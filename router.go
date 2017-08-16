// Package marmoset ,reinventing the wheel
package marmoset

import (
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// NewRouter ...
func NewRouter() *Router {
	return &Router{
		routes:     map[HTTPMethod]map[string]http.HandlerFunc{},
		regexps:    map[HTTPMethod]map[*regexp.Regexp]http.HandlerFunc{},
		subrouters: []*Router{},
	}
}

// HTTPMethod ...
type HTTPMethod string

const (
	// MethodGet ...
	MethodGet HTTPMethod = http.MethodGet
	// MethodPost ...
	MethodPost HTTPMethod = http.MethodPost
	// MethodAny represents any method type
	MethodAny HTTPMethod = "*"
)

const pathParameterExpressionString = "\\(\\?P\\<[^>]+\\>[^)]+\\)"

// Resolver only resolves
type Resolver interface {
	FindHandler(*http.Request) (http.HandlerFunc, bool)
}

// Router ...
type Router struct {
	static     *static
	routes     map[HTTPMethod]map[string]http.HandlerFunc
	regexps    map[HTTPMethod]map[*regexp.Regexp]http.HandlerFunc
	subrouters []*Router
	subrouter  *Router
	resolver   Resolver
	filters    []Filterable
}

type static struct {
	Path   string
	Server http.Handler
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if router.static != nil && strings.HasPrefix(req.URL.Path, router.static.Path) {
		router.static.Server.ServeHTTP(w, req)
		return
	}

	handlerFunc, ok := router.FindHandler(req)
	if ok && handlerFunc != nil {
		router.handleFilteredFunc(handlerFunc, w, req)
		return
	}

	http.NotFound(w, req)
}

// FindHandler ...
func (router *Router) FindHandler(req *http.Request) (http.HandlerFunc, bool) {

	// Ask child first
	for _, child := range router.subrouters {
		handler, ok := child.FindHandler(req)
		if ok && handler != nil {
			return func(w http.ResponseWriter, req *http.Request) {
				child.handleFilteredFunc(handler, w, req)
			}, ok
		}
	}

	// Ask subroute
	if handlers, ok := router.routes[MethodAny]; ok {
		for key, handler := range handlers {
			if strings.HasPrefix(req.URL.Path, key) {
				// Match!!
				rest := strings.Replace(req.URL.Path, key, "", -1)
				if !strings.HasPrefix(rest, "/") {
					rest = "/" + rest
				}
				return handler, true
			}
		}
	}

	// Ask normal routes
	if handlers, ok := router.routes[HTTPMethod(req.Method)]; ok {
		if handler, ok := handlers[req.URL.Path]; ok {
			return handler, true
		}
	}

	// Ask regex routes
	if handlers, ok := router.regexps[HTTPMethod(req.Method)]; ok {
		for exp, handler := range handlers {
			if exp.MatchString(req.URL.Path) {
				matched := exp.FindAllStringSubmatch(req.URL.Path, -1)
				if req.Form == nil {
					req.Form = url.Values{}
				}
				for i, name := range exp.SubexpNames() {
					req.Form.Add(name, matched[0][i])
				}
				return handler, true
			}
		}
	}
	return nil, false
}

// add ...
func (router *Router) add(method HTTPMethod, path string, handler http.HandlerFunc) *Router {
	if ok, compiled := isRegexpPath(path); ok {
		return router.addRegexpRoute(method, path, compiled, handler)
	}
	if _, ok := router.routes[method]; !ok {
		router.routes[method] = map[string]http.HandlerFunc{}
	}
	if _, ok := router.routes[method][path]; ok {
		log.Fatalf("route duplicated on `%s %s`", method, path)
	}
	router.routes[method][path] = handler
	return router
}

// isRegexpPath ...
func isRegexpPath(path string) (bool, *regexp.Regexp) {
	compiled, err := regexp.Compile(path)
	if err != nil || compiled == nil {
		return false, nil
	}
	pathParameterExpression := regexp.MustCompile(pathParameterExpressionString)

	all := strings.Split(path, "/")
	for i, part := range all {
		if pathParameterExpression.MatchString(part) {
			// If the expression is the last part of URL path,
			// The end expression should be appended in Regexp.
			if i == len(all)-1 {
				compiled = regexp.MustCompile(compiled.String() + "$")
			}
			return true, compiled
		}
	}
	return false, nil
}

// addRegexpRoute ...
func (router *Router) addRegexpRoute(method HTTPMethod, path string, pathCompiled *regexp.Regexp, handler http.HandlerFunc) *Router {
	if _, ok := router.regexps[method]; !ok {
		router.regexps[method] = map[*regexp.Regexp]http.HandlerFunc{}
	}
	router.regexps[method][pathCompiled] = handler
	return router
}

// GET ...
func (router *Router) GET(path string, handler http.HandlerFunc) *Router {
	return router.add("GET", path, handler)
}

// POST ...
func (router *Router) POST(path string, handler http.HandlerFunc) *Router {
	return router.add("POST", path, handler)
}

// Handle ...
func (router *Router) Handle(path string, handler http.Handler) *Router {
	return router.add("*", path, handler.ServeHTTP)
}

// Subrouter ...
func (router *Router) Subrouter(child *Router) *Router {
	router.subrouters = append(router.subrouters, child)
	return router
}

// Static ...
func (router *Router) Static(p, dir string) *Router {

	if !filepath.IsAbs(dir) {
		_, f, _, _ := runtime.Caller(1)
		dir = path.Join(path.Dir(f), dir)
	}

	fs := http.FileServer(http.Dir(dir))
	router.static = &static{
		Path:   p,
		Server: http.StripPrefix(p, fs),
	}
	return router
}

// StaticRelative ...
func (router *Router) StaticRelative(p string, relativeDir string) *Router {
	// TODO: is this needed?? ;)
	_, f, _, _ := runtime.Caller(1)
	return router.Static(p, path.Join(path.Dir(f), relativeDir))
}

// Apply applies Filters
func (router *Router) Apply(filters ...Filterable) error {
	router.filters = filters
	return nil
}

func (router *Router) handleFilteredFunc(fn http.HandlerFunc, w http.ResponseWriter, req *http.Request) {
	if len(router.filters) == 0 {
		fn(w, req)
		return
	}
	// TODO: Fix
	// TODO: This if fuxxxn' buggy to indirect function to single refrence for slice
	// filters := []Filterable{}
	// copy(filters, router.filters)
	handler := http.HandlerFunc(fn)
	router.filters[len(router.filters)-1].SetNext(handler)
	chained := chainFilters(router.filters)
	chained.ServeHTTP(w, req)
	return
}

func chainFilters(filters []Filterable) Filterable {
	if len(filters) == 1 {
		return filters[0]
	}
	filters[len(filters)-2].SetNext(filters[len(filters)-1])
	return chainFilters(filters[:len(filters)-1])
}

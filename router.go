// Package marmoset ,reinventing the wheel
package marmoset

import (
	"log"
	"net/http"
	"path"
	"runtime"
	"strings"
)

// NewRouter ...
func NewRouter() *Router {
	return &Router{
		routes:   map[string]map[string]http.HandlerFunc{},
		notfound: http.NotFound,
	}
}

// Router ...
type Router struct {
	static   *static
	routes   map[string]map[string]http.HandlerFunc
	notfound http.HandlerFunc
}

type static struct {
	Path   string
	Server http.Handler
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if router.static != nil && strings.HasPrefix(r.URL.Path, router.static.Path) {
		router.static.Server.ServeHTTP(w, r)
		return
	}
	methodes, ok := router.routes[r.Method]
	if !ok {
		router.notfound(w, r)
		return
	}
	handler, ok := methodes[r.URL.Path]
	if !ok {
		router.notfound(w, r)
		return
	}
	handler.ServeHTTP(w, r)
}

// add ...
func (router *Router) add(method string, path string, handler http.HandlerFunc) *Router {
	if _, ok := router.routes[method]; !ok {
		router.routes[method] = map[string]http.HandlerFunc{}
	}
	if _, ok := router.routes[method][path]; ok {
		log.Fatalf("route duplicated on `%s %s`", method, path)
	}
	router.routes[method][path] = handler
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

// NotFound ...
func (router *Router) NotFound(handler http.HandlerFunc) *Router {
	router.notfound = handler
	return router
}

// Static ...
func (router *Router) Static(p string, dir string) *Router {
	fs := http.FileServer(http.Dir(dir))
	router.static = &static{
		Path:   p,
		Server: http.StripPrefix(p, fs),
	}
	return router
}

// StaticRelative ...
func (router *Router) StaticRelative(p string, relativeDir string) *Router {
	_, f, _, _ := runtime.Caller(1)
	return router.Static(p, path.Join(path.Dir(f), relativeDir))
}

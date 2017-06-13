package marmoset

import "net/http"

// Filterable represents struct which can be a filter
type Filterable interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	SetNext(http.Handler)
}

// Filter ...
// Remember "Last added, First called"
type Filter struct {
	http.Handler
	Next http.Handler
}

// SetNext ...
func (f *Filter) SetNext(next http.Handler) {
	f.Next = next
}

// NewFilter ...
// func NewFilter(root *Router) *FilteredChain {
// 	return &FilteredChain{
// 		root:    root,
// 		current: root,
// 	}
// }

// Add ...
// func (chain *FilteredChain) Add(filter http.Handler) *FilteredChain {
// 	v := reflect.ValueOf(filter)
// 	switch v.Kind() {
// 	case reflect.Interface, reflect.Ptr:
// 		// pass
// 	default:
// 		log.Fatalf("type `%s` is not addressable", v.Type().String())
// 		// return chain
// 	}
// 	if !v.Elem().FieldByName("Next").CanSet() {
// 		log.Fatalf("type `%s` must have `Next` field", v.Type().String())
// 		// return chain
// 	}
// 	v.Elem().FieldByName("Next").Set(reflect.ValueOf(chain.current))
// 	chain.current = filter
// 	return chain
// }
//
// // Router ...
// func (chain *FilteredChain) Router() *Router {
// 	child := NewRouter()
// 	child.Handle("/", chain.current)
// 	router := &Router{
// 		resolver:   chain.root,
// 		subrouters: []*Router{child},
// 	}
// 	return router
// }

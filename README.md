marmoset
========

[![Go](https://github.com/otiai10/marmoset/actions/workflows/go.yml/badge.svg)](https://github.com/otiai10/marmoset/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/otiai10/marmoset/branch/main/graph/badge.svg?token=jQrl6xAFNa)](https://codecov.io/gh/otiai10/marmoset)
[![Go Report Card](https://goreportcard.com/badge/github.com/otiai10/marmoset)](https://goreportcard.com/report/github.com/otiai10/marmoset)
[![Maintainability](https://api.codeclimate.com/v1/badges/ba1ab060a31d7ff00c8a/maintainability)](https://codeclimate.com/github/otiai10/marmoset/maintainability)
[![GoDoc](https://godoc.org/github.com/otiai10/marmoset?status.svg)](https://godoc.org/github.com/otiai10/marmoset)

less than "web framework", just make your code a bit DRY.

```go
func main() {
	router := marmoset.NewRouter()
	router.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello!"))
	})
	http.ListenAndServe(":8080", router)
}
```

# Features

- Path parameters
- Static directory
- Request filters
	- Request-Context accessor, if you want

```go
func main() {

	router := marmoset.NewRouter()

	router.GET("/foo", your.FooHttpHandlerFunc)
	router.POST("/bar", your.BarHttpHandlerFunc)

	// Path parameters available with regex like declaration
	router.GET("/users/(?P<name>[a-zA-Z0-9]+)/hello", func(w http.ResponseWriter, r *http.Request) {
		marmoset.Render(w).HTML("hello", map[string]string{
			// Path parameters can be accessed by req.FromValue()
			"message": fmt.Printf("Hello, %s!", r.FormValue("name")),
		})
	})

	// Set static file path
	router.Static("/public", "/your/assets/path")

	// Last added filter will be called first.
	server := marmoset.NewFilter(router).
		Add(&your.Filter{}).
		Add(&your.AnotherFilter{}).
		Add(&marmoset.ContextFilter{}). // if you want
		Server()

	http.ListenAndServe(":8080", server)
}
```

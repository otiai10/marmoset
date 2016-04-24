marmoset
========

less than "web framework", just make your code a bit DRY.

```go
func main() {

	r := marmoset.NewRouter()

	r.GET("/foo", your.FooHttpHandlerFunc)
	r.POST("/bar", your.BarHttpHandlerFunc)

	// Use path parameters
	r.GET("/users/(?P<name>[a-zA-Z0-9]+)/hello", func(w http.ResponseWriter, req *http.Request) {
		marmoset.Render(w).HTML("hello", map[string]string{
			// Path parameters can be accessed by req.FromValue()
			"name": req.FormValue("name"),
		})
	})

	// Set static file path
	r.Static("/public", "/your/assets/path")

	s := marmoset.NewFilter(r).Add(&your.Filter{}).Server()

	http.ListenAndServe(":8080", s)
}
```

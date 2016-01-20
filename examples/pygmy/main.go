package main

import (
	"net/http"

	"github.com/otiai10/marmoset"
)

func main() {
	r := marmoset.NewRouter()
	r.GET("/api", func(w http.ResponseWriter, r *http.Request) {
		marmoset.RenderJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Hello, this is pygmy marmoset!",
		})
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		marmoset.RenderJSON(w, http.StatusNotFound, map[string]interface{}{
			"message": "not found :(",
		})
	})
	http.ListenAndServe(":8080", r)
}

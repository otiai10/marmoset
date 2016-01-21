package main

import (
	"net/http"

	"github.com/otiai10/marmoset"
)

func main() {

	marmoset.LoadViews("./")

	r := marmoset.NewRouter()

	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		marmoset.Render(w).HTML("index", map[string]interface{}{
			"message": "Hello, this is pygmy marmoset!",
		})
	})
	r.GET("/api", func(w http.ResponseWriter, r *http.Request) {
		marmoset.Render(w).JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello, this is pygmy marmoset API!",
		})
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		marmoset.RenderJSON(w, http.StatusNotFound, map[string]interface{}{
			"message": "not found :(",
		})
	})

	r.StaticRelative("/public", "./")

	http.ListenAndServe(":8080", r)
}

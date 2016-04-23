package example

import (
	"net/http"

	m "github.com/otiai10/marmoset"
)

func init() {

	m.LoadViews("./views")

	r := m.NewRouter()

	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		m.Render(w).HTML("index", nil)
	})
	r.GET("/api", func(w http.ResponseWriter, r *http.Request) {
		m.Render(w).JSON(http.StatusOK, map[string]interface{}{
			"message":     "Hello, this is pygmy marmoset API!",
			"remote_addr": r.RemoteAddr,
		})
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		m.RenderJSON(w, http.StatusNotFound, map[string]interface{}{
			"message": "not found :(",
		})
	})

	r.StaticRelative("/public", "./assets")

	http.Handle("/", r)
	// ListenAndServe will be invoked by GAE SDK.
	// http.ListenAndServe(":"+os.Getenv("PORT"), r)
}

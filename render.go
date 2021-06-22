package marmoset

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Renderer ...
type Renderer struct {
	Pretty     bool
	EscapeHTML bool
	writer     http.ResponseWriter
}

// Render ...
func Render(w http.ResponseWriter, pretty ...bool) Renderer {
	pretty = append(pretty, true)
	return Renderer{
		writer:     w,
		Pretty:     pretty[0],
		EscapeHTML: true,
	}
}

// JSON ...
func (r Renderer) JSON(status int, data interface{}) error {

	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(status)

	enc := json.NewEncoder(r.writer)
	if r.Pretty {
		enc.SetIndent("", "\t")
	}
	enc.SetEscapeHTML(r.EscapeHTML)
	err := enc.Encode(data)

	return err
}

// RenderJSON ...
func RenderJSON(w http.ResponseWriter, status int, data interface{}) error {
	return Render(w, true).JSON(status, data)
}

// HTML ...
func (r Renderer) HTML(name string, data interface{}) error {
	if templates == nil {
		return fmt.Errorf("templates not loaded")
	}
	r.writer.Header().Add("Content-Type", "text/html")
	return templates.Lookup(name).Execute(r.writer, data)
}

package marmoset

import (
	"encoding/json"
	"net/http"
)

// Renderer ...
type Renderer interface {
	JSON(http.ResponseWriter, int, interface{}) error
}

// PrettyRenderer ...
type PrettyRenderer struct{}

// JSON ...
func (pr PrettyRenderer) JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

// RenderJSON ...
func RenderJSON(w http.ResponseWriter, status int, data interface{}) error {
	return PrettyRenderer{}.JSON(w, status, data)
}

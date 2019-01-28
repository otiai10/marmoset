package marmoset

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/otiai10/mint"
)

func TestRender(t *testing.T) {
	w := httptest.NewRecorder()
	renderer := Render(w, true)
	Expect(t, renderer).TypeOf("marmoset.Renderer")
}

func TestRenderer_EscapeHTML(t *testing.T) {

	message := "<h1>Test</h1>"

	w := httptest.NewRecorder()
	r := Render(w, false)
	err := r.JSON(http.StatusOK, P{"html": message})
	Expect(t, err).ToBe(nil)
	b, err := ioutil.ReadAll(w.Body)
	Expect(t, err).ToBe(nil)
	Expect(t, bytes.Trim(b, "\n")).ToBe([]byte(`{"html":"\u003ch1\u003eTest\u003c/h1\u003e"}`))

	Because(t, "EscapeHTML can be disabled", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := Render(w, false)
		r.EscapeHTML = false
		err := r.JSON(http.StatusOK, P{"html": message})
		Expect(t, err).ToBe(nil)
		b, err := ioutil.ReadAll(w.Body)
		Expect(t, err).ToBe(nil)
		Expect(t, bytes.Trim(b, "\n")).ToBe([]byte(`{"html":"<h1>Test</h1>"}`))
	})

}

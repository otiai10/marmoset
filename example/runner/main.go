package main

import (
	"net/http"

	_ "github.com/otiai10/marmoset/example"
)

func main() {
	http.ListenAndServe(":8080", nil)
}

// +build !appengine

package marmoset

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// AuthFilter can detect request header and define current request user.
// For example, you can restore User model from auth-token in request header.
type AuthFilter struct {
	Filter
}

// ServeHTTP will be called before the root router's ServeHTTP.
func (f *AuthFilter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	user := new(User)
	err := json.NewDecoder(strings.NewReader(req.Header.Get("X-User"))).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if user.ID == 20 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ctx := Context().Get(req)
	Context().Set(req, context.WithValue(ctx, "user", user))

	f.Next.ServeHTTP(w, req)
}

var SampleController = func(w http.ResponseWriter, r *http.Request) {
	user := Context().Get(r).Value("user").(*User)
	json.NewEncoder(w).Encode(map[string]interface{}{"request_user": user})
}

func ExampleFilter() {

	// Define routings.
	router := NewRouter()
	router.GET("/test", SampleController)

	// Add Filters.
	// If you want to use `Context`, `ContextFilter` must be added for the last.
	// Remember "Last added, First called"
	router.Apply(new(AuthFilter))

	// Use `http.ListenAndServe` in real case, instead of httptest.
	server := httptest.NewServer(router)

	req, _ := http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Add("X-User", `{"id":10,"name":"otiai10"}`)
	res, _ := http.DefaultClient.Do(req)
	fmt.Println("StatusCode:", res.StatusCode)

	req, _ = http.NewRequest("GET", server.URL+"/test", nil)
	req.Header.Add("X-User", `{"id":20,"name":"otiai20"}`)
	res, _ = http.DefaultClient.Do(req)
	fmt.Println("StatusCode:", res.StatusCode)

	// Output:
	// StatusCode: 200
	// StatusCode: 403
}

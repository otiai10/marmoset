package marmoset

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/otiai10/mint"
)

func TestNewRouter(t *testing.T) {
	Expect(t, NewRouter()).TypeOf("*marmoset.Router")
}

func TestRouter_GET(t *testing.T) {

	router := NewRouter()
	router.GET("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is /foo handler"))
	})
	server := httptest.NewServer(router)

	res, err := http.Get(server.URL + "/foo")
	Expect(t, err).ToBe(nil)
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	Expect(t, err).ToBe(nil)
	Expect(t, string(b)).ToBe("This is /foo handler")
}

func TestRouter_POST(t *testing.T) {

	p := make([]byte, 2)

	router := NewRouter()
	router.POST("/bar", func(w http.ResponseWriter, r *http.Request) {
		r.Body.Read(p)
		r.Body.Close()
		w.Write([]byte("This is /bar handler"))
	})
	server := httptest.NewServer(router)

	res, err := http.Post(server.URL+"/bar", "text/plain", strings.NewReader("Hi"))
	Expect(t, err).ToBe(nil)
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	Expect(t, err).ToBe(nil)
	Expect(t, string(b)).ToBe("This is /bar handler")
	Expect(t, string(p)).ToBe("Hi")
}

func TestRouter_Subrouter(t *testing.T) {

	sub1 := NewRouter()
	sub1.GET("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is /foo handler from subrouter 1"))
	})
	sub2 := NewRouter()
	sub2.GET("/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is /bar handler from subrouter 2"))
	})
	root := NewRouter()
	root.Subrouter(sub1)
	root.Subrouter(sub2)
	server := httptest.NewServer(root)

	res, err := http.Get(server.URL + "/bar")
	Expect(t, err).ToBe(nil)
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	Expect(t, err).ToBe(nil)
	Expect(t, string(b)).ToBe("This is /bar handler from subrouter 2")

	res, err = http.Get(server.URL + "/foo")
	Expect(t, err).ToBe(nil)
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	Expect(t, err).ToBe(nil)
	Expect(t, string(b)).ToBe("This is /foo handler from subrouter 1")

}

type MyPlainFilter struct {
	count int
	Filter
}

func (f *MyPlainFilter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f.count++
	f.Next.ServeHTTP(w, req)
}

func TestNewFilter_Add(t *testing.T) {

	instance := new(MyPlainFilter)
	root := NewRouter()
	root.POST("/", func(w http.ResponseWriter, req *http.Request) {})
	err := root.Apply(instance)
	Expect(t, err).ToBe(nil)
	server := httptest.NewServer(root)

	http.Post(server.URL, "text/plain", nil)
	Expect(t, instance.count).ToBe(1)
	http.Post(server.URL, "text/plain", nil)
	Expect(t, instance.count).ToBe(2)
}

type SampleUser struct {
	Name string `json:"name"`
}
type MyAuthFilter struct {
	context map[*http.Request]*SampleUser
	Filter
}

func (f *MyAuthFilter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if f.context == nil {
		f.context = map[*http.Request]*SampleUser{}
	}
	user := new(SampleUser)
	if err := json.NewDecoder(req.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if user.Name == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	f.context[req] = user
	f.Next.ServeHTTP(w, req)
}

func TestRouter_Apply(t *testing.T) {

	requesters := map[string]*SampleUser{}
	instance := new(MyAuthFilter)

	authorized := NewRouter()
	authorized.POST("/user/1", func(w http.ResponseWriter, req *http.Request) {
		requesters["/user/1"] = instance.context[req]
	})
	authorized.POST("/user/2", func(w http.ResponseWriter, req *http.Request) {
		requesters["/user/2"] = instance.context[req]
	})
	authorized.Apply(instance)

	unauthorized := NewRouter()
	unauthorized.POST("/welcome", func(w http.ResponseWriter, req *http.Request) {
		requesters["/welcome"] = instance.context[req]
	})
	unauthorized.POST("/login", func(w http.ResponseWriter, req *http.Request) {
		requesters["/login"] = instance.context[req]
	})

	root := NewRouter()
	root.Subrouter(authorized)
	root.Subrouter(unauthorized)
	server := httptest.NewServer(root)

	var res *http.Response
	var err error

	res, err = http.Post(server.URL+"/user/1", "application/json", nil)
	Expect(t, err).ToBe(nil)
	Expect(t, res.StatusCode).ToBe(http.StatusForbidden)

	res, err = http.Post(server.URL+"/user/1", "application/json", strings.NewReader(`{"name":"otiai10"}`))
	Expect(t, err).ToBe(nil)
	Expect(t, res.StatusCode).ToBe(http.StatusOK)

	res, err = http.Post(server.URL+"/welcome", "application/json", nil)
	Expect(t, err).ToBe(nil)
	Expect(t, res.StatusCode).ToBe(http.StatusOK)

	res, err = http.Post(server.URL+"/user/2", "application/json", nil)
	Expect(t, err).ToBe(nil)
	Expect(t, res.StatusCode).ToBe(http.StatusForbidden)

	res, err = http.Post(server.URL+"/user/2", "application/json", strings.NewReader(`{"name":"otiai20"}`))
	Expect(t, err).ToBe(nil)
	Expect(t, res.StatusCode).ToBe(http.StatusOK)

	res, err = http.Post(server.URL+"/login", "application/json", nil)
	Expect(t, err).ToBe(nil)
	Expect(t, res.StatusCode).ToBe(http.StatusOK)

}

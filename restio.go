package restio

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Initializer ...
func Initializer(fs ...interface{}) error {
	for _, f := range fs {
		if e := dynamicFunc(f); e != nil {
			return e
		}
	}
	return nil
}

// Run the app on it's router
func Run(r http.Handler, host string) {
	Panic("%s", http.ListenAndServe(host, r))
}

// Endpoint ...
func Endpoint(i interface{}) http.Handler {
	return &Handler{F: i}
}

// Params ...
func Params(r *http.Request) map[string]string {
	m := mux.Vars(r)
	for k, v := range r.URL.Query() {
		m[k] = v[0]
	}
	return m
}

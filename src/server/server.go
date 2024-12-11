package server

import (
	"fmt"
	"net/http"
	"strings"
)

type Server struct {
	Addr   string
	Routes []Route
}

type Route struct {
	URI     string
	Methods []Method
	Handler http.HandlerFunc
}

// Http method
type Method int

const (
	GET    Method = 1
	POST   Method = 2
	PUT    Method = 4
	DELETE Method = 8
	ANY    Method = GET | POST | PUT | DELETE
)

func (s Server) ListenAndServe() {
	for _, route := range s.Routes {
		http.HandleFunc(route.URI, func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI != route.URI || !containsMethod(route.Methods, r.Method) {
				fmt.Printf("[SERVER] Request failed at %s, method: %s\n", r.RequestURI, r.Method)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			route.Handler(w, r)
		})
	}

	http.ListenAndServe(s.Addr, nil)
}

func getMethodFromString(method string) (Method, error) {
	s := strings.ToUpper(method)
	switch {
	case s == "" || s == "GET":
		return GET, nil
	case s == "POST":
		return POST, nil
	case s == "PUT":
		return PUT, nil
	case s == "DELETE":
		return DELETE, nil
	}

	return -1, fmt.Errorf("invalid method as input: %s", method)
}

func getFlagValue(methods []Method) int {
	value := 0
	for _, method := range methods {
		value = value | int(method)
	}

	return value
}

func containsMethod(methods []Method, method string) bool {
	m, err := getMethodFromString(method)
	if err != nil {
		return false
	}

	return getFlagValue(methods)&int(m) == int(m)
}

package sql

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
	Methods []HttpMethod
	Handler http.HandlerFunc
}

// Http method
type HttpMethod int

const (
	HTTP_GET    HttpMethod = 1
	HTTP_POST   HttpMethod = 2
	HTTP_PUT    HttpMethod = 4
	HTTP_DELETE HttpMethod = 8
	HTTP_ALL    HttpMethod = HTTP_GET | HTTP_POST | HTTP_PUT | HTTP_DELETE
)

func (s *Server) ListenAndServe() {
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

func getMethodFromString(method string) (HttpMethod, error) {
	s := strings.ToUpper(method)
	switch {
	case s == "" || s == "GET":
		return HTTP_GET, nil
	case s == "POST":
		return HTTP_POST, nil
	case s == "PUT":
		return HTTP_PUT, nil
	case s == "DELETE":
		return HTTP_DELETE, nil
	}

	return -1, fmt.Errorf("invalid method as input: %s", method)
}

func GetFlag(methods []HttpMethod) int {
	value := 0
	for _, method := range methods {
		value = value | int(method)
	}

	return value
}

func containsMethod(methods []HttpMethod, method string) bool {
	m, err := getMethodFromString(method)
	if err != nil {
		return false
	}

	flag := GetFlag(methods)
	return flag|int(m) == flag
}

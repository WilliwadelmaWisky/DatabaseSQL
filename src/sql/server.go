package sql

import (
	"fmt"
	"net/http"
	"strings"
)

// Represents an http server
type Server struct {
	Addr   string
	Routes []Route
}

// Represents a single route on the server
type Route struct {
	URI     string
	Methods []HttpMethod
	Handler http.HandlerFunc
}

// Enum to represent an http method, values are prefixed with HTTP.
// Enum is flagged so many methods can be combined by bitwise or |
type HttpMethod int

const (
	// Represents an http get request
	HTTP_GET HttpMethod = 1
	// Represents an http post request
	HTTP_POST HttpMethod = 2
	// Represents an http put request
	HTTP_PUT HttpMethod = 4
	// Represents an http delete request
	HTTP_DELETE HttpMethod = 8
	// Represents any http request
	HTTP_ALL HttpMethod = HTTP_GET | HTTP_POST | HTTP_PUT | HTTP_DELETE
)

// Starts a server
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

// Get an http method from a string.
// Empty string is considered a GET request
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

// Calculates a flag value of http method array
func GetFlag(methods []HttpMethod) int {
	value := 0
	for _, method := range methods {
		value = value | int(method)
	}

	return value
}

// Check if method is contained in the flag
func containsMethod(methods []HttpMethod, method string) bool {
	m, err := getMethodFromString(method)
	if err != nil {
		return false
	}

	flag := GetFlag(methods)
	return flag|int(m) == flag
}

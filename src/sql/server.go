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
	URI        string
	MethodFlag HttpMethod
	Handler    http.HandlerFunc
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
			if r.RequestURI != route.URI || !route.IsAllowedMethodString(r.Method) {
				fmt.Printf("[SERVER] Request failed at %s, method: %s\n", r.RequestURI, r.Method)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			route.Handler(w, r)
		})
	}

	http.ListenAndServe(s.Addr, nil)
}

// Check if route allows http request of a certain method
func (route *Route) IsAllowedMethod(method HttpMethod) bool {
	return route.MethodFlag|method == route.MethodFlag
}

// Check if route allows http requests of a certain method string (not casesensitive).
// Empty string is interpreted as a GET request
func (route *Route) IsAllowedMethodString(method string) bool {
	m, err := GetMethod(method)
	if err != nil {
		return false
	}

	return route.IsAllowedMethod(m)
}

// Get an http method from a string (not casesensitive).
// Empty string is interpreted as a GET request
func GetMethod(method string) (HttpMethod, error) {
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

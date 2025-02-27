package sql

import (
	"fmt"
	"net/http"
	"strings"
)

// Represents an http server
type Server struct {
	Addr   string  // Server address eg. localhost:8080
	Routes []Route // Server routes
}

// Represents a single route on the server
type Route struct {
	URI        string           // Route path eg. /path/to/resource
	MethodFlag HttpMethod       // Http methods allowed in the route, can be multiple
	Handler    http.HandlerFunc // Route handler function
}

// Enum to represent an http method, values are prefixed with HTTP.
// Enum is flagged so many methods can be combined by bitwise or |
type HttpMethod int

const (
	HTTP_GET    HttpMethod = 1                                             // Represents an http get request
	HTTP_POST   HttpMethod = 2                                             // Represents an http post request
	HTTP_PUT    HttpMethod = 4                                             // Represents an http put request
	HTTP_DELETE HttpMethod = 8                                             // Represents an http delete request
	HTTP_ALL    HttpMethod = HTTP_GET | HTTP_POST | HTTP_PUT | HTTP_DELETE // Represents any http request
)

// Starts a server
func (s *Server) ListenAndServe() {
	for _, route := range s.Routes {
		http.HandleFunc(route.URI, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
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
func (route *Route) IsAllowedMethodString(s string) bool {
	method, err := GetMethod(s)
	if err != nil {
		return false
	}

	return route.IsAllowedMethod(method)
}

// Get an http method from a string (not casesensitive).
// Empty string is interpreted as a GET request
func GetMethod(s string) (HttpMethod, error) {
	switch strings.ToUpper(s) {
	case "", "GET":
		return HTTP_GET, nil
	case "POST":
		return HTTP_POST, nil
	case "PUT":
		return HTTP_PUT, nil
	case "DELETE":
		return HTTP_DELETE, nil
	}

	return -1, fmt.Errorf("invalid method as input: %s", s)
}

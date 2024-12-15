package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/WilliwadelmaWisky/DatabaseSQL/sql"
)

// Main function
func main() {
	port := 9000
	if len(os.Args) > 1 {
		value, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		port = value
	}

	database := &sql.Database{}

	server := &sql.Server{
		Addr: fmt.Sprintf("localhost:%d", port),
		Routes: []sql.Route{
			{URI: "/", Methods: []sql.HttpMethod{sql.HTTP_POST}, Handler: func(w http.ResponseWriter, r *http.Request) {
				serverHandler(w, r, database)
			}},
		},
	}

	fmt.Printf("Server starting at http://localhost:%d/\nPress <CTRL+C> to terminate!\n", port)
	server.ListenAndServe()
}

// h
//
// # Arguments
//   - responseWriter - Hello
//   - request - Request
//   - database - Database
func serverHandler(responseWriter http.ResponseWriter, request *http.Request, database *sql.Database) {
	bytes, _ := io.ReadAll(request.Body)
	fmt.Printf("[SQL]: %s\n", string(bytes))

	tokens := sql.Tokenize(bytes)
	if len(tokens) == 0 {
		fmt.Print("No tokens received from request\n")
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	// TOKEN DEBUGGING
	// values := util.Map(tokens, func(token sql.Token) string {
	// 	return fmt.Sprintf("'%s'", token.Value)
	// })
	// fmt.Printf("Tokens found: %s\n", strings.Join(values, " "))

	operation, err := sql.Parse(tokens)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	result := operation.Call(database)
	if result != nil {
		responseWriter.Header().Add("Content-Length", strconv.Itoa(len(result)))
		responseWriter.Header().Add("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write(result)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}

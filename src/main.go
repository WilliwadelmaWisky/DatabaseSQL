package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/WilliwadelmaWisky/DatabaseSQL/sql"
	"github.com/WilliwadelmaWisky/DatabaseSQL/util"
)

// Main function
func main() {
	database := sql.Database{}

	server := sql.Server{
		Addr: "localhost:9000",
		Routes: []sql.Route{
			{URI: "/", Methods: []sql.HttpMethod{sql.HTTP_POST}, Handler: func(w http.ResponseWriter, r *http.Request) {
				serverHandler(w, r, &database)
			}},
		},
	}

	fmt.Println("Server starting at http://localhost:9000/\nPress <CTRL+C> to terminate!")
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
	tokens := sql.Tokenize(bytes)
	values := util.Map(tokens, func(token sql.Token) string {
		return fmt.Sprintf("'%s'", token.Value)
	})

	fmt.Printf("Request SQL: %s\n", string(bytes))
	fmt.Printf("Tokens found: %s\n", strings.Join(values, " "))

	if len(tokens) == 0 {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	operation, _ := sql.Parse(tokens)
	result := operation.Call(database)

	responseWriter.Write(result)
	responseWriter.WriteHeader(http.StatusOK)
}

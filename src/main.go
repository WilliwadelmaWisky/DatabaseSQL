package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/WilliwadelmaWisky/DatabaseSQL/server"
	"github.com/WilliwadelmaWisky/DatabaseSQL/sql"
	"github.com/WilliwadelmaWisky/DatabaseSQL/util"
)

// Main function
func main() {
	s := server.Server{
		Addr: "localhost:9000",
		Routes: []server.Route{
			{URI: "/", Methods: []server.Method{server.POST}, Handler: serverHandler},
		},
	}

	fmt.Println("Server starting at http://localhost:9000/\nPress <CTRL+C> to terminate!")
	s.ListenAndServe()
}

// h
//
// # Arguments
//   - w - Hello
//   - r - Request
func serverHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := io.ReadAll(r.Body)
	tokens := sql.Tokenize(bytes)
	values := util.Map(tokens, func(token sql.Token) string {
		return fmt.Sprintf("'%s'", token.Value)
	})

	fmt.Printf("Request SQL: %s\n", string(bytes))
	fmt.Printf("Tokens found: %s\n", strings.Join(values, " "))

	if len(tokens) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//sql.Evaluate(tokens)
	sql.Parse(tokens)
	w.WriteHeader(http.StatusOK)
}

type Table struct {
	Table   string   `json:"table"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Column string        `json:"column"`
	Values []interface{} `json:"values"`
}

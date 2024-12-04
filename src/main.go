package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/WilliwadelmaWisky/DatabaseSQL/nlp"
	"github.com/WilliwadelmaWisky/DatabaseSQL/util"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.RequestURI != "/" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bytes, _ := io.ReadAll(r.Body)
		tokens := nlp.Tokenize(bytes)
		values := util.Map(tokens, func(token nlp.Token) string {
			return fmt.Sprintf("'%s'", token.Value)
		})

		fmt.Printf("Request SQL: %s\n", string(bytes))
		fmt.Printf("Tokens found: %s\n", strings.Join(values, " "))

		if len(tokens) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		nlp.Evaluate(tokens)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Server starting at http://localhost:9000\nPress <CTRL+C> to terminate!")
	http.ListenAndServe("localhost:9000", nil)
}

type Table struct {
	Table   string   `json:"table"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Column string        `json:"column"`
	Values []interface{} `json:"values"`
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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
		tokens := tokenize(bytes)

		fmt.Printf("Request SQL: %s\n", string(bytes))
		fmt.Printf("Tokens found: %s\n", strings.Join(tokens, ", "))

		if len(tokens) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch tokens[0] {
		case "CREATE":
			w.WriteHeader(http.StatusOK)
		case "SELECT":
			evaluateSelect(tokens[1:])
			w.WriteHeader(http.StatusOK)
		case "UPDATE":
			w.WriteHeader(http.StatusOK)
		case "DELETE":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	fmt.Println("Server starting at http://localhost:9000\nPress <CTRL+C> to terminate!")
	http.ListenAndServe("localhost:9000", nil)
}

func tokenize(b []byte) []string {
	start := 0
	tokens := make([]string, 0)

	for i := 0; i < len(b); i++ {
		if (!isAlphaNumeric(b[i]) && !isSpecial(b[i])) || i == len(b)-1 {
			// End of token
			if start != i {
				token := string(b[start:i])
				tokens = append(tokens, token)
			}

			start = i + 1
		}
	}

	return tokens
}

func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func isSpecial(c byte) bool {
	return c == '=' || c == '(' || c == ')' || c == '*'
}

func evaluateSelect(tokens []string) {
	columns := make([]string, 0)

	for _, token := range tokens {
		if token == columns[0] {
			break
		}
	}
}

type Table struct {
	Table   string   `json:"table"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Column string        `json:"column"`
	Values []interface{} `json:"values"`
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/WilliwadelmaWisky/DatabaseSQL/sql"
)

// Main function
func main() {
	if len(os.Args) <= 1 {
		fmt.Println("No arguments were given to the script! Use --help to get more information.")
		return
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Print help")
		return
	}

	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Println("Print version")
		return
	}

	url := os.Args[1]
	port := 9000
	if len(os.Args) > 2 {
		value, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		port = value
	}

	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(append([]string{homeDir, ".WilliwadelmaWisky", "DatabaseSQL"}, strings.Split(url, "/")...)...)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	database := sql.NewDatabase(path)
	err = database.Load()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	server := &sql.Server{
		Addr: fmt.Sprintf("localhost:%d", port),
		Routes: []sql.Route{
			{
				URI:        "/",
				MethodFlag: sql.HTTP_POST,
				Handler:    func(w http.ResponseWriter, r *http.Request) { sqlRequestHandler(w, r, database) },
			},
			{
				URI:        "/information_schema",
				MethodFlag: sql.HTTP_GET,
				Handler:    func(w http.ResponseWriter, r *http.Request) { informationSchemaRequestHandler(w, r, database) },
			},
		},
	}

	fmt.Printf("Server starting at http://localhost:%d/\nPress <CTRL+C> to terminate!\n", port)
	server.ListenAndServe()
}

// HttpServer request handler for sql requests
func sqlRequestHandler(w http.ResponseWriter, r *http.Request, database *sql.Database) {
	bytes, _ := io.ReadAll(r.Body)
	fmt.Printf("[SQL]: %s\n", string(bytes))

	tokens := sql.Tokenize(bytes)
	if len(tokens) == 0 {
		fmt.Print("[ERROR]: No tokens received from request\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TOKEN DEBUGGING
	// values := sql.Map(tokens, func(token *sql.Token) string {
	// 	return fmt.Sprintf("'%s'", token.Value)
	// })
	// fmt.Printf("Tokens found: %s\n", strings.Join(values, " "))

	operation, err := sql.Parse(tokens)
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := operation.Call(database)
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if result != nil {
		w.Header().Add("Content-Length", strconv.Itoa(len(result)))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HttpServer request handler for information_schema requests
func informationSchemaRequestHandler(w http.ResponseWriter, r *http.Request, database *sql.Database) {
	informationSchema := sql.NewInformationSchema(database)
	bytes, err := json.Marshal(informationSchema)

	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if bytes != nil {
		w.Header().Add("Content-Length", strconv.Itoa(len(bytes)))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
		return
	}

	w.WriteHeader(http.StatusOK)
}

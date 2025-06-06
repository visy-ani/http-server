package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")              // logging client message
	io.WriteString(w, "This is my Website!\n") // sending response
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")  // logging client message
	io.WriteString(w, "Hello, HTTP!\n") // sending response
}

func main() {
	// seting up Multiplexer and register handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	// set router
	// http.HandleFunc("/", getRoot)
	// http.HandleFunc("/hello", getHello)

	// starting server and if anything goes wrong store error
	// err := http.ListenAndServe(":3333", nil) // without mux
	err := http.ListenAndServe(":3333", mux) // with mux

	// check if server closed
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

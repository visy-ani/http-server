package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
	// fmt.Printf("got / request\n")              // logging client message
	io.WriteString(w, "This is my Website!\n") // sending response
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	// fmt.Printf("got /hello request\n")  // logging client message
	io.WriteString(w, "Hello, HTTP!\n") // sending response
}

// func main() {
// 	// seting up Multiplexer and register handlers
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", getRoot)
// 	mux.HandleFunc("/hello", getHello)

// 	// set router
// 	// http.HandleFunc("/", getRoot)
// 	// http.HandleFunc("/hello", getHello)

// 	// starting server and if anything goes wrong store error
// 	// err := http.ListenAndServe(":3333", nil) // without mux
// 	err := http.ListenAndServe(":3333", mux) // with mux

// 	// check if server closed
// 	if errors.Is(err, http.ErrServerClosed) {
// 		fmt.Printf("server closed\n")
// 	} else if err != nil {
// 		fmt.Printf("error starting server: %s\n", err)
// 		os.Exit(1)
// 	}
// }

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	ctx, cancelCtx := context.WithCancel(context.Background())

	serveOne := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	serveTwo := &http.Server{
		Addr:    ":4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := serveOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Println("error starting server one: ", err)
		}
		cancelCtx()
	}()

	go func() {
		err := serveTwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server two closed\n")
		} else if err != nil {
			fmt.Println("error starting server two: ", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}

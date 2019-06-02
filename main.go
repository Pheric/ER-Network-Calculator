package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port int

func main() {
	setupClFlags()

	mux := http.NewServeMux()
	mux.Handle("/stylesheets/", http.StripPrefix("/stylesheets", http.FileServer(http.FileSystem(http.Dir("./www/stylesheets")))))
	mux.Handle("/hello", http.HandlerFunc(serveExamplePage))
	mux.Handle("/", http.HandlerFunc(serveIndex))

	fmt.Printf("Now listening on port %d!\n", port)
	log.Fatalf("error while listening on port %d: %v\nExiting...\n", port, http.ListenAndServe(fmt.Sprintf("127.1:%d", port), mux))
}

func setupClFlags() {
	flag.IntVar(&port, "port", 2500, "the port to run the website on")

	flag.Parse()
}

func serveExamplePage(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello, World!")); err != nil {
		log.Printf("error serving example page: %v\n", err)
	}
}

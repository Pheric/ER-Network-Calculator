package main

import (
	"flag"
	"fmt"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"log"
	"net/http"
)

var port int

func main() {
	setupClFlags()
	setupReadMe()

	mux := http.NewServeMux()
	mux.Handle("/stylesheets/", http.StripPrefix("/stylesheets", http.FileServer(http.FileSystem(http.Dir("./www/stylesheets")))))
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.FileSystem(http.Dir("./www/assets")))))
	mux.Handle("/hello", http.HandlerFunc(ServeReadMe))
	mux.Handle("/", http.HandlerFunc(serveIndex))

	fmt.Printf("Now listening on port %d!\n", port)
	log.Fatalf("error while listening on port %d: %v\nExiting...\n", port, http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func setupClFlags() {
	flag.IntVar(&port, "port", 2500, "the port to run the website on")

	flag.Parse()
}

var readme []byte
func setupReadMe() {
	md, err := ioutil.ReadFile("./readme.md")
	if err != nil {
		fmt.Printf("error reading readme.md: %v\n", err)
		return
	}

	readme = blackfriday.Run(md)
}

func ServeReadMe(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write(readme); err != nil {
		fmt.Printf("error writing readme.md: %v\n", err)
	}
}
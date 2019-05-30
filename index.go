package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./www/index.html", "./www/templates/header.html")
	if err != nil {
		e := fmt.Sprintf("error parsing index page: %v", err)
		if _, err = w.Write([]byte("503 Internal Server Error")); err != nil {
			e = fmt.Sprintf("error writing error to index page:\ninitial error:%s\ncurrent error: %v", e, err)
		}

		log.Println(e)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		log.Printf("error executing index template: %v\n", err)
	}
}
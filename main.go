package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {

	port := ":8080"

	fs := http.FileServer(http.Dir("./public/"))
	http.Handle("/public/", http.StripPrefix("/public", fs))

	http.HandleFunc("/", indexHandler)

	fmt.Printf("listening on %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index.html").ParseFiles("public/templates/header.tmpl", "public/html/index.html"))
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

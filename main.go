package main

import (
	"html/template"
	"net/http"
)

//var templates = template.Must()

func main() {
	http.HandleFunc("/", index)

	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	html := template.Must(template.ParseFiles("public/index.html"))
	html.ExecuteTemplate(res,"index.html",nil)
}

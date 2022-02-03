package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	type ViewData struct {
		Title   string
		Message string
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := ViewData{
			Title:   "Our dream cars",
			Message: "Tesla will come",
		}
		tmpl, _ := template.ParseFiles("templates/index.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			return
		}
	})

	fmt.Println("Server is listening...")
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		return
	}
}

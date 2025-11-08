package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "page not found", 404)
		return
	}

	tmpl, err := template.ParseFiles("./frontend/index.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "./frontend" + r.URL.Path
	// Serve the file directly
	http.ServeFile(w, r, path)
}

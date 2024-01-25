package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/jouer", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "jouer.html")
	})

	http.HandleFunc("/options", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "options.html")
	})

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "style.css")
	})

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

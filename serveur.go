package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strings"
    "dylan/play"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
    http.HandleFunc("/hangman", hangmanHandler)
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/guess", guessHandler)
    http.HandleFunc("/jouer", jouerHandler)
    http.HandleFunc("/abandon", abandonHandler)
    http.HandleFunc("/style.css", styleHandler)

    fs := http.FileServer(http.Dir("pictures"))
    http.Handle("/pictures/", http.StripPrefix("/pictures/", fs))

    fmt.Println("Server started on port 8080")
    http.ListenAndServe(":8080", nil)
}

func hangmanHandler(w http.ResponseWriter, r *http.Request) {
    if play.GetGame() == nil {
        play.InitGame()
    }
    game := play.GetGame()
    fmt.Println("Page refreshed successfully")
    err := templates.ExecuteTemplate(w, "hangman.html", game)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    err := templates.ExecuteTemplate(w, "index.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func jouerHandler(w http.ResponseWriter, r *http.Request) {
    play.InitGame()
    err := templates.ExecuteTemplate(w, "jouer.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
    letter := strings.ToLower(r.FormValue("letter"))
    play.GuessLetter(letter)
    http.Redirect(w, r, "/hangman", http.StatusFound)
}

func abandonHandler(w http.ResponseWriter, r *http.Request) {
    play.ResetInitialWord()
    http.Redirect(w, r, "/", http.StatusFound)
}

func styleHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/style.css")
}
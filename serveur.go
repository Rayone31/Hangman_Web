package main

import (
    "fmt"
    "html/template"
    "dylan/play"
    "net/http"
    "strings"
)

var templates = template.Must(template.ParseGlob("template/*.html"))

func main() {
    http.HandleFunc("/hangman", hangmanHandler)
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/guess", guessHandler)
    http.HandleFunc("/jouer", jouerHandler)
    http.HandleFunc("/options", optionsHandler)
    http.HandleFunc("/abandon", abandonHandler) // Nouvelle route pour abandonner la partie
    http.HandleFunc("/style.css", styleHandler)

    fs := http.FileServer(http.Dir("pictures"))
    http.Handle("/pictures/", http.StripPrefix("/pictures/", fs))

    fmt.Println("Server started on port 8080")
    http.ListenAndServe(":8080", nil)
}

func hangmanHandler(w http.ResponseWriter, r *http.Request) {
    play.InitGame()                                           // Initialise le jeu
    game := play.GetGame()
    err := templates.ExecuteTemplate(w, "hangman.html", game) // Passe l'objet game à la template
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
    err := templates.ExecuteTemplate(w, "jouer.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
    err := templates.ExecuteTemplate(w, "options.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func abandonHandler(w http.ResponseWriter, r *http.Request) {
    // Réinitialiser le mot initialisé lorsque le joueur abandonne la partie
    play.ResetInitialWord()
    http.Redirect(w, r, "/", http.StatusFound)
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
    letter := strings.ToLower(r.FormValue("letter"))
    play.GuessLetter(letter)
    http.Redirect(w, r, "/hangman", http.StatusFound)
}

func styleHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/style.css")
}

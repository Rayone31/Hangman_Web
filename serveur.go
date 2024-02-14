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
    http.HandleFunc("/finish", StatusHandler)
    http.HandleFunc("/difficulty", difficultyHandler)

    fs := http.FileServer(http.Dir("assets"))
    http.Handle("/assets/", http.StripPrefix("/assets/", fs))

    fmt.Println("Server started on port 8080")
    http.ListenAndServe(":8080", nil)
}

func hangmanHandler(w http.ResponseWriter, r *http.Request) {
    if play.GetGame() == nil {
        // Si le jeu n'est pas initialisé, rediriger vers la page de sélection de difficulté
        http.Redirect(w, r, "/jouer", http.StatusFound)
        return
    }

    game := play.GetGame()
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
    play.ResetGame()
    http.Redirect(w, r, "/", http.StatusFound)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
    play.ResetGame()
    http.Redirect(w, r, "/", http.StatusFound)
}

func difficultyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        err := templates.ExecuteTemplate(w, "difficulty.html", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == "POST" {
        difficulty := r.FormValue("difficulty")
        language := r.FormValue("language")
        if difficulty == "" || language == "" {
            http.Error(w, "Veuillez sélectionner une difficulté et une langue", http.StatusBadRequest)
            return
        }
        play.InitGame(difficulty, language)
        http.Redirect(w, r, "/hangman", http.StatusFound)
    }
}

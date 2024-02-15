package main

import (
    "fmt"
    "html/template"
    "net/http"
    "sort"
    "dylan/play"
    "strconv"
    "strings"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
    http.HandleFunc("/hangman", hangmanHandler)
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/guess", guessHandler)
    http.HandleFunc("/jouer", jouerHandler)
    http.HandleFunc("/abandon", abandonHandler)
    http.HandleFunc("/victoire", victoryHandler)
    http.HandleFunc("/defaite", defeatHandler)
    http.HandleFunc("/difficulty", difficultyHandler)
    http.HandleFunc("/scoreboard", scoreboardHandler)
    http.HandleFunc("/sauvegarder", saveScoreHandler)

    fs := http.FileServer(http.Dir("assets"))
    http.Handle("/assets/", http.StripPrefix("/assets/", fs))

    fmt.Println("Server started on port 8080")
    http.ListenAndServe(":8080", nil)
}

func hangmanHandler(w http.ResponseWriter, r *http.Request) {
    if play.GetGame() == nil {
        http.Redirect(w, r, "/jouer", http.StatusFound)
        return
    }

    game := play.GetGame()

    if game.GameStatus == "victory" {
        http.Redirect(w, r, "/victoire", http.StatusFound)
        return
    }

    if game.GameStatus == "game over" {
        http.Redirect(w, r, "/defaite", http.StatusFound)
        return
    }

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

func victoryHandler(w http.ResponseWriter, r *http.Request) {
    game := play.GetGame()
    if game == nil {
        http.Redirect(w, r, "/", http.StatusFound)
        return
    }

    data := struct {
        Score int
    }{
        Score: game.Score,
    }

    err := templates.ExecuteTemplate(w, "victoire.html", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func defeatHandler(w http.ResponseWriter, r *http.Request) {
    err := templates.ExecuteTemplate(w, "defaite.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
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

func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
    scores := play.GetScores()
    sortedScores := sortScores(scores)

    showNameInput := false
    if len(sortedScores) < 8 || sortedScores[len(sortedScores)-1].Value < play.GetGame().Score {
        showNameInput = true
    }

    data := struct {
        Scores       []play.Score
        ShowNameInput bool
        NewScore      int
    }{
        Scores:       sortedScores,
        ShowNameInput: showNameInput,
        NewScore:     play.GetGame().Score,
    }

    err := templates.ExecuteTemplate(w, "scoreboard.html", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func sortScores(scores []play.Score) []play.Score {
    sort.Slice(scores, func(i, j int) bool {
        return scores[i].Value > scores[j].Value
    })
    return scores
}

func saveScoreHandler(w http.ResponseWriter, r *http.Request) {
    playerName := r.FormValue("player_name")
    score, _ := strconv.Atoi(r.FormValue("score"))
    if playerName == "" {
        playerName = "Inconnu"
    }
    play.AddScore(play.Score{PlayerName: playerName, Value: score})
    http.Redirect(w, r, "/scoreboard", http.StatusFound)
}

// Package main est le point d'entrée de l'application.
package main

import (
    "fmt"
    "html/template"
    "net/http"
    "sort"
    "dylan/play" // Importation du package play pour la fonctionnalité du jeu.
    "strconv"
    "strings"
)

// Variable globale pour contenir les modèles HTML.
var templates = template.Must(template.ParseGlob("templates/*.html"))

// La fonction main configure les routes HTTP et démarre le serveur.
func main() {
    // Configuration des routes HTTP pour différents gestionnaires.
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

    // Servir les fichiers statiques depuis le répertoire "assets".
    fs := http.FileServer(http.Dir("assets"))
    http.Handle("/assets/", http.StripPrefix("/assets/", fs))

    // Démarrage du serveur.
    fmt.Println("Serveur démarré sur le port 8080")
    http.ListenAndServe(":8080", nil)
}

// hangmanHandler gère les demandes pour jouer au jeu du Pendu.
func hangmanHandler(w http.ResponseWriter, r *http.Request) {
    // Si aucun jeu n'est en cours, redirige vers la page de départ.
    if play.GetGame() == nil {
        http.Redirect(w, r, "/jouer", http.StatusFound)
        return
    }

    game := play.GetGame()

    // Redirige vers la page de victoire si le jeu est gagné.
    if game.GameStatus == "victory" {
        http.Redirect(w, r, "/victoire", http.StatusFound)
        return
    }

    // Redirige vers la page de défaite si le jeu est perdu.
    if game.GameStatus == "game over" {
        http.Redirect(w, r, "/defaite", http.StatusFound)
        return
    }

    // Rend la page du jeu du Pendu.
    err := templates.ExecuteTemplate(w, "hangman.html", game)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// indexHandler gère les demandes pour la page d'accueil.
func indexHandler(w http.ResponseWriter, r *http.Request) {
    // Rend la page d'accueil.
    err := templates.ExecuteTemplate(w, "index.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// jouerHandler gère les demandes pour commencer un nouveau jeu.
func jouerHandler(w http.ResponseWriter, r *http.Request) {
    // Rend la page pour commencer un nouveau jeu.
    err := templates.ExecuteTemplate(w, "jouer.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// guessHandler gère les demandes pour deviner une lettre dans le jeu du Pendu.
func guessHandler(w http.ResponseWriter, r *http.Request) {
    // Obtenir la lettre devinée du formulaire et la convertir en minuscules.
    letter := strings.ToLower(r.FormValue("letter"))
    // Traiter la lettre devinée dans la logique du jeu.
    play.GuessLetter(letter)
    // Rediriger vers la page du jeu du Pendu.
    http.Redirect(w, r, "/hangman", http.StatusFound)
}

// abandonHandler gère les demandes pour abandonner le jeu en cours.
func abandonHandler(w http.ResponseWriter, r *http.Request) {
    // Réinitialiser l'état du jeu.
    play.ResetGame()
    // Rediriger vers la page d'accueil.
    http.Redirect(w, r, "/", http.StatusFound)
}

// victoryHandler gère les demandes pour afficher la page de victoire.
func victoryHandler(w http.ResponseWriter, r *http.Request) {
    // Obtenir l'état du jeu actuel.
    game := play.GetGame()
    if game == nil {
        // Si aucun jeu n'est en cours, rediriger vers la page d'accueil.
        http.Redirect(w, r, "/", http.StatusFound)
        return
    }

    // Préparer les données pour le rendu de la page de victoire.
    data := struct {
        Score int
    }{
        Score: game.Score,
    }

    // Rend la page de victoire.
    err := templates.ExecuteTemplate(w, "victoire.html", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// defeatHandler gère les demandes pour afficher la page de défaite.
func defeatHandler(w http.ResponseWriter, r *http.Request) {
    // Rend la page de défaite.
    err := templates.ExecuteTemplate(w, "defaite.html", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// difficultyHandler gère les demandes pour définir la difficulté du jeu.
func difficultyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Rend la page pour choisir la difficulté.
        err := templates.ExecuteTemplate(w, "difficulty.html", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    } else if r.Method == "POST" {
        // Gérer la soumission du formulaire pour commencer un nouveau jeu avec la difficulté choisie.
        difficulty := r.FormValue("difficulty")
        language := r.FormValue("language")
        if difficulty == "" || language == "" {
            http.Error(w, "Veuillez sélectionner une difficulté et une langue", http.StatusBadRequest)
            return
        }
        play.InitGame(difficulty, language)
        // Rediriger vers la page du jeu du Pendu.
        http.Redirect(w, r, "/hangman", http.StatusFound)
    }
}

// scoreboardHandler gère les demandes pour afficher le tableau de bord.
func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
    // Obtenir les scores du jeu.
    scores := play.GetScores()
    // Trier les scores.
    sortedScores := sortScores(scores)

    // Déterminer s'il faut afficher le champ de saisie du nom en fonction du score.
    showNameInput := false
    if len(sortedScores) < 8 || sortedScores[len(sortedScores)-1].Value < play.GetGame().Score {
        showNameInput = true
    }

    // Préparer les données pour le rendu du tableau de bord.
    data := struct {
        Scores       []play.Score
        ShowNameInput bool
        NewScore      int
    }{
        Scores:       sortedScores,
        ShowNameInput: showNameInput,
        NewScore:     play.GetGame().Score,
    }

    // Rend la page du tableau de bord.
    err := templates.ExecuteTemplate(w, "scoreboard.html", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// sortScores trie les scores par ordre décroissant.
func sortScores(scores []play.Score) []play.Score {
    sort.Slice(scores, func(i, j int) bool {
        return scores[i].Value > scores[j].Value
    })
    return scores
}

// saveScoreHandler gère les demandes pour enregistrer le score du joueur.
func saveScoreHandler(w http.ResponseWriter, r *http.Request) {
    // Obtenir le nom du joueur et le score du formulaire.
    playerName := r.FormValue("player_name")
    score, _ := strconv.Atoi(r.FormValue("score"))
    // Si aucun nom n'est fourni, utiliser "Inconnu".
    if playerName == "" {
        playerName = "Inconnu"
    }
    // Ajouter le score au tableau de bord.
    play.AddScore(play.Score{PlayerName: playerName, Value: score})
    // Rediriger vers la page du tableau de bord.
    http.Redirect(w, r, "/scoreboard", http.StatusFound)
}

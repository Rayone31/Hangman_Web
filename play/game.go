// Définition du package
package play

// Importation des packages nécessaires
import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "strings"
    "time"
)

// Définition de la structure Game qui représente une partie de jeu
type Game struct {
    Word            string
    PartialWord     string
    GuessedLetters  []string
    LivesRemaining  int
    GameStatus      string
    ProposedLetters []string
    HangmanImage    string
    Score           int
    Streak          int
}

// Définition de la structure Score qui représente le score d'un joueur
type Score struct {
    PlayerName string
    Value      int
}

// Déclaration des variables globales
var (
    game   Game
    scores []Score
)

// InitGame initialise une nouvelle partie de jeu
func InitGame(difficulty string, language string) {
    // Sélection du fichier de mots en fonction de la difficulté et de la langue
    var filePath string
    switch difficulty {
    case "easy":
        filePath = "Ressources/" + language + "_words_easy.txt"
    case "medium":
        filePath = "Ressources/" + language + "_words_medium.txt"
    case "hard":
        filePath = "Ressources/" + language + "_words_hard.txt"
    default:
        filePath = "Ressources/french_words_easy.txt" 
    }

    // Obtention d'un mot aléatoire à partir du fichier
    word := getRandomWordFromFile(filePath)

    // Initialisation de la partie de jeu
    game = Game{
        Word:            word,
        PartialWord:     strings.Repeat("-", len(word)),
        GuessedLetters:  []string{},
        LivesRemaining:  10,
        GameStatus:      "",
        ProposedLetters: []string{},
        HangmanImage:    "/assets/Hangman_0.png",
    }
}

// GetGame renvoie la partie de jeu actuelle
func GetGame() *Game {
    return &game
}

// GuessLetter traite une lettre proposée par le joueur
func GuessLetter(letter string) {
    // Si le jeu est terminé, on ne fait rien
    if game.GameStatus != "" {
        return
    }

    // Ajout de la lettre proposée à la liste des lettres proposées
    game.ProposedLetters = append(game.ProposedLetters, letter)

    // Si la lettre a déjà été devinée, on ne fait rien
    for _, guessedLetter := range game.GuessedLetters {
        if guessedLetter == letter {
            return
        }
    }

    // Ajout de la lettre à la liste des lettres devinées
    game.GuessedLetters = append(game.GuessedLetters, letter)

    // Si la lettre n'est pas dans le mot, on décrémente le nombre de vies restantes
    // et on met à jour l'image du pendu
    if !strings.Contains(game.Word, letter) {
        game.LivesRemaining--
        game.HangmanImage = fmt.Sprintf("/assets/Hangman_%d.png", 10-game.LivesRemaining)
        game.Streak = 0
    } else {
        // Si la lettre est dans le mot, on augmente le score et la série
        game.Score += 100
        game.Streak++
        if game.Streak > 1 {
            game.Score *= game.Streak
        }
    }

    // Mise à jour du mot partiellement deviné
    game.PartialWord = updatePartialWord(letter)

    // Si le mot a été entièrement deviné, on déclare une victoire
    // Sinon, si le nombre de vies restantes est nul, on déclare une défaite
    if game.PartialWord == game.Word {
        game.GameStatus = "victory"
        return
    } else if game.LivesRemaining == 0 {
        game.GameStatus = "game over"
        return
    }
}

// updatePartialWord met à jour le mot partiellement deviné avec la nouvelle lettre
func updatePartialWord(letter string) string {
    updatedWord := ""
    for i, char := range game.Word {
        if strings.ContainsRune(letter, char) {
            updatedWord += string(char)
        } else {
            updatedWord += string(game.PartialWord[i])
        }
    }
    return updatedWord
}

// getRandomWordFromFile obtient un mot aléatoire à partir d'un fichier
func getRandomWordFromFile(filePath string) string {
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var words []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        words = append(words, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }

    rand.Seed(time.Now().UnixNano())
    return words[rand.Intn(len(words))]
}

// ResetGame réinitialise la partie de jeu
func ResetGame() {
    game = Game{}
    game.Score = 0
}

// AddScore ajoute un score à la liste des scores
func AddScore(score Score) {
    scores = append(scores, score)
}

// GetScores renvoie la liste des scores
func GetScores() []Score {
    return scores
}
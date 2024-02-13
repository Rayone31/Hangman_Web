package play

import (
    "bufio"
    "math/rand"
    "os"
    "strings"
    "time"
    "fmt"
)

type Game struct {
    Word           string
    PartialWord    string
    GuessedLetters []string
    LivesRemaining int
    GameStatus     string
    ProposedLetters []string
    HangmanImage    string
    Score              int 
    Streak             int
}

var (
    initialWord string
    game        Game
)

func InitGame(difficulty string) {
    var filePath string
    switch difficulty {
    case "easy":
        filePath = "Ressources/french_words_easy.txt"
    case "medium":
        filePath = "Ressources/french_words_medium.txt"
    case "hard":
        filePath = "Ressources/french_words_hard.txt"
    default:
        filePath = "Ressources/french_words_easy.txt" // Par défaut, choisir facile
    }

    word := getRandomWordFromFile(filePath)

    game = Game{
        Word:           word,
        PartialWord:    strings.Repeat("-", len(word)),
        GuessedLetters: []string{},
        LivesRemaining: 10,
        GameStatus:     "",
        ProposedLetters: []string{},
        HangmanImage:   "/assets/Hangman_0.png",
    }
}

func GetGame() *Game {
    return &game
}

func GuessLetter(letter string) {
    if game.GameStatus != "" {
        ResetGame()
        return
    }

    // Ajouter la lettre proposée à la liste des lettres proposées
    game.ProposedLetters = append(game.ProposedLetters, letter)

    for _, guessedLetter := range game.GuessedLetters {
        if guessedLetter == letter {
            return
        }
    }

    game.GuessedLetters = append(game.GuessedLetters, letter)

    if !strings.Contains(game.Word, letter) {
        game.LivesRemaining--
        game.HangmanImage = fmt.Sprintf("/assets/Hangman_%d.png", 10-game.LivesRemaining)
        game.Streak = 0 // Réinitialiser la série si le joueur se trompe
    } else {
        // Si la lettre est correcte, augmenter le score et mettre à jour la série
        game.Score += 100
        game.Streak++
        if game.Streak > 1 {
            game.Score *= game.Streak // Si le joueur trouve des lettres d'affilée, doubler le score
        }
    }

    game.PartialWord = updatePartialWord(letter)

    if game.PartialWord == game.Word {
        game.GameStatus = "victory"
    } else if game.LivesRemaining == 0 {
        game.GameStatus = "game over"
    }
}

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

func ResetGame() {
    initialWord = ""
    game = Game{}
}

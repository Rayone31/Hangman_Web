package play

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "strings"
    "time"
)

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

type Score struct {
    PlayerName string
    Value      int
}

var (
    game   Game
    scores []Score
)

func InitGame(difficulty string, language string) {
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

    word := getRandomWordFromFile(filePath)

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

func GetGame() *Game {
    return &game
}

func GuessLetter(letter string) {
    if game.GameStatus != "" {
        return
    }

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
        game.Streak = 0
    } else {
        game.Score += 100
        game.Streak++
        if game.Streak > 1 {
            game.Score *= game.Streak
        }
    }

    game.PartialWord = updatePartialWord(letter)

    if game.PartialWord == game.Word {
        game.GameStatus = "victory"
        return
    } else if game.LivesRemaining == 0 {
        game.GameStatus = "game over"
        return
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
    game = Game{}
    game.Score = 0
}

func AddScore(score Score) {
    scores = append(scores, score)
}

func GetScores() []Score {
    return scores
}

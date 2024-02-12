package play

import (
    "bufio"
    "math/rand"
    "os"
    "strings"
    "time"
)

type Game struct {
    Word           string
    PartialWord    string
    GuessedLetters []string
    LivesRemaining int
    GameStatus     string
    ProposedLetters []string
}

var game Game

func InitGame() {
    // Choix aléatoire du fichier de mots
    files := []string{"Ressources/french_words1.txt", "Ressources/french_words2.txt", "Ressources/french_words3.txt"}
    rand.Seed(time.Now().UnixNano())
    fileIndex := rand.Intn(len(files))
    filePath := files[fileIndex]

    // Lecture du fichier et choix aléatoire d'un mot
    word := getRandomWordFromFile(filePath)

    game = Game{
        Word:           initialWord,
        PartialWord:    strings.Repeat("-", len(initialWord)),
        GuessedLetters: []string{},
        LivesRemaining: 10,
        GameStatus:     "",
        ProposedLetters: []string{},
    }
}

func GetGame() *Game {
    return &game
}

func GuessLetter(letter string) {
    if game.GameStatus != "" {
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
    }

    game.PartialWord = updatePartialWord(letter)
}

// Fonction pour mettre à jour le mot partiel avec la lettre proposée si elle est correcte
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

// Fonction pour réinitialiser le mot initialisé
func ResetInitialWord() {
    initialWord = ""
}

func ResetInitialWord() {
    initialWord = ""
}
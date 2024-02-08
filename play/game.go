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
}

var game Game

// Fonction pour initialiser le jeu
func InitGame() {
    // Choix aléatoire du fichier de mots
    files := []string{"Ressources/french_words1.txt", "Ressources/french_words2.txt", "Ressources/french_words3.txt"}
    rand.Seed(time.Now().UnixNano())
    fileIndex := rand.Intn(len(files))
    filePath := files[fileIndex]

    // Lecture du fichier et choix aléatoire d'un mot
    word := getRandomWordFromFile(filePath)

    // Initialisation du jeu avec le mot choisi
    game = Game{
        Word:           word,
        PartialWord:    strings.Repeat("-", len(word)),
        GuessedLetters: []string{},
        LivesRemaining: 10,
        GameStatus:     "",
    }
}

// Fonction pour obtenir l'objet de jeu actuel
func GetGame() Game {
    return game
}

// Fonction pour proposer une lettre
func GuessLetter(letter string) {
    // Vérifier si le jeu est terminé
    if game.GameStatus != "" {
        return
    }

    // Vérifier si la lettre a déjà été devinée
    for _, guessedLetter := range game.GuessedLetters {
        if guessedLetter == letter {
            return
        }
    }

    // Ajouter la lettre à la liste des lettres devinées
    game.GuessedLetters = append(game.GuessedLetters, letter)

    // Réduire le nombre de vies restantes si la lettre proposée n'est pas dans le mot
    if !strings.Contains(game.Word, letter) {
        game.LivesRemaining--
    }

    // Mettre à jour le mot partiel avec la lettre proposée si elle est correcte
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

// Fonction pour obtenir un mot aléatoire à partir d'un fichier
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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const exit = "STOP"
const startAttempts = 10
const saveFilename = "save.txt"

var hasWin bool

func main() {
	var attempts int
	var letters []string
	var submittedLetters []string
	var word string
	word, attempts, letters, submittedLetters = initGame()

	if os.Args[2] == "--startsWith" {
		if len(os.Args) > 3 {
			newSaveFilename := os.Args[3]
			var err error
			attempts, letters, submittedLetters, word, err = recoverFromSave(newSaveFilename)
			if err != nil {
				word, attempts, letters, submittedLetters = initGame()
			}
		}
	}


	fmt.Printf("Good Luck, you have %v attempts.\n", attempts)
	for {
		printWord(letters)
		submission, doExit := getLetter()
		if doExit {
			break
		}

		if submission == word {
			win()
			break
		}

		if len(submission) == 1 {
			if !isLetter(rune(submission[0])) {
				fmt.Println("You can only submit letters")
				continue
			}

			if isLetterSubmitted(submission, submittedLetters) {
				fmt.Println("You already submitted this letter")
				continue
			}

			submittedLetters = append(submittedLetters, submission)
			if !isLetterInWord(submission, word) {
				attempts--
				fmt.Printf("Not present in the word, %v attempts remaining\n", attempts)
				continue
			}
		} else {
			attempts -= 2
		}

		for i := 0; i < len(word); i++ {
			if string(word[i]) == submission {
				letters[i] = submission
			}
		}
		if isWordGuessed(letters, word) {
			win()
			break
		}
	}

	if !hasWin {
		fmt.Printf("Saving data in %v\n", saveFilename)
		data := Data{
			Attempts:         attempts,
			ActualWord:       strings.Join(letters, ""),
			LettersSubmitted: submittedLetters,
			Word:             word,
		}
		if submittedLetters == nil {
			data.LettersSubmitted = []string{}
		}

		err := data.SaveInJSONFile(saveFilename)
		if err != nil {
			fmt.Printf("Error saving data: %v\nSave is empty.", err)
		}
	}
}

func chooseWordFromFile(selectedFile string) string {
	file, err := ioutil.ReadFile(selectedFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("File does not exist")
		} else {
			log.Fatalf("Error reading file : %v", err)
		}
	}
	split := strings.Split(string(file), "\n")
	randIndex := rand.Intn(len(split))
	return strings.ToUpper(split[randIndex])
}

func getLetter() (result string, doExit bool) {
	fmt.Print("Choose: ")
	letter := strings.ToUpper(readLine())
	return letter, letter == exit
}

func initGame() (word string, attempts int, letters []string, submittedLetters []string) {
	rand.Seed(time.Now().Unix())
	word = chooseWordFromFile(os.Args[1])
	lettersToReveal := len(word)/2 - 1
	return word, startAttempts, setVisibleLetters(word, []string{}, lettersToReveal), []string{}
}

func isLetter(letter rune) bool {
	return letter >= 'A' && letter <= 'Z'
}

func isLetterInWord(letter string, word string) bool {
	for i := 0; i < len(word); i++ {
		if string(word[i]) == letter {
			return true
		}
	}
	return false
}

func isLetterSubmitted(submission string, letter []string) bool {
	for _, l := range letter {
		if l == submission {
			return true
		}
	}
	return false
}

func isWordGuessed(letters []string, word string) bool {
	for i := 0; i < len(word); i++ {
		if string(word[i]) != letters[i] {
			return false
		}
	}
	return true
}

func printWord(letters []string) {
	fmt.Println(strings.Join(letters, " "))
}

func setVisibleLetters(word string, letters []string, lettersToReveal int) []string {
	for i := 0; i < len(word); i++ {
		letters = append(letters, "_")
	}

	for i := 0; i < lettersToReveal; i++ {
		index := rand.Intn(len(word))
		letters[index] = string(word[index])
	}
	return letters
}

func readLine() string {
	var line string
	_, err := fmt.Scanln(&line)
	if err != nil {
		log.Fatalf("Error reading line : %v", err)
	}
	return line
}

func recoverFromSave(saveFilename string) (attempts int, letters, submittedLetters []string, word string, err error) {
	file, err := ioutil.ReadFile(saveFilename)
	if err != nil {
		log.Fatal(err)
	}

	data := Data{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Cannot read save file, data of this game will recreated.")
		return startAttempts, []string{}, []string{}, "", err
	}

	fmt.Println("Game recovered from save.")
	return data.Attempts, strings.Split(data.ActualWord, ""), data.LettersSubmitted, data.Word, nil
}

func win() {
	hasWin = true
	fmt.Printf("You won!")
}
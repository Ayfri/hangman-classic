package main

import (
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
	rand.Seed(time.Now().Unix())
	word := chooseWordFromFile(os.Args[1])

	var letters []string
	var submittedLetters []string
	lettersToReveal := len(word)/2 - 1
	attempts := startAttempts
	letters = setVisibleLetters(word, letters, lettersToReveal)

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

		if len(submission) == 1 && isLetter(submission) {
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

func printWord(letters []string) {
	println(strings.Join(letters, " "))
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

func isLetter(letter string) bool {
	return letter >= "a" && letter <= "z"
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

func win() {
	hasWin = true
	fmt.Printf("You won!")
}

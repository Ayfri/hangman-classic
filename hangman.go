package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	exit          = "STOP"
	startAttempts = 10
)

var asciiLettersFile = flag.String("letterFile", "", "Use ASCII letters from file")
var asciiLetters = map[rune]string{}
var hasWin bool
var saveFilename = flag.String("startWith", "", "File to save the game")

func main() {
	var attempts int
	var err error
	var letters []string
	var submittedLetters []string
	var word string
	flag.Parse()

	if len(os.Args) > 2 {
		attempts, letters, submittedLetters, word, err = recoverFromSave(*saveFilename)
		if err != nil {
			word, attempts, letters, submittedLetters = initGame()
		} else {
			fmt.Printf("Welcome back, you have %v attempts remaining.\n", attempts)
		}
	} else {
		word, attempts, letters, submittedLetters = initGame()
		fmt.Printf("Good Luck, you have %v attempts.\n", attempts)
	}

	if *saveFilename == "" {
		*saveFilename = "save.txt"
	}

	if *asciiLettersFile != "" {
		for i := 'A'; i <= 'Z'; i++ {
			asciiLetters[i] = getAsciiLetter(*asciiLettersFile, i - ('a' - 'A'))
		}
		asciiLetters['_'] = getAsciiLetter(*asciiLettersFile, 63)
	}

	for {
		if attempts <= 0 {
			fmt.Println("You lost, the word was:", word)
			break
		}

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
				fmt.Println(getHangmanPosition(9 - attempts))
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

		err := data.SaveInJSONFile(*saveFilename)
		if err != nil {
			fmt.Printf("Error saving data: %v\nSave is empty.", err)
		}
	}
}

func chooseWordFromFile(selectedFile string) string {
	content, err := readFile(selectedFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		win()
		return ""
	}
	split := strings.Split(content, "\n")
	randIndex := rand.Intn(len(split))
	return strings.ToUpper(split[randIndex])
}

func getHangmanPosition(position int) string {
	file, err := ioutil.ReadFile("resources/hangman.txt")
	if err != nil {
		log.Fatalf("Error reading hangman file : %v", err)
	}
	content := string(file)
	split := strings.Split(content, "\n\n")
	return split[position]
}

func getAsciiLetter(file string, letter rune) string {
	content, err := readFile(file)
	if err != nil {
		log.Fatalf("Error reading file : %v", err)
	}
	split := strings.Split(content, "\n\n")
	return split[letter]
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
	if len(asciiLetters) > 1 {
		str := strings.Join(letters, " ") + "\n"
		var array []string
		for i := range str {
			array = append(array, asciiLetters[rune(str[i])])
		}
		result := ""
		limit := len(strings.Split(array[0], "\n"))
		for i := 0; i < limit; i++ {
			for _, letter := range array {
				for index, line := range strings.Split(letter, "\n") {
					if index == i {
						result += line + "  "
					}
				}
			}
			result += "\n"
		}

		fmt.Println(result)
	} else {
		fmt.Println(strings.Join(letters, " ") + "\n")
	}
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

func readFile(selectedFile string) (string, error) {
	file, err := ioutil.ReadFile(selectedFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("File does not exist")
		} else {
			log.Fatalf("Error reading file : %v", err)
		}
		return "", err
	}
	return string(file), nil
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
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Save file %v has not been found, a new game will start.\n", saveFilename)
			return startAttempts, []string{}, []string{}, "", err
		} else {
			log.Fatalf("Error reading file : %v", err)
		}
	}

	data := Data{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Cannot read save file, data of this game will recreated.")
		return startAttempts, []string{}, []string{}, "", err
	}

	fmt.Printf("Game recovered from save file %v.\n", saveFilename)
	return data.Attempts, strings.Split(data.ActualWord, ""), data.LettersSubmitted, data.Word, nil
}

func win() {
	hasWin = true
	fmt.Printf("Congrats !")
}

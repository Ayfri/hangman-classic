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

func main() {
	rand.Seed(time.Now().Unix())
	selectedFile := os.Args[1]
	file, err := ioutil.ReadFile(selectedFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("File does not exist")
		} else {
			log.Fatalf("Error reading file : %v", err)
		}
	}
	split := strings.Split(string(file), "\n")
	word := split[rand.Intn(len(split))]

	var letters []string
	lettersToReveal := len(word)/2 - 1
	for i := 0; i < len(word); i++ {
		letters = append(letters, "_")
	}

	for i := 0; i < lettersToReveal; i++ {
		index := rand.Intn(len(word))
		letters[index] = string(word[index])
	}

	attempts := 10
	for {
		println(strings.Join(letters, " "))
		submission, doExit := getLetter()
		if doExit {
			break
		}

		if submission == word {
			break
		}

		if len(submission) == 1 && isLetter(submission) {
			if !isLetterInWord(submission, word) {
				attempts--
				fmt.Printf("Wrong submission, attempts: %v\n", attempts)
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
			break
		}
	}
	fmt.Printf("You won!")
}

func getLetter() (result string, doExit bool) {
	letter := strings.ToLower(readLine())
	return letter, letter == exit
}

func isLetterInWord(letter string, word string) bool {
    for i := 0; i < len(word); i++ {
        if string(word[i]) == letter {
            return true
        }
    }
    return false
}

func isLetter(letter string) bool {
	return letter >= "a" && letter <= "z"
}

func isWordGuessed(letters []string, word string) bool {
	for i := 0; i < len(word); i++ {
		if string(word[i]) != letters[i] {
			return false
		}
	}
	return true
}

func readLine() string {
	var line string
	_, err := fmt.Scanln(&line)
	if err != nil {
		log.Fatalf("Error reading line : %v", err)
	}
	return line
}

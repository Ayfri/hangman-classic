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
		letter := getLetter()
		if letter == "exit" {
			break
		}
		if !isLetterInWord(letter, word) {
			attempts--
			fmt.Printf("Wrong letter, attempts: %v\n", attempts)
			continue
		}
		for i := 0; i < len(word); i++ {
			if string(word[i]) == letter {
				letters[i] = letter
			}
		}
		if isWordGuessed(letters, word) {
			fmt.Printf("You won!")
			break
		}
	}

	fmt.Printf(strings.Join(letters, ""))
}

func getLetter() string {
	var letter string
	for {
		letter = strings.ToLower(readLine())
		if len(letter) == 1 && isLetter(letter) {
			break
		}
		fmt.Println("Please enter a single letter")
	}
	return letter
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

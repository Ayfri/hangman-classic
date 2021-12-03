package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Data struct {
	ActualWord       string   `json:"actual_word"`
	Attempts         int      `json:"attempts"`
	LettersSubmitted []string `json:"letters_submitted"`
	Word             string   `json:"word"`
}

func (d *Data) MarshalJSON() (result []byte, err error) {
	return json.MarshalIndent(*d, "", "")
}

func (d *Data) SaveInJSONFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	marshalJSON, err := d.MarshalJSON()
	_ = ioutil.WriteFile(fileName, marshalJSON, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

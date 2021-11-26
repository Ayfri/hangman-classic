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

func NewData() *Data {
	return &Data{
		ActualWord:       "",
		Attempts:         0,
		LettersSubmitted: []string{},
		Word:             "",
	}
}

func (d *Data) MarshalJSON() (result []byte, err error) {
	return json.MarshalIndent(*d, "", "")
}

func (d *Data) UnmarshalJSON(b []byte) error {
	var data Data
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	return nil
}

func (d *Data) GetFromJSONFile(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, d); err != nil {
		return err
	}
	return nil
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
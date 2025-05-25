package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// use godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

type Prompt struct {
	Prompt         string `json:"prompt"`
	ExpectedOutput any    `json:"expectedOutput"`
}

func GoGetJsonValue(filepath string, key string) (Prompt, error) {
	// add filepath checker
	// if file is not found, return empty Prompt and error or path not ending with .json\
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return Prompt{}, err
	}
	defer jsonFile.Close()

	jsonData, _ := ioutil.ReadFile(filepath)
	var rawData map[string]json.RawMessage
	var prompt Prompt

	err = json.Unmarshal(jsonData, &rawData)
	if err != nil || len(rawData) == 0 {
		return Prompt{}, err
	}

	err = json.Unmarshal(rawData[key], &prompt)
	if err != nil {
		return Prompt{}, err
	}

	return prompt, nil
}

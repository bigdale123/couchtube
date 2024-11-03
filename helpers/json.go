package helpers

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func LoadJSONFromFile[T any](filePath string) (T, error) {
	var result T

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get working directory: %v", err)
		return result, err
	}

	jsonFile, err := os.Open(wd + filePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filePath, err)
		return result, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Printf("Failed to read file %s: %v", filePath, err)
		return result, err
	}

	if err := json.Unmarshal(byteValue, &result); err != nil {
		log.Printf("Failed to parse JSON from file %s: %v", filePath, err)
		return result, err
	}

	return result, nil
}

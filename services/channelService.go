package services

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/ozencb/couchtube/models"
)

func GetChannels() ([]models.Channel, error) {
	wd, err := os.Getwd()

	if err != nil {
		log.Printf("Failed to get working directory: %v", err)
		return nil, err
	}

	filePath := wd + "/channels.json"
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filePath, err)
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Printf("Failed to read file %s: %v", filePath, err)
		return nil, err
	}

	var channelsWrapper models.Channels
	if err := json.Unmarshal(byteValue, &channelsWrapper); err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		return nil, err
	}

	return channelsWrapper.Channels, nil
}

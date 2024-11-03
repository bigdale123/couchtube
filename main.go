package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

type Video struct {
	Url string `json:"url"`
}

type Channel struct {
	Id     int8    `json:"id"`
	Type   string  `json:"type"`
	Videos []Video `json:"videos"`
}

type Channels struct {
	Channels []Channel `json:"channels"`
}

func main() {
	http.Handle("/", corsMiddleware(http.FileServer(http.Dir("./static"))))
	http.Handle("/channels", corsMiddleware(http.HandlerFunc(GetChannels)))

	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func corsMiddleware(next http.Handler) http.Handler {
	whitelist := map[string]bool{
		"http://play.google.com":  true,
		"http://youtube.com":      true,
		"https://play.google.com": true,
		"https://youtube.com":     true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if _, ok := whitelist[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	channels, err := getChannels()
	if err != nil {
		http.Error(w, "Failed to load channels", http.StatusInternalServerError)
		log.Printf("Error loading channels.json: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(channels)
}

func getChannels() ([]byte, error) {
	filePath := "channels.json"
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

	return byteValue, nil
}

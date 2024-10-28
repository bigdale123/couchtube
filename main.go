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
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := getChannels()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(channels)
}

func getChannels() ([]byte, error) {
	jsonFile, err := os.Open("channels.json")

	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

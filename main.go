package main

import (
	"log"
	"net/http"

	"github.com/ozencb/couchtube/handlers"
	"github.com/ozencb/couchtube/middleware"
)

func main() {
	http.Handle("/", middleware.CORSMiddleware(http.FileServer(http.Dir("./static"))))
	http.Handle("/channels", middleware.CORSMiddleware(http.HandlerFunc(handlers.GetChannels)))

	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

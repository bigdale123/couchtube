package main

import (
	"log"
	"net/http"

	"github.com/ozencb/couchtube/db"
	"github.com/ozencb/couchtube/handlers"
	"github.com/ozencb/couchtube/middleware"
)

func main() {
	http.Handle("/", middleware.CORSMiddleware(http.FileServer(http.Dir("./static"))))
	http.Handle("/channels", middleware.CORSMiddleware(http.HandlerFunc(handlers.GetChannels)))
	http.Handle("/current-video", middleware.CORSMiddleware(http.HandlerFunc(handlers.GetCurrentVideo)))
	http.Handle("/submit-list", middleware.CORSMiddleware(http.HandlerFunc(handlers.SubmitList)))

	db.InitTables()
	db.PopulateDatabase()

	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

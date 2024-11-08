package main

import (
	"log"
	"net/http"

	"github.com/ozencb/couchtube/db"
	"github.com/ozencb/couchtube/handlers"
	"github.com/ozencb/couchtube/middleware"

	repo "github.com/ozencb/couchtube/repositories"
	"github.com/ozencb/couchtube/services"
)

func main() {
	// Initialize database
	dbInstance, err := db.GetConnector()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer dbInstance.Close()

	// Initialize Repositories
	channelRepo := repo.NewChannelRepository(dbInstance)
	videoRepo := repo.NewVideoRepository(dbInstance)

	// Initialize Services
	channelService := services.NewChannelService(channelRepo, videoRepo)
	submitListService := services.NewSubmitListService(channelRepo, videoRepo)

	// Initialize Handlers with services
	channelHandler := handlers.NewChannelHandler(channelService)
	submitListHandler := handlers.NewSubmitListHandler((submitListService))

	http.Handle("/", middleware.CORSMiddleware(http.FileServer(http.Dir("./static"))))
	http.Handle("/channels", middleware.CORSMiddleware(http.HandlerFunc(channelHandler.GetChannels)))
	http.Handle("/current-video", middleware.CORSMiddleware(http.HandlerFunc(channelHandler.GetCurrentVideo)))
	http.Handle("/submit-list", middleware.CORSMiddleware(http.HandlerFunc(submitListHandler.SubmitList)))

	db.InitTables()
	db.PopulateDatabase()

	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

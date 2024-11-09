package main

import (
	"log"
	"net/http"

	"github.com/ozencb/couchtube/db"
	"github.com/ozencb/couchtube/handlers"

	repo "github.com/ozencb/couchtube/repositories"
	"github.com/ozencb/couchtube/services"
)

type Route struct {
	Path    string
	Handler http.HandlerFunc
}

func registerRoutes(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		handler := route.Handler
		mux.Handle(route.Path, handler)
	}
}

func main() {
	// Initialize the database
	dbInstance, err := db.GetDbConnection()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.CloseConnector()

	db.InitDatabase(dbInstance)

	// Initialize Repositories
	channelRepo := repo.NewChannelRepository(dbInstance)
	videoRepo := repo.NewVideoRepository(dbInstance)

	// Initialize Services
	mediaService := services.NewMediaService(channelRepo, videoRepo)

	// Initialize Handlers with services
	mediaHandler := handlers.NewMediaHandler(mediaService)

	routes := []Route{
		{Path: "/", Handler: http.FileServer(http.Dir("./static")).ServeHTTP},
		{Path: "/channels", Handler: mediaHandler.FetchAllChannels},
		{Path: "/current-video", Handler: mediaHandler.GetCurrentVideo},
		{Path: "/submit-list", Handler: mediaHandler.SubmitList},
		{Path: "/invalidate-video", Handler: mediaHandler.InvalidateVideo},
	}
	registerRoutes(http.DefaultServeMux, routes)

	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

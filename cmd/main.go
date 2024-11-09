package main

import (
	"log"
	"net/http"

	"github.com/ozencb/couchtube/config"
	"github.com/ozencb/couchtube/db"
	"github.com/ozencb/couchtube/handlers"

	"github.com/ozencb/couchtube/middleware"
	repo "github.com/ozencb/couchtube/repositories"
	"github.com/ozencb/couchtube/services"
)

type Route struct {
	Path     string
	Handler  http.HandlerFunc
	Readonly bool
}

func registerRoutes(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		handler := route.Handler

		if route.Readonly {
			handler = middleware.ReadOnlyGuard(handler)
		}

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

	readonlyEnabled := config.GetReadonlyMode()

	routes := []Route{
		{Path: "/", Handler: http.FileServer(http.Dir("./static")).ServeHTTP, Readonly: false},
		{Path: "/channels", Handler: mediaHandler.FetchAllChannels, Readonly: false},
		{Path: "/current-video", Handler: mediaHandler.GetCurrentVideo, Readonly: false},
		{Path: "/submit-list", Handler: mediaHandler.SubmitList, Readonly: readonlyEnabled},
		{Path: "/invalidate-video", Handler: mediaHandler.InvalidateVideo, Readonly: readonlyEnabled},
	}
	registerRoutes(http.DefaultServeMux, routes)

	port := config.GetPort()

	log.Println("Server starting on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

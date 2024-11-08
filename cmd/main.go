package main

import (
	"log"
	"net/http"

	"github.com/ozencb/couchtube/db"
	"github.com/ozencb/couchtube/handlers"
	middleware "github.com/ozencb/couchtube/middleware"

	repo "github.com/ozencb/couchtube/repositories"
	"github.com/ozencb/couchtube/services"
)

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
	Cors    bool
}

func registerRoutes(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		handler := route.Handler
		if route.Cors {
			handler = middleware.CORSMiddleware(handler).(http.HandlerFunc)
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

	routes := []Route{
		{Path: "/", Method: "GET", Handler: http.FileServer(http.Dir("./static")).ServeHTTP, Cors: true},
		{Path: "/channels", Method: "GET", Handler: mediaHandler.FetchAllChannels, Cors: true},
		{Path: "/current-video", Method: "GET", Handler: mediaHandler.GetCurrentVideo, Cors: true},
		{Path: "/submit-list", Method: "POST", Handler: mediaHandler.SubmitList, Cors: true},
	}
	registerRoutes(http.DefaultServeMux, routes)

	log.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

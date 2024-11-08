package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	dbmodels "github.com/ozencb/couchtube/models/db"
	"github.com/ozencb/couchtube/services"
)

func GetChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	channels, err := services.GetChannels()
	if err != nil {
		http.Error(w, "Failed to load channels", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"channels": channels})
}

func GetCurrentVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	channelID := r.URL.Query().Get("channel-id")
	if channelID == "" {
		http.Error(w, "channel-id is required", http.StatusBadRequest)
		return
	}
	channelIDInt, err := strconv.Atoi(channelID)
	if err != nil {
		http.Error(w, "Invalid channel-id", http.StatusBadRequest)
		return
	}

	videoID := r.URL.Query().Get("video-id")

	var video *dbmodels.Video
	// if videoId is provided, call GetNextVideo
	if videoID != "" {
		videoIDInt, err := strconv.Atoi(videoID)
		if err != nil {
			http.Error(w, "Invalid video-id", http.StatusBadRequest)
			return
		}
		video = services.GetNextVideo(channelIDInt, videoIDInt)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"video": video})
		return
	} else {
		// if videoId is not provided, call GetCurrentVideoByChannelId
		video, err = services.GetCurrentVideoByChannelId(channelIDInt)
		if err != nil {
			http.Error(w, "Failed to load video", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"video": video})
}

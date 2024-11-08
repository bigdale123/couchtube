package handlers

import (
	"encoding/json"
	"net/http"

	jsonmodels "github.com/ozencb/couchtube/models/json"
	"github.com/ozencb/couchtube/services"
)

type SubmitListHandler struct {
	Service *services.SubmitListService
}

func NewSubmitListHandler(service *services.SubmitListService) *SubmitListHandler {
	return &SubmitListHandler{Service: service}
}

func (h *SubmitListHandler) SubmitList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var list jsonmodels.SubmitListRequestJson
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		http.Error(w, "Failed to parse list", http.StatusBadRequest)
		return
	}

	success, err := h.Service.SubmitList(list)
	if err != nil {
		http.Error(w, "Failed to submit list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": success})
}

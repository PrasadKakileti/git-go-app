package handlers

import (
	"encoding/json"
	"net/http"
)

type updateProfileRequest struct {
	Email      string `json:"email"`
	Location   string `json:"location"`
	Domain     string `json:"domain"`
	Experience string `json:"experience"`
	Frequency  string `json:"notification_frequency"`
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	if req.Location == "" || req.Domain == "" || req.Experience == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Location, domain and experience are required"})
		return
	}
	if req.Frequency == "" {
		req.Frequency = "daily"
	}

	if err := h.userRepo.UpdateProfile(user.ID, req.Location, req.Domain, req.Experience, req.Frequency); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update profile"})
		return
	}

	// Update the user object and fetch fresh jobs for new preferences
	user.Location = req.Location
	user.Domain = req.Domain
	user.Experience = req.Experience
	user.NotificationFrequency = req.Frequency
	go h.triggerImmediateFetch(user) // Fetches + emails immediately

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Profile updated! Fetching new job matches for your updated preferences.",
	})
}

// RefreshJobs lets the user manually trigger a fresh job fetch from the dashboard.
func (h *Handler) RefreshJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "email required"})
		return
	}

	user, err := h.userRepo.GetByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	go h.triggerImmediateFetch(user)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Fetching fresh jobs — you'll receive them via email within 1 minute.",
	})
}

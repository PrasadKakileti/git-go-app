package handlers

import (
	"encoding/json"
	"job-portal/repository"
	"net/http"
	"time"
)

type JobHandler struct {
	jobRepo  *repository.JobRepository
	userRepo *repository.UserRepository
}

func NewJobHandler(jobRepo *repository.JobRepository, userRepo *repository.UserRepository) *JobHandler {
	return &JobHandler{
		jobRepo:  jobRepo,
		userRepo: userRepo,
	}
}

func (h *JobHandler) GetJobsForUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parameter required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	since := time.Now().AddDate(0, 0, -30) // Last 30 days
	jobs, err := h.jobRepo.GetUnsentJobsForUser(user.ID, user.Location, user.Domain, since)
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jobs": jobs,
	})
}

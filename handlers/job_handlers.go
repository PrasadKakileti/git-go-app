package handlers

import (
	"encoding/json"
	"job-portal/models"
	"job-portal/repository"
	"net/http"
	"time"
)

type JobHandler struct {
	jobRepo  *repository.JobRepository
	userRepo *repository.UserRepository
}

func NewJobHandler(jobRepo *repository.JobRepository, userRepo *repository.UserRepository) *JobHandler {
	return &JobHandler{jobRepo: jobRepo, userRepo: userRepo}
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func (h *JobHandler) GetJobsForUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email parameter required"})
		return
	}

	user, err := h.userRepo.GetByEmail(email)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found — please log in again"})
		return
	}

	since := time.Now().AddDate(0, 0, -30)
	jobs, err := h.jobRepo.GetUnsentJobsForUser(user.ID, user.Location, user.Domain, user.Experience, since)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "could not load jobs — database query failed",
		})
		return
	}

	// Return empty slice (not null) so frontend never sees null
	if jobs == nil {
		jobs = []*models.Job{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jobs":  jobs,
		"total": len(jobs),
	})
}

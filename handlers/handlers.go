package handlers

import (
	"encoding/json"
	"job-portal/config"
	"job-portal/models"
	"job-portal/repository"
	"net/http"
)

// JobFetcher is the subset of services.JobService the handler needs.
// Declared here to avoid an import cycle.
type JobFetcher interface {
	FetchAndStoreJobs(location, domain, experience string) error
	FetchJobsAndNotify(user *models.User) error // Fetch jobs + immediately send email
}

type Handler struct {
	userRepo     *repository.UserRepository
	emailService EmailServiceInterface
	jobFetcher   JobFetcher
	cfg          *config.Config
}

func NewHandler(userRepo *repository.UserRepository) *Handler {
	return &Handler{userRepo: userRepo}
}

func (h *Handler) SetEmailService(svc EmailServiceInterface) { h.emailService = svc }
func (h *Handler) SetConfig(cfg *config.Config)              { h.cfg = cfg }
func (h *Handler) SetJobFetcher(jf JobFetcher)               { h.jobFetcher = jf }

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	if req.Email == "" || req.Location == "" || req.Domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email, location, and domain are required"})
		return
	}

	if req.NotificationFrequency == "" {
		req.NotificationFrequency = "daily"
	}

	user := &models.User{
		Email:                 req.Email,
		Location:              req.Location,
		Domain:                req.Domain,
		NotificationFrequency: req.NotificationFrequency,
		IsActive:              true,
	}

	if err := h.userRepo.Create(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to register: " + err.Error()})
		return
	}

	go h.sendWelcomeEmail(user)
	go h.triggerImmediateFetch(user)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Registration successful! Job alerts are being fetched and will be emailed within 1 minute.",
		"user":    user,
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email is required"})
		return
	}

	user, err := h.userRepo.GetByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) sendWelcomeEmail(user *models.User) {
	if h.emailService != nil {
		h.emailService.SendWelcomeEmail(user.Email, user.Location, user.Domain, user.NotificationFrequency)
	}
}

// triggerImmediateFetch fetches jobs AND sends them as an email immediately.
// New users get job notifications within seconds of signup, not waiting for the 1-hour scheduler.
func (h *Handler) triggerImmediateFetch(user *models.User) {
	if h.jobFetcher != nil {
		// FetchJobsAndNotify does: fetch → store → email → mark as sent
		h.jobFetcher.FetchJobsAndNotify(user)
	}
}

type EmailServiceInterface interface {
	SendWelcomeEmail(email, location, domain, frequency string) error
}

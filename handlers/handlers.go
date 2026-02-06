package handlers

import (
	"encoding/json"
	"job-portal/models"
	"job-portal/repository"
	"net/http"
)

type Handler struct {
	userRepo     *repository.UserRepository
	emailService EmailServiceInterface
}

func NewHandler(userRepo *repository.UserRepository) *Handler {
	return &Handler{userRepo: userRepo}
}

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
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to register user: " + err.Error()})
		return
	}

	// Send welcome email
	go h.sendWelcomeEmail(user)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Registration successful! Check your email for confirmation.",
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

func (h *Handler) SetEmailService(emailService EmailServiceInterface) {
	h.emailService = emailService
}

type EmailServiceInterface interface {
	SendWelcomeEmail(email, location, domain, frequency string) error
}

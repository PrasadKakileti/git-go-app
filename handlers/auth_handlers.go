package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"job-portal/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	if (req.Email == "" && req.Mobile == "") || req.Password == "" || req.Location == "" || req.Domain == "" || req.Experience == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "All fields are required"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to process password"})
		return
	}

	user := &models.User{
		Email:                 req.Email,
		Mobile:                req.Mobile,
		Password:              string(hashedPassword),
		Location:              req.Location,
		Domain:                req.Domain,
		Experience:            req.Experience,
		NotificationFrequency: req.NotificationFrequency,
		IsActive:              true,
		IsVerified:            false,
		VerificationCode:      generateCode(),
	}
	if user.NotificationFrequency == "" {
		user.NotificationFrequency = "daily"
	}

	if err := h.userRepo.CreateWithAuth(user); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email or mobile already registered"})
		return
	}

	go h.sendWelcomeEmail(user)
	go h.triggerImmediateFetch(user) // fetch matching jobs right away, no 6h wait

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Signup successful! Job alerts are being fetched and will be emailed to you within 1 minute.",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid request"})
		return
	}

	// Admin login: credentials come from environment, never hardcoded.
	if h.cfg != nil && req.EmailOrMobile == h.cfg.AdminEmail && h.cfg.AdminPassword != "" && req.Password == h.cfg.AdminPassword {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Login successful",
			"token":   generateSecureToken(),
			"user": map[string]interface{}{
				"id":       0,
				"email":    h.cfg.AdminEmail,
				"role":     "admin",
				"location": "Admin",
				"domain":   "Admin",
			},
		})
		return
	}

	user, err := h.userRepo.GetByEmailOrMobile(req.EmailOrMobile)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid credentials"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"token":   generateSecureToken(),
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"mobile":     user.Mobile,
			"location":   user.Location,
			"domain":     user.Domain,
			"experience": user.Experience,
		},
	})
}

// generateSecureToken returns a cryptographically random 32-byte hex token.
func generateSecureToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateCode() string {
	b := make([]byte, 3) // 6 hex chars
	rand.Read(b)
	return hex.EncodeToString(b)
}


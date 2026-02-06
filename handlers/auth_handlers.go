package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"job-portal/models"
	"net/http"
	"time"

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

	// Validate required fields
	if (req.Email == "" && req.Mobile == "") || req.Password == "" || req.Location == "" || req.Domain == "" || req.Experience == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "All fields are required"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to process password"})
		return
	}

	// Generate verification code
	verificationCode := generateVerificationCode()

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
		VerificationCode:      verificationCode,
	}

	if err := h.userRepo.CreateWithAuth(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email or mobile already registered"})
		return
	}

	// Send welcome email
	go h.sendWelcomeEmail(user)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Signup successful! Please login to continue.",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request",
		})
		return
	}

	// Hardcoded admin login
	if req.EmailOrMobile == "Test123@gmail.com" && req.Password == "Test@123" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Login successful",
			"token":   "admin-token",
			"user": map[string]interface{}{
				"id":       0,
				"email":    "Test123@gmail.com",
				"mobile":   "",
				"location": "Admin",
				"domain":   "Admin",
			},
		})
		return
	}

	// Get user by email or mobile
	user, err := h.userRepo.GetByEmailOrMobile(req.EmailOrMobile)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid credentials",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid credentials",
		})
		return
	}

	// Generate session token (simple implementation)
	token := generateToken(user.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"token":   token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"email":    user.Email,
			"mobile":   user.Mobile,
			"location": user.Location,
			"domain":   user.Domain,
		},
	})
}

func generateVerificationCode() string {
	const digits = "0123456789"
	b := make([]byte, 6)
	rand.Read(b)
	for i := range b {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b)
}

func generateToken(userID int) string {
	return fmt.Sprintf("%d-%d", userID, time.Now().Unix())
}

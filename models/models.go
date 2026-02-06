package models

import "time"

type User struct {
	ID                    int       `json:"id"`
	Email                 string    `json:"email"`
	Mobile                string    `json:"mobile"`
	Password              string    `json:"-"`
	Location              string    `json:"location"`
	Domain                string    `json:"domain"`
	Experience            string    `json:"experience"`
	NotificationFrequency string    `json:"notification_frequency"`
	IsActive              bool      `json:"is_active"`
	IsVerified            bool      `json:"is_verified"`
	VerificationCode      string    `json:"-"`
	CreatedAt             time.Time `json:"created_at"`
}

type Job struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Company      string    `json:"company"`
	Location     string    `json:"location"`
	Domain       string    `json:"domain"`
	EmailContact string    `json:"email_contact"`
	Description  string    `json:"description"`
	PostedAt     time.Time `json:"posted_at"`
	ScrapedAt    time.Time `json:"scraped_at"`
	SourceURL    string    `json:"source_url"`
}

type RegisterRequest struct {
	Email                 string `json:"email"`
	Location              string `json:"location"`
	Domain                string `json:"domain"`
	NotificationFrequency string `json:"notification_frequency"`
}


type SignupRequest struct {
	Email                 string `json:"email"`
	Mobile                string `json:"mobile"`
	Password              string `json:"password"`
	Location              string `json:"location"`
	Domain                string `json:"domain"`
	Experience            string `json:"experience"`
	NotificationFrequency string `json:"notification_frequency"`
}

type LoginRequest struct {
	EmailOrMobile string `json:"emailOrMobile"`
	Password      string `json:"password"`
}

package models

import "time"

// SupportedDomains is the canonical list used in the signup form and job fetching.
var SupportedDomains = []string{
	"Java",
	"Node.js",
	"PLM Teamcenter",
	"OPcenter",
}

// SupportedExperience maps display labels to the stored value.
var SupportedExperience = []struct {
	Value string
	Label string
}{
	{"0-1", "0–1 years (Fresher)"},
	{"1-3", "1–3 years"},
	{"3-5", "3–5 years"},
	{"5-10", "5–10 years"},
	{"10+", "10+ years"},
}

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
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	Domain      string    `json:"domain"`
	Description string    `json:"description"`
	PostedAt    time.Time `json:"posted_at"`
	ScrapedAt   time.Time `json:"scraped_at"`
	SourceURL   string    `json:"source_url"` // Real apply link (LinkedIn, Indeed, Glassdoor, etc.)
	Source      string    `json:"source"`     // Platform name: "LinkedIn", "Indeed", etc.
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

// sendtest — one-shot program that:
//   1. Calls the JSearch API for the test user's preferences
//   2. Stores real jobs in the database
//   3. Immediately emails them to prasadkakileti143@gmail.com
//
// Usage:  go run cmd/sendtest/main.go
package main

import (
	"job-portal/config"
	"job-portal/database"
	"job-portal/providers"
	"job-portal/repository"
	"job-portal/services"
	"log"
	"time"
)

const testEmail = "prasadkakileti143@gmail.com"

func main() {
	cfg := config.Load()

	if cfg.JSearchAPIKey == "" {
		log.Fatal("JSEARCH_API_KEY is not set in .env — cannot fetch real jobs")
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	jobRepo  := repository.NewJobRepository(db)

	// Load the test user
	user, err := userRepo.GetByEmail(testEmail)
	if err != nil {
		log.Fatalf("User %s not found in DB: %v", testEmail, err)
	}
	log.Printf("User: %s | Location: %s | Domain: %s | Experience: %s",
		user.Email, user.Location, user.Domain, user.Experience)

	// Fetch real jobs from JSearch
	p := providers.NewJSearchProvider(cfg.JSearchAPIKey)
	log.Printf("Fetching %s jobs in %s for %s yrs experience...", user.Domain, user.Location, user.Experience)

	jobs, err := p.FetchJobs(user.Location, user.Domain, user.Experience)
	if err != nil {
		log.Fatalf("JSearch error: %v", err)
	}
	log.Printf("Fetched %d real jobs from API", len(jobs))

	// Store in DB
	stored := 0
	for _, job := range jobs {
		if err := jobRepo.Upsert(job); err != nil {
			log.Printf("  skip (dup or error): %s — %v", job.Title, err)
		} else {
			stored++
			log.Printf("  ✓ %s @ %s [%s]", job.Title, job.Company, job.Source)
		}
	}
	log.Printf("Stored %d new jobs", stored)

	// Get unsent jobs for this user (last 30 days)
	since := time.Now().AddDate(0, 0, -30)
	unsent, err := jobRepo.GetUnsentJobsForUser(user.ID, user.Location, user.Domain, user.Experience, since)
	if err != nil {
		log.Fatalf("GetUnsentJobs: %v", err)
	}
	log.Printf("Found %d unsent jobs ready to email", len(unsent))

	if len(unsent) == 0 {
		log.Println("No new jobs to send — all already emailed, or none matched. Try again after clearing user_job_sent.")
		return
	}

	// Send email
	emailSvc := services.NewUnifiedEmailService(cfg)
	log.Printf("Sending email to %s...", user.Email)
	if err := emailSvc.SendJobNotifications(user.Email, unsent); err != nil {
		log.Fatalf("Email send failed: %v", err)
	}

	// Mark as sent
	for _, job := range unsent {
		_ = jobRepo.MarkJobAsSent(user.ID, job.ID)
	}
	log.Printf("✅ Done! Email with %d jobs sent to %s", len(unsent), user.Email)
	log.Println("   Check your inbox — sender will appear as 'JobHub'")
}

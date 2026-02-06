package main

import (
	"fmt"
	"job-portal/config"
	"job-portal/database"
	"job-portal/models"
	"job-portal/repository"
	"job-portal/scraper"
	"job-portal/services"
	"log"
)

func main() {
	cfg := config.Load()
	
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup services
	jobRepo := repository.NewJobRepository(db)
	userRepo := repository.NewUserRepository(db)
	naukriScraper := scraper.NewNaukriScraper()
	emailService := services.NewUnifiedEmailService(cfg)

	fmt.Println("========================================")
	fmt.Println("   Email Test - Job Alert")
	fmt.Println("========================================")
	fmt.Println()

	// Get or create test user
	testEmail := "babusurendra500@gmail.com"
	user, err := userRepo.GetByEmail(testEmail)
	if err != nil {
		fmt.Println("Creating test user...")
		user = &models.User{
			Email:                 testEmail,
			Location:              "Bangalore",
			Domain:                "Java",
			NotificationFrequency: "daily",
			IsActive:              true,
		}
		if err := userRepo.Create(user); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("✅ User created: %s\n", testEmail)
	} else {
		fmt.Printf("✅ User found: %s\n", testEmail)
	}

	fmt.Println()
	fmt.Println("Scraping jobs for Bangalore...")
	
	// Scrape jobs
	jobs, err := naukriScraper.ScrapeJobs("Bangalore")
	if err != nil {
		log.Fatalf("Failed to scrape jobs: %v", err)
	}
	
	fmt.Printf("✅ Found %d jobs\n", len(jobs))
	
	// Store jobs
	for _, job := range jobs {
		jobRepo.Create(job)
	}
	
	// Get Java jobs for user
	javaJobs := []*models.Job{}
	for _, job := range jobs {
		if job.Domain == "Java" {
			javaJobs = append(javaJobs, job)
		}
	}
	
	if len(javaJobs) == 0 {
		fmt.Println("⚠️  No Java jobs found, using all jobs for test")
		javaJobs = jobs
	}
	
	fmt.Println()
	fmt.Printf("Sending %d jobs to %s...\n", len(javaJobs), testEmail)
	
	// Send email
	err = emailService.SendJobNotifications(testEmail, javaJobs)
	if err != nil {
		log.Fatalf("❌ Failed to send email: %v", err)
	}
	
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("✅ Email sent successfully!")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Printf("Check inbox: %s\n", testEmail)
	fmt.Println("(Check spam folder if not in inbox)")
	fmt.Println()
}

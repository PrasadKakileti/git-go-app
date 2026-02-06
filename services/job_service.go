package services

import (
	"job-portal/repository"
	"job-portal/scraper"
	"log"
	"time"
)

type JobService struct {
	jobRepo  *repository.JobRepository
	userRepo *repository.UserRepository
	scraper  *scraper.NaukriScraper
	email    *UnifiedEmailService
}

func NewJobService(jobRepo *repository.JobRepository, userRepo *repository.UserRepository, 
	scraper *scraper.NaukriScraper, email *UnifiedEmailService) *JobService {
	return &JobService{
		jobRepo:  jobRepo,
		userRepo: userRepo,
		scraper:  scraper,
		email:    email,
	}
}

func (s *JobService) ScrapeAndStoreJobs(locations []string) error {
	for _, location := range locations {
		jobs, err := s.scraper.ScrapeJobs(location)
		if err != nil {
			log.Printf("Error scraping jobs for %s: %v", location, err)
			continue
		}

		for _, job := range jobs {
			if err := s.jobRepo.Create(job); err != nil {
				log.Printf("Error storing job: %v", err)
			}
		}
		log.Printf("Scraped and stored %d jobs for %s", len(jobs), location)
	}
	return nil
}

func (s *JobService) SendNotifications(frequency string) error {
	users, err := s.userRepo.GetActiveUsers(frequency)
	if err != nil {
		return err
	}

	var since time.Time
	if frequency == "daily" {
		since = time.Now().Add(-24 * time.Hour)
	} else {
		since = time.Now().Add(-7 * 24 * time.Hour)
	}

	for _, user := range users {
		jobs, err := s.jobRepo.GetUnsentJobsForUser(user.ID, user.Location, user.Domain, since)
		if err != nil {
			log.Printf("Error getting jobs for user %s: %v", user.Email, err)
			continue
		}

		if len(jobs) > 0 {
			if err := s.email.SendJobNotifications(user.Email, jobs); err != nil {
				log.Printf("Error sending email to %s: %v", user.Email, err)
				continue
			}

			for _, job := range jobs {
				s.jobRepo.MarkJobAsSent(user.ID, job.ID)
			}
			log.Printf("Sent %d jobs to %s", len(jobs), user.Email)
		}
	}
	return nil
}

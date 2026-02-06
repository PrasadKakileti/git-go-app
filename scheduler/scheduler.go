package scheduler

import (
	"job-portal/services"
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	jobService *services.JobService
	cron       *cron.Cron
}

func NewScheduler(jobService *services.JobService) *Scheduler {
	return &Scheduler{
		jobService: jobService,
		cron:       cron.New(),
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("0 */6 * * *", func() {
		log.Println("Starting job scraping...")
		locations := []string{"Bangalore", "Mumbai", "Delhi", "Hyderabad", "Pune"}
		if err := s.jobService.ScrapeAndStoreJobs(locations); err != nil {
			log.Printf("Error in scraping: %v", err)
		}
	})

	s.cron.AddFunc("*/10 * * * *", func() {
		log.Println("Sending job notifications (every 10 minutes)...")
		if err := s.jobService.SendNotifications("daily"); err != nil {
			log.Printf("Error sending notifications: %v", err)
		}
	})

	s.cron.Start()
	log.Println("Scheduler started - Notifications every 10 minutes")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

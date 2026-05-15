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
	// Every 6 hours: fetch fresh jobs from JSearch for every user's location+domain+experience.
	s.cron.AddFunc("0 */6 * * *", func() {
		log.Println("[Scheduler] Fetching jobs for all user preferences...")
		if err := s.jobService.FetchForAllUsers(); err != nil {
			log.Printf("[Scheduler] FetchForAllUsers error: %v", err)
		}
	})

	// Every hour: send job notification emails to daily subscribers.
	s.cron.AddFunc("0 * * * *", func() {
		log.Println("[Scheduler] Sending hourly job notifications...")
		if err := s.jobService.SendNotifications("daily"); err != nil {
			log.Printf("[Scheduler] hourly notifications error: %v", err)
		}
	})

	// Every Sunday at 9 AM: send weekly digest.
	s.cron.AddFunc("0 9 * * 0", func() {
		log.Println("[Scheduler] Sending weekly digest...")
		if err := s.jobService.SendNotifications("weekly"); err != nil {
			log.Printf("[Scheduler] weekly notifications error: %v", err)
		}
	})

	s.cron.Start()
	log.Println("[Scheduler] Started — job fetch every 6h, email every 1h")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

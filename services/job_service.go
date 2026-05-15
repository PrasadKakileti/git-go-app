package services

import (
	"job-portal/models"
	"job-portal/providers"
	"job-portal/repository"
	"log"
	"time"
)

type JobService struct {
	jobRepo   *repository.JobRepository
	userRepo  *repository.UserRepository
	providers []providers.JobProvider
	email     *UnifiedEmailService
}

func NewJobService(
	jobRepo *repository.JobRepository,
	userRepo *repository.UserRepository,
	providerList []providers.JobProvider,
	email *UnifiedEmailService,
) *JobService {
	return &JobService{
		jobRepo:   jobRepo,
		userRepo:  userRepo,
		providers: providerList,
		email:     email,
	}
}

// FetchAndStoreJobs pulls jobs from all registered providers for the given location+domain+experience
// and persists them to the database (deduplicating by source_url).
func (s *JobService) FetchAndStoreJobs(location, domain, experience string) error {
	for _, p := range s.providers {
		jobs, err := p.FetchJobs(location, domain, experience)
		if err != nil {
			log.Printf("[%s] fetch error for %s/%s: %v", p.Name(), domain, location, err)
			continue
		}
		for _, job := range jobs {
			if err := s.jobRepo.Upsert(job); err != nil {
				log.Printf("store job error: %v", err)
			}
		}
		log.Printf("[%s] stored %d jobs — domain=%s location=%s", p.Name(), len(jobs), domain, location)
	}
	return nil
}

// FetchJobsAndNotify fetches fresh jobs for a user AND immediately emails them.
// Used for new signup — ensures immediate job notification, not waiting for scheduler.
func (s *JobService) FetchJobsAndNotify(user *models.User) error {
	// 1. Fetch fresh jobs
	if err := s.FetchAndStoreJobs(user.Location, user.Domain, user.Experience); err != nil {
		return err
	}

	// 2. Get unsent jobs (all jobs we just added are unsent)
	since := time.Now().AddDate(0, 0, -30)
	jobs, err := s.jobRepo.GetUnsentJobsForUser(user.ID, user.Location, user.Domain, user.Experience, since)
	if err != nil {
		return err
	}

	// 3. Send them immediately as a notification email
	if len(jobs) > 0 {
		if err := s.email.SendJobNotifications(user.Email, jobs); err != nil {
			// Log but don't fail — at least jobs are in the DB
			log.Printf("[FetchJobsAndNotify] email send failed for %s: %v", user.Email, err)
			return nil
		}
		// 4. Mark as sent so they don't email again in 1 hour
		for _, job := range jobs {
			_ = s.jobRepo.MarkJobAsSent(user.ID, job.ID)
		}
		log.Printf("[FetchJobsAndNotify] Sent %d jobs to %s immediately", len(jobs), user.Email)
	}
	return nil
}

// FetchForAllUsers fetches jobs for every unique (location, domain, experience) combination
// across all active users.  Called by the scheduler.
func (s *JobService) FetchForAllUsers() error {
	combos, err := s.userRepo.GetDistinctUserPreferences()
	if err != nil {
		return err
	}
	for _, c := range combos {
		if err := s.FetchAndStoreJobs(c.Location, c.Domain, c.Experience); err != nil {
			log.Printf("fetch error for %+v: %v", c, err)
		}
	}
	return nil
}

// SendNotifications emails matched jobs to all active users for the given frequency.
func (s *JobService) SendNotifications(frequency string) error {
	users, err := s.userRepo.GetActiveUsers(frequency)
	if err != nil {
		return err
	}

	var since time.Time
	switch frequency {
	case "weekly":
		since = time.Now().Add(-7 * 24 * time.Hour)
	default:
		since = time.Now().Add(-24 * time.Hour)
	}

	for _, user := range users {
		jobs, err := s.jobRepo.GetUnsentJobsForUser(user.ID, user.Location, user.Domain, user.Experience, since)
		if err != nil {
			log.Printf("jobs fetch error for %s: %v", user.Email, err)
			continue
		}

		if len(jobs) == 0 {
			continue
		}

		if err := s.email.SendJobNotifications(user.Email, jobs); err != nil {
			log.Printf("email error for %s: %v", user.Email, err)
			continue
		}

		for _, job := range jobs {
			_ = s.jobRepo.MarkJobAsSent(user.ID, job.ID)
		}
		log.Printf("Sent %d jobs to %s", len(jobs), user.Email)
	}
	return nil
}

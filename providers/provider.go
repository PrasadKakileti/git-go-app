package providers

import "job-portal/models"

// JobProvider is the interface every job-source integration must implement.
// Add LinkedIn, Indeed, Naukri etc. as separate structs implementing this.
type JobProvider interface {
	FetchJobs(location, domain, experience string) ([]*models.Job, error)
	Name() string
}

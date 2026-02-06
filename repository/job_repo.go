package repository

import (
	"database/sql"
	"job-portal/models"
	"time"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(job *models.Job) error {
	query := `INSERT INTO jobs (title, company, location, domain, email_contact, description, posted_at, source_url) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, job.Title, job.Company, job.Location, job.Domain,
		job.EmailContact, job.Description, job.PostedAt, job.SourceURL)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	job.ID = int(id)
	return nil
}

func (r *JobRepository) GetUnsentJobsForUser(userID int, location, domain string, since time.Time) ([]*models.Job, error) {
	query := `SELECT j.id, j.title, j.company, j.location, j.domain, j.email_contact, j.description, j.posted_at, j.source_url
			  FROM jobs j
			  WHERE j.location LIKE ? 
			  AND j.domain = ?
			  AND j.posted_at >= ?
			  AND NOT EXISTS (
				  SELECT 1 FROM user_job_sent ujs 
				  WHERE ujs.user_id = ? AND ujs.job_id = j.id
			  )
			  ORDER BY j.posted_at DESC
			  LIMIT 50`
	
	rows, err := r.db.Query(query, "%"+location+"%", domain, since, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		job := &models.Job{}
		if err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location, &job.Domain,
			&job.EmailContact, &job.Description, &job.PostedAt, &job.SourceURL); err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *JobRepository) MarkJobAsSent(userID, jobID int) error {
	query := `INSERT INTO user_job_sent (user_id, job_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE sent_at = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, userID, jobID)
	return err
}

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

// Upsert inserts a job; if source_url already exists it updates the title/description only.
// This prevents duplicate rows for the same real job posting.
func (r *JobRepository) Upsert(job *models.Job) error {
	query := `INSERT INTO jobs (title, company, location, domain, description, posted_at, source_url, source)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			  ON DUPLICATE KEY UPDATE title=VALUES(title), description=VALUES(description), scraped_at=CURRENT_TIMESTAMP`
	result, err := r.db.Exec(query,
		job.Title, job.Company, job.Location, job.Domain,
		job.Description, job.PostedAt, job.SourceURL, job.Source,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	job.ID = int(id)
	return nil
}

// GetUnsentJobsForUser returns jobs matching the user's location, domain, and experience
// that haven't already been emailed to them.
func (r *JobRepository) GetUnsentJobsForUser(userID int, location, domain, experience string, since time.Time) ([]*models.Job, error) {
	query := `SELECT j.id, j.title, j.company, j.location, j.domain, j.description,
			         j.posted_at, COALESCE(j.source,''), j.source_url
			  FROM jobs j
			  WHERE j.location LIKE ?
			    AND j.domain = ?
			    AND j.posted_at >= ?
			    AND NOT EXISTS (
			        SELECT 1 FROM user_job_sent ujs WHERE ujs.user_id = ? AND ujs.job_id = j.id
			    )
			  ORDER BY j.posted_at DESC
			  LIMIT 20`

	rows, err := r.db.Query(query, "%"+location+"%", domain, since, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		job := &models.Job{}
		if err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location,
			&job.Domain, &job.Description, &job.PostedAt, &job.Source, &job.SourceURL); err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *JobRepository) MarkJobAsSent(userID, jobID int) error {
	query := `INSERT INTO user_job_sent (user_id, job_id) VALUES (?, ?)
			  ON DUPLICATE KEY UPDATE sent_at = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, userID, jobID)
	return err
}

// GetRecentJobs is used by the dashboard API.
func (r *JobRepository) GetRecentJobs(location, domain string, limit int) ([]*models.Job, error) {
	query := `SELECT id, title, company, location, domain, description, posted_at, COALESCE(source,''), source_url
			  FROM jobs
			  WHERE location LIKE ? AND domain = ?
			  ORDER BY posted_at DESC
			  LIMIT ?`
	rows, err := r.db.Query(query, "%"+location+"%", domain, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		job := &models.Job{}
		if err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location,
			&job.Domain, &job.Description, &job.PostedAt, &job.Source, &job.SourceURL); err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

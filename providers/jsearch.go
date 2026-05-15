package providers

// JSearch aggregates jobs from LinkedIn, Indeed, Glassdoor, ZipRecruiter, and more.
// Free tier on RapidAPI: 200 requests/month — enough for development and small scale use.
// Sign up at: https://rapidapi.com/letscrape-6bRBa3QguO5/api/jsearch
// Every job returned has a real job_apply_link pointing to the actual application page.

import (
	"encoding/json"
	"fmt"
	"job-portal/models"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type JSearchProvider struct {
	apiKey string
	client *http.Client
}

func NewJSearchProvider(apiKey string) *JSearchProvider {
	return &JSearchProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *JSearchProvider) Name() string { return "JSearch (LinkedIn/Indeed/Glassdoor)" }

// experienceToRequirement maps the user's stored experience string to the JSearch API parameter.
func experienceToRequirement(exp string) string {
	switch exp {
	case "0-1":
		return "no_experience"
	case "1-3":
		return "under_3_years_experience"
	case "3-5", "5-10", "10+":
		return "more_than_3_years_experience"
	default:
		return ""
	}
}

type jsearchResponse struct {
	Status string      `json:"status"`
	Data   []jsearchJob `json:"data"`
}

type jsearchJob struct {
	JobID          string  `json:"job_id"`
	EmployerName   string  `json:"employer_name"`
	JobTitle       string  `json:"job_title"`
	JobDescription string  `json:"job_description"`
	JobApplyLink   string  `json:"job_apply_link"`
	JobGoogleLink  string  `json:"job_google_link"`
	JobCity        string  `json:"job_city"`
	JobState       string  `json:"job_state"`
	JobCountry     string  `json:"job_country"`
	JobPostedAt    int64   `json:"job_posted_at_timestamp"`
	JobSource      string  `json:"job_publisher"`
	JobIsRemote    bool    `json:"job_is_remote"`
}

func (p *JSearchProvider) FetchJobs(location, domain, experience string) ([]*models.Job, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("JSearch API key not configured — set JSEARCH_API_KEY in .env (free at rapidapi.com/letscrape-6bRBa3QguO5/api/jsearch)")
	}

	query := fmt.Sprintf("%s jobs in %s", domain, location)
	params := url.Values{}
	params.Set("query", query)
	params.Set("page", "1")
	params.Set("num_pages", "1")
	params.Set("date_posted", "week")

	if req := experienceToRequirement(experience); req != "" {
		params.Set("job_requirements", req)
	}

	reqURL := "https://jsearch.p.rapidapi.com/search?" + params.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-RapidAPI-Key", p.apiKey)
	req.Header.Set("X-RapidAPI-Host", "jsearch.p.rapidapi.com")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("JSearch request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("JSearch rate limit exceeded — free tier allows 200 req/month")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("JSearch API returned status %d", resp.StatusCode)
	}

	var result jsearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode JSearch response: %w", err)
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("JSearch API status: %s", result.Status)
	}

	var jobs []*models.Job
	for _, j := range result.Data {
		applyLink := j.JobApplyLink
		if applyLink == "" {
			applyLink = j.JobGoogleLink
		}
		if applyLink == "" {
			continue // skip jobs with no valid apply link
		}

		jobLocation := strings.TrimSpace(fmt.Sprintf("%s %s %s", j.JobCity, j.JobState, j.JobCountry))
		if jobLocation == "" {
			jobLocation = location
		}

		desc := j.JobDescription
		if len(desc) > 500 {
			desc = desc[:497] + "..."
		}

		postedAt := time.Now()
		if j.JobPostedAt > 0 {
			postedAt = time.Unix(j.JobPostedAt, 0)
		}

		source := j.JobSource
		if source == "" {
			source = "JSearch"
		}

		jobs = append(jobs, &models.Job{
			Title:       j.JobTitle,
			Company:     j.EmployerName,
			Location:    jobLocation,
			Domain:      domain,
			Description: desc,
			PostedAt:    postedAt,
			SourceURL:   applyLink, // REAL apply link — no fake paths
			Source:      source,
		})
	}

	log.Printf("[JSearch] Fetched %d real jobs for '%s' in '%s'", len(jobs), domain, location)
	return jobs, nil
}

package services

import (
	"fmt"
	"job-portal/config"
	"job-portal/models"
	"log"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridService struct {
	cfg    *config.Config
	client *sendgrid.Client
}

func NewSendGridService(cfg *config.Config) *SendGridService {
	return &SendGridService{
		cfg:    cfg,
		client: sendgrid.NewSendClient(cfg.SendGridAPIKey),
	}
}

func (s *SendGridService) SendJobNotifications(userEmail string, jobs []*models.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	from := mail.NewEmail("JobHub Alerts", s.cfg.SMTPUser)
	to := mail.NewEmail("", userEmail)
	subject := fmt.Sprintf("🎯 %d New Job Opportunities Matching Your Preferences", len(jobs))
	
	htmlContent := s.buildEmailHTML(jobs)
	plainContent := s.buildEmailPlain(jobs)
	
	message := mail.NewSingleEmail(from, subject, to, plainContent, htmlContent)
	
	response, err := s.client.Send(message)
	if err != nil {
		log.Printf("SendGrid error: %v", err)
		return err
	}
	
	if response.StatusCode >= 400 {
		log.Printf("SendGrid status: %d, body: %s", response.StatusCode, response.Body)
		return fmt.Errorf("sendgrid error: %d", response.StatusCode)
	}
	
	log.Printf("Email sent successfully to %s (Status: %d)", userEmail, response.StatusCode)
	return nil
}

func (s *SendGridService) buildEmailHTML(jobs []*models.Job) string {
	var sb strings.Builder
	
	sb.WriteString(`
<!DOCTYPE html>
<html>
<head>
<style>
body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
.container { max-width: 600px; margin: 0 auto; padding: 20px; }
.header { background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
.job-card { background: #f9fafb; border: 1px solid #e5e7eb; padding: 20px; margin: 15px 0; border-radius: 8px; }
.job-title { color: #6366f1; font-size: 18px; font-weight: bold; margin: 0 0 10px 0; }
.job-meta { color: #6b7280; font-size: 14px; margin: 5px 0; }
.job-meta strong { color: #374151; }
.apply-btn { display: inline-block; background: #6366f1; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin-top: 10px; }
.footer { text-align: center; color: #6b7280; font-size: 12px; margin-top: 30px; padding-top: 20px; border-top: 1px solid #e5e7eb; }
</style>
</head>
<body>
<div class="container">
<div class="header">
<h1>🎯 New Job Opportunities</h1>
<p>We found jobs matching your preferences!</p>
</div>
`)
	
	for _, job := range jobs {
		sb.WriteString(fmt.Sprintf(`
<div class="job-card">
<h2 class="job-title">%s</h2>
<div class="job-meta"><strong>Company:</strong> %s</div>
<div class="job-meta"><strong>Location:</strong> %s</div>
<div class="job-meta"><strong>Domain:</strong> %s</div>
<div class="job-meta"><strong>Posted:</strong> %s</div>
<p>%s</p>
<a href="%s" class="apply-btn">View & Apply →</a>
</div>
`, job.Title, job.Company, job.Location, job.Domain, 
			job.PostedAt.Format("Jan 02, 2006"), job.Description, job.SourceURL))
	}
	
	sb.WriteString(`
<div class="footer">
<p>You're receiving this because you subscribed to JobHub alerts.</p>
<p>To update preferences or unsubscribe, visit your dashboard.</p>
</div>
</div>
</body>
</html>
`)
	
	return sb.String()
}

func (s *SendGridService) buildEmailPlain(jobs []*models.Job) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("New Job Opportunities - %d Jobs Found\n\n", len(jobs)))
	
	for i, job := range jobs {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, job.Title))
		sb.WriteString(fmt.Sprintf("   Company: %s\n", job.Company))
		sb.WriteString(fmt.Sprintf("   Location: %s\n", job.Location))
		sb.WriteString(fmt.Sprintf("   Domain: %s\n", job.Domain))
		sb.WriteString(fmt.Sprintf("   Apply: %s\n\n", job.SourceURL))
	}
	
	return sb.String()
}

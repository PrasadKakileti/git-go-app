package services

import (
	"fmt"
	"job-portal/config"
	"job-portal/models"
	"net/smtp"
	"strings"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) SendJobNotifications(userEmail string, jobs []*models.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	subject := fmt.Sprintf("Latest Job Updates - %d New Opportunities", len(jobs))
	body := s.buildEmailBody(jobs)

	return s.sendEmail(userEmail, subject, body)
}

func (s *EmailService) buildEmailBody(jobs []*models.Job) string {
	var sb strings.Builder
	
	sb.WriteString("<html><body style='font-family: Arial, sans-serif;'>")
	sb.WriteString("<h2 style='color: #2c3e50;'>Latest Job Opportunities</h2>")
	
	for _, job := range jobs {
		sb.WriteString("<div style='border: 1px solid #ddd; padding: 15px; margin: 10px 0; border-radius: 5px;'>")
		sb.WriteString(fmt.Sprintf("<h3 style='color: #1a4968ff; margin: 0;'>%s</h3>", job.Title))
		sb.WriteString(fmt.Sprintf("<p style='margin: 5px 0;'><strong>Company:</strong> %s</p>", job.Company))
		sb.WriteString(fmt.Sprintf("<p style='margin: 5px 0;'><strong>Location:</strong> %s</p>", job.Location))
		sb.WriteString(fmt.Sprintf("<p style='margin: 5px 0;'><strong>Domain:</strong> %s</p>", job.Domain))
		sb.WriteString(fmt.Sprintf("<p style='margin: 5px 0;'><strong>Contact:</strong> %s</p>", job.EmailContact))
		sb.WriteString(fmt.Sprintf("<p style='margin: 5px 0;'><strong>Posted:</strong> %s</p>", job.PostedAt.Format("Jan 02, 2006 15:04")))
		sb.WriteString(fmt.Sprintf("<p style='margin: 10px 0;'>%s</p>", job.Description))
		sb.WriteString(fmt.Sprintf("<a href='%s' style='color: #3498db;'>View Job</a>", job.SourceURL))
		sb.WriteString("</div>")
	}
	
	sb.WriteString("</body></html>")
	return sb.String()
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.SMTPUser, []string{to}, msg)
}

package services

import (
	"fmt"
	"job-portal/config"
	"job-portal/models"
	"log"
	"net/smtp"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type UnifiedEmailService struct {
	cfg            *config.Config
	useSendGrid    bool
	sendGridClient *sendgrid.Client
}

func NewUnifiedEmailService(cfg *config.Config) *UnifiedEmailService {
	useSendGrid := cfg.SendGridAPIKey != "" && cfg.SendGridAPIKey != "your-sendgrid-api-key-here"
	
	var client *sendgrid.Client
	if useSendGrid {
		client = sendgrid.NewSendClient(cfg.SendGridAPIKey)
		log.Println("Email service: Using SendGrid")
	} else {
		log.Println("Email service: Using SMTP")
	}
	
	return &UnifiedEmailService{
		cfg:            cfg,
		useSendGrid:    useSendGrid,
		sendGridClient: client,
	}
}

func (s *UnifiedEmailService) SendJobNotifications(userEmail string, jobs []*models.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	if s.useSendGrid {
		return s.sendViaSendGrid(userEmail, jobs)
	}
	return s.sendViaSMTP(userEmail, jobs)
}

func (s *UnifiedEmailService) sendViaSendGrid(userEmail string, jobs []*models.Job) error {
	from := mail.NewEmail("JobHub Alerts", s.cfg.SMTPUser)
	to := mail.NewEmail("", userEmail)
	subject := fmt.Sprintf("🎯 %d New Job Opportunities", len(jobs))
	
	htmlContent := s.buildEmailHTML(jobs)
	plainContent := s.buildEmailPlain(jobs)
	
	message := mail.NewSingleEmail(from, subject, to, plainContent, htmlContent)
	
	response, err := s.sendGridClient.Send(message)
	if err != nil {
		return err
	}
	
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: %d", response.StatusCode)
	}
	
	log.Printf("✅ Email sent via SendGrid to %s", userEmail)
	return nil
}

func (s *UnifiedEmailService) sendViaSMTP(userEmail string, jobs []*models.Job) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	
	subject := fmt.Sprintf("🎯 %d New Job Opportunities", len(jobs))
	htmlContent := s.buildEmailHTML(jobs)
	
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", userEmail, subject, htmlContent))

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	err := smtp.SendMail(addr, auth, s.cfg.SMTPUser, []string{userEmail}, msg)
	
	if err == nil {
		log.Printf("✅ Email sent via SMTP to %s", userEmail)
	}
	return err
}

func (s *UnifiedEmailService) buildEmailHTML(jobs []*models.Job) string {
	var sb strings.Builder
	
	sb.WriteString(`<!DOCTYPE html><html><head><style>
body{font-family:Arial,sans-serif;line-height:1.6;color:#333;background:#f9fafb;margin:0;padding:20px}
.container{max-width:600px;margin:0 auto;background:white;border-radius:10px;overflow:hidden;box-shadow:0 4px 6px rgba(0,0,0,0.1)}
.header{background:linear-gradient(135deg,#6366f1 0%,#8b5cf6 100%);color:white;padding:30px;text-align:center}
.header h1{margin:0;font-size:24px}
.content{padding:20px}
.job-card{background:#f9fafb;border-left:4px solid #6366f1;padding:15px;margin:15px 0;border-radius:5px}
.job-title{color:#6366f1;font-size:18px;font-weight:bold;margin:0 0 10px 0}
.job-meta{color:#6b7280;font-size:14px;margin:5px 0}
.apply-btn{display:inline-block;background:#6366f1;color:white;padding:10px 20px;text-decoration:none;border-radius:5px;margin-top:10px}
.footer{text-align:center;color:#6b7280;font-size:12px;padding:20px;border-top:1px solid #e5e7eb}
</style></head><body><div class="container">
<div class="header"><h1>🎯 New Job Opportunities</h1><p>Jobs matching your preferences</p></div>
<div class="content">`)
	
	for _, job := range jobs {
		sb.WriteString(fmt.Sprintf(`<div class="job-card">
<h2 class="job-title">%s</h2>
<div class="job-meta">🏢 <strong>%s</strong></div>
<div class="job-meta">📍 %s | 💼 %s</div>
<div class="job-meta">📅 %s</div>
<p>%s</p>
<a href="%s" class="apply-btn">View & Apply →</a>
</div>`, job.Title, job.Company, job.Location, job.Domain,
			job.PostedAt.Format("Jan 02, 2006"), job.Description, job.SourceURL))
	}
	
	sb.WriteString(`</div><div class="footer">
<p>You're receiving this because you subscribed to JobHub alerts.</p>
<p>Powered by JobHub</p></div></div></body></html>`)
	
	return sb.String()
}

func (s *UnifiedEmailService) buildEmailPlain(jobs []*models.Job) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("New Job Opportunities - %d Jobs\n\n", len(jobs)))
	
	for i, job := range jobs {
		sb.WriteString(fmt.Sprintf("%d. %s at %s\n   Location: %s | Domain: %s\n   %s\n\n",
			i+1, job.Title, job.Company, job.Location, job.Domain, job.SourceURL))
	}
	
	return sb.String()
}


func (s *UnifiedEmailService) SendWelcomeEmail(email, location, domain, frequency string) error {
	subject := "🎉 Welcome to JobHub - Your Job Alerts Are Active!"
	
	htmlContent := fmt.Sprintf(`<!DOCTYPE html><html><head><style>
body{font-family:Arial,sans-serif;line-height:1.6;color:#333;background:#f9fafb;margin:0;padding:20px}
.container{max-width:600px;margin:0 auto;background:white;border-radius:10px;overflow:hidden;box-shadow:0 4px 6px rgba(0,0,0,0.1)}
.header{background:linear-gradient(135deg,#6366f1 0%%,#8b5cf6 100%%);color:white;padding:40px;text-align:center}
.content{padding:30px}
.info-box{background:#f3f4f6;padding:15px;border-radius:8px;margin:20px 0}
.footer{text-align:center;color:#6b7280;font-size:12px;padding:20px}
</style></head><body><div class="container">
<div class="header"><h1>🎉 Welcome to JobHub!</h1><p>Your job alerts are now active</p></div>
<div class="content">
<h2>Registration Confirmed</h2>
<p>Thank you for registering! You'll now receive job alerts matching your preferences.</p>
<div class="info-box">
<strong>Your Preferences:</strong><br>
📧 Email: %s<br>
📍 Location: %s<br>
💼 Domain: %s<br>
⏰ Frequency: %s
</div>
<h3>What Happens Next?</h3>
<ul>
<li>We scrape Naukri every 6 hours for new jobs</li>
<li>You'll receive %s emails with matching jobs</li>
<li>Only jobs you haven't seen before</li>
<li>Direct links to apply on Naukri</li>
</ul>
<p><strong>First email:</strong> You'll receive your first job alert at the next scheduled time (9 AM).</p>
</div>
<div class="footer">
<p>Powered by JobHub | Your Career Companion</p>
</div></div></body></html>`, email, location, domain, frequency, frequency)
	
	plainContent := fmt.Sprintf(`Welcome to JobHub!

Your job alerts are now active.

Your Preferences:
- Email: %s
- Location: %s
- Domain: %s
- Frequency: %s

You'll receive %s emails with jobs matching your preferences.

First email will arrive at the next scheduled time (9 AM).

Powered by JobHub`, email, location, domain, frequency, frequency)

	if s.useSendGrid {
		return s.sendWelcomeViaSendGrid(email, subject, plainContent, htmlContent)
	}
	return s.sendWelcomeViaSMTP(email, subject, htmlContent)
}

func (s *UnifiedEmailService) sendWelcomeViaSendGrid(email, subject, plainContent, htmlContent string) error {
	from := mail.NewEmail("JobHub", s.cfg.SMTPUser)
	to := mail.NewEmail("", email)
	message := mail.NewSingleEmail(from, subject, to, plainContent, htmlContent)
	
	response, err := s.sendGridClient.Send(message)
	if err != nil {
		return err
	}
	
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: %d", response.StatusCode)
	}
	
	log.Printf("✅ Welcome email sent via SendGrid to %s", email)
	return nil
}

func (s *UnifiedEmailService) sendWelcomeViaSMTP(email, subject, htmlContent string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", email, subject, htmlContent))

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	err := smtp.SendMail(addr, auth, s.cfg.SMTPUser, []string{email}, msg)
	
	if err == nil {
		log.Printf("✅ Welcome email sent via SMTP to %s", email)
	}
	return err
}

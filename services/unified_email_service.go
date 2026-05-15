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
		log.Println("[Email] Using SendGrid")
	} else {
		log.Println("[Email] Using SMTP (Gmail)")
	}

	return &UnifiedEmailService{cfg: cfg, useSendGrid: useSendGrid, sendGridClient: client}
}

// ─────────────────────────────────────────────────────────────
// JOB NOTIFICATION EMAIL
// ─────────────────────────────────────────────────────────────

func (s *UnifiedEmailService) SendJobNotifications(userEmail string, jobs []*models.Job) error {
	if len(jobs) == 0 {
		return nil
	}
	subject := fmt.Sprintf("🎯 %d New Job Matches For You", len(jobs))
	html := s.buildJobEmailHTML(jobs)
	plain := s.buildJobEmailPlain(jobs)

	if s.useSendGrid {
		return s.sendViaSendGrid(userEmail, subject, plain, html)
	}
	return s.sendViaSMTP(userEmail, subject, html)
}

func (s *UnifiedEmailService) buildJobEmailHTML(jobs []*models.Job) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html><html><head><meta charset="UTF-8">
<style>
body{margin:0;padding:0;background:#f3f4f6;font-family:Arial,sans-serif}
.wrap{max-width:620px;margin:30px auto;background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 4px 20px rgba(0,0,0,.1)}
.hdr{background:linear-gradient(135deg,#4f46e5,#7c3aed);color:#fff;padding:32px 24px;text-align:center}
.hdr h1{margin:0 0 6px;font-size:22px}
.hdr p{margin:0;opacity:.85;font-size:14px}
.body{padding:20px 24px}
.card{background:#f9fafb;border-left:4px solid #4f46e5;border-radius:6px;padding:18px;margin:14px 0}
.title{color:#1e1b4b;font-size:17px;font-weight:700;margin:0 0 8px}
.meta{font-size:13px;color:#6b7280;margin:3px 0}
.badge{display:inline-block;background:#ede9fe;color:#5b21b6;font-size:11px;font-weight:600;padding:2px 8px;border-radius:20px;margin-bottom:8px}
.desc{font-size:13px;color:#374151;margin:10px 0;line-height:1.5}
.apply{display:inline-block;background:#4f46e5;color:#fff;text-decoration:none;padding:10px 22px;border-radius:6px;font-size:14px;font-weight:600;margin-top:8px}
.apply:hover{background:#4338ca}
.footer{text-align:center;padding:20px;color:#9ca3af;font-size:12px;border-top:1px solid #e5e7eb}
</style></head><body><div class="wrap">
<div class="hdr"><h1>🎯 New Job Opportunities</h1><p>Fresh matches based on your profile</p></div>
<div class="body">`)

	for _, job := range jobs {
		source := job.Source
		if source == "" {
			source = "Job Board"
		}
		sb.WriteString(fmt.Sprintf(`<div class="card">
<span class="badge">%s</span>
<div class="title">%s</div>
<div class="meta">🏢 <strong>%s</strong></div>
<div class="meta">📍 %s &nbsp;|&nbsp; 💼 %s</div>
<div class="meta">📅 Posted: %s</div>
<div class="desc">%s</div>
<a href="%s" class="apply" target="_blank">Apply Now →</a>
</div>`,
			htmlEscape(source),
			htmlEscape(job.Title),
			htmlEscape(job.Company),
			htmlEscape(job.Location),
			htmlEscape(job.Domain),
			job.PostedAt.Format("Jan 02, 2006"),
			htmlEscape(job.Description),
			job.SourceURL, // real link — LinkedIn / Indeed / Glassdoor URL
		))
	}

	sb.WriteString(`</div>
<div class="footer">
  <p>You're receiving this because you subscribed to JobHub alerts.</p>
  <p>© JobHub — Your Career Companion</p>
</div></div></body></html>`)

	return sb.String()
}

func (s *UnifiedEmailService) buildJobEmailPlain(jobs []*models.Job) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("New Job Opportunities (%d jobs)\n\n", len(jobs)))
	for i, job := range jobs {
		sb.WriteString(fmt.Sprintf("%d. %s — %s\n   %s | %s\n   Apply: %s\n\n",
			i+1, job.Title, job.Company, job.Location, job.Domain, job.SourceURL))
	}
	return sb.String()
}

// ─────────────────────────────────────────────────────────────
// WELCOME EMAIL
// ─────────────────────────────────────────────────────────────

func (s *UnifiedEmailService) SendWelcomeEmail(email, location, domain, frequency string) error {
	subject := "🎉 Welcome to JobHub — Your Job Alerts Are Active!"
	html := s.buildWelcomeHTML(email, location, domain, frequency)
	plain := fmt.Sprintf("Welcome to JobHub!\n\nYour alerts are active.\nLocation: %s | Domain: %s | Frequency: %s\n\nYou'll start receiving matching jobs within 10 minutes.", location, domain, frequency)

	if s.useSendGrid {
		return s.sendViaSendGrid(email, subject, plain, html)
	}
	return s.sendViaSMTP(email, subject, html)
}

func (s *UnifiedEmailService) buildWelcomeHTML(email, location, domain, frequency string) string {
	return fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="UTF-8">
<style>
body{margin:0;background:#f3f4f6;font-family:Arial,sans-serif}
.wrap{max-width:600px;margin:30px auto;background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 4px 20px rgba(0,0,0,.1)}
.hdr{background:linear-gradient(135deg,#4f46e5,#7c3aed);color:#fff;padding:36px 24px;text-align:center}
.hdr h1{margin:0 0 6px;font-size:24px}
.body{padding:28px 24px}
.box{background:#f5f3ff;border-radius:8px;padding:16px 20px;margin:18px 0}
.box p{margin:4px 0;font-size:14px;color:#374151}
ul{padding-left:20px;color:#374151;font-size:14px;line-height:1.8}
.footer{text-align:center;padding:20px;color:#9ca3af;font-size:12px;border-top:1px solid #e5e7eb}
</style></head><body><div class="wrap">
<div class="hdr"><h1>🎉 Welcome to JobHub!</h1><p>Your job alerts are now active</p></div>
<div class="body">
<h2 style="color:#1e1b4b">Registration Confirmed</h2>
<p style="color:#374151">Thank you for joining! You'll receive job alerts matching your profile.</p>
<div class="box">
<p>📧 <strong>Email:</strong> %s</p>
<p>📍 <strong>Location:</strong> %s</p>
<p>💼 <strong>Domain:</strong> %s</p>
<p>⏰ <strong>Frequency:</strong> %s</p>
</div>
<h3 style="color:#1e1b4b">What Happens Next?</h3>
<ul>
<li>We fetch real jobs from <strong>LinkedIn, Indeed, Glassdoor</strong> and more</li>
<li>First email arrives within the hour (or instantly if you clicked Refresh Jobs)</li>
<li>Every job link takes you <strong>directly to the apply page</strong></li>
<li>Only new jobs you haven't seen before</li>
</ul>
</div>
<div class="footer"><p>© JobHub — Your Career Companion</p></div>
</div></body></html>`, email, location, domain, frequency)
}

// ─────────────────────────────────────────────────────────────
// TRANSPORT
// ─────────────────────────────────────────────────────────────

func (s *UnifiedEmailService) sendViaSendGrid(to, subject, plain, html string) error {
	from := mail.NewEmail("JobHub", s.cfg.SMTPUser) // display name shown in inbox: "JobHub"
	msg := mail.NewSingleEmail(mail.NewEmail("", to), subject, mail.NewEmail("", to), plain, html)
	msg.SetFrom(from)

	resp, err := s.sendGridClient.Send(msg)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error %d: %s", resp.StatusCode, resp.Body)
	}
	log.Printf("[Email/SendGrid] sent to %s", to)
	return nil
}

func (s *UnifiedEmailService) sendViaSMTP(to, subject, html string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	// "From: JobHub <address>" sets the display name visible in the inbox
	fromHeader := fmt.Sprintf("\"JobHub\" <%s>", s.cfg.SMTPUser)
	msg := []byte("From: " + fromHeader + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		html + "\r\n")

	addr := s.cfg.SMTPHost + ":" + s.cfg.SMTPPort
	if err := smtp.SendMail(addr, auth, s.cfg.SMTPUser, []string{to}, msg); err != nil {
		return err
	}
	log.Printf("[Email/SMTP] sent to %s", to)
	return nil
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&#34;")
	return s
}

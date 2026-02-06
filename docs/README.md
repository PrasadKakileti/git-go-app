# JobHub - Job Portal Application Documentation

## 📚 Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Application Flow](#application-flow)
4. [Setup Guide](#setup-guide)
5. [API Documentation](#api-documentation)
6. [Email System](#email-system)
7. [Database Schema](#database-schema)
8. [Deployment](#deployment)

---

## Overview

JobHub is an automated job alert system that:
- Scrapes job listings from Naukri.com
- Sends personalized email alerts to registered users
- Filters jobs by location and domain
- Supports daily and weekly notification frequencies

### Key Features
- ✅ Real-time job scraping from Naukri
- ✅ Domain-based filtering (Java, Golang, Chartered Accountant)
- ✅ Location-based filtering (10+ Indian cities)
- ✅ Email notifications (SMTP/SendGrid)
- ✅ Automated scheduling (cron jobs)
- ✅ Admin dashboard
- ✅ Mobile-responsive UI

---

## Architecture

### Technology Stack
- **Backend:** Golang 1.21+
- **Database:** MySQL 8.0+
- **Email:** SMTP (Gmail) / SendGrid
- **Frontend:** HTML, CSS, JavaScript
- **Scheduler:** robfig/cron
- **Scraping:** goquery

### Project Structure
```
Go-app/
├── main.go                 # Application entry point
├── config/                 # Configuration management
│   └── config.go
├── database/               # Database connection
│   └── db.go
├── handlers/               # HTTP request handlers
│   ├── handlers.go
│   └── list_users.go
├── models/                 # Data models
│   └── models.go
├── repository/             # Database operations
│   ├── user_repo.go
│   └── job_repo.go
├── services/               # Business logic
│   ├── email_service.go
│   ├── unified_email_service.go
│   └── job_service.go
├── scraper/                # Job scraping
│   └── naukri_scraper.go
├── scheduler/              # Cron jobs
│   └── scheduler.go
├── frontend/               # UI files
│   ├── index.html
│   ├── admin.html
│   ├── css/
│   └── js/
├── docs/                   # Documentation
└── .env                    # Environment variables
```

---

## Application Flow

### 1. User Registration Flow
```
User fills form → POST /api/register → Validate data → 
Save to database → Send welcome email → Return success
```

**Steps:**
1. User visits http://localhost:8080
2. Fills registration form (email, location, domain, frequency)
3. Clicks "Start Receiving Jobs"
4. Backend validates and saves to MySQL
5. Welcome email sent immediately
6. Success message displayed

### 2. Job Scraping Flow
```
Cron triggers (every 6 hours) → Scrape Naukri → 
Parse job data → Store in database → Log results
```

**Steps:**
1. Scheduler runs at: 00:00, 06:00, 12:00, 18:00
2. For each location (Bangalore, Mumbai, Delhi, etc.)
3. For each domain (Java, Golang, Chartered Accountant)
4. Scrape Naukri.com using goquery
5. Extract: title, company, location, domain, URL
6. Store in `jobs` table
7. Log: "Scraped X jobs for Y location"

### 3. Email Notification Flow
```
Cron triggers (9 AM daily/weekly) → Get active users → 
For each user: Get unsent jobs → Send email → Mark as sent
```

**Steps:**
1. Daily: 9 AM every day
2. Weekly: 9 AM every Monday
3. Query users by notification frequency
4. For each user:
   - Get jobs matching location + domain
   - Filter out already-sent jobs
   - Send email with job list
   - Mark jobs as sent in `user_job_sent` table
5. Log: "Sent X jobs to user@email.com"

### 4. Admin Dashboard Flow
```
User visits /admin → GET /api/users → 
Fetch all users → Display in table → Auto-refresh every 30s
```

---

## Setup Guide

### Prerequisites
```bash
# Check versions
go version        # 1.21+
mysql --version   # 8.0+
```

### Installation Steps

**1. Clone and Setup**
```bash
cd Go-app
go mod download
```

**2. Database Setup**
```bash
mysql -u root -p < schema.sql
```

**3. Configure Environment**
```bash
cp .env.example .env
# Edit .env with your credentials
```

**4. Run Application**
```bash
./start.sh
```

**5. Access**
- Main: http://localhost:8080
- Admin: http://localhost:8080/admin
- Network: http://YOUR_IP:8080

---

## API Documentation

### Endpoints

#### 1. Register User
```http
POST /api/register
Content-Type: application/json

{
  "email": "user@example.com",
  "location": "Bangalore",
  "domain": "Java",
  "notification_frequency": "daily"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Registration successful! Check your email for confirmation.",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "location": "Bangalore",
    "domain": "Java",
    "notification_frequency": "daily",
    "is_active": true
  }
}
```

#### 2. Get User
```http
GET /api/user?email=user@example.com
```

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "location": "Bangalore",
  "domain": "Java",
  "notification_frequency": "daily",
  "is_active": true,
  "created_at": "2024-01-30T12:00:00Z"
}
```

#### 3. List All Users
```http
GET /api/users
```

**Response:**
```json
{
  "total": 10,
  "users": [...]
}
```

#### 4. Health Check
```http
GET /api/health
```

**Response:**
```json
{
  "status": "ok"
}
```

---

## Email System

### Configuration

**Option 1: Gmail SMTP**
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

**Option 2: SendGrid**
```env
SENDGRID_API_KEY=SG.your-key-here
SMTP_USER=verified-sender@gmail.com
```

### Email Types

**1. Welcome Email**
- Sent immediately on registration
- Confirms preferences
- Explains what to expect

**2. Job Alert Email**
- Sent daily (9 AM) or weekly (Monday 9 AM)
- Lists matching jobs
- Includes: title, company, location, domain, apply link

### Email Features
- HTML formatted
- Mobile responsive
- Plain text fallback
- Automatic retry on failure
- Detailed logging

---

## Database Schema

### Tables

**users**
```sql
id                      INT PRIMARY KEY AUTO_INCREMENT
email                   VARCHAR(255) UNIQUE NOT NULL
location                VARCHAR(100)
domain                  VARCHAR(100)
notification_frequency  ENUM('daily', 'weekly')
is_active               BOOLEAN DEFAULT TRUE
created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

**jobs**
```sql
id              INT PRIMARY KEY AUTO_INCREMENT
title           VARCHAR(255) NOT NULL
company         VARCHAR(255)
location        VARCHAR(100)
domain          VARCHAR(100)
email_contact   VARCHAR(255)
description     TEXT
posted_at       TIMESTAMP
scraped_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
source_url      VARCHAR(500)
```

**user_job_sent**
```sql
id          INT PRIMARY KEY AUTO_INCREMENT
user_id     INT FOREIGN KEY → users(id)
job_id      INT FOREIGN KEY → jobs(id)
sent_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
UNIQUE(user_id, job_id)
```

---

## Deployment

### Production Checklist
- [ ] Configure SendGrid API key
- [ ] Set strong database password
- [ ] Enable HTTPS
- [ ] Set up monitoring
- [ ] Configure backup
- [ ] Test email delivery
- [ ] Verify cron jobs running

### Monitoring
```bash
# Check logs
tail -f /tmp/jobhub.log

# Check database
mysql -u root -p job_portal

# Check running processes
lsof -i:8080
```

---

## Additional Documentation

- [Setup Guide](./SETUP.md)
- [Email Configuration](./EMAIL_SETUP.md)
- [API Reference](./API.md)
- [Troubleshooting](./TROUBLESHOOTING.md)
- [FAQ](./FAQ.md)

---

**Version:** 1.0.0  
**Last Updated:** 2024-01-30  
**Maintainer:** JobHub Team

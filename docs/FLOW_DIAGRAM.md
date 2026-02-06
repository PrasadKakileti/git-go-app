# Application Flow Diagram

## Complete System Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                         USER INTERACTION                         │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  1. USER REGISTRATION                                            │
│                                                                  │
│  Browser → http://localhost:8080                                │
│  ├─ Fill Form (email, location, domain, frequency)             │
│  ├─ Click "Start Receiving Jobs"                               │
│  └─ POST /api/register                                          │
│                                                                  │
│  Backend (handlers/handlers.go)                                 │
│  ├─ Validate input                                              │
│  ├─ Save to MySQL (users table)                                 │
│  ├─ Send welcome email (immediate)                              │
│  └─ Return success JSON                                         │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. WELCOME EMAIL (Immediate)                                    │
│                                                                  │
│  services/unified_email_service.go                              │
│  ├─ Build HTML email                                            │
│  ├─ Include user preferences                                    │
│  ├─ Send via SMTP/SendGrid                                      │
│  └─ Log: "✅ Welcome email sent to user@email.com"             │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. JOB SCRAPING (Every 6 hours: 00:00, 06:00, 12:00, 18:00)   │
│                                                                  │
│  scheduler/scheduler.go                                         │
│  └─ Cron: "0 */6 * * *"                                         │
│                                                                  │
│  services/job_service.go → ScrapeAndStoreJobs()                │
│  ├─ For each location (Bangalore, Mumbai, Delhi...)            │
│  │   └─ For each domain (Java, Golang, CA)                     │
│  │       ├─ scraper/naukri_scraper.go                          │
│  │       ├─ Fetch: https://www.naukri.com/Java-jobs-in-Bangalore│
│  │       ├─ Parse HTML with goquery                            │
│  │       ├─ Extract: title, company, location, URL             │
│  │       └─ Store in jobs table                                │
│  └─ Log: "Scraped X jobs for Y location"                       │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  4. DAILY NOTIFICATIONS (Every day at 9 AM)                     │
│                                                                  │
│  scheduler/scheduler.go                                         │
│  └─ Cron: "0 9 * * *"                                           │
│                                                                  │
│  services/job_service.go → SendNotifications("daily")          │
│  ├─ Get all users with frequency="daily"                       │
│  ├─ For each user:                                              │
│  │   ├─ Get jobs from last 24 hours                            │
│  │   ├─ Filter by user's location                              │
│  │   ├─ Filter by user's domain                                │
│  │   ├─ Exclude already-sent jobs (user_job_sent table)        │
│  │   ├─ If jobs found:                                         │
│  │   │   ├─ Send email with job list                           │
│  │   │   ├─ Mark jobs as sent                                  │
│  │   │   └─ Log: "Sent X jobs to user@email.com"              │
│  │   └─ Else: Skip (no new jobs)                               │
│  └─ Repeat for all daily users                                 │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  5. WEEKLY NOTIFICATIONS (Every Monday at 9 AM)                 │
│                                                                  │
│  scheduler/scheduler.go                                         │
│  └─ Cron: "0 9 * * 1"                                           │
│                                                                  │
│  services/job_service.go → SendNotifications("weekly")         │
│  ├─ Get all users with frequency="weekly"                      │
│  ├─ For each user:                                              │
│  │   ├─ Get jobs from last 7 days                              │
│  │   ├─ Filter by location + domain                            │
│  │   ├─ Exclude already-sent jobs                              │
│  │   └─ Send email if jobs found                               │
│  └─ Same process as daily                                      │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  6. ADMIN DASHBOARD (Real-time monitoring)                      │
│                                                                  │
│  Browser → http://localhost:8080/admin                         │
│  ├─ GET /api/users                                              │
│  ├─ Display all registered users                               │
│  ├─ Show: email, location, domain, frequency, status           │
│  ├─ Auto-refresh every 30 seconds                              │
│  └─ Statistics: Total users, Daily, Weekly                     │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Browser │────▶│ Handlers │────▶│Repository│────▶│  MySQL   │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                       │                                  │
                       │                                  │
                       ▼                                  ▼
                 ┌──────────┐                      ┌──────────┐
                 │ Services │◀─────────────────────│  Jobs    │
                 └──────────┘                      │  Users   │
                       │                           │user_job_ │
                       │                           │  sent    │
                       ▼                           └──────────┘
                 ┌──────────┐
                 │  Email   │
                 │ Service  │
                 └──────────┘
                       │
                       ▼
                 ┌──────────┐
                 │   SMTP   │
                 │SendGrid  │
                 └──────────┘
```

## Component Interaction

### 1. Registration Flow
```
User Input → Validation → Database → Email → Response
   ↓            ↓            ↓         ↓        ↓
 Form      handlers.go   user_repo  email    JSON
                                    service
```

### 2. Scraping Flow
```
Scheduler → Job Service → Scraper → Parser → Database
    ↓           ↓            ↓         ↓         ↓
  Cron    ScrapeAndStore  Naukri   goquery   jobs
                                              table
```

### 3. Notification Flow
```
Scheduler → Job Service → Repository → Email Service → User
    ↓           ↓             ↓            ↓           ↓
  Cron    SendNotifications  Get Jobs   Build HTML  Inbox
                            Filter      Send SMTP
```

## File Responsibilities

| File | Responsibility |
|------|---------------|
| `main.go` | Initialize app, setup routes, start server |
| `config/config.go` | Load environment variables |
| `database/db.go` | MySQL connection |
| `handlers/handlers.go` | HTTP request handling |
| `models/models.go` | Data structures |
| `repository/user_repo.go` | User database operations |
| `repository/job_repo.go` | Job database operations |
| `services/job_service.go` | Business logic |
| `services/unified_email_service.go` | Email sending |
| `scraper/naukri_scraper.go` | Web scraping |
| `scheduler/scheduler.go` | Cron jobs |

## Timing Diagram

```
Time    | Event
--------|--------------------------------------------------
00:00   | Job scraping starts
00:05   | Jobs stored in database
06:00   | Job scraping starts
06:05   | Jobs stored in database
09:00   | Daily notifications sent
09:00   | Weekly notifications sent (if Monday)
12:00   | Job scraping starts
12:05   | Jobs stored in database
18:00   | Job scraping starts
18:05   | Jobs stored in database
```

## Error Handling Flow

```
Error Occurs
    │
    ├─ Log error message
    │
    ├─ Return error to caller
    │
    ├─ Continue with next item (don't crash)
    │
    └─ User sees friendly error message
```

## Security Flow

```
Request → Validate Input → Sanitize → Process → Response
   ↓           ↓             ↓          ↓         ↓
 HTTPS      Required      SQL Safe   Business   JSON
           Fields        Queries     Logic
```

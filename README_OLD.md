# Job Portal Platform

A job portal that scrapes Naukri for latest job listings and sends email notifications to registered users.

## Features

- User registration with email
- Location-based job filtering
- Daily/Weekly email notifications
- Real-time job data from Naukri
- Automated job scraping
- MySQL database storage

## Setup

### Prerequisites

- Go 1.21+
- MySQL 8.0+
- SMTP credentials (Gmail recommended)

### Installation

1. Clone and navigate to project:
```bash
cd Go-app
```

2. Install dependencies:
```bash
go mod download
```

3. Setup MySQL database:
```bash
mysql -u root -p < schema.sql
```

4. Configure environment:
```bash
cp .env.example .env
# Edit .env with your credentials
```

5. Run the application:
```bash
go run main.go
```

6. Access the application:
```
http://localhost:8080
```

## Configuration

Edit `.env` file:

- `DB_USER`, `DB_PASSWORD`: MySQL credentials
- `SMTP_USER`, `SMTP_PASSWORD`: Email service credentials (use Gmail App Password)
- `SERVER_PORT`: Application port (default: 8080)

## API Endpoints

- `POST /api/register` - Register new user
- `GET /api/user?email=<email>` - Get user details
- `GET /api/health` - Health check

## Scheduler

- Job scraping: Every 6 hours
- Daily notifications: 9 AM daily
- Weekly notifications: 9 AM every Monday

## Production Notes

Replace mock scraper in `scraper/naukri_scraper.go` with actual Naukri integration using:
- Naukri API (if available)
- Web scraping with goquery/colly libraries

## License

MIT

# Complete Job Alert System - Integration Guide

## 🎯 System Overview

Your job portal now includes:
1. **Real Naukri Scraping** with goquery (fallback to mock data)
2. **SendGrid Integration** (100 free emails/day)
3. **SMTP Fallback** (Gmail support)
4. **Automated Scheduling** (scrape every 6 hours, send daily/weekly)
5. **Domain-based Filtering** (Chartered Accountant, Java, Golang)

## 📋 Features Implemented

### ✅ Job Data Extraction
- Web scraping from Naukri.com using goquery
- Searches by domain + location
- Extracts: title, company, location, job link
- Automatic fallback to mock data if scraping fails
- Rate limiting and error handling

### ✅ User Preferences
- Email, Location, Domain selection
- Daily/Weekly notification frequency
- Stored in MySQL database
- Admin dashboard to view all users

### ✅ Email Alert System
**Two Options:**
1. **SendGrid** (Recommended)
   - 100 emails/day free
   - Professional delivery
   - Better inbox placement
   
2. **Gmail SMTP** (Backup)
   - Unlimited (with limits)
   - Requires app password
   - May hit spam filters

### ✅ Automated Scheduling
```
Every 6 hours  → Scrape new jobs from Naukri
Every day 9 AM → Send daily notifications
Every Mon 9 AM → Send weekly notifications
```

### ✅ Job Matching
- Filters by user's location
- Filters by user's domain
- Only sends unsent jobs
- Tracks sent jobs to avoid duplicates

## 🚀 Quick Start

### 1. Setup SendGrid (5 minutes)
```bash
# See SENDGRID_SETUP.md for detailed steps
# 1. Create account: https://signup.sendgrid.com/
# 2. Get API key
# 3. Verify sender email
# 4. Add to .env:
SENDGRID_API_KEY=SG.your-key-here
SMTP_USER=your-verified-email@gmail.com
```

### 2. Start Application
```bash
./start.sh
```

### 3. Register Users
- Visit: http://localhost:8080
- Or network: http://10.21.12.100:8080
- Fill form: email, location, domain, frequency

### 4. Monitor
```bash
# View registered users
bash view_users.sh

# Admin dashboard
http://localhost:8080/admin
```

## 📧 Email Template

Users receive beautiful HTML emails with:
- Job title and company
- Location and domain
- Posted date
- Job description
- "View & Apply" button linking to Naukri

## 🔧 Configuration

### Environment Variables (.env)
```bash
# Database
DB_USER=root
DB_PASSWORD=password
DB_NAME=job_portal

# Email (Choose one)
SENDGRID_API_KEY=SG.xxx  # Recommended
SMTP_USER=email@gmail.com
SMTP_PASSWORD=app-password

# Server
SERVER_PORT=8080
```

### Scheduler Settings
Edit `scheduler/scheduler.go`:
```go
// Job scraping frequency
"0 */6 * * *"  // Every 6 hours

// Daily notifications
"0 9 * * *"    // 9 AM daily

// Weekly notifications
"0 9 * * 1"    // 9 AM Monday
```

## 🎨 Customization

### Add More Domains
Edit `frontend/index.html`:
```html
<option value="Python">Python</option>
<option value="DevOps">DevOps</option>
```

### Add More Cities
Edit `frontend/index.html`:
```html
<option value="Jaipur">Jaipur</option>
```

### Change Email Template
Edit `services/unified_email_service.go`:
- Modify `buildEmailHTML()` function
- Customize colors, layout, content

## 📊 Database Schema

```sql
users:
- id, email, location, domain
- notification_frequency (daily/weekly)
- is_active, created_at

jobs:
- id, title, company, location, domain
- email_contact, description
- posted_at, source_url

user_job_sent:
- Tracks which jobs sent to which users
- Prevents duplicate notifications
```

## 🔍 How It Works

1. **User Registers**
   - Selects: email, location, domain, frequency
   - Stored in database

2. **Scheduler Runs**
   - Every 6 hours: Scrapes Naukri for all domains/locations
   - Stores new jobs in database

3. **Notification Time**
   - Gets active users by frequency (daily/weekly)
   - For each user:
     - Finds unsent jobs matching location + domain
     - Sends email with job list
     - Marks jobs as sent

4. **Email Delivery**
   - Uses SendGrid if configured
   - Falls back to SMTP
   - Logs success/failure

## 🛡️ Security & Privacy

- Email addresses encrypted in transit
- No passwords stored (email-only registration)
- GDPR compliant (users can unsubscribe)
- Rate limiting on scraping
- SQL injection protection

## 📈 Scaling

**Current Setup:**
- Handles 100 users/day (SendGrid free tier)
- Scrapes 3 domains × 10 cities = 30 searches every 6 hours
- ~10 jobs per search = 300 jobs/day

**To Scale:**
- Upgrade SendGrid plan ($15/month = 40K emails)
- Add more domains/cities
- Increase scraping frequency
- Add job caching layer

## 🐛 Troubleshooting

**Emails not sending?**
```bash
# Check logs
tail -f logs/app.log

# Test SendGrid
curl -X POST https://api.sendgrid.com/v3/mail/send \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Scraping not working?**
- Naukri may block requests
- Check User-Agent header
- Add delays between requests
- Use proxy/VPN if needed

**Database errors?**
```bash
# Reset database
mysql -u root -ppassword < schema.sql
```

## 📱 Access URLs

**Local:**
- Main: http://localhost:8080
- Admin: http://localhost:8080/admin

**Network (same WiFi):**
- Main: http://10.21.12.100:8080
- Admin: http://10.21.12.100:8080/admin

## 🎓 Next Steps

1. **Get SendGrid API key** (5 min)
2. **Test with your email** (register yourself)
3. **Wait for next scheduled run** (or trigger manually)
4. **Check your inbox!** 📬

## 💡 Tips

- Start with daily notifications for testing
- Monitor first few sends carefully
- Check spam folder initially
- Ask users to whitelist your sender email
- Use professional sender name: "JobHub Alerts"

---

**Need Help?**
- Check logs: `tail -f logs/app.log`
- Run diagnostics: `bash fix.sh`
- View users: `bash view_users.sh`

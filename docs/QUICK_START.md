# Quick Start Guide

## ✅ Email System Now Working!

### What Changed
- **Welcome Email**: Sent immediately on registration
- **Job Alerts**: Sent daily/weekly at 9 AM
- **Both emails**: Fully functional with your SMTP credentials

---

## Test Email Immediately

### Option 1: Register New User (Gets Welcome Email)
```bash
# Start server
./start.sh

# Visit in browser
http://localhost:8080

# Register with any email
# Welcome email sent immediately!
```

### Option 2: Send Test Job Alert
```bash
# Send job alert to babusurendra500@gmail.com
bash send_test_email.sh
```

---

## Email Types

### 1. Welcome Email (Immediate)
**Trigger:** User registers  
**Content:**
- Confirms registration
- Shows preferences (location, domain, frequency)
- Explains what happens next

**Example:**
```
Subject: 🎉 Welcome to JobHub - Your Job Alerts Are Active!

Your Preferences:
📧 Email: user@example.com
📍 Location: Bangalore
💼 Domain: Java
⏰ Frequency: daily

You'll receive daily emails with matching jobs.
First email will arrive at 9 AM tomorrow.
```

### 2. Job Alert Email (Scheduled)
**Trigger:** Daily 9 AM or Weekly Monday 9 AM  
**Content:**
- List of matching jobs
- Job title, company, location, domain
- "View & Apply" buttons

**Example:**
```
Subject: 🎯 5 New Job Opportunities

1. Senior Java Developer - Tech Corp
   Location: Bangalore | Domain: Java
   [View & Apply →]

2. Java Specialist - StartupXYZ
   Location: Bangalore | Domain: Java
   [View & Apply →]
```

---

## Current Configuration

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=prasadkakileti105@gmail.com
SMTP_PASSWORD=cpbx hgkj anzu vpsk
```

✅ **Status:** Working perfectly!

---

## Test Results

### Test 1: Welcome Email
```
Email: testuser123@example.com
Status: ✅ Sent successfully
Log: "✅ Welcome email sent via SMTP to testuser123@example.com"
```

### Test 2: Job Alert Email
```
Email: babusurendra500@gmail.com
Jobs: 2 Java jobs in Bangalore
Status: ✅ Sent successfully
Log: "✅ Email sent via SMTP to babusurendra500@gmail.com"
```

---

## Verify Email Delivery

### Check Logs
```bash
# Real-time logs
tail -f /tmp/jobhub.log

# Look for:
# "✅ Welcome email sent via SMTP to user@email.com"
# "✅ Email sent via SMTP to user@email.com"
```

### Check Inbox
1. **Welcome email**: Check immediately after registration
2. **Job alerts**: Check at 9 AM (daily/weekly)
3. **Spam folder**: Check if not in inbox

---

## Scheduled Email Times

| Frequency | Schedule | Cron Expression |
|-----------|----------|-----------------|
| Daily | Every day 9 AM | `0 9 * * *` |
| Weekly | Every Monday 9 AM | `0 9 * * 1` |

---

## Troubleshooting

### Email Not Received?

**1. Check logs:**
```bash
tail -f /tmp/jobhub.log
```

**2. Look for errors:**
- "Failed to send email" → SMTP issue
- "No jobs found" → No matching jobs

**3. Check spam folder**

**4. Verify SMTP credentials:**
```bash
# Test SMTP connection
telnet smtp.gmail.com 587
```

**5. Check user is active:**
```bash
mysql -u root -ppassword job_portal -e \
  "SELECT email, is_active FROM users WHERE email='your@email.com';"
```

---

## Manual Testing

### Send Welcome Email
```bash
# Register new user via API
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"test@example.com",
    "location":"Bangalore",
    "domain":"Java",
    "notification_frequency":"daily"
  }'

# Check logs
tail -f /tmp/jobhub.log
```

### Send Job Alert
```bash
# Use test script
bash send_test_email.sh

# Or run test program
go run test_email.go
```

---

## Production Checklist

- [x] SMTP configured
- [x] Welcome email working
- [x] Job alert email working
- [x] Scheduler running
- [x] Database connected
- [x] Users can register
- [ ] SendGrid configured (optional upgrade)
- [ ] Monitoring setup
- [ ] Backup configured

---

## Next Steps

1. **Register yourself**: http://localhost:8080
2. **Check welcome email**: Should arrive immediately
3. **Wait for job alert**: Will arrive at 9 AM
4. **Monitor logs**: `tail -f /tmp/jobhub.log`
5. **View users**: http://localhost:8080/admin

---

## Support

**Documentation:**
- Main docs: `docs/README.md`
- Flow diagram: `docs/FLOW_DIAGRAM.md`
- Email setup: `SENDGRID_SETUP.md`

**Scripts:**
- Start server: `./start.sh`
- Test email: `bash send_test_email.sh`
- View users: `bash view_users.sh`
- Fix issues: `bash fix.sh`

---

**Status: ✅ FULLY OPERATIONAL**

Both welcome emails and job alerts are working!

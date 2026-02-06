# SendGrid Integration Guide

## Free Email Service Setup (100 emails/day free)

### Step 1: Create SendGrid Account
1. Go to: https://signup.sendgrid.com/
2. Sign up for FREE account
3. Verify your email address

### Step 2: Get API Key
1. Login to SendGrid dashboard
2. Go to: Settings → API Keys
3. Click "Create API Key"
4. Name: `JobHub-Production`
5. Permissions: Select "Full Access"
6. Click "Create & View"
7. **COPY THE API KEY** (you won't see it again!)

### Step 3: Verify Sender Email
1. Go to: Settings → Sender Authentication
2. Click "Verify a Single Sender"
3. Fill in your details:
   - From Name: JobHub Alerts
   - From Email: your-email@gmail.com
   - Reply To: your-email@gmail.com
4. Check your email and verify

### Step 4: Configure Application
Edit `.env` file:
```bash
# Add your SendGrid API key
SENDGRID_API_KEY=SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Set sender email (must match verified email)
SMTP_USER=your-verified-email@gmail.com
```

### Step 5: Test
```bash
./start.sh
```

The app will automatically use SendGrid if API key is configured!

## Alternative: Gmail SMTP (Backup)

If you don't want to use SendGrid:

1. Enable 2-Factor Authentication on Gmail
2. Generate App Password:
   - Google Account → Security → 2-Step Verification → App passwords
3. Update `.env`:
```bash
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
```

## Email Service Features

✅ **Automatic Selection**: App uses SendGrid if configured, falls back to SMTP
✅ **Beautiful HTML Emails**: Professional job alert templates
✅ **Plain Text Fallback**: Works in all email clients
✅ **Batch Sending**: Efficient notification delivery
✅ **Error Handling**: Graceful fallback and logging

## SendGrid Free Tier Limits
- 100 emails/day forever free
- No credit card required
- Perfect for testing and small user base

## Monitoring
Check logs for email status:
```bash
tail -f logs/app.log
```

You'll see:
- `Email service: Using SendGrid` - SendGrid active
- `Email service: Using SMTP` - SMTP fallback
- `✅ Email sent via SendGrid to user@example.com` - Success!

# View Registered Users - Quick Guide

## ✅ MySQL Access Fixed
- Username: `root`
- Password: `password`
- Database: `job_portal`
- Host: `localhost:3306`

## 3 Ways to View Users

### 1. Admin Dashboard (Web UI) ⭐ RECOMMENDED
Open in browser:
```
http://localhost:8080/admin
```
Features:
- Real-time user list
- Statistics (Total, Daily, Weekly)
- Auto-refresh every 30 seconds
- Clean table view

### 2. Terminal Script
```bash
bash view_users.sh
```
Shows:
- All registered users
- Statistics summary

### 3. DBeaver
1. Open DBeaver
2. Connect to "job_portal" database
   - Host: localhost
   - Port: 3306
   - Database: job_portal
   - Username: root
   - Password: password
3. Navigate to: job_portal → Tables → users
4. Right-click → View Data

### 4. API Endpoint
```bash
curl http://localhost:8080/api/users
```

## Current Registered User
- Email: prasadkakileti105@gmail.com
- Location: Hyderabad
- Frequency: Daily
- Status: Active
- Registered: 2026-01-29 17:11

## SQL Queries
Run in DBeaver or terminal:

```sql
-- View all users
SELECT * FROM users;

-- Count by location
SELECT location, COUNT(*) FROM users GROUP BY location;

-- Count by frequency
SELECT notification_frequency, COUNT(*) FROM users GROUP BY notification_frequency;
```

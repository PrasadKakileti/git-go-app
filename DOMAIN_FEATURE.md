# Domain-Based Job Filtering Feature

## Overview

The platform implements intelligent job filtering based on professional domains, enabling users to receive highly targeted job recommendations aligned with their career specialization.

## Supported Domains

| Domain | Description | Target Audience |
|--------|-------------|-----------------|
| **Chartered Accountant** | Finance, accounting, and audit roles | CA professionals, finance experts |
| **Java** | Java development and enterprise applications | Java developers, backend engineers |
| **Golang** | Go programming and microservices | Go developers, cloud engineers |

## Technical Implementation

### 1. Data Model

**Users Table Enhancement**
```sql
ALTER TABLE users ADD COLUMN domain VARCHAR(100);
CREATE INDEX idx_user_domain ON users(domain);
```

**Jobs Table Enhancement**
```sql
ALTER TABLE jobs ADD COLUMN domain VARCHAR(100);
CREATE INDEX idx_job_domain ON jobs(domain);
```

### 2. API Integration

**Registration Endpoint**
```http
POST /api/register
Content-Type: application/json

{
  "email": "professional@example.com",
  "location": "Bangalore",
  "domain": "Java",
  "notification_frequency": "daily"
}
```

**Response Schema**
```json
{
  "success": true,
  "message": "Registration successful",
  "user": {
    "id": 1,
    "email": "professional@example.com",
    "location": "Bangalore",
    "domain": "Java",
    "notification_frequency": "daily",
    "is_active": true
  }
}
```

### 3. Job Matching Algorithm

**Filtering Logic**
```go
// Pseudo-code
func GetMatchingJobs(user User) []Job {
    return jobs.Where(
        location LIKE user.Location AND
        domain = user.Domain AND
        posted_at >= user.LastNotification AND
        NOT IN user_job_sent
    )
}
```

**Matching Criteria**
1. **Location Match**: Jobs in user's preferred location
2. **Domain Match**: Jobs matching user's professional domain
3. **Temporal Filter**: Only new jobs since last notification
4. **Deduplication**: Excludes previously sent jobs

### 4. Email Notification Enhancement

**Email Template Structure**
- Job Title
- Company Name
- Location
- **Domain** (highlighted)
- Contact Information
- Posted Date
- Job Description
- Application Link

## User Experience Flow

```
1. User Registration
   ├─ Select Domain from dropdown
   ├─ Choose Location
   └─ Set Notification Frequency

2. Job Scraping (Every 6 hours)
   ├─ Fetch jobs from Naukri.com
   ├─ Classify by domain
   └─ Store with domain tag

3. Notification Dispatch (9 AM daily/weekly)
   ├─ Query users by frequency
   ├─ Filter jobs by domain + location
   ├─ Send personalized email
   └─ Mark jobs as sent
```

## Frontend Implementation

**Domain Selector Component**
```html
<select id="domain" name="domain" required>
    <option value="">Select your domain</option>
    <option value="Chartered Accountant">Chartered Accountant</option>
    <option value="Java">Java</option>
    <option value="Golang">Golang</option>
</select>
```

**Validation**
- Required field
- Client-side validation
- Server-side validation
- Error handling

## Admin Dashboard

**Enhanced User View**
- Email
- Location
- **Domain** (new column)
- Notification Frequency
- Status
- Registration Date

## Testing

### Unit Tests
```bash
# Test domain filtering
go test ./repository -run TestGetJobsByDomain

# Test user registration with domain
go test ./handlers -run TestRegisterWithDomain
```

### Integration Tests
```bash
# End-to-end registration flow
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "location": "Bangalore",
    "domain": "Java",
    "notification_frequency": "daily"
  }'
```

### Database Verification
```sql
-- Verify domain data
SELECT email, location, domain, notification_frequency 
FROM users 
WHERE domain IS NOT NULL;

-- Check job distribution by domain
SELECT domain, COUNT(*) as job_count 
FROM jobs 
GROUP BY domain;
```

## Production Considerations

### Scalability
- **Domain Expansion**: Easy to add new domains via configuration
- **Index Optimization**: Composite indexes on (location, domain)
- **Query Performance**: Optimized filtering queries

### Data Quality
- **Domain Validation**: Enum-based validation
- **Data Consistency**: Foreign key constraints
- **Audit Trail**: Track domain changes

### Future Enhancements
- Multi-domain selection per user
- Domain-specific job sources
- Machine learning for domain classification
- Domain popularity analytics

## Migration Guide

### Existing Users
```sql
-- Set default domain for existing users
UPDATE users 
SET domain = 'Java' 
WHERE domain IS NULL AND location LIKE '%Bangalore%';
```

### Rollback Plan
```sql
-- Remove domain feature
ALTER TABLE users DROP COLUMN domain;
ALTER TABLE jobs DROP COLUMN domain;
```

## Monitoring

**Key Metrics**
- Domain distribution among users
- Job availability per domain
- Match rate (jobs sent vs available)
- User engagement by domain

**Queries**
```sql
-- User distribution by domain
SELECT domain, COUNT(*) as users 
FROM users 
GROUP BY domain;

-- Job availability by domain
SELECT domain, COUNT(*) as jobs 
FROM jobs 
WHERE posted_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
GROUP BY domain;
```

## Documentation

- [API Documentation](docs/API.md)
- [Database Schema](docs/DATABASE.md)
- [User Guide](docs/USER_GUIDE.md)

---

**Feature Status**: ✅ Production Ready  
**Version**: 1.0.0  
**Last Updated**: 2024-01-30

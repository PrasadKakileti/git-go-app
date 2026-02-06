# Project Structure & Architecture

## Directory Layout

```
Go-app/
├── main.go                          # Application entry point
├── go.mod                           # Go module dependencies
├── go.sum                           # Dependency checksums
├── .env                             # Environment configuration
├── .env.example                     # Environment template
├── schema.sql                       # Database schema
│
├── cmd/                             # Command-line utilities
│   └── test_email/                  # Email testing utility
│       └── test_email.go
│
├── config/                          # Configuration management
│   └── config.go                    # Environment loader
│
├── database/                        # Database connectivity
│   └── db.go                        # MySQL connection handler
│
├── handlers/                        # HTTP request handlers
│   ├── handlers.go                  # Core API handlers
│   └── list_users.go                # User listing handler
│
├── models/                          # Data models
│   └── models.go                    # User, Job, Request DTOs
│
├── repository/                      # Data access layer
│   ├── user_repo.go                 # User CRUD operations
│   └── job_repo.go                  # Job CRUD operations
│
├── services/                        # Business logic layer
│   ├── job_service.go               # Job processing logic
│   ├── email_service.go             # Legacy email service
│   ├── unified_email_service.go     # Multi-provider email
│   └── sendgrid_service.go          # SendGrid integration
│
├── scraper/                         # Web scraping
│   └── naukri_scraper.go            # Naukri.com scraper
│
├── scheduler/                       # Task scheduling
│   └── scheduler.go                 # Cron job manager
│
├── frontend/                        # Web interface
│   ├── index.html                   # Landing page
│   ├── admin.html                   # Admin dashboard
│   ├── css/
│   │   └── style.css                # Stylesheet
│   └── js/
│       └── app.js                   # Client-side logic
│
├── docs/                            # Documentation
│   ├── README.md                    # Main documentation
│   ├── FLOW_DIAGRAM.md              # System flow diagrams
│   ├── QUICK_START.md               # Quick start guide
│   └── API.md                       # API reference
│
└── scripts/                         # Utility scripts
    ├── start.sh                     # Application launcher
    ├── fix.sh                       # Troubleshooting script
    ├── view_users.sh                # User viewer
    ├── send_test_email.sh           # Email tester
    └── network_info.sh              # Network configuration
```

## Layer Architecture

### 1. Presentation Layer
**Location**: `frontend/`, `handlers/`

**Responsibilities**:
- HTTP request/response handling
- Input validation
- Response formatting
- Static file serving

**Components**:
- REST API endpoints
- HTML/CSS/JS frontend
- Admin dashboard

### 2. Business Logic Layer
**Location**: `services/`

**Responsibilities**:
- Job matching algorithm
- Notification scheduling
- Email composition
- Business rules enforcement

**Components**:
- JobService: Job processing
- EmailService: Notification delivery
- Scraper integration

### 3. Data Access Layer
**Location**: `repository/`

**Responsibilities**:
- Database operations
- Query optimization
- Transaction management
- Data mapping

**Components**:
- UserRepository: User CRUD
- JobRepository: Job CRUD
- Audit trail management

### 4. Infrastructure Layer
**Location**: `database/`, `config/`, `scheduler/`

**Responsibilities**:
- Database connectivity
- Configuration management
- Task scheduling
- External service integration

**Components**:
- MySQL connection pool
- Environment configuration
- Cron scheduler
- SMTP/SendGrid clients

## Design Patterns

### 1. Repository Pattern
**Purpose**: Abstract data access logic

```go
type UserRepository interface {
    Create(user *User) error
    GetByEmail(email string) (*User, error)
    GetActiveUsers(frequency string) ([]*User, error)
}
```

### 2. Service Layer Pattern
**Purpose**: Encapsulate business logic

```go
type JobService struct {
    jobRepo  *JobRepository
    userRepo *UserRepository
    scraper  *Scraper
    email    *EmailService
}
```

### 3. Dependency Injection
**Purpose**: Loose coupling, testability

```go
func NewJobService(
    jobRepo *JobRepository,
    userRepo *UserRepository,
    scraper *Scraper,
    email *EmailService
) *JobService {
    return &JobService{...}
}
```

### 4. Strategy Pattern
**Purpose**: Multiple email providers

```go
type EmailService interface {
    SendEmail(to, subject, body string) error
}

// Implementations: SMTPService, SendGridService
```

## Data Flow

### Request Flow
```
Client Request
    ↓
HTTP Handler (handlers/)
    ↓
Service Layer (services/)
    ↓
Repository Layer (repository/)
    ↓
Database (MySQL)
```

### Job Processing Flow
```
Scheduler (cron)
    ↓
JobService.ScrapeAndStoreJobs()
    ↓
Scraper.ScrapeJobs()
    ↓
JobRepository.Create()
    ↓
Database
```

### Notification Flow
```
Scheduler (cron)
    ↓
JobService.SendNotifications()
    ↓
UserRepository.GetActiveUsers()
    ↓
JobRepository.GetUnsentJobs()
    ↓
EmailService.SendEmail()
    ↓
SMTP/SendGrid
```

## Module Responsibilities

| Module | Purpose | Key Functions |
|--------|---------|---------------|
| **main.go** | Bootstrap | Initialize, wire dependencies, start server |
| **config** | Configuration | Load env vars, validate config |
| **database** | Connectivity | Manage DB connections, connection pooling |
| **handlers** | API | Handle HTTP requests, validate input |
| **models** | Data | Define data structures, DTOs |
| **repository** | Data Access | CRUD operations, queries |
| **services** | Business Logic | Job matching, notifications |
| **scraper** | Data Collection | Web scraping, data extraction |
| **scheduler** | Automation | Cron jobs, task scheduling |
| **frontend** | UI | User interface, admin dashboard |

## Configuration Management

### Environment-based Configuration
```go
type Config struct {
    DBUser         string
    DBPassword     string
    DBHost         string
    DBPort         string
    DBName         string
    SMTPHost       string
    SMTPPort       string
    SMTPUser       string
    SMTPPass       string
    SendGridAPIKey string
    ServerPort     string
}
```

### Configuration Loading
1. Load from `.env` file
2. Override with environment variables
3. Validate required fields
4. Provide sensible defaults

## Error Handling Strategy

### Levels
1. **Handler Level**: Return HTTP error codes
2. **Service Level**: Log and return errors
3. **Repository Level**: Wrap database errors
4. **Infrastructure Level**: Retry with backoff

### Error Types
- **Validation Errors**: 400 Bad Request
- **Not Found**: 404 Not Found
- **Server Errors**: 500 Internal Server Error
- **Database Errors**: Logged and wrapped

## Testing Strategy

### Unit Tests
```bash
# Test individual components
go test ./handlers -v
go test ./services -v
go test ./repository -v
```

### Integration Tests
```bash
# Test component interaction
go test ./... -tags=integration
```

### End-to-End Tests
```bash
# Test complete flows
./scripts/test_registration.sh
./scripts/test_email.sh
```

## Deployment Architecture

### Single Server Deployment
```
┌─────────────────────────────────┐
│      Application Server         │
│  ┌──────────┐  ┌──────────┐    │
│  │  JobHub  │  │  MySQL   │    │
│  │  (Go)    │──│  (8.0+)  │    │
│  └──────────┘  └──────────┘    │
│       │                         │
│       ▼                         │
│  ┌──────────┐                  │
│  │   SMTP   │                  │
│  │ SendGrid │                  │
│  └──────────┘                  │
└─────────────────────────────────┘
```

### Scalable Deployment
```
┌──────────────┐     ┌──────────────┐
│ Load Balancer│────▶│  App Server  │
└──────────────┘     │   (Primary)  │
                     └──────────────┘
                            │
                     ┌──────────────┐
                     │  App Server  │
                     │  (Secondary) │
                     └──────────────┘
                            │
                     ┌──────────────┐
                     │    MySQL     │
                     │  (Master)    │
                     └──────────────┘
                            │
                     ┌──────────────┐
                     │    MySQL     │
                     │   (Replica)  │
                     └──────────────┘
```

## Security Architecture

### Layers
1. **Network**: Firewall, HTTPS
2. **Application**: Input validation, SQL injection prevention
3. **Data**: Encrypted connections, secure storage
4. **Access**: Rate limiting, authentication

### Best Practices
- Parameterized queries (SQL injection prevention)
- Input sanitization
- Environment-based secrets
- HTTPS in production
- Regular security audits

---

**Document Version**: 1.0.0  
**Last Updated**: 2024-01-30  
**Maintained By**: Development Team

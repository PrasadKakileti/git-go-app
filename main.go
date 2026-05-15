package main

import (
	"job-portal/config"
	"job-portal/database"
	"job-portal/handlers"
	"job-portal/middleware"
	"job-portal/providers"
	"job-portal/repository"
	"job-portal/scheduler"
	"job-portal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)

	// Job providers — JSearch aggregates LinkedIn, Indeed, Glassdoor, ZipRecruiter.
	// Add more providers here as needed (e.g. a Naukri-specific one).
	jobProviders := []providers.JobProvider{
		providers.NewJSearchProvider(cfg.JSearchAPIKey),
	}

	// Services
	emailService := services.NewUnifiedEmailService(cfg)
	jobService := services.NewJobService(jobRepo, userRepo, jobProviders, emailService)

	// Handlers
	handler := handlers.NewHandler(userRepo)
	handler.SetEmailService(emailService)
	handler.SetConfig(cfg)
	handler.SetJobFetcher(jobService)
	jobHandler := handlers.NewJobHandler(jobRepo, userRepo)

	// Scheduler
	sched := scheduler.NewScheduler(jobService)
	sched.Start()
	defer sched.Stop()

	// Router
	r := mux.NewRouter()
	r.Use(middleware.CORS) // applies to every route

	// API routes
	r.HandleFunc("/api/signup", handler.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/login", handler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/register", handler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/users", handler.ListUsers).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user", handler.GetUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/jobs", jobHandler.GetJobsForUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/profile", handler.UpdateProfile).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/refresh-jobs", handler.RefreshJobs).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/health", handler.Health).Methods("GET", "OPTIONS")

	// Static frontend
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./frontend/css"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./frontend/js"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./frontend/images"))))

	r.HandleFunc("/signup.html", serveFile("./frontend/signup.html"))
	r.HandleFunc("/login.html", serveFile("./frontend/login.html"))
	r.HandleFunc("/dashboard.html", serveFile("./frontend/dashboard.html"))
	r.HandleFunc("/profile.html", serveFile("./frontend/profile.html"))
	r.HandleFunc("/admin", serveFile("./frontend/admin.html"))
	r.HandleFunc("/", serveFile("./frontend/login.html"))

	log.Printf("Server listening on :%s", cfg.ServerPort)
	log.Printf("  Local:   http://localhost:%s", cfg.ServerPort)
	log.Printf("  Network: http://<your-ip>:%s (find your IP below)", cfg.ServerPort)

	if cfg.JSearchAPIKey == "" {
		log.Println("WARNING: JSEARCH_API_KEY not set. Job fetching is disabled until you add it to .env")
		log.Println("         Get a free key at: https://rapidapi.com/letscrape-6bRBa3QguO5/api/jsearch")
	}

	log.Fatal(http.ListenAndServe("0.0.0.0:"+cfg.ServerPort, r))
}

func serveFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

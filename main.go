package main

import (
	"job-portal/config"
	"job-portal/database"
	"job-portal/handlers"
	"job-portal/repository"
	"job-portal/scheduler"
	"job-portal/scraper"
	"job-portal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	naukriScraper := scraper.NewNaukriScraper()
	emailService := services.NewUnifiedEmailService(cfg)
	jobService := services.NewJobService(jobRepo, userRepo, naukriScraper, emailService)

	handler := handlers.NewHandler(userRepo)
	handler.SetEmailService(emailService)
	jobHandler := handlers.NewJobHandler(jobRepo, userRepo)

	scheduler := scheduler.NewScheduler(jobService)
	scheduler.Start()
	defer scheduler.Stop()

	r := mux.NewRouter()
	r.HandleFunc("/api/signup", handler.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/login", handler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/register", handler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/users", handler.ListUsers).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user", handler.GetUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/jobs", jobHandler.GetJobsForUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/health", handler.Health).Methods("GET", "OPTIONS")
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./frontend/css"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./frontend/js"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./frontend/images"))))
	r.HandleFunc("/signup.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/signup.html")
	})
	r.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/login.html")
	})
	r.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/admin.html")
	})
	r.HandleFunc("/dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/dashboard.html")
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/login.html")
	})

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("Local access: http://localhost:%s", cfg.ServerPort)
	log.Printf("Network access: http://10.21.12.100:%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+cfg.ServerPort, r))
}

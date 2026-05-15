package repository

import (
	"database/sql"
	"job-portal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// UserPreference is used by the scheduler to know what to fetch.
type UserPreference struct {
	Location   string
	Domain     string
	Experience string
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, location, domain, notification_frequency) VALUES (?, ?, ?, ?)
			  ON DUPLICATE KEY UPDATE location=VALUES(location), domain=VALUES(domain), notification_frequency=VALUES(notification_frequency)`
	result, err := r.db.Exec(query, user.Email, user.Location, user.Domain, user.NotificationFrequency)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	if id > 0 {
		user.ID = int(id)
	}
	return nil
}

func (r *UserRepository) CreateWithAuth(user *models.User) error {
	query := `INSERT INTO users (email, mobile, password, location, domain, experience, notification_frequency, is_verified, verification_code)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query,
		user.Email, user.Mobile, user.Password,
		user.Location, user.Domain, user.Experience,
		user.NotificationFrequency, user.IsVerified, user.VerificationCode,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, COALESCE(mobile,''), location, domain, COALESCE(experience,''), notification_frequency, is_active, created_at
			  FROM users WHERE email = ?`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Mobile,
		&user.Location, &user.Domain, &user.Experience,
		&user.NotificationFrequency, &user.IsActive, &user.CreatedAt,
	)
	return user, err
}

func (r *UserRepository) GetByEmailOrMobile(emailOrMobile string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, COALESCE(mobile,''), password, location, domain, COALESCE(experience,''),
			  notification_frequency, is_active, is_verified, created_at
			  FROM users WHERE email = ? OR mobile = ?`
	err := r.db.QueryRow(query, emailOrMobile, emailOrMobile).Scan(
		&user.ID, &user.Email, &user.Mobile, &user.Password,
		&user.Location, &user.Domain, &user.Experience,
		&user.NotificationFrequency, &user.IsActive, &user.IsVerified, &user.CreatedAt,
	)
	return user, err
}

func (r *UserRepository) GetActiveUsers(frequency string) ([]*models.User, error) {
	query := `SELECT id, email, COALESCE(mobile,''), location, domain, COALESCE(experience,''), notification_frequency
			  FROM users WHERE is_active = TRUE AND notification_frequency = ?`
	rows, err := r.db.Query(query, frequency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.Mobile, &u.Location, &u.Domain, &u.Experience, &u.NotificationFrequency); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, email, COALESCE(mobile,''), location, domain, COALESCE(experience,''), notification_frequency, is_active, created_at
			  FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.Mobile, &u.Location, &u.Domain, &u.Experience,
			&u.NotificationFrequency, &u.IsActive, &u.CreatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) UpdateProfile(userID int, location, domain, experience, frequency string) error {
	query := `UPDATE users SET location=?, domain=?, experience=?, notification_frequency=? WHERE id=?`
	_, err := r.db.Exec(query, location, domain, experience, frequency, userID)
	return err
}

// GetDistinctUserPreferences returns unique (location, domain, experience) combos from active users.
// The scheduler uses this to know exactly what to fetch — no hardcoded city lists.
func (r *UserRepository) GetDistinctUserPreferences() ([]UserPreference, error) {
	query := `SELECT DISTINCT location, domain, COALESCE(experience,'') FROM users WHERE is_active = TRUE`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prefs []UserPreference
	for rows.Next() {
		var p UserPreference
		if err := rows.Scan(&p.Location, &p.Domain, &p.Experience); err != nil {
			continue
		}
		prefs = append(prefs, p)
	}
	return prefs, nil
}

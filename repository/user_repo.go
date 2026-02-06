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

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, location, domain, notification_frequency) VALUES (?, ?, ?, ?)
			  ON DUPLICATE KEY UPDATE location=?, domain=?, notification_frequency=?`
	result, err := r.db.Exec(query, user.Email, user.Location, user.Domain, user.NotificationFrequency,
		user.Location, user.Domain, user.NotificationFrequency)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	if id > 0 {
		user.ID = int(id)
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, location, domain, notification_frequency, is_active, created_at FROM users WHERE email = ?`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Location, 
		&user.Domain, &user.NotificationFrequency, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetActiveUsers(frequency string) ([]*models.User, error) {
	query := `SELECT id, email, location, domain, notification_frequency FROM users 
			  WHERE is_active = TRUE AND notification_frequency = ?`
	rows, err := r.db.Query(query, frequency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Location, &user.Domain, &user.NotificationFrequency); err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, email, location, domain, notification_frequency, is_active, created_at 
			  FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Location, &user.Domain,
			&user.NotificationFrequency, &user.IsActive, &user.CreatedAt); err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}


func (r *UserRepository) CreateWithAuth(user *models.User) error {
	query := `INSERT INTO users (email, mobile, password, location, domain, experience, notification_frequency, is_verified, verification_code) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, user.Email, user.Mobile, user.Password, user.Location, user.Domain, 
		user.Experience, user.NotificationFrequency, user.IsVerified, user.VerificationCode)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetByEmailOrMobile(emailOrMobile string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, mobile, password, location, domain, experience, notification_frequency, is_active, is_verified, created_at 
			  FROM users WHERE email = ? OR mobile = ?`
	err := r.db.QueryRow(query, emailOrMobile, emailOrMobile).Scan(
		&user.ID, &user.Email, &user.Mobile, &user.Password, &user.Location, &user.Domain,
		&user.Experience, &user.NotificationFrequency, &user.IsActive, &user.IsVerified, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

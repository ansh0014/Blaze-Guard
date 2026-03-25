package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	FirebaseUID string    `json:"firebase_uid" db:"auth0_id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Picture     string    `json:"picture"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Session struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) FindByFirebaseUID(firebaseUID string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, auth0_id, email, name, picture, role, created_at, updated_at
		FROM users
		WHERE auth0_id = $1
	`
	err := r.db.QueryRow(query, firebaseUID).Scan(
		&user.ID,
		&user.FirebaseUID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, auth0_id, email, name, picture, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirebaseUID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (auth0_id, email, name, picture, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		user.FirebaseUID,
		user.Email,
		user.Name,
		user.Picture,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users
		SET email = $1, name = $2, picture = $3, role = $4
		WHERE auth0_id = $5
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		user.Email,
		user.Name,
		user.Picture,
		user.Role,
		user.FirebaseUID,
	).Scan(&user.UpdatedAt)
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

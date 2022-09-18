package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDuplicateAdminUsername = errors.New("duplicate username")
)

type Admin struct {
	Id       int64    `json:"id"`
	Username string   `json:"username"`
	Password password `json:"-"`
}

type AdminModel struct {
	DB *sql.DB
}

func (a AdminModel) Insert(admin *Admin) error {
	query := `
        INSERT INTO admin (username, password) 
        VALUES ($1, $2) RETURNING id`

	args := []interface{}{admin.Username, admin.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(&admin.Id)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (a AdminModel) GetByUsername(username string) (*Admin, error) {
	query := `
        SELECT * FROM admin
        WHERE username = $1`

	var admin Admin

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, username).Scan(
		&admin.Id,
		&admin.Username,
		&admin.Password.hash,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &admin, nil
}

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousUser = &User{}

type User struct {
	ID          int64    `json:"id"`
	Username    string   `json:"username"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	Password    password `json:"-"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Insert(user *User) error {
	query := `
        INSERT INTO users (username, phone_number,email, password) 
        VALUES ($1, $2, $3, $4) RETURNING id`

	args := []interface{}{user.Username, user.PhoneNumber, user.Email, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		fmt.Println(err.Error(), query)
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (u UserModel) GetUsers() ([]User, error) {
	query := `
        SELECT id,  username, email, phone_number
        FROM users`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PhoneNumber,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	fmt.Println(users)
	return users, nil
}

//func (u *UserModel) GetUserById(id int) (*User,error)

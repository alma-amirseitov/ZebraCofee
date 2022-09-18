package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
)

type Cashier struct {
	ID       int64    `json:"id"`
	BranchID int64    `json:"branch_id"`
	Username string   `json:"username"`
	Password password `json:"-"`
}

type CashierModel struct {
	DB *sql.DB
}

func (c CashierModel) Insert(cashier *Cashier) error {
	query := `
        INSERT INTO cashier (branch_id,username, password) 
        VALUES ($1, $2,$3) RETURNING id`

	args := []interface{}{cashier.BranchID, cashier.Username, cashier.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&cashier.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (c CashierModel) GetByUsername(username string) (*Cashier, error) {
	query := `
        SELECT * FROM cashier
        WHERE username = $1`

	var cashier Cashier

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, username).Scan(
		&cashier.ID,
		&cashier.BranchID,
		&cashier.Username,
		&cashier.Password.hash,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &cashier, nil
}

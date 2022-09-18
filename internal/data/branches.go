package data

import (
	"context"
	"database/sql"
	"time"
)

type Branches struct {
	ID          int64   `json:"ID"`
	Address     string  `json:"address"`
	XCoordinate float64 `json:"x_coordinate"`
	YCoordinate float64 `json:"y_coordinate"`
}

type BranchesModel struct {
	DB *sql.DB
}

func (branch BranchesModel) AddBranches(branches Branches) error {
	query := `Insert Into branches (address,x_coordinate,y_coordinate) values ($1,$2,$3) returning id`

	args := []interface{}{branches.Address, branches.XCoordinate, branches.YCoordinate}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := branch.DB.QueryRowContext(ctx, query, args...).Scan(&branches.ID)
	if err != nil {
		return err
	}
	return nil
}

func (branch BranchesModel) GetBranches() ([]Branches, error) {
	query := `SELECT * FROM branches`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := branch.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var branches []Branches

	for rows.Next() {
		var b Branches

		err := rows.Scan(
			&b.ID,
			&b.Address,
			&b.XCoordinate,
			&b.YCoordinate,
		)

		if err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}
	return branches, nil
}

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Products struct {
	ID          int64   `json:"id"`
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
	NetCost     float64 `json:"net_cost"`
	Price       float64 `json:"price"`
	PictureUrl  string  `json:"picture_url"`
	Description string  `json:"description"`
}

type ProductModel struct {
	DB *sql.DB
}

func (p ProductModel) Insert(product *Products) error {
	query := `
        INSERT INTO products (product_name, category, net_cost, price,picture_url,description) 
        VALUES ($1, $2, $3, $4,$5,$6)
        RETURNING id`

	args := []interface{}{product.ProductName, product.Category, product.NetCost, product.Price, product.PictureUrl, product.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID)
}

func (p ProductModel) Get(id int64) (*Products, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT * FROM products
        WHERE id = $1`

	var product Products

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.ProductName,
		&product.Category,
		&product.NetCost,
		&product.Price,
		&product.PictureUrl,
		&product.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (p ProductModel) Update(product *Products) error {
	query := `UPDATE Products 
			  SET product_name = $1, category = $2, net_cost = $3, price = $4, picture_url = $5, description = $6 
			  WHERE id = $7`

	args := []interface{}{
		product.ProductName,
		product.Category,
		product.NetCost,
		product.Price,
		product.PictureUrl,
		product.Description,
		product.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID)
	if err != nil {
		fmt.Println(err.Error(), query, product.ID)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (p ProductModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM products
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (p ProductModel) GetAll() ([]*Products, error) {
	query := fmt.Sprintf(`SELECT * FROM products`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var products []*Products

	for rows.Next() {
		var product Products

		err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Category,
			&product.NetCost,
			&product.Price,
			&product.PictureUrl,
			&product.Description,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

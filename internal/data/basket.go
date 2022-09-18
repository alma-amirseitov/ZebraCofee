package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type Basket struct {
	ID           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	BasketStatus string `json:"status"`
}

type BasketProducts struct {
	BaskedId  int64    `json:"basked_id"`
	ProductId int64    `json:"product_id"`
	Amount    int64    `json:"amount"`
	Addition  []string `json:"addition"`
}

type BasketModel struct {
	DB *sql.DB
}

func (b BasketModel) CreateBasket(userId int64) (int64, error) {
	queryBasket := `INSERT INTO basket (user_id,status) values ($1,'active') returning id,user_id,status`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var basket Basket

	err := b.DB.QueryRowContext(ctx, queryBasket, userId).Scan(&basket.ID, &basket.UserId, &basket.BasketStatus)
	if err != nil {
		return -1, err
	}
	return basket.ID, nil
}

func (b BasketModel) GetActiveBasket(userId int64) (int64, error) {
	query := `SELECT id FROM basket where user_id=$1 and status = 'active'`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int64

	err := b.DB.QueryRowContext(ctx, query, userId).Scan(&id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return -1, ErrRecordNotFound
		default:
			return -1, err
		}
	}

	return id, nil
}

func (b BasketModel) AddBasketProduct(basketProducts *BasketProducts) error {
	queryBasketOrder := `INSERT INTO basket_products (basket_id, product_id, amount, addition) VALUES ($1,$2,$3,$4)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{basketProducts.BaskedId, basketProducts.ProductId, basketProducts.Amount, pq.Array(basketProducts.Addition)}
	_, err := b.DB.ExecContext(ctx, queryBasketOrder, args...)

	if err != nil {
		return err
	}
	return nil
}

func (b BasketModel) ChangeBasketProduct(basketId int64, productId int64, i int64) error {
	query := `UPDATE basket_products SET amount = amount + $1 where basket_id = $2 and product_id = $3`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := b.DB.ExecContext(ctx, query, i, basketId, productId)
	if err != nil {
		return err
	}
	return nil
}

func (b BasketModel) CheckBasketProduct(basketId int64, productId int64, addition []string) (int64, error) {
	query := `SELECT amount FROM basket_products WHERE product_id=$1 and basket_id = $2 and addition = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var amount int64

	err := b.DB.QueryRowContext(ctx, query, productId, basketId, pq.Array(addition)).Scan(&amount)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, nil
		default:
			return -1, err
		}
	}
	return amount, err
}

func (b BasketModel) GetProductsByBasket(id int64) ([]BasketProducts, error) {
	query := `SELECT basket_id,product_id,amount,addition FROM basket_products where basket_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	var basketProducts []BasketProducts

	for rows.Next() {
		var basketProduct BasketProducts

		err := rows.Scan(
			&basketProduct.BaskedId,
			&basketProduct.ProductId,
			&basketProduct.Amount,
			pq.Array(&basketProduct.Addition),
		)
		if err != nil {
			return nil, err
		}
		basketProducts = append(basketProducts, basketProduct)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return basketProducts, nil
}

func (b BasketModel) CloseBasketStatus(basketId int64) error {
	query := `UPDATE basket SET status = 'inactive' where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := b.DB.ExecContext(ctx, query, basketId)
	if err != nil {
		return err
	}
	return nil
}

package data

import (
	"context"
	"database/sql"
	"time"
)

type Order struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CashierID   int64     `json:"cashier_id"`
	Time        time.Time `json:"time"`
	Status      string    `json:"status"`
	Preferences string    `json:"preferences"`
	OrderType   string    `json:"order_type"`
	BasketId    int64     `json:"basket_id"`
	Discount    int64     `json:"discount"`
	TotalPrice  float64   `json:"total_price"`
	OrderPrice  float64   `json:"order_price"`
	BranchesId  int64     `json:"branches_id"`
}

type OrderModel struct {
	DB *sql.DB
}

func (o OrderModel) MakeOrder(order *Order) error {
	queryOrder := `INSERT INTO orders (user_id, cashier_id, time, status, preferences, order_type, basket_id, total_price,order_price,discount) 
			  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		order.UserID,
		order.CashierID,
		order.Time,
		order.Status,
		order.Preferences,
		order.OrderType,
		order.BasketId,
		order.TotalPrice,
		order.OrderPrice,
		order.Discount,
	}

	return o.DB.QueryRowContext(ctx, queryOrder, args...).Scan(&order.ID)
}

func (o OrderModel) CloseOrderStatus(id int64) error {
	query := `UPDATE orders SET status = 'done' where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := o.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users   UserModel
	Admin   AdminModel
	Cashier CashierModel
	Tokens  TokenModel

	Products ProductModel
	Order    OrderModel
	Basket   BasketModel

	Branches BranchesModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:   UserModel{DB: db},
		Admin:   AdminModel{DB: db},
		Cashier: CashierModel{DB: db},
		Tokens:  TokenModel{DB: db},

		Products: ProductModel{DB: db},
		Order:    OrderModel{DB: db},
		Basket:   BasketModel{DB: db},

		Branches: BranchesModel{DB: db},
	}
}

package main

import (
	"ZebraCofee/internal/data"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) addBasketProductHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserId    int64    `json:"user_id"`
		ProductId int64    `json:"product_id"`
		Amount    int64    `json:"amount"`
		Addition  []string `json:"addition"`
	}

	basketId, err := app.models.Basket.GetActiveBasket(input.UserId)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			basketId, err = app.models.Basket.CreateBasket(input.UserId)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	basketProducts := &data.BasketProducts{
		BaskedId:  basketId,
		ProductId: input.ProductId,
		Amount:    input.Amount,
		Addition:  input.Addition,
	}
	amount, err := app.models.Basket.CheckBasketProduct(basketId, input.ProductId, input.Addition)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if amount == 0 {
		err = app.models.Basket.AddBasketProduct(basketProducts)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		err = app.models.Basket.ChangeBasketProduct(basketId, input.ProductId, input.Amount)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		basketProducts.Amount += amount
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"Added Product": basketProducts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) MakeOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId      int64  `json:"user_id"`
		Preferences string `json:"preferences"`
		OrderType   string `json:"order_type"`
		BranchesId  int64  `json:"branches_id"`
		CashierId   int64  `json:"cashier_id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		fmt.Println("here1")
		app.badRequestResponse(w, r, err)
		return
	}
	basketId, err := app.models.Basket.GetActiveBasket(input.UserId)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.badRequestResponse(w, r, errors.New("basket is empty"))
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	busketProducts, err := app.models.Basket.GetProductsByBasket(basketId)
	if err != nil {
		fmt.Println("here3")
		app.badRequestResponse(w, r, err)
		return
	}
	orderPrice := 0.0

	for _, v := range busketProducts {
		product, err := app.models.Products.Get(v.ProductId)
		if err != nil {
			fmt.Println("here4")
			app.badRequestResponse(w, r, err)
			return
		}
		orderPrice += product.Price*float64(v.Amount) + 50.0*float64(len(product.Category))
	}

	err = app.models.Basket.CloseBasketStatus(basketId)
	if err != nil {
		fmt.Println("here5")
		app.badRequestResponse(w, r, err)
		return
	}

	order := &data.Order{
		UserID:      input.UserId,
		Status:      "in progress",
		Preferences: input.Preferences,
		OrderType:   input.OrderType,
		BasketId:    basketId,
		OrderPrice:  orderPrice,
		TotalPrice:  orderPrice * 0.85,
		Discount:    15,
		CashierID:   input.CashierId,
		BranchesId:  input.BranchesId,
		Time:        time.Now(),
	}

	err = app.models.Order.MakeOrder(order)
	if err != nil {
		fmt.Println("here6")
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"order": order}, nil)
	if err != nil {
		fmt.Println("here7")
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) CloseOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id int64 `json:"id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = app.models.Order.CloseOrderStatus(input.Id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"order_id": input.Id}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

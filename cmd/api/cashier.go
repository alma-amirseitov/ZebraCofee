package main

import (
	"ZebraCofee/internal/data"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) registerCashierHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BranchesId int64  `json:"branches_id"`
		Username   string `json:"username"`
		Password   string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cashier := &data.Cashier{
		BranchID: input.BranchesId,
		Username: input.Username,
	}

	err = cashier.Password.Set(input.Password)
	if err != nil {
		fmt.Println(err.Error())
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Cashier.Insert(cashier)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(cashier.ID, 3*24*time.Hour, data.RoleCashier)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"cashier": cashier, "authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//loginAdmin
func (app *application) loginCashierHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cashier, err := app.models.Cashier.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := cashier.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(cashier.ID, 24*time.Hour, data.RoleCashier)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

package main

import (
	"ZebraCofee/internal/data"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) registerAdminHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	admin := &data.Admin{
		Username: input.Username,
	}

	err = admin.Password.Set(input.Password)
	if err != nil {
		fmt.Println(err.Error())
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Admin.Insert(admin)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(admin.Id, 3*24*time.Hour, data.RoleAdmin)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"admin": admin, "authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//loginAdmin
func (app *application) loginAdminHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	admin, err := app.models.Admin.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := admin.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(admin.Id, 24*time.Hour, data.RoleAdmin)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//Add New Branch
func (app *application) AddBranch(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Address     string  `json:"address"`
		XCoordinate float64 `json:"x_coordinate"`
		YCoordinate float64 `json:"x_coordinate"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	branch := &data.Branches{
		Address:     input.Address,
		XCoordinate: input.XCoordinate,
		YCoordinate: input.YCoordinate,
	}

	err = app.models.Branches.AddBranches(*branch)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"branch": branch}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

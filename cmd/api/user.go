package main

import (
	"ZebraCofee/internal/data"
	"ZebraCofee/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) getUsers(w http.ResponseWriter, r *http.Request) {
	u := struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
	}{}
	user, err := app.models.Users.GetUsers()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	var users []struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
	}

	for _, v := range user {
		u.Username = v.Username
		u.Email = v.Email
		u.PhoneNumber = v.PhoneNumber
		users = append(users, u)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserByEmail(w http.ResponseWriter, r *http.Request) {

}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		fmt.Println("here3")
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Username:    input.Username,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		fmt.Println("here2")
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		fmt.Println("here1")
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {

}

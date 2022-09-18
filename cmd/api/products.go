package main

import (
	"ZebraCofee/internal/data"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductName string  `json:"product_name"`
		Category    string  `json:"category"`
		NetCost     float64 `json:"net_cost"`
		Price       float64 `json:"price"`
		PictureUrl  string  `json:"picture_url"`
		Description string  `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &data.Products{
		ProductName: input.ProductName,
		Category:    input.Category,
		NetCost:     input.NetCost,
		Price:       input.Price,
		PictureUrl:  input.PictureUrl,
		Description: input.Description,
	}

	err = app.models.Products.Insert(product)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/products/%d", product.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"product": product}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := app.models.Products.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"products": products}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Products.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "product successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	product, err := app.models.Products.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	product, err := app.models.Products.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		ProductName *string  `json:"product_name"`
		Category    *string  `json:"category"`
		NetCost     *float64 `json:"net_cost"`
		Price       *float64 `json:"price"`
		PictureUrl  *string  `json:"picture_url"`
		Description *string  `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.ProductName != nil {
		product.ProductName = *input.ProductName
	}

	if input.Category != nil {
		product.Category = *input.Category
	}
	if input.NetCost != nil {
		product.NetCost = *input.NetCost
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.PictureUrl != nil {
		product.PictureUrl = *input.PictureUrl
	}
	if input.Description != nil {
		product.Description = *input.Description
	}

	err = app.models.Products.Update(product)
	if err != nil {
		fmt.Println(err.Error())
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

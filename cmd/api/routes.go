package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//User
	router.HandlerFunc(http.MethodPost, "/api/users/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/api/users/login", app.loginUserHandler)

	//Admin
	router.HandlerFunc(http.MethodPost, "/api/admin/register", app.registerAdminHandler)
	router.HandlerFunc(http.MethodPost, "/api/admin/login", app.loginAdminHandler)

	//Cashier
	router.HandlerFunc(http.MethodPost, "/api/cashier/register", app.registerCashierHandler)
	router.HandlerFunc(http.MethodPost, "/api/cashier/login", app.loginCashierHandler)

	//Products
	router.HandlerFunc(http.MethodGet, "/api/products", app.ListProductsHandler)
	router.HandlerFunc(http.MethodPost, "/api/products", app.createProductHandler)
	router.HandlerFunc(http.MethodGet, "/api/products/:id", app.GetProductHandler)
	router.HandlerFunc(http.MethodPatch, "/api/products/:id", app.updateProductHandler)
	router.HandlerFunc(http.MethodDelete, "/api/products/:id", app.deleteProductHandler)

	//BasketProducts
	router.HandlerFunc(http.MethodPost, "/api/basket/add", app.addBasketProductHandler)

	//Order
	router.HandlerFunc(http.MethodPost, "/api/order/make", app.MakeOrderHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}

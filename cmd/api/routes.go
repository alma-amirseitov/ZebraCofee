package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/api/createUser", app.createUser)
	router.HandlerFunc(http.MethodGet, "/api/users", app.getUsers)

	router.HandlerFunc(http.MethodGet, "/api/products", app.ListProductsHandler)
	router.HandlerFunc(http.MethodPost, "/api/products", app.createProductHandler)
	router.HandlerFunc(http.MethodGet, "/api/products/:id", app.GetProductHandler)
	router.HandlerFunc(http.MethodPatch, "/api/products/:id", app.updateProductHandler)
	router.HandlerFunc(http.MethodDelete, "/api/products/:id", app.deleteProductHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}

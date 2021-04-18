// Package Classification Of Product API
//
// Documentation for Product API demo
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta


package Handlers

import (
	"Microservices/NotMod/data"
	"context"
	"fmt"

	//"regexp"
	"strconv"

	//"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// A list of products
// swagger:response productResponse
type productResponse struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:response noContent
type productsNoContent struct {

}

// swagger:parameters updateProduct
type productIDParameterWrapper struct {
	// The ID of the product to be updated
	// in: path
	// required: true
	ID int `json:"id"`
}

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// Returns a lft of products
// responses:
//   200: productResponse 

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route PUT /products/{id} products updateProducts
// Returns a lft of products
// responses:
//   201: noContent

// Updates an existing product
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {


	prod := r.Context().Value(KeyProduct{}).(data.Product)
	//err := data.AddProduct(&prod)
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.l.Fatal(err)
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductuctNotFound {
		http.Error(w, "Product Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product Not Found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{
	
}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		err = prod.Validate()
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error validating product %s: ", err),
				http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(w , req)
	})


}
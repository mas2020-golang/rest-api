/*
Package to handle all the request that comes from / path.
*/

package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/mas2020-golang/rest-api/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// test it with:
// curl -s  http://localhost:9090/ | jq
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle GET Products")
	lp := data.Products.GetProducts()
	// return JSON to the caller
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST Product")
	// take the product from the request context. The product has been inserted into the context from the middleware function
	// call
	prod, _ := r.Context().Value("prod").(*data.Product) // cast the interface{} to *data.Product
	p.l.Printf("add product %#v", prod)
	prod.Add()
}

// UpdateProduct is the handler for the update of a single product
func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// retrieve the id from the path
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "{id} not found in the path", http.StatusBadRequest)
		return
	}
	p.l.Printf("Handle PUT Product with id %d", id)
	// take the product from the request context. The product has been inserted into the context from the middleware function
	// call
	prod, _ := r.Context().Value("prod").(*data.Product) // cast the interface{} to *data.Product
	p.l.Printf("update product %#v", prod)
	err = prod.Update()
	// error check
	switch err {
	case data.RecordNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(`"message": "update product done!"`))
}

// MiddlewareProductValidation is a function call before the effective function. Its scope is to unmarshall the json
// object in the body of the request in a valid Product object, save this object into a new context, inject the new
// context in the request and serve the next handler in the chain
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.l.Println("MiddlewareProductValidation func execution")
		prod := &data.Product{}
		if err := prod.FromJSON(r.Body); err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// context creation to store product
		ctx := context.WithValue(r.Context(), "prod", prod)
		// create a new request with the new context
		req := r.WithContext(ctx)
		// server the next handler
		next.ServeHTTP(w, req)
	})
}

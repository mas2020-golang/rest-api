/*
Package to handle all the request that comes from / path.
*/

package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mas2020-golang/rest-api/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l    *log.Logger
	pool *pgxpool.Pool
}

func NewProducts(l *log.Logger, pool *pgxpool.Pool) *Products {
	return &Products{l, pool}
}

// test it with:
// curl -s  http://localhost:9090/ | jq
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle GET Products")
	lp, err := data.Products.GetProducts(p.pool)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// return JSON to the caller
	err = lp.ToJSON(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Products) GetProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle GET Product")
	// retrieve the id from the path
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "{id} not found in the path", http.StatusBadRequest)
		return
	}
	prod := &data.Product{ID: id}
	err = prod.Get(p.pool)
	if err != nil {
		http.Error(w, "product not found: "+err.Error(), http.StatusNotFound)
		return
	}
	json, err := prod.ToJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST Product")
	// take the product from the request context. The product has been inserted into the context from the middleware function
	// call
	prod, _ := r.Context().Value("prod").(*data.Product) // cast the interface{} to *data.Product
	p.l.Printf("product from http body is %#v", prod)
	err := prod.Add(p.pool)
	if err != nil {
		http.Error(w, "error occurred during product creation", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonResp, err := prod.ToJSON()
	//p.l.Printf("json marshal of product is: %s", string(jsonResp))
	if err != nil {
		http.Error(w, "failed to represent json object", http.StatusInternalServerError)
		p.l.Printf("failed to represent json object: %s", err.Error())
		return
	}
	w.Write(jsonResp)
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
	p.l.Printf("product data received are %#v", prod)
	prod.ID = id
	err = prod.Update(p.pool)
	// error check
	switch err {
	case data.RecordNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// MiddlewareProductValidation is a function call before the effective function. Its scope is to unmarshall the json
// object in the body of the request in a valid Product object, save this object into a new context, inject the new
// context in the request and serve the next handler in the chain
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.l.Println("MiddlewareProductValidation func execution")
		prod := &data.Product{}
		if err := prod.FromJSON(r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p.l.Println("product middleware validation")
		if err := prod.Validate(); err != nil {
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

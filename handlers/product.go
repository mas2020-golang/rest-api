/*
Package to handle all the request that comes from / path.
*/

package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mas2020-golang/rest-api/data"
	"github.com/mas2020-golang/rest-api/utils"
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

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle GET Products")
	// get the claims
	//claims, _ := r.Context().Value("claims").(jwt.MapClaims) // cast the interface{} to jwt.MapClaims
	//p.l.Printf("claims data in the context are %#v", claims)

	lp, err := data.Products.GetAll(p.pool)
	if err != nil {
		utils.ReturnError(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return JSON to the caller
	json, err := lp.ToJSON()
	if err != nil {
		utils.ReturnError(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (p *Products) GetProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle GET Product")
	// retrieve the id from the path
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(&w, "{id} not found in the path", http.StatusNotFound)
		return
	}
	prod, err := data.Products.Get(p.pool, id)
	if err != nil {
		utils.ReturnError(&w, fmt.Sprintf("product not found (%s)", err.Error()), http.StatusNotFound)
		return
	}
	json, err := prod.ToJSON()
	if err != nil {
		utils.ReturnError(&w, err.Error(), http.StatusNotFound)
		return
	}
	w.Write(json)
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST /products")
	// take the product from the request context. The product has been inserted into the context from the middleware function
	// call
	prod, _ := r.Context().Value("prod").(*data.Product) // cast the interface{} to *data.Product
	p.l.Printf("product from http body is %#v", prod)
	err := data.Products.Add(p.pool, prod)
	if err != nil {
		utils.ReturnError(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonBody, err := prod.ToJSON()
	if err != nil {
		utils.ReturnError(&w, fmt.Sprintf("error to marshall json response, %s", err.Error()),
			http.StatusInternalServerError)
		return
	}
	w.Write(jsonBody)
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
	err = data.Products.Update(p.pool, prod)
	if err != nil {
		// error check
		switch err {
		case data.RecordNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

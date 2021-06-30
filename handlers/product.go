/*
Package to handle all the request that comes from / path.
*/

package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/rest-api/models"
	"github.com/mas2020-golang/rest-api/utils"
	"net/http"
	"strconv"
)

type Products struct {
	pool *pgxpool.Pool
}

func NewProducts(pool *pgxpool.Pool) *Products {
	return &Products{pool}
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	output.InfoLog("", "GET /products")
	// get the claims
	//claims, _ := r.Context().Value("claims").(jwt.MapClaims) // cast the interface{} to jwt.MapClaims
	//p.l.Printf("claims models in the context are %#v", claims)

	lp, err := models.Products.GetAll(p.pool)
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
	// retrieve the id from the path
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ReturnError(&w, "{id} not found in the path", http.StatusNotFound)
		return
	}
	output.InfoLog("", fmt.Sprintf("GET /products/%d", id))
	prod, err := models.Products.Get(p.pool, id)
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
	output.InfoLog("", "POST /products")
	// take the product from the request context. The product has been inserted into the context from the middleware function
	// call
	prod, _ := r.Context().Value("prod").(*models.Product) // cast the interface{} to *models.Product
	output.DebugLog("", fmt.Sprintf("product content in http body: %#v", prod))
	err := models.Products.Add(p.pool, prod)
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
	output.InfoLog("", fmt.Sprintf("PUT /products/%d", id))
	// take the product from the request context. The product has been inserted into the context from the middleware function
	// call
	prod, _ := r.Context().Value("prod").(*models.Product) // cast the interface{} to *models.Product
	output.DebugLog("", fmt.Sprintf("product content in http body: %#v", prod))
	prod.ID = id
	err = models.Products.Update(p.pool, prod)
	if err != nil {
		// error check
		switch err {
		case models.RecordNotFound:
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
		prod := &models.Product{}
		if err := prod.FromJSON(r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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

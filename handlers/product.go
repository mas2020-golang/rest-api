package handlers

import (
	"encoding/json"
	"github.com/mas2020-golang/rest-api/data"
	"log"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products{
	return &Products{l}
}

func (p *Products) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	// transform the list of products into a slice of byte using json package
	b, err := json.Marshal(lp)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(b)
}



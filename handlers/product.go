/*
Package to handle all the request that comes from / path.
*/

package handlers

import (
	"github.com/mas2020-golang/rest-api/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

/*
Single router that checks the call and decide wich internal methods to
call to serve the request.
This method is what is automatically call from the NewServeMux http multi plexer when
a / path is called. ServeHTTP is called because it implements the Handler interface.
*/
func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// routing the HTTP method
	switch r.Method {
	case http.MethodGet: // handle GET method
		p.getProducts(w, r)
		return
	case http.MethodPost: // handle POST method
		p.addProduct(w, r)
		return
	case http.MethodPut: // handle PUT method
		// regex to catch the id
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 || len(g[0]) != 2 {
			http.Error(w, `{"message": "Invalid URI"}`, http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(g[0][1])
		if err != nil {
			http.Error(w, `{"message": "{id} parameter is not an integer"}`, http.StatusBadRequest)
			return
		}
		p.updateProduct(id, w, r)
		return
	}

	// catch all
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// test it with:
// curl -s  http://localhost:9090/ | jq
func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET ProductsType")
	lp := data.Products.GetProducts()
	// return JSON to the caller
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	prod, err := new(data.ProductsType).FromJSON(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		p.l.Printf("add product %#v", prod)
		data.Products.Add(prod)
	}
}

// updateProduct is the handler for the update of a single product
func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Printf("Handle PUT Product with id %d", id)
	prod, err := new(data.ProductsType).FromJSON(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		p.l.Printf("update product %#v", prod)
		err := data.Products.Update(id, prod)
		// error check
		switch err {
		case data.RecordNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(`"message": "update product done!"`))
	}
}

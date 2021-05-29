/*
Package data the structure (modeling) of the data.
To avoid reading from a database the data are also inserted here.
*/
package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// custom errors
var RecordNotFound = fmt.Errorf("product not found")

// ProductsType type to return directly a slice of Product
type ProductsType []*Product

// ToJSON encode the ProductsType object and write it into the writer
// passed as argument
func (p *ProductsType) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// FromJSON returns a new Product decoding the JSON read from the reader
func (p *ProductsType) FromJSON(r io.Reader) (*Product, error) {
	prod := &Product{}
	e := json.NewDecoder(r)
	err := e.Decode(prod)
	if err != nil {
		return nil, err
	} else {
		return prod, nil
	}
}

// Add a new product to the internal slice of ProductsType
func (p *ProductsType) Add(new *Product) {
	productList = append(productList, new)
}

// Update the product in the collection
func (p *ProductsType) Update(id int, prod *Product) error {
	idx := p.findProduct(id)
	if idx > -1 {
		productList[idx] = prod
		return nil
	} else {
		return RecordNotFound
	}
}

func (p *ProductsType) GetProducts() ProductsType {
	return productList
}

// findProduct returns the position of the product in the list or -1 in case of not found
func (p *ProductsType) findProduct(id int) int {
	idx := -1
	for i, prodItem := range productList {
		if prodItem.ID == id {
			idx = i
			break
		}
	}
	return idx
}

var Products ProductsType = ProductsType{}

var productList = ProductsType{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "dfad",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee",
		Price:       1.99,
		SKU:         "dfadds",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

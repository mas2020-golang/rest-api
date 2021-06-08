/*
Package data the structure (modeling) of the data.
To avoid reading from a database the data are also inserted here.
*/
package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"io"
	"regexp"
	"time"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// custom errors
var RecordNotFound = fmt.Errorf("product not found")

// Validate the structure
func (p *Product) Validate() error {
	validate := validator.New()
	// this methods is used to validate in a custom way the field SKU
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

// validateSKU is a custom validation func for the SKU field
func validateSKU(fl validator.FieldLevel) bool {
	// sku is of sada-ads-adas
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	result := re.FindAllString(fl.Field().String(), -1)
	return len(result) == 1
}

// FromJSON fills Product decoding the JSON read from the reader.
// Returns an error in case of any.
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// ToJSON encode the Product object into a json representation in []byte
func (p *Product) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// Add a new product to the internal slice of ProductsType
func (p *Product) Add(pool *pgxpool.Pool) error {
	err := pool.QueryRow(context.Background(),
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

// Update the product in the collection
// TODO: write to the database and code the corresponding test method
func (p *Product) Update() error {
	// search the product in the internal collection (in a real application it would update the database with the value
	// of p searching the row with the corresponding product.ID)
	idx := findProduct(p.ID)
	if idx > -1 {
		productList[idx] = p
		return nil
	} else {
		return RecordNotFound
	}
}

// ProductsType type to return directly a slice of Product
type ProductsType []*Product

// ToJSON encode the ProductsType object and write it into the writer
// passed as argument
func (p *ProductsType) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *ProductsType) GetProducts(pool *pgxpool.Pool) (ProductsType, error) {
	var productList = ProductsType{}
	rows, err := pool.Query(context.Background(), "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate through the result set
	for rows.Next() {
		p := Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		productList = append(productList, &p)
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	// return the list of products
	return productList, nil
}

// findProduct returns the position of the product in the list or -1 in case of not found
func findProduct(id int) int {
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

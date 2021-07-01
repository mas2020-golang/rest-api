/*
Package models the structure (modeling) of the models.
To avoid reading from a database the models are also inserted here.
*/
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"io"
	"regexp"
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
var (
	RecordNotFound = fmt.Errorf("product not found")
	Products ProductsT = ProductsT{}
)

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

// ProductsT type to return directly a slice of Product
type ProductsT []*Product

// ToJSON encode the Products object into a json representation in []byte
func (p *ProductsT) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// GetAll returns a slice of *Product.
func (p *ProductsT) GetAll(pool *pgxpool.Pool) (ProductsT, error) {
	var productList ProductsT
	rows, err := pool.Query(context.Background(), "SELECT id, name, description, price, sku FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	/*
	    id SERIAL,
	    name TEXT NOT NULL,
	    description TEXT NOT NULL,
	    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
	    sku varchar(100),
	 */
	// iterate through the result set
	for rows.Next() {
		p := Product{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.SKU)
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

// Get the product reading db. All the fields are stored in the *Product object returned by the method.
func (p *ProductsT) Get(pool *pgxpool.Pool, id int) (prod *Product, err error) {
	row := pool.QueryRow(context.Background(), "SELECT id, name, description, price, sku FROM products WHERE id=$1",
		id)
	prod = new(Product)
	err = row.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.SKU)
	return prod, err
}

// Add a new product to the products table
func (p *ProductsT) Add(pool *pgxpool.Pool, new *Product) error {
	err := pool.QueryRow(context.Background(),
		"INSERT INTO products(name, price, description, sku) VALUES($1, $2, $3, $4) RETURNING id",
		(*new).Name, (*new).Price, (*new).Description, (*new).SKU).Scan(&new.ID)

	if err != nil {
		return err
	}
	return nil
}

// Update the product in the collection
func (p *ProductsT) Update(pool *pgxpool.Pool, prod *Product) error {
	tag, err := pool.Exec(context.Background(), "UPDATE products SET name = $1, price = $2, " +
		"description = $3, sku = $4 WHERE id = $5",
		prod.Name, prod.Price, prod.Description, prod.SKU, prod.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return RecordNotFound
	}
	return nil
}
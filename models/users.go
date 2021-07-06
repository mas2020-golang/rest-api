/*
Package models the structure (modeling) of the models.
To avoid reading from a database the models are also inserted here.
*/
package models

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"io"
	"time"
)

var (
	Users UsersT = UsersT{}
)

// User defines the structure for an API user
type User struct {
	ID          int     `json:"user-id"`
	Username    string  `json:"username"`
	Description string  `json:"description"`
	Email       string  `json:"email"`
	ApiKey      string `json:"api-key"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
	Disabled    bool    `json:"disabled"`
}

// FromJSON fills User decoding the JSON read from the reader.
// Returns an error in case of any.
func (p *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// ToJSON encode the User object into a json representation in []byte
func (p *User) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// UsersT type to return directly a slice of Product
type UsersT []*User

// ToJSON encode the Products object into a json representation in []byte
func (p *UsersT) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// SearchByUserPwd return a User object if the username and password match with the specified argument, otherwise
// it returns nil
func (p *UsersT) SearchByUserPwd(pool *pgxpool.Pool, username, password string) (user *User, err error) {
	row := pool.QueryRow(context.Background(), `
	SELECT user_id, username, created FROM users WHERE username=$1 and api_key=$2`,
		username, password)

	user = new(User)
	err = row.Scan(&user.ID, &user.Username, &user.Created)
	if err != nil {
		return nil, err
	}

	return user, err
}

// GetAll returns a slice of *Product.
func (p *UsersT) GetAll(pool *pgxpool.Pool) (ProductsT, error) {
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
func (p *UsersT) Get(pool *pgxpool.Pool, id int) (prod *Product, err error) {
	row := pool.QueryRow(context.Background(), "SELECT id, name, description, price, sku FROM products WHERE id=$1",
		id)
	prod = new(Product)
	err = row.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.SKU)
	return prod, err
}

// Add a new product to the products table
func (p *UsersT) Add(pool *pgxpool.Pool, new *Product) error {
	err := pool.QueryRow(context.Background(),
		"INSERT INTO products(name, price, description, sku) VALUES($1, $2, $3, $4) RETURNING id",
		(*new).Name, (*new).Price, (*new).Description, (*new).SKU).Scan(&new.ID)

	if err != nil {
		return err
	}
	return nil
}

// Update the product in the collection
func (p *UsersT) Update(pool *pgxpool.Pool, prod *Product) error {
	tag, err := pool.Exec(context.Background(), "UPDATE products SET name = $1, price = $2, "+
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

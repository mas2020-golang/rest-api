package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mas2020-golang/rest-api/models"
	"github.com/mas2020-golang/rest-api/server"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var a server.App
var token string

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    sku varchar(100),
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`

func ensureTableExists() {
	if _, err := a.DBPool.Exec(context.Background(), tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DBPool.Exec(context.Background(), "DELETE FROM products")
	a.DBPool.Exec(context.Background(), "ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func generateToken() string {
	buf := bytes.Buffer{}
	body := `{
"username": "andrea",
"password": "my-andrea-pwd"
}`
	buf.Write([]byte(body))
	req, _ := http.NewRequest("POST", "/login", &buf)
	req.Header.Set("Content-Type", "multipart/form-models")
	response := executeRequest(req)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if response.Code != 201 {
		log.Printf("Expected response code %d. Got %d (body: %s)\n", 201, response.Code, response.Body.String())
		os.Exit(1)
	}
	return m["token"]
}

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_NAME"))

	ensureTableExists()
	// get the token
	token = generateToken()
	// all the test are executed by calling m.Run()
	code := m.Run()
	// after test execution the table test models are deleted
	clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
	if body := response.Body.String(); body == "" {
		t.Errorf("Expected an error response %s", body)
	}
	// do the call using the token
	//t.Logf("token is %s", token)
	req.Header.Set("Authorization", "Bearer "+token)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if strings.Trim(response.Body.String(), "\n") != "null" {
		t.Errorf("Expected body == null, got %v", response.Body.String())
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`
{
	"name":"test product", 
	"price": 11.22, 
	"sku": "dfr-fadf-adfa"
}
`)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProducts(t *testing.T) {
	clearTable()
	// insert 5 fake products
	for i := 0; i < 5; i++ {
		// add new Product to test the update
		p := models.Product{
			Name:        fmt.Sprintf("test-%d", i),
			Description: fmt.Sprintf("test-%d", i),
			Price:       100 + float32(i),
			SKU:         "dsda-asd-asd",
		}
		err := models.Products.Add(a.DBPool, &p)
		if err != nil {
			t.Error("error occurred during the product creation")
		}
	}

	// first store info on the existing resource
	req, _ := http.NewRequest("GET", "/products", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// check the content
	var products []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &products)
	if len(products) != 5 {
		t.Errorf("Expected 5 products. Got %d", len(products))
	}
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	// add new Product to test the update
	p := models.Product{
		Name:        "test",
		Description: "test",
		Price:       100,
		SKU:         "dsda-asd-asd",
	}
	err := models.Products.Add(a.DBPool, &p)
	if err != nil {
		t.Error("error occurred during the product creation")
	}
	// first store info on the existing resource
	req, _ := http.NewRequest("GET", "/products/1", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	// update call
	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22,"sku": "dfr-fadf-adfa"}`)
	req, _ = http.NewRequest("PUT", "/products/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	// get call to read the existing resource values
	req, _ = http.NewRequest("GET", "/products/1", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"],
			"test product - updated name", m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"],
			11.22, m["price"])
	}
}

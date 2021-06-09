package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/mas2020-golang/rest-api/data"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
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

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_NAME"))

	ensureTableExists()
	// all the tests are executed by calling m.Run()
	code := m.Run()
	// after test execution the table test data are deleted
	clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]\n" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`{"name":"test product", "price": 11.22, "sku": "dfr-fadf-adfa"}`)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

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

func TestUpdateProduct(t *testing.T) {

	clearTable()
	// add new Product to test the update
	p := data.Product{
		ID:          0,
		Name:        "test",
		Description: "test",
		Price:       100,
		SKU:         "dsda-asd-asd",
	}
	err := p.Add(a.DBPool)
	if err != nil {
		t.Error("error occurred during the product creation")
	}
	req, _ := http.NewRequest("GET", "/products/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22,"sku": "dfr-fadf-adfa"}`)
	req, _ = http.NewRequest("PUT", "/products/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

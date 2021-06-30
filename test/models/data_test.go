package models

import (
	"github.com/mas2020-golang/rest-api/models"
	"testing"
)

// TestValidation test the Product validation
func TestValidation(t *testing.T){
	p := &models.Product{}
	p.Name = "Test"
	p.Price = 1.99
	p.SKU = "ads-fdsd-sdas"
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}

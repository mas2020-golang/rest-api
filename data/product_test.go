package data

import "testing"

// TestValidation tests the Product validation
func TestValidation(t *testing.T){
	p := &Product{}
	p.Name = "Test"
	p.Price = 1.99
	p.SKU = "ads-fdsd-sdas"
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}

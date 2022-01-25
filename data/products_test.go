package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "nics",
		Price: 20.0,
		SKU: "abs-abc-def",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
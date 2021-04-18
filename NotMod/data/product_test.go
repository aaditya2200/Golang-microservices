package data

import (
	"testing"
)

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "Masala chai",
		Price: 3.40,
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
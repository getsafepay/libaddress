//go:generate go run generator/main.go

// Package libaddress is a library that validates address using data
// generated from Google's Address Data Service
package libaddress

import "fmt"

// Address represents a valid address made up of its
// child components
type Address struct {
	Country            string
	Name               string
	Organization       string
	StreetAddress      []string
	DependentLocality  string
	Locality           string
	AdministrativeArea string
	PostCode           string
	SortingCode        string
}

// New creates a new unvalidated address. The validity of the address
// should be checked using the validator.
func New(fields ...func(*Address)) Address {
	address := Address{}
	for _, f := range fields {
		f(&address)
	}
	return address
}

// NewValid creates a new Address. If the address is invalid, an error
// is returned. In the case where an error is returned, the error is a
// hashicorp/go-multierror (https://github.com/hashicorp/go-multierror).
// You can use a type switch to get a list of validation errors for
// the address.
func NewValid(fields ...func(*Address)) (Address, error) {

	address := New(fields...)

	err := Validate(address)

	if err != nil {
		return address, fmt.Errorf("invalid address: %s", err)
	}

	return address, nil
}

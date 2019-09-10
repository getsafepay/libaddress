//go:generate go run generator/main.go

// Package libaddress is a library that validates address using data
// generated from Google's Address Data Service
package libaddress

type formatData struct {
	Country                     string
	CountryEn                   string
	Name                        string
	Organization                string
	StreetAddress               []string
	DependentLocality           string
	Locality                    string
	AdministrativeArea          string
	AdministrativeAreaPostalKey string
	PostCode                    string
	SortingCode                 string
}

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

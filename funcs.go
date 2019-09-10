package libaddress

import "strings"

// WithCountry sets the country code of an address.
// The country code must be an ISO 3166-1 country code.
func WithCountry(cc string) func(*Address) {
	return func(a *Address) {
		a.Country = strings.ToUpper(cc)
	}
}

// WithName sets the addressee's name of an address.
func WithName(name string) func(*Address) {
	return func(a *Address) {
		a.Name = name
	}
}

// WithOrganization sets teh addressee's organization
// of an address
func WithOrganization(organization string) func(*Address) {
	return func(a *Address) {
		a.Organization = organization
	}
}

// WithStreetAddress sets the street address of an address.
// The street address is a slice of strings, with each
// element representing an address line.
func WithStreetAddress(streetAddress []string) func(*Address) {
	return func(a *Address) {
		a.StreetAddress = streetAddress
	}
}

// WithDependentLocality sets the dependent locality (commonly
// know as the suburb) of an address. If the country of the
// address has a list of dependent localities, then the key of
// the dependent locality should be used, otherwise the validation
// will fail.
func WithDependentLocality(dependentLocality string) func(*Address) {
	return func(a *Address) {
		a.DependentLocality = dependentLocality
	}
}

// WithLocality sets the locality (commonly known as the city) of an
// address. If the country of the address has a list of localities,
// then the keys of the locality should be used, otherwise validation
// will fail.
func WithLocality(locality string) func(*Address) {
	return func(a *Address) {
		a.Locality = locality
	}
}

// WithAdministrativeArea sets the administrative area (commonly known
// as the state) of an address. If the country of the address has a list
// of administrative areas, then the key of the administrative area
// should be used, otherwise, the validation will fail.
func WithAdministrativeArea(administrativeArea string) func(*Address) {
	return func(a *Address) {
		a.AdministrativeArea = administrativeArea
	}
}

// WithPostCode sets the post code of an address.
func WithPostCode(postCode string) func(*Address) {
	return func(a *Address) {
		a.PostCode = postCode
	}
}

// WithSortingCode sets the sorting code of an address.
func WithSortingCode(sortingCode string) func(*Address) {
	return func(a *Address) {
		a.SortingCode = sortingCode
	}
}

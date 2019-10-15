package libaddress

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidCountryCode indicates that the country code
	// used to create an address is invalid
	ErrInvalidCountryCode = errors.New("invalid:Country")

	// ErrInvalidDependentLocality indicates that the dependent
	// locality is invalid. This is usually due to the country
	// having a pre-determined list of dependent localities and
	// the value does not match any of the keys in the list of
	// dependent localities.
	ErrInvalidDependentLocality = errors.New("invalid:DependentLocality")

	// ErrInvalidLocality indicates that the locality is invalid.
	// This is usually due to the country having a pre-determined
	// list of localities and the value does not match any of the
	// keys in the list of localities
	ErrInvalidLocality = errors.New("invalid:Locality")

	// ErrInvalidAdministrativeArea indicates that the administrative
	// area is invalid. This is usually due to the country having a
	// pre-determined list of administrative areas and the value
	// does not match any of the keys in the list of administrative
	// areas.
	ErrInvalidAdministrativeArea = errors.New("invalid:AdministrativeArea")

	// ErrInvalidPostCode indicates that the post code did not validate
	// using the regular expression of the country
	ErrInvalidPostCode = errors.New("invalid:PostCode")
)

// ErrMissingRequiredFields indicates that a required address
// field is missing. The Fields field can be used can be used
// to get a list of missing field.
type ErrMissingRequiredFields struct {
	country string
	Fields  []Field
}

func (e ErrMissingRequiredFields) Error() string {
	var fieldStr []string
	for _, field := range e.Fields {
		fieldStr = append(fieldStr, field.String())
	}

	return fmt.Sprintf(
		"missing required fields:%s",
		strings.Join(fieldStr, ","),
	)
}

// ErrUnsupportedFields indicates that an address field is provided,
// but it is not supported by the address format of the country. The
// Fields field can be used to get a list of unsupported fields.
type ErrUnsupportedFields struct {
	country string
	Fields  []Field
}

func (e ErrUnsupportedFields) Error() string {
	var fieldStr []string
	for _, field := range e.Fields {
		fieldStr = append(fieldStr, field.String())
	}

	return fmt.Sprintf(
		"unsupported fields for:%s",
		strings.Join(fieldStr, ","),
	)
}

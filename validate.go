package libaddress

import (
	"github.com/hashicorp/go-multierror"
	"regexp"
	"strings"
)

// Validate checks an address to determine if it
// is valid. To create a valid address, the
// `address.NewValid()` function can do it
// in one call.
func Validate(address Address) error {
	var result *multierror.Error

	if !generated.hasCountry(address.Country) {
		result = multierror.Append(result, ErrInvalidCountryCode)
		return result
	}

	var data country = generated.getCountry(address.Country)

	if err := checkRequiredFields(address, data.RequiredFields); err != nil {
		result = multierror.Append(result, err)
	}

	if err := checkAllowedFields(address, data.AllowedFields); err != nil {
		result = multierror.Append(result, err)
	}

	if len(data.AdministrativeAreas) > 0 {
		if adminAreaData, ok := data.AdministrativeAreas[data.DefaultLanguage]; ok {
			if err := checkSubdivisions(address, adminAreaData); err != nil {
				result = multierror.Append(result, err.(*multierror.Error).Errors...)
			}
		}
	}

	if address.PostCode != "" {
		if err := checkPostCode(address, data.PostCodeRegex); err != nil {
			result = multierror.Append(result, err.(*multierror.Error).Errors...)
		}

	}

	return result.ErrorOrNil()
}

func checkRequiredFields(address Address, required map[Field]struct{}) error {
	errors := ErrMissingRequiredFields{
		country: address.Country,
	}

	for field := range required {
		switch field {
		case Name:
			if strings.TrimSpace(address.Name) == "" {
				errors.Fields = append(errors.Fields, Name)
			}
		case Organization:
			if strings.TrimSpace(address.Organization) == "" {
				errors.Fields = append(errors.Fields, Organization)
			}

		case StreetAddress:
			isEmpty := true
			for _, line := range address.StreetAddress {
				if strings.TrimSpace(line) != "" {
					isEmpty = false
					break
				}
			}
			if isEmpty {
				errors.Fields = append(errors.Fields, StreetAddress)
			}

		case DependentLocality:
			if strings.TrimSpace(address.DependentLocality) == "" {
				errors.Fields = append(errors.Fields, DependentLocality)
			}

		case Locality:
			if strings.TrimSpace(address.Locality) == "" {
				errors.Fields = append(errors.Fields, Locality)
			}

		case AdministrativeArea:
			if strings.TrimSpace(address.AdministrativeArea) == "" {
				errors.Fields = append(errors.Fields, AdministrativeArea)
			}

		case PostCode:
			if strings.TrimSpace(address.PostCode) == "" {
				errors.Fields = append(errors.Fields, PostCode)
			}

		case SortingCode:
			if strings.TrimSpace(address.SortingCode) == "" {
				errors.Fields = append(errors.Fields, SortingCode)
			}
		}
	}

	if len(errors.Fields) <= 0 {
		return nil
	}

	return errors
}

func checkAllowedFields(address Address, allowed map[Field]struct{}) error {
	errors := ErrUnsupportedFields{
		country: address.Country,
	}

	if _, ok := allowed[Name]; address.Name != "" && !ok {
		errors.Fields = append(errors.Fields, Name)
	}

	if _, ok := allowed[Organization]; address.Organization != "" && !ok {
		errors.Fields = append(errors.Fields, Organization)
	}

	if _, ok := allowed[StreetAddress]; len(address.StreetAddress) > 0 && !ok {
		errors.Fields = append(errors.Fields, StreetAddress)
	}

	if _, ok := allowed[DependentLocality]; address.DependentLocality != "" && !ok {
		errors.Fields = append(errors.Fields, DependentLocality)
	}

	if _, ok := allowed[Locality]; address.Locality != "" && !ok {
		errors.Fields = append(errors.Fields, Locality)
	}

	if _, ok := allowed[AdministrativeArea]; address.AdministrativeArea != "" && !ok {
		errors.Fields = append(errors.Fields, AdministrativeArea)
	}

	if _, ok := allowed[PostCode]; address.PostCode != "" && !ok {
		errors.Fields = append(errors.Fields, PostCode)
	}

	if _, ok := allowed[SortingCode]; address.SortingCode != "" && !ok {
		errors.Fields = append(errors.Fields, SortingCode)
	}

	if len(errors.Fields) <= 0 {
		return nil
	}

	return errors
}

func checkSubdivisions(address Address, data administrativeAreaSlice) error {
	var err *multierror.Error

	if address.AdministrativeArea != "" {
		adminAreaIndex := -1

		for i, area := range data {
			if area.ID == address.AdministrativeArea {
				adminAreaIndex = i
			}
		}

		if adminAreaIndex == -1 {
			err = multierror.Append(err, ErrInvalidAdministrativeArea)
			return err.ErrorOrNil()
		}

		localities := data[adminAreaIndex].Localities
		localityIndex := -1

		if address.Locality == "" || len(localities) <= 0 {
			return err.ErrorOrNil()
		}

		for i, locality := range localities {
			if locality.ID == address.Locality {
				localityIndex = i
			}
		}

		if localityIndex == -1 {
			err = multierror.Append(err, ErrInvalidLocality)
			return err.ErrorOrNil()
		}

		dependentLocalities := localities[localityIndex].DependentLocalities
		dependentLocalityIndex := -1

		if address.DependentLocality == "" || len(dependentLocalities) <= 0 {
			return err.ErrorOrNil()
		}

		for i, dl := range dependentLocalities {
			if dl.ID == address.DependentLocality {
				dependentLocalityIndex = i
			}
		}

		if dependentLocalityIndex == -1 {
			err = multierror.Append(err, ErrInvalidDependentLocality)
			return err.ErrorOrNil()
		}
	}

	return err.ErrorOrNil()
}

func checkPostCode(address Address, regex postCodeRegex) error {
	var err *multierror.Error

	if address.PostCode != "" && regex.regex != "" {
		country := regex
		countryRegex := regexp.MustCompile(country.regex)

		if !countryRegex.MatchString(address.PostCode) {
			err = multierror.Append(err, ErrInvalidPostCode)
			return err.ErrorOrNil()
		}

		if area, ok := country.subdivisionRegex[address.AdministrativeArea]; ok {
			areaRegex := regexp.MustCompile(area.regex)

			if !areaRegex.MatchString(address.PostCode) {
				err = multierror.Append(err, ErrInvalidPostCode)
				return err.ErrorOrNil()
			}

			if locality, ok := area.subdivisionRegex[address.Locality]; ok {
				localityRegex := regexp.MustCompile(locality.regex)

				if !localityRegex.MatchString(address.PostCode) {
					err = multierror.Append(err, ErrInvalidPostCode)
					return err.ErrorOrNil()
				}

				if dependentLocality, ok := locality.subdivisionRegex[address.DependentLocality]; ok {
					dependentLocalityRegex := regexp.MustCompile(dependentLocality.regex)

					if !dependentLocalityRegex.MatchString(address.PostCode) {
						err = multierror.Append(err, ErrInvalidPostCode)
						return err.ErrorOrNil()
					}
				}
			}
		}
	}

	return err.ErrorOrNil()
}

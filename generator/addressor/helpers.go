package addressor

import (
	"fmt"
	"github.com/ziyadparekh/libaddress"
	"regexp"
)

func getAllowedFields(format string) map[libaddress.Field]struct{} {
	allowed := make(map[libaddress.Field]struct{})
	fields := ADDRESS_FORMAT_REGEX.FindAllString(format, -1)

	for _, field := range fields {
		switch field {
		case "%N":
			allowed[libaddress.Name] = struct{}{}
		case "%O":
			allowed[libaddress.Organization] = struct{}{}
		case "%A":
			allowed[libaddress.StreetAddress] = struct{}{}
		case "%D":
			allowed[libaddress.DependentLocality] = struct{}{}
		case "%C":
			allowed[libaddress.Locality] = struct{}{}
		case "%Z":
			allowed[libaddress.PostCode] = struct{}{}
		case "%X":
			allowed[libaddress.SortingCode] = struct{}{}
		}
	}

	return allowed
}

func getFields(fields string) map[libaddress.Field]struct{} {
	upper := make(map[libaddress.Field]struct{})

	for _, field := range fields {
		switch string(field) {
		case "N":
			upper[libaddress.Name] = struct{}{}
		case "O":
			upper[libaddress.Organization] = struct{}{}
		case "A":
			upper[libaddress.StreetAddress] = struct{}{}
		case "D":
			upper[libaddress.DependentLocality] = struct{}{}
		case "C":
			upper[libaddress.Locality] = struct{}{}
		case "S":
			upper[libaddress.AdministrativeArea] = struct{}{}
		case "Z":
			upper[libaddress.PostCode] = struct{}{}
		case "X":
			upper[libaddress.SortingCode] = struct{}{}
		}
	}

	return upper
}

func fieldNameToConstant(fieldname string) (libaddress.FieldName, error) {
	switch fieldname {
	case "area":
		return libaddress.Area, nil
	case "city":
		return libaddress.City, nil
	case "county":
		return libaddress.County, nil
	case "department":
		return libaddress.Department, nil
	case "district":
		return libaddress.District, nil
	case "do_si":
		return libaddress.DoSi, nil
	case "eircode":
		return libaddress.Eircode, nil
	case "emirate":
		return libaddress.Emirate, nil
	case "island":
		return libaddress.Island, nil
	case "neighborhood":
		return libaddress.Neighborhood, nil
	case "oblast":
		return libaddress.Oblast, nil
	case "pin":
		return libaddress.PINCode, nil
	case "parish":
		return libaddress.Parish, nil
	case "post_town":
		return libaddress.PostTown, nil
	case "postal":
		return libaddress.PostalCode, nil
	case "prefecture":
		return libaddress.Prefecture, nil
	case "province":
		return libaddress.Province, nil
	case "state":
		return libaddress.State, nil
	case "suburb":
		return libaddress.Suburb, nil
	case "townland":
		return libaddress.Townland, nil
	case "village_township":
		return libaddress.VillageTownship, nil
	case "zip":
		return libaddress.ZipCode, nil
	default:
		return libaddress.FieldName(-1), fmt.Errorf("unknown field name: %s", fieldname)
	}
}

func checkPostalCodeRegex(regex string, postalcodes []string) error {
	pcr, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("unable to compile regex: %s", regex)
	}

	for _, pc := range postalcodes {
		if !pcr.MatchString(pc) {
			return fmt.Errorf(
				"sample postcode %s could not be validated by regex %s",
				pc, regex,
			)
		}
	}

	return nil
}

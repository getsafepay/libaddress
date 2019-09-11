package libaddress

import (
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"sort"
)

// DependentLocalityData contains the name and ID of a
// dependent locality. The ID must be passed to
// WithDependentLocalityData() when creating an address.
// The name is useful for displaying to the end user.
type DependentLocalityData struct {
	ID   string
	Name string
}
type DependentLocalityDataSlice []DependentLocalityData

// LocalityData contains the name and ID of a
// locality. The ID must be passed to
// WithLocalityData() when creating an address.
// The name is useful for displaying to the end user.
type LocalityData struct {
	ID   string
	Name string

	DependentLocalities DependentLocalityDataSlice
}
type LocalitySlice []LocalityData

// AdministrativeAreaData contains the name and ID of an
// administrative area. The ID must be passed to
// WithAdministrativeArea() when creating an address.
// The name is useful for displaying to the end user.
type AdministrativeAreaData struct {
	ID   string
	Name string

	Localities LocalitySlice
}
type AdministrativeAreaSlice []AdministrativeAreaData

// PostCodeRegexData contains regular expressions for
// validating post codes for a given country.
// If the country has subdivisions
// (administrative areas, localities and dependent localities),
// the SubdivisionRegex field may contain further regular
// expressions to Validate the post code.
type PostCodeRegexData struct {
	Regex            string
	SubdivisionRegex map[string]PostCodeRegexData
}

// CountryData contains the address data for a country.
// The AdministrativeAreas field contains a list of nested
// subdivisions (administrative areas, localities and dependent
// localities) grouped by their translated languages. They
// are also sorted according to the sort order of the languages
// they are in.
type CountryData struct {
	Format                     string
	LatinizedFormat            string
	Required                   []Field
	Allowed                    []Field
	DefaultLanguage            string
	AdministrativeAreaNameType FieldName
	LocalityNameType           FieldName
	DependentLocalityNameType  FieldName
	PostCodeNameType           FieldName
	PostCodeRegex              PostCodeRegexData
	AdministrativeAreas        map[string]AdministrativeAreaSlice
}

// CountryListItem represents a single country
// containing the ISO 3166-1 code and the name
// of the country.
type CountryListItem struct {
	Code string
	Name string
}

// CountryList contains a list of countries that can be used
// to create addresses.
type CountryList []CountryListItem

// Len returns the number of countries in the list. This is used
// for sorting the countries and would not be generally used in
// client code.
func (c CountryList) Len() int {
	return len(c)
}

// Swap swaps 2 countries in the list. This is used for sorting
// the countries and would not be generally used in client code.
func (c CountryList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Bytes returns a country name in bytes. This is used for sorting
// the countries and would not be generally used in client code.
func (c CountryList) Bytes(i int) []byte {
	return []byte(c[i].Name)
}

// ListCountries returns a list of countries that can be used to create
// addresses. Language must be a valid ISO 639-1 language code such as
// jp, zh etc. If the language does not have any translations or is invalid
// then english is used as the fallback. The returned list of countries is
// sorted according to the chosen language.
func ListCountries(l string) CountryList {
	lang, err := language.Parse(l)
	if err != nil {
		lang = language.English
	}

	c := collate.New(lang)
	n := display.Regions(lang)

	if n == nil {
		n = display.Regions(language.English)
	}

	var list CountryList
	for country := range generated {
		if country == "ZZ" {
			continue
		}
		cc := language.MustParseRegion(country)
		list = append(list, CountryListItem{
			Code: country,
			Name: n.Name(cc),
		})
	}

	c.Sort(list)
	return list
}

// Get country returns address information for a given country.
func GetCountry(cc string) CountryData {
	country := generated.getCountry(cc)
	return internalToExternalCountry(country)
}

func internalToExternalCountry(c country) CountryData {
	data := CountryData{
		Format:                     c.Format,
		LatinizedFormat:            c.LatinizedFormat,
		DefaultLanguage:            c.DefaultLanguage,
		AdministrativeAreaNameType: c.AdministrativeAreaNameType,
		LocalityNameType:           c.LocalityNameType,
		DependentLocalityNameType:  c.DependentLocalityNameType,
		PostCodeNameType:           c.PostCodeNameType,
		PostCodeRegex:              internalToExternalPostCodeRegex(c.PostCodeRegex),
	}

	var required []Field
	for field := range c.RequiredFields {
		required = append(required, field)
	}

	sort.Slice(required, func(i, j int) bool {
		return required[i].String() < required[j].String()
	})

	data.Required = required

	var allowed []Field

	for field := range c.AllowedFields {
		allowed = append(allowed, field)
	}

	sort.Slice(allowed, func(i, j int) bool {
		return allowed[i].String() < allowed[j].String()
	})

	data.Allowed = allowed

	administrativeAreas := map[string]AdministrativeAreaSlice{}

	for lang, adminAreas := range c.AdministrativeAreas {
		administrativeAreas[lang] = internalToExternalAdministrativeArea(adminAreas)
	}

	if len(administrativeAreas) > 0 {
		data.AdministrativeAreas = administrativeAreas
	}

	return data

}

func internalToExternalPostCodeRegex(regex postCodeRegex) PostCodeRegexData {
	result := PostCodeRegexData{
		Regex: regex.regex,
	}

	for subID, regex := range regex.subdivisionRegex {
		if result.SubdivisionRegex == nil {
			result.SubdivisionRegex = map[string]PostCodeRegexData{}
		}
		result.SubdivisionRegex[subID] = internalToExternalPostCodeRegex(regex)
	}

	return result
}

func internalToExternalAdministrativeArea(areas administrativeAreaSlice) AdministrativeAreaSlice {
	var result AdministrativeAreaSlice
	for _, area := range areas {
		var localities LocalitySlice
		for _, locality := range area.Localities {

			var dependentLocalities DependentLocalityDataSlice
			for _, dependentLocality := range locality.DependentLocalities {
				dependentLocalities = append(dependentLocalities, DependentLocalityData{
					ID:   dependentLocality.ID,
					Name: dependentLocality.Name,
				})
			}

			ld := LocalityData{
				ID:   locality.ID,
				Name: locality.Name,
			}

			if len(dependentLocalities) > 0 {
				ld.DependentLocalities = dependentLocalities
			}

			localities = append(localities, ld)
		}

		aad := AdministrativeAreaData{
			ID:   area.ID,
			Name: area.Name,
		}

		if len(localities) > 0 {
			aad.Localities = localities
		}

		result = append(result, aad)
	}

	return result
}

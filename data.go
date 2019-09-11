package libaddress

// DependentLocalityData contains the name and ID of a
// dependent locality. The ID must be passed to
// WithDependentLocalityData() when creating an address.
// The name is useful for displaying to the end user.
type DependentLocalityData struct {
	ID string
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
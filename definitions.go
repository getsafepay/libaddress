package libaddress

type country struct {
	ID              string
	Name            string
	DefaultLanguage string
	PostCodePrefix  string
	PostCodeRegex   postCodeRegex
	Format          string
	LatinizedFormat string

	AdministrativeAreaNameType FieldName
	LocalityNameType           FieldName
	DependentLocalityNameType  FieldName
	PostCodeNameType           FieldName

	AllowedFields  map[Field]struct{}
	RequiredFields map[Field]struct{}
	Upper          map[Field]struct{}

	AdministrativeAreas map[string]administrativeAreaSlice
}

type postCodeRegex struct {
	regex            string
	subdivisionRegex map[string]postCodeRegex
}

type administrativeArea struct {
	ID         string
	Name       string
	PostalKey  string
	Localities localitySlice
}

type locality struct {
	ID                  string
	Name                string
	DependentLocalities dependentLocalitySlice
}

type dependentLocality struct {
	ID   string
	Name string
}

type data map[string]country
type administrativeAreaSlice []administrativeArea
type localitySlice []locality
type dependentLocalitySlice []dependentLocality

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

	AdministrativeAreas map[string][]administrativeArea
}

type postCodeRegex struct {
	regex            string
	subdivisionRegex map[string]postCodeRegex
}

type administrativeArea struct {
	ID         string
	Name       string
	PostalKey  string
	Localities []locality
}

type locality struct {
	ID                  string
	Name                string
	DependentLocalities []dependentLocality
}

type dependentLocality struct {
	ID   string
	Name string
}

type data map[string]country
type administrativeAreaSlice []administrativeArea
type dependentLocalitySlice []dependentLocality

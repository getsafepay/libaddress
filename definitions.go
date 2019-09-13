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

// var generated data

func (d data) getCountry(cc string) country {
	var data country = generated[cc]
	var defaults country = generated["ZZ"]

	if data.Format == "" {
		data.Format = defaults.Format
	}

	if data.AdministrativeAreaNameType == 0 {
		data.AdministrativeAreaNameType = defaults.AdministrativeAreaNameType
	}

	if data.LocalityNameType == 0 {
		data.LocalityNameType = defaults.LocalityNameType
	}

	if data.DependentLocalityNameType == 0 {
		data.DependentLocalityNameType = defaults.DependentLocalityNameType
	}

	if data.PostCodeNameType == 0 {
		data.PostCodeNameType = defaults.PostCodeNameType
	}

	if len(data.AllowedFields) <= 0 {
		data.AllowedFields = defaults.AllowedFields
	}

	if len(data.RequiredFields) <= 0 {
		data.RequiredFields = defaults.RequiredFields
	}

	if len(data.Upper) <= 0 {
		data.Upper = defaults.Upper
	}

	return data
}

func (d data) hasCountry(cc string) bool {
	if _, ok := generated[cc]; ok {
		return true
	}
	return false
}

func (d data) getAdministrativeAreaName(cc, areaID, language string) string {
	data := d.getCountry(cc)
	lang := d.normalizeLanguage(cc, language)

	for _, area := range data.AdministrativeAreas[lang] {
		if area.ID == areaID {
			return area.Name
		}
	}

	return ""
}

func (d data) getAdministrativeAreaPostalKey(cc, areaID string) string {
	data := d.getCountry(cc)
	lang := d.normalizeLanguage(cc, "")

	for _, area := range data.AdministrativeAreas[lang] {
		if area.ID == areaID {
			return area.PostalKey
		}
	}

	return ""
}

func (d data) getLocalityName(cc, areaID, localityID, language string) string {
	data := d.getCountry(cc)
	lang := d.normalizeLanguage(cc, language)

	for _, area := range data.AdministrativeAreas[lang] {
		if area.ID == areaID {
			for _, l := range area.Localities {
				if l.ID == localityID {
					return l.Name
				}
			}
		}
	}

	return ""
}

func (d data) getDependentLocalityName(cc, areaID, lID, dlID, language string) string {
	data := d.getCountry(cc)
	lang := d.normalizeLanguage(cc, language)

	for _, area := range data.AdministrativeAreas[lang] {
		if area.ID == areaID {
			for _, l := range area.Localities {
				if l.ID == lID {
					for _, d := range l.DependentLocalities {
						if d.ID == dlID {
							return d.Name
						}
					}
				}
			}
		}
	}

	return ""
}

func (d data) normalizeLanguage(cc, language string) string {
	country := d.getCountry(cc)
	if _, ok := country.AdministrativeAreas[language]; ok {
		return language
	}
	return country.DefaultLanguage
}

package libaddress

import "sort"

type DisplayFieldName struct {
	ID      FieldName `json:"id"`
	Display string    `json:"display"`
}

func fieldNameToDisplay(fn FieldName) *DisplayFieldName {
	return &DisplayFieldName{
		ID:      fn,
		Display: fn.Readable(),
	}
}

type DisplayField struct {
	ID      Field  `json:"id"`
	Display string `json:"display"`
}

func fieldToDisplay(fn Field) *DisplayField {
	return &DisplayField{
		ID:      fn,
		Display: fn.Readable(),
	}
}

type ExternalCountry struct {
	Required                   []*DisplayField                    `json:"required"`
	Allowed                    []*DisplayField                    `json:"allowed"`
	Language                   string                             `json:"language"`
	AdministrativeAreaNameType *DisplayFieldName                  `json:"administrative_area_name_type"`
	LocalityNameType           *DisplayFieldName                  `json:"locality_name_type"`
	DependentLocalityNameType  *DisplayFieldName                  `json:"dependent_locality_name_type"`
	PostCodeNameType           *DisplayFieldName                  `json:"post_code_name_type"`
	AdministrativeAreas        map[string]AdministrativeAreaSlice `json:"administrative_areas"`
}

func Externalize(c country) ExternalCountry {
	data := ExternalCountry{
		Language:                   c.DefaultLanguage,
		AdministrativeAreaNameType: fieldNameToDisplay(c.AdministrativeAreaNameType),
		LocalityNameType:           fieldNameToDisplay(c.LocalityNameType),
		DependentLocalityNameType:  fieldNameToDisplay(c.DependentLocalityNameType),
		PostCodeNameType:           fieldNameToDisplay(c.PostCodeNameType),
	}

	var required []*DisplayField
	for field := range c.RequiredFields {
		required = append(required, fieldToDisplay(field))
	}

	sort.Slice(required, func(i, j int) bool {
		return required[i].ID.String() < required[j].ID.String()
	})

	data.Required = required

	var allowed []*DisplayField
	for field := range c.AllowedFields {
		allowed = append(allowed, fieldToDisplay(field))
	}

	sort.Slice(allowed, func(i, j int) bool {
		return allowed[i].ID.String() < allowed[j].ID.String()
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

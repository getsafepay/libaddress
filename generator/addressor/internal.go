package addressor

import (
	"fmt"
	"github.com/ziyadparekh/libaddress"
	"sort"
)

type postCodeRegex struct {
	regex            string
	subdivisionRegex map[string]postCodeRegex
}

func (p postCodeRegex) toCode() string {
	// Generate postcode regex in order
	// to avoid huge diffs when updating
	// the data

	var ids []string

	for id := range p.subdivisionRegex {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	str := fmt.Sprintf(`{
		regex: `+"`%s`,", p.regex,
	)

	if len(p.subdivisionRegex) > 0 {
		str += `
subdivisionRegex: map[string]postCodeRegex{
`

		for _, id := range ids {
			str += fmt.Sprintf(
				`"%s": %s,`,
				id,
				p.subdivisionRegex[id].toCode(),
			)
		}

		str += `}`
	}

	str += `}`

	return str
}

type dependentLocalitySlice []dependentLocality
type dependentLocality struct {
	ID   string
	Name string
}

func (d dependentLocality) toCode() string {
	return fmt.Sprintf(`{
		ID: "%s",
		Name: "%s",
	}`, d.ID, d.Name)
}

type localitySlice []locality
type locality struct {
	ID   string
	Name string

	DependentLocalities dependentLocalitySlice
}

func (l locality) toCode() string {
	str := fmt.Sprintf(`{
		ID: "%s",
		Name: "%s",`,
		l.ID, l.Name,
	)

	if len(l.DependentLocalities) > 0 {
		str += `
DependentLocalities: dependentLocalitySlice{
`

		for _, dl := range l.DependentLocalities {
			str += dl.toCode() + ",\n"
		}

		str += `
},`
	}

	str += `
}`

	return str
}

type administrativeAreaSlice []administrativeArea
type administrativeArea struct {
	ID        string
	Name      string
	PostalKey string

	Localities localitySlice
}

func (a administrativeArea) toCode() string {
	str := fmt.Sprintf(`{
		ID: "%s",
		Name: "%s",
		PostalKey: "%s",`,
		a.ID, a.Name, a.PostalKey,
	)

	if len(a.Localities) > 0 {
		str += `
Localities: localitySlice{
`

		for _, l := range a.Localities {
			str += l.toCode() + ",\n"
		}

		str += `
},`
	}

	str += `
}`

	return str
}

type Country struct {
	ID   string
	Name string

	DefaultLanguage string

	PostCodePrefix string
	PostCodeRegex  postCodeRegex

	Format          string
	LatinizedFormat string

	AdministrativeAreaNameType libaddress.FieldName
	LocalityNameType           libaddress.FieldName
	DependentLocalityNameType  libaddress.FieldName
	PostCodeNameType           libaddress.FieldName

	AllowedFields  map[libaddress.Field]struct{}
	RequiredFields map[libaddress.Field]struct{}
	Upper          map[libaddress.Field]struct{}

	AdministrativeAreas map[string]administrativeAreaSlice
}

func (c Country) ToCode() string {
	str := fmt.Sprintf(`{
		ID: "%s",
		Name: "%s",`,
		c.ID, c.Name,
	)

	if c.DefaultLanguage != "" {
		str += fmt.Sprintf(`
DefaultLanguage: "%s",`,
			c.DefaultLanguage,
		)
	} else {
		fmt.Println(c.ID)
	}

	if c.PostCodePrefix != "" {
		str += fmt.Sprintf(`
PostCodePrefix: "%s",`, c.PostCodePrefix)
	}

	if c.PostCodeRegex.regex != "" || len(c.PostCodeRegex.subdivisionRegex) > 0 {
		str += fmt.Sprintf(`
PostCodeRegex: postCodeRegex%s,`, c.PostCodeRegex.toCode())
	}

	if c.Format != "" {
		str += fmt.Sprintf(`
Format: "%s",`, c.Format)
	}

	if c.LatinizedFormat != "" {
		str += fmt.Sprintf(`
LatinizedFormat: "%s",`, c.LatinizedFormat)
	}

	if c.AdministrativeAreaNameType != 0 {
		str += fmt.Sprintf(`
AdministrativeAreaNameType: %s,`, c.AdministrativeAreaNameType.String())
	}

	if c.LocalityNameType != 0 {
		str += fmt.Sprintf(`
LocalityNameType: %s,`, c.LocalityNameType.String())
	}

	if c.DependentLocalityNameType != 0 {
		str += fmt.Sprintf(`
DependentLocalityNameType: %s,`, c.DependentLocalityNameType.String())
	}

	if c.PostCodeNameType != 0 {
		str += fmt.Sprintf(`
PostCodeNameType: %s,`, c.PostCodeNameType.String())
	}

	if len(c.AllowedFields) > 0 {
		// Generate fields in order to avoid huge diffs when
		// updating the data.
		var fields []string
		for field := range c.AllowedFields {
			fields = append(fields, field.String())
		}

		sort.Strings(fields)

		str += fmt.Sprintf(`
AllowedFields: map[Field]struct{}{`,
		)

		for _, field := range fields {
			str += fmt.Sprintf(`
%s: {},`, field)
		}

		str += `},`
	}

	if len(c.RequiredFields) > 0 {
		// Generate fields in order to avoid huge diffs when
		// updating the data.
		var fields []string
		for field := range c.RequiredFields {
			fields = append(fields, field.String())
		}

		sort.Strings(fields)

		str += fmt.Sprintf(`
RequiredFields: map[Field]struct{}{`,
		)

		for _, field := range fields {
			str += fmt.Sprintf(`
%s: {},`, field)
		}

		str += `},`
	}

	if len(c.Upper) > 0 {
		var fields []string
		for field := range c.Upper {
			fields = append(fields, field.String())
		}

		sort.Strings(fields)

		str += fmt.Sprintf(`
Upper: map[Field]struct{}{`,
		)

		for _, field := range fields {
			str += fmt.Sprintf(`
%s: {},`, field)
		}

		str += `
},`
	}

	if len(c.AdministrativeAreas) > 0 {
		// Generate languages first in order to avoid
		// huge diffs when updating the address data
		var languages []string
		for language := range c.AdministrativeAreas {
			languages = append(languages, language)
		}

		sort.Strings(languages)

		str += fmt.Sprintf(`
AdministrativeAreas: map[string]administrativeAreaSlice{`,
		)

		for _, language := range languages {
			areas := c.AdministrativeAreas[language]

			str += fmt.Sprintf(`
"%s": {`, language)

			for _, area := range areas {
				str += fmt.Sprintf(`
%s,`, area.toCode())
			}

			str += `
},`
		}

		str += `
},`
	}

	str += `
}`

	return str
}

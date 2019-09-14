//go:generate stringer -type=Field,FieldName -output=constant_string.go
package libaddress

// Field is an address field type.
type Field int

const (
	Country Field = iota + 1
	Name
	Organization
	StreetAddress
	DependentLocality
	Locality
	AdministrativeArea
	PostCode
	SortingCode
)

func (i Field) Readable() string {
	switch i {
	case Country:
		return "Country"
	case Name:
		return "Name"
	case Organization:
		return "Organization"
	case StreetAddress:
		return "StreetAddress"
	case DependentLocality:
		return "DependentLocality"
	case Locality:
		return "Locality"
	case AdministrativeArea:
		return "AdministrativeArea"
	case PostCode:
		return "PostCode"
	case SortingCode:
		return "SortingCode"
	default:
		return ""
	}
}

// Key returns the corresponding one-letter abbreviation used by Google to refer to address
// fields. This is useful for parsing the address format for a country.
// See https://github.com/googlei18n/libaddressinput/wiki/AddressValidationMetadata
// for more information
func (i Field) Key() string {
	switch i {
	case Country:
		return "country"
	case Name:
		return "N"
	case Organization:
		return "O"
	case StreetAddress:
		return "A"
	case DependentLocality:
		return "D"
	case Locality:
		return "C"
	case AdministrativeArea:
		return "S"
	case PostCode:
		return "Z"
	case SortingCode:
		return "X"
	default:
		return ""
	}
}

// FieldName is the name to be used when referring to a field.
// For example, in India the postal code is called PIN Code instead of PostalCode.
// The field name allows you to display the appropriate form labels to the user.
type FieldName int

const (
	Area FieldName = iota + 1
	City
	County
	Department
	District
	DoSi
	Eircode
	Emirate
	Island
	Neighborhood
	Oblast
	PINCode
	Parish
	PostTown
	PostalCode
	Prefecture
	Province
	State
	Suburb
	Townland
	VillageTownship
	ZipCode
)

func (i FieldName) Readable() string {
	switch i {
	case Area:
		return "Area"
	case City:
		return "City"
	case County:
		return "County"
	case Department:
		return "Department"
	case District:
		return "District"
	case DoSi:
		return "DoSi"
	case Eircode:
		return "Eircode"
	case Emirate:
		return "Emirate"
	case Island:
		return "Island"
	case Neighborhood:
		return "Neighborhood"
	case Oblast:
		return "Oblast"
	case PINCode:
		return "PINCode"
	case Parish:
		return "Parish"
	case PostTown:
		return "PostTown"
	case PostalCode:
		return "PostalCode"
	case Prefecture:
		return "Prefecture"
	case Province:
		return "Province"
	case Suburb:
		return "Suburb"
	case Townland:
		return "Townland"
	case VillageTownship:
		return "VillageTownship"
	case ZipCode:
		return "ZipCode"
	default:
		return ""

	}
}

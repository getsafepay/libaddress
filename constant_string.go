// Code generated by "stringer -type=Field,FieldName -output=constant_string.go"; DO NOT EDIT.

package libaddress

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Country-1]
	_ = x[Name-2]
	_ = x[Organization-3]
	_ = x[StreetAddress-4]
	_ = x[DependentLocality-5]
	_ = x[Locality-6]
	_ = x[AdministrativeArea-7]
	_ = x[PostCode-8]
	_ = x[SortingCode-9]
}

const _Field_name = "CountryNameOrganizationStreetAddressDependentLocalityLocalityAdministrativeAreaPostCodeSortingCode"

var _Field_index = [...]uint8{0, 7, 11, 23, 36, 53, 61, 79, 87, 98}

func (i Field) String() string {
	i -= 1
	if i < 0 || i >= Field(len(_Field_index)-1) {
		return "Field(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Field_name[_Field_index[i]:_Field_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Area-1]
	_ = x[City-2]
	_ = x[County-3]
	_ = x[Department-4]
	_ = x[District-5]
	_ = x[DoSi-6]
	_ = x[Eircode-7]
	_ = x[Emirate-8]
	_ = x[Island-9]
	_ = x[Neighborhood-10]
	_ = x[Oblast-11]
	_ = x[PINCode-12]
	_ = x[Parish-13]
	_ = x[PostTown-14]
	_ = x[PostalCode-15]
	_ = x[Prefecture-16]
	_ = x[Province-17]
	_ = x[State-18]
	_ = x[Suburb-19]
	_ = x[Townland-20]
	_ = x[VillageTownship-21]
	_ = x[ZipCode-22]
}

const _FieldName_name = "AreaCityCountyDepartmentDistrictDoSiEircodeEmirateIslandNeighborhoodOblastPINCodeParishPostTownPostalCodePrefectureProvinceStateSuburbTownlandVillageTownshipZipCode"

var _FieldName_index = [...]uint8{0, 4, 8, 14, 24, 32, 36, 43, 50, 56, 68, 74, 81, 87, 95, 105, 115, 123, 128, 134, 142, 157, 164}

func (i FieldName) String() string {
	i -= 1
	if i < 0 || i >= FieldName(len(_FieldName_index)-1) {
		return "FieldName(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _FieldName_name[_FieldName_index[i]:_FieldName_index[i+1]]
}

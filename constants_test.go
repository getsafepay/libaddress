package libaddress

import "testing"

func TestConstantKeys(t *testing.T) {

	testCases := []struct {
		Field    Field
		Expected string
	}{
		{
			Field:    Country,
			Expected: "country",
		},
		{
			Field:    Name,
			Expected: "N",
		},
		{
			Field:    Organization,
			Expected: "O",
		},
		{
			Field:    StreetAddress,
			Expected: "A",
		},
		{
			Field:    DependentLocality,
			Expected: "D",
		},
		{
			Field:    Locality,
			Expected: "C",
		},
		{
			Field:    AdministrativeArea,
			Expected: "S",
		},
		{
			Field:    PostCode,
			Expected: "Z",
		},
		{
			Field:    SortingCode,
			Expected: "X",
		},
	}

	for _, c := range testCases {
		if c.Field.Key() != c.Expected {
			t.Errorf(
				"Expected key for %s to be %s, gor %s",
				c.Field,
				c.Expected,
				c.Field.Key(),
			)
		}
	}
}

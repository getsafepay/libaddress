package address

import "github.com/getsafepay/communist/payments/common"

type Address struct {
	Country            string   `json:"country"`
	Name               string   `json:"name"`
	StreetAddress      []string `json:"street_address"`
	DependentLocality  string   `json:"dependent_locality"`
	Locality           string   `json:"locality"`
	AdministrativeArea string   `json:"administrative_area"`
	PostCode           string   `json:"post_code"`
	SortingCode        string   `json:"sorting_code"`
}

type ValidateRequest struct {
	common.Head

	Address *Address `json:"address"`
}

func NewValidateRequest(token string, address *Address) *ValidateRequest {
	vr := new(ValidateRequest)
	vr.AccessToken = token
	vr.Address = address
	return vr
}

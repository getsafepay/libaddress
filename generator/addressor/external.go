package addressor

import (
	"encoding/json"
	"io"
)

type externalCountries struct {
	Countries string `json:"countries"`
}

func DecodeCountries(reader io.Reader) (externalCountries, error) {
	var ecs = externalCountries{}
	if err := json.NewDecoder(reader).Decode(&ecs); err != nil {
		return ecs, err
	}

	return ecs, nil
}

type externalCountry struct {
	ID  string `json:"id"`
	Key string `json:"key"`

	Lang      string `json:"lang"`
	Languages string `json:"languages"`
	Name      string `json:"name"`

	Fmt  string `json:"fmt"`
	Lfmt string `json:"lfmt"`

	StateNameType       string `json:"state_name_type"`
	LocalityNameType    string `json:"locality_name_type"`
	SubLocalityNameType string `json:"sublocality_name_type"`
	ZipNameType         string `json:"zip_name_type"`

	Require string `json:"require"`
	Upper   string `json:"upper"`

	SubISOIDs string `json:"sub_isoids"`
	SubKeys   string `json:"sub_keys"`
	SubLNames string `json:"sub_lnames"`
	SubNames  string `json:"sub_names"`

	SubMores string `json:"sub_mores"`

	SubXRequires string `json:"sub_xrequires"`
	SubXZips     string `json:"sub_xzips"`

	SubZips   string `json:"sub_zips"`
	SubZipExs string `json:"sub_zipexs"`

	PostPrefix string `json:"post_prefix"`
	Zip        string `json:"zip"`
	Zipex      string `json:"zipex"`
}

func decodeCountry(reader io.Reader) (externalCountry, error) {
	var ec = externalCountry{}
	if err := json.NewDecoder(reader).Decode(&ec); err != nil {
		return ec, err
	}

	return ec, nil
}

type externalSubdivision struct {
	ID  string `json:"id"`
	Key string `json:"key"`

	Name  string `json:"name"`
	LName string `json:"lname"`

	Lang string `json:"lang"`

	ISOID   string `json:"isoid"`
	SubKeys string `json:"sub_keys"`

	SubNames   string `json:"sub_names"`
	SubMores   string `json:"sub_mores"`
	SubLNames  string `json:"sub_lnames"`
	SubLFNames string `json:"sub_lfnames"`

	Zip       string `json:"zip"`
	ZipEx     string `json:"zipex"`
	SubZips   string `json:"sub_zips"`
	SubZipExs string `json:"sub_zipexs"`
}

func decodeSubdivision(reader io.Reader) (externalSubdivision, error) {
	var esd = externalSubdivision{}
	if err := json.NewDecoder(reader).Decode(&esd); err != nil {
		return esd, err
	}

	return esd, nil
}

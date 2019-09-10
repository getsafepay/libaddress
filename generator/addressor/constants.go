package addressor

import "regexp"

const (
	GOOGLE_ADDRESS_URL string = "https://chromium-i18n.appspot.com/ssl-address"
	NUM_WORKERS        int    = 25
)

var (
	ADDRESS_FORMAT_REGEX = regexp.MustCompile(`%[NOADCSZX]`)
	REMOVE_LANG_REGEX    = regexp.MustCompile(`--.*`)

	POST_PREFIX_FIXES = map[string]string{
		"PR": "PR ",
	}

	LANGUAGE_OVERRIDES = map[string]string{
		"AQ": "en",
		"AS": "en",
		"BQ": "nl",
		"BV": "nb",
		"CW": "nl",
		"DJ": "fr",
		"GS": "en",
		"HM": "en",
		"MV": "en",
		"PG": "en",
		"PW": "en",
		"TK": "en",
		"VU": "fr",
		"WS": "en",
	}

	LOCAL_NAME_OVERRIDES = map[string]string{
		"TV": "Tuvalu",
	}
)

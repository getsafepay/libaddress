package addressor

import (
	"fmt"
	"golang.org/x/text/language"
	"net/http"
	"strings"
)

type Result struct {
	Error   error
	Country Country
}

type Worker struct {
	CC   chan string
	Stop chan struct{}
	Res  chan Result
}

func (w *Worker) Start() {
	go func() {
		for {
		exit:
			select {
			case <-w.Stop:
				return
			case cc := <-w.CC:
				url := fmt.Sprintf("%s/data/%s", GOOGLE_ADDRESS_URL, cc)

				data, err := http.Get(url)
				if err != nil {
					res := Result{
						Error: fmt.Errorf(
							"error getting data using url(%s): %s",
							url, err.Error(),
						),
					}
					w.Res <- res
					break
				}

				ec, err := decodeCountry(data.Body)
				if err != nil {
					res := Result{
						Error: fmt.Errorf(
							"error unmarhaling data for url(%s): %s",
							url, err.Error(),
						),
					}
					w.Res <- res
					break
				}

				fmtAllowedFields := getAllowedFields(ec.Fmt)
				lFmtAllowedFields := getAllowedFields(ec.Lfmt)
				requiredFields := getFields(ec.Require)
				upperFields := getFields(ec.Upper)

				// Sanity check latinized format
				if ec.Lfmt != "" && len(fmtAllowedFields) != len(lFmtAllowedFields) {
					res := Result{
						Error: fmt.Errorf(
							`number of fields in address format and latinized 
									address format do not match for %s`,
							ec.Key,
						),
					}
					w.Res <- res
					break
				}

				// Sanity check post code regex
				if ec.Zip != "" {
					if err := checkPostalCodeRegex(ec.Zip, strings.Split(ec.Zipex, ",")); err != nil {
						res := Result{
							Error: fmt.Errorf(
								"error validating post code regex for %s: %s",
								ec.Key, err.Error(),
							),
						}
						w.Res <- res
						break
					}
				}

				c := Country{
					ID:   cc,
					Name: ec.Name,

					PostCodeRegex: postCodeRegex{
						regex: ec.Zip,
					},

					Format:          ec.Fmt,
					LatinizedFormat: ec.Lfmt,

					AllowedFields:  fmtAllowedFields,
					RequiredFields: requiredFields,
					Upper:          upperFields,
				}

				if ec.Lang != "" {
					c.DefaultLanguage = ec.Lang
				} else if l, ok := LANGUAGE_OVERRIDES[cc]; ok {
					c.DefaultLanguage = l
				} else {
					l, _ := language.Make(fmt.Sprintf("und-%s", cc)).Base()
					c.DefaultLanguage = l.String()
				}

				if ec.StateNameType != "" {
					aant, err := fieldNameToConstant(ec.StateNameType)
					if err != nil {
						res := Result{
							Error: fmt.Errorf(
								"error converting administrative area name type for %s: %s",
								ec.Key, err.Error(),
							),
						}
						w.Res <- res
						break
					}

					c.AdministrativeAreaNameType = aant
				}

				if ec.LocalityNameType != "" {
					lnt, err := fieldNameToConstant(ec.LocalityNameType)
					if err != nil {
						res := Result{
							Error: fmt.Errorf(
								"error converting locality name type for %s: %s",
								ec.Key, err.Error(),
							),
						}
						w.Res <- res
						break
					}

					c.LocalityNameType = lnt
				}

				if ec.SubLocalityNameType != "" {
					dlnt, err := fieldNameToConstant(ec.SubLocalityNameType)
					if err != nil {
						if err != nil {
							res := Result{
								Error: fmt.Errorf(
									"error converting dependent locality name type for %s: %s",
									ec.Key, err.Error(),
								),
							}
							w.Res <- res
							break
						}
					}

					c.DependentLocalityNameType = dlnt
				}

				if ec.ZipNameType != "" {
					pcnt, err := fieldNameToConstant(ec.ZipNameType)
					if err != nil {
						res := Result{
							Error: fmt.Errorf(
								"error converting postcode name type for %s: %s",
								ec.Key, err.Error(),
							),
						}
						w.Res <- res
						break
					}

					c.PostCodeNameType = pcnt
				}

				if prefix, ok := POST_PREFIX_FIXES[ec.Key]; ok {
					c.PostCodePrefix = prefix
				} else {
					c.PostCodePrefix = ec.PostPrefix
				}

				// Process subdivisions
				if ec.SubKeys != "" {
					// Sanity check
					if ec.Languages == "" {
						res := Result{
							Error: fmt.Errorf(
								"%s has subkeys but not any languages",
								ec.Key,
							),
						}
						w.Res <- res
						break
					}

					c.AdministrativeAreas = make(map[string]administrativeAreaSlice)

					// Get languages
					languages := strings.Split(ec.Languages, "~")
					if len(languages) > 1 {
						for _, l := range languages {
							if l != ec.Lang {
								data, err := http.Get(url + "--" + l)
								if err != nil {
									res := Result{
										Error: fmt.Errorf(
											"err getting language %s for country %s: %s",
											l, ec.Key, err.Error(),
										),
									}
									w.Res <- res
									break exit
								}
								ecl, err := decodeCountry(data.Body)
								if err != nil {
									res := Result{
										Error: fmt.Errorf(
											"err decoding language %s for country %s: %s",
											l, ec.Key, err.Error(),
										),
									}
									w.Res <- res
									break exit
								}

								langAdminAreas, _, err := processAdministrativeAreas(ecl, l)
								if err != nil {
									res := Result{
										Error: fmt.Errorf(
											"error processing admin areas in language"+
												"%s for country %s: %s",
											l, ec.Key, err.Error(),
										),
									}
									w.Res <- res
									break exit
								}

								for l, aa := range langAdminAreas {
									c.AdministrativeAreas[l] = aa
								}
							} else {
								aam, pcrm, err := processAdministrativeAreas(ec, "")
								if err != nil {
									res := Result{
										Error: fmt.Errorf(
											"error processing admin areas in default"+
												"language for country %s: %s",
											ec.Key, err.Error(),
										),
									}
									w.Res <- res
									break exit
								}

								c.PostCodeRegex.subdivisionRegex = pcrm
								for l, aa := range aam {
									c.AdministrativeAreas[l] = aa
								}
							}
						}
					} else {
						aam, pcrm, err := processAdministrativeAreas(ec, "")
						if err != nil {
							res := Result{
								Error: fmt.Errorf(
									"error processing admin areas in default"+
										"language for country %s: %s",
									ec.Key, err.Error(),
								),
							}
							w.Res <- res
							break exit
						}

						c.PostCodeRegex.subdivisionRegex = pcrm
						for l, aa := range aam {
							c.AdministrativeAreas[l] = aa
						}
					}
				}

				res := Result{
					Country: c,
				}
				w.Res <- res
			}
		}
	}()
}

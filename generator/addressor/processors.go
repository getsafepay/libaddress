package addressor

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func processAdministrativeAreas(ec externalCountry, lang string) (
	map[string]administrativeAreaSlice,
	map[string]postCodeRegex,
	error,
) {
	administrativeAreaMap := make(map[string]administrativeAreaSlice)
	postCodeRegexMap := make(map[string]postCodeRegex)

	subIsoIds := strings.Split(ec.SubISOIDs, "~")
	subNames := strings.Split(ec.SubNames, "~")
	subZips := strings.Split(ec.SubZips, "~")
	subMores := strings.Split(ec.SubMores, "~")
	subKeys := strings.Split(ec.SubKeys, "~")
	subZipExs := strings.Split(ec.SubZipExs, "~")
	subLNames := strings.Split(ec.SubLNames, "~")

	// Countries like China include places like Taiwan and Hong Kong in their
	// list of administrative divisions. However, these places are already
	// in the list of countries, so we check to see if they have special post
	// code regexes or required fields to filter them out
	skippedSubdivisions := make(map[string]struct{})
	if ec.SubXRequires != "" {
		for idx, requires := range strings.Split(ec.SubXRequires, "~") {
			if requires != "" {
				skippedSubdivisions[subIsoIds[idx]] = struct{}{}
			}
		}
	}

	if ec.SubXZips != "" {
		for idx, xzip := range strings.Split(ec.SubXZips, "~") {
			if xzip != "" {
				skippedSubdivisions[subIsoIds[idx]] = struct{}{}
			}
		}
	}

	var processedAdministrativeAreas administrativeAreaSlice
	var latinizedAdministrativeAreas administrativeAreaSlice
	var ids []string

	// Deal with the case where a country has subkeys, but the list
	// of ISO ids is blank (eg: ES)
	if ec.SubISOIDs != "" {
		ids = subIsoIds
	} else if ec.SubKeys != "" {
		ids = subKeys
	}

	for i, isoID := range ids {
		// Skip administrative areas without ISO Ids due to regions being
		// contested or not recognized (e.g Crimea and Sevastopol in Russia)
		if isoID == "" {
			continue
		}

		if _, ok := skippedSubdivisions[isoID]; ok {
			continue
		}

		// Sanity check
		if ec.SubZips != "" && ec.SubZipExs != "" && subZips[i] != "" && subZipExs[i] != "" {
			regexp := "^" + subZips[i]
			pcs := strings.Split(subZipExs[i], ",")
			if e := checkPostalCodeRegex(regexp, pcs); e != nil {
				err := fmt.Errorf(
					"error checking administrative area post "+
						"code regex for %s / %s against sample: %s",
					isoID, ec.Key, e.Error(),
				)

				return administrativeAreaMap, postCodeRegexMap, err
			}
		}

		adminArea := administrativeArea{
			ID:        isoID,
			PostalKey: subKeys[i],
		}

		if ec.SubNames != "" {
			adminArea.Name = subNames[i]
		} else {
			adminArea.Name = subKeys[i]
		}

		if ec.SubZips != "" && subZips[i] != "" {
			postCodeRegexMap[isoID] = postCodeRegex{
				regex: fmt.Sprintf("^%s", subZips[i]),
			}
		}

		var latinizedAdministrativeArea administrativeArea
		if ec.SubLNames != "" {
			latinizedAdministrativeArea.ID = isoID
			latinizedAdministrativeArea.Name = subLNames[i]
			latinizedAdministrativeArea.PostalKey = subKeys[i]
		}

		if ec.SubMores != "" && subMores[i] == "true" {
			sanitized := REMOVE_LANG_REGEX.ReplaceAllString(ec.ID, "")
			url := fmt.Sprintf("%s/%s", GOOGLE_ADDRESS_URL, sanitized)
			if lang != "" {
				url += fmt.Sprintf("--%s", lang)
			}

			data, e := http.Get(url)
			if e != nil {
				err := fmt.Errorf(
					"error fetching administrative area"+
						"data for %s: %s",
					url, e.Error(),
				)
				return administrativeAreaMap, postCodeRegexMap, err
			}

			esd, e := decodeSubdivision(data.Body)
			if e != nil {
				err := fmt.Errorf(
					"error decoding subdivision" +
						"for %s/%s: %s",
					ec.Key, subKeys[i], e.Error(),
				)

				return administrativeAreaMap, postCodeRegexMap, err
			}

			lm, pcrm, e := processLocalities(esd, lang)
			if e != nil {
				err := fmt.Errorf(
					"error processing localities "+
						"for %s/%s: %s",
					ec.Key, subKeys[i], e.Error(),
				)

				return administrativeAreaMap, postCodeRegexMap, err
			}

			// Sanity check
			if _, ok := postCodeRegexMap[isoID]; !ok && len(pcrm) > 0 {
				err := fmt.Errorf(
					"locality %s has postcode regexes"+
						"but the parent locality does not",
					ec.ID,
				)

				return administrativeAreaMap, postCodeRegexMap, err
			}

			if len(pcrm) > 0 {
				pcr := postCodeRegexMap[isoID]
				pcr.subdivisionRegex = pcrm
				postCodeRegexMap[isoID] = pcr
			}

			if len(lm[ec.Lang]) > 0 {
				adminArea.Localities = lm[ec.Lang]
			}

			// Consider latinized names to be english
			if ec.SubLNames != "" {
				// sanity check
				if _, ok := lm["en"]; !ok {
					err := fmt.Errorf(
						"%s has latinized admin areas but does not"+
							"have any latinized localities for %s",
						ec.Key, esd.ID,
					)

					return administrativeAreaMap, postCodeRegexMap, err
				}

				latinizedAdministrativeArea.Localities = lm["en"]
			}
		}

		processedAdministrativeAreas = append(processedAdministrativeAreas, adminArea)
		if latinizedAdministrativeArea.ID != "" {
			latinizedAdministrativeAreas = append(latinizedAdministrativeAreas, latinizedAdministrativeArea)
		}
	}

	administrativeAreaMap[ec.Lang] = processedAdministrativeAreas

	if len(latinizedAdministrativeAreas) > 0 {
		// sanity check
		if len(latinizedAdministrativeAreas) != len(processedAdministrativeAreas) {
			err := fmt.Errorf(
				"number of latinized areas (%d) does not"+
					"match number of admin areas (%d) for %s",
				len(latinizedAdministrativeAreas), len(processedAdministrativeAreas), ec.ID,
			)

			return administrativeAreaMap, postCodeRegexMap, err
		}

		sort.Slice(latinizedAdministrativeAreas, func(i, j int) bool {
			return latinizedAdministrativeAreas[i].Name < latinizedAdministrativeAreas[j].Name
		})

		administrativeAreaMap["en"] = latinizedAdministrativeAreas
	}

	return administrativeAreaMap, postCodeRegexMap, nil
}

func processLocalities(esd externalSubdivision, lang string) (
	map[string]localitySlice,
	map[string]postCodeRegex,
	error,
) {
	localityMap := make(map[string]localitySlice)
	postCodeRegexMap := make(map[string]postCodeRegex)

	subKeys := strings.Split(esd.SubKeys, "~")
	subNames := strings.Split(esd.SubNames, "~")
	subMores := strings.Split(esd.SubMores, "~")
	subZips := strings.Split(esd.SubZips, "~")
	subZipExs := strings.Split(esd.SubZipExs, "~")
	subLNames := strings.Split(esd.SubLNames, "~")

	var latinizedLocalities localitySlice
	var processedLocalities localitySlice

	for i, key := range subKeys {
		// Sanity check
		if esd.SubZips != "" && esd.SubZipExs != "" && subZips[i] != "" && subZipExs[i] != "" {
			regexp := "^" + subZips[i]
			pcs := strings.Split(subZipExs[i], ",")
			if e := checkPostalCodeRegex(regexp, pcs); e != nil {
				err := fmt.Errorf(
					"error checking default locality post "+
						"code regex for %s against sample: %s",
					esd.ID, e.Error(),
				)

				return localityMap, postCodeRegexMap, err
			}
		}
		// dependent locality
		dl := locality{
			// No ISO ID at this level,
			// so we use the key from
			// Google's data set.
			ID: key,
		}

		if esd.SubNames != "" {
			dl.Name = subNames[i]
		} else {
			dl.Name = subKeys[i]
		}

		if esd.SubZips != "" && subZips[i] != "" {
			postCodeRegexMap[key] = postCodeRegex{
				regex: fmt.Sprintf("^%s", subZips[i]),
			}
		}
		// latinized locality
		var ll locality
		if esd.SubLNames != "" {
			ll.ID = key
			ll.Name = subLNames[i]
		}

		if esd.SubMores != "" && subMores[i] == "true" {
			sanitized := REMOVE_LANG_REGEX.ReplaceAllString(esd.ID, "")
			url := fmt.Sprintf("%s/%s", GOOGLE_ADDRESS_URL, sanitized)
			if lang != "" {
				url += fmt.Sprintf("--%s", lang)
			}

			data, e := http.Get(url)
			if e != nil {
				err := fmt.Errorf(
					"error fetching default locality"+
						"data for %s: %s",
					url, e.Error(),
				)

				return localityMap, postCodeRegexMap, err
			}

			externalLocality, e := decodeSubdivision(data.Body)
			if e != nil {
				err := fmt.Errorf(
					"error unmarhaling data for url(%s): %s",
					url, e.Error(),
				)

				return localityMap, postCodeRegexMap, err
			}

			dlm, pcrm, err := processDependentLocalities(externalLocality)
			if err != nil {
				err := fmt.Errorf(
					"error processing dependent localities for %s/%s: %s",
					esd.ID, subKeys[i], e.Error(),
				)

				return localityMap, postCodeRegexMap, err
			}

			if _, ok := postCodeRegexMap[key]; !ok && len(pcrm) > 0 {
				err := fmt.Errorf(
					"dependent locality %s/%s has postcode regexes"+
						"but the parent locality does not",
					esd.ID, subKeys[i],
				)

				return localityMap, postCodeRegexMap, err
			}

			if len(pcrm) > 0 {
				postCodeRegex := postCodeRegexMap[key]
				postCodeRegex.subdivisionRegex = pcrm
				postCodeRegexMap[key] = postCodeRegex
			}

			if len(dlm[esd.Lang]) > 0 {
				dl.DependentLocalities = dlm[esd.Lang]
			}

			// Consider latinized names to be english
			if esd.SubLNames != "" {
				if _, ok := dlm["en"]; !ok {
					err := fmt.Errorf(
						"%s has latinized localities but does not "+
							"have any latinized dependent localities for %s",
						esd.ID, subKeys[i],
					)

					return localityMap, postCodeRegexMap, err
				}

				ll.DependentLocalities = dlm["en"]
			}
		}

		processedLocalities = append(processedLocalities, dl)
		if esd.SubLNames != "" {
			latinizedLocalities = append(latinizedLocalities, ll)
		}
	}

	localityMap[esd.Lang] = processedLocalities

	if len(latinizedLocalities) > 0 {
		// sanity check
		if len(latinizedLocalities) != len(processedLocalities) {
			err := fmt.Errorf(
				"number of latinized localities (%d) does not"+
					"match number of localities (%d) for %s",
				len(latinizedLocalities), len(processedLocalities), esd.ID,
			)

			return localityMap, postCodeRegexMap, err
		}

		sort.Slice(latinizedLocalities, func(i, j int) bool {
			return latinizedLocalities[i].Name < latinizedLocalities[j].Name
		})

		localityMap["en"] = latinizedLocalities
	}

	return localityMap, postCodeRegexMap, nil
}

func processDependentLocalities(esd externalSubdivision) (
	map[string]dependentLocalitySlice,
	map[string]postCodeRegex,
	error,
) {
	dependentLocalityMap := make(map[string]dependentLocalitySlice)
	postCodeRegexMap := make(map[string]postCodeRegex)

	subKeys := strings.Split(esd.SubKeys, "~")
	subNames := strings.Split(esd.SubNames, "~")
	subZips := strings.Split(esd.SubZips, "~")
	subZipExs := strings.Split(esd.SubZipExs, "~")

	var processedDependentLocalities dependentLocalitySlice
	for i, key := range subKeys {
		// sanity check
		if esd.SubZips != "" && esd.SubZipExs != "" && subZips[i] != "" && subZipExs[i] != "" {
			regexp := "^" + subZips[i]
			pcs := strings.Split(subZipExs[i], ",")
			if e := checkPostalCodeRegex(regexp, pcs); e != nil {
				err := fmt.Errorf(
					"error checking dependent locality post "+
						"code regex for %s against sample: %s",
					esd.Key, e.Error(),
				)

				return dependentLocalityMap, postCodeRegexMap, err
			}
		}

		dl := dependentLocality{
			// No ISO ID at this level,
			// so we use the key from
			// Google's data set.
			ID: key,
		}

		if esd.SubNames != "" {
			dl.Name = subNames[i]
		} else {
			dl.Name = subKeys[i]
		}

		if esd.SubZips != "" && subZips[i] != "" {
			postCodeRegexMap[key] = postCodeRegex{
				regex: fmt.Sprintf("^%s", subZips[i]),
			}
		}

		processedDependentLocalities = append(processedDependentLocalities, dl)
	}

	dependentLocalityMap[esd.Lang] = processedDependentLocalities

	// We consider latinized names to be english
	if esd.SubLNames != "" {
		subLNames := strings.Split(esd.SubLNames, "~")

		var ldl dependentLocalitySlice
		for i, key := range subKeys {
			dl := dependentLocality{
				// No ISO ID at this level
				// so we use the key from
				// Google's data set.
				ID:   key,
				Name: subLNames[i],
			}

			ldl = append(ldl, dl)
		}

		sort.Slice(ldl, func(i, j int) bool {
			return ldl[i].Name < ldl[j].Name
		})

		dependentLocalityMap["en"] = ldl
	}

	return dependentLocalityMap, postCodeRegexMap, nil
}

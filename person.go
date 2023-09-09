package main

import (
	"errors"
	"slices"
)

type Person struct {
	uid                  string
	nationalityPrimary   string
	nationalitySecondary string
	ethnicValue          int
}

func (person Person) GetEthnic(nationToEthnic map[string]string) (string, error) {
	primaryEthnic, hasPrimaryEthnic := nationToEthnic[person.nationalityPrimary]
	if !hasPrimaryEthnic {
		return "", errors.New("no primary nationality")
	}
	secondaryEthnic := nationToEthnic[person.nationalitySecondary]
	nationalityEthnics := []string{primaryEthnic, secondaryEthnic}

	containsEthnic := func(ethnic string) bool {
		return slices.Contains(nationalityEthnics, ethnic)
	}

	if person.ethnicValue > 10 {
		return "", errors.New("ethnic value out of bounds")
	}

	switch person.ethnicValue {
	case 1:
		if containsEthnic(Scandinavian) ||
			containsEthnic(SouthEastAsian) ||
			containsEthnic(CentralEuropean) ||
			containsEthnic(Caucasion) ||
			containsEthnic(African) ||
			containsEthnic(Asian) ||
			containsEthnic(MiddleEastNorthAfrican) ||
			containsEthnic(MiddleEastSouthAmerican) ||
			containsEthnic(EasternEuropeanCentralAsian) {
			return SouthAmerican, nil
		}
		if primaryEthnic != "" {
			return primaryEthnic, nil
		}
		return secondaryEthnic, nil
	case 3, 6, 7, 8, 9:
		if person.ethnicValue == 7 {
			if primaryEthnic == SouthAmericanMediterranean {
				return SouthAmericanMediterranean, nil
			}
			if primaryEthnic == SouthAmerican {
				return SouthAmerican, nil
			}
		}
		return African, nil
	case 10:
		if primaryEthnic == SouthAmerican {
			return SouthAmerican, nil
		}
		return Asian, nil
	case 2:
		if containsEthnic(MiddleEastSouthAmerican) {
			return MiddleEastSouthAmerican, nil
		}
		return MiddleEastNorthAfrican, nil
	case 5:
		return SouthEastAsian, nil
	case 0:
		if containsEthnic(Scandinavian) {
			return Scandinavian, nil
		}
		if containsEthnic(Caucasion) {
			return Caucasion, nil
		}
		return CentralEuropean, nil
	case 4:
		return MiddleEastSouthAmerican, nil
	default:
		return "", errors.New("ethnic not found")
	}
}

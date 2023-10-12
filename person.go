package main

import (
	"errors"
)

type Person struct {
	uid             string
	ethnicPrimary   string
	ethnicSecondary string
	ethnicValue     int
}

func (person Person) GetEthnic() (string, error) {
	hasEthnic := func(ethnic string) bool {
		return person.ethnicPrimary == ethnic || person.ethnicSecondary == ethnic
	}

	switch person.ethnicValue {
	case 0:
		if hasEthnic(Scandinavian) {
			return Scandinavian, nil
		}
		if hasEthnic(Caucasion) {
			return Caucasion, nil
		}
		return CentralEuropean, nil
	case 1:
		if hasEthnic(Scandinavian) ||
			hasEthnic(SouthEastAsian) ||
			hasEthnic(CentralEuropean) ||
			hasEthnic(Caucasion) ||
			hasEthnic(African) ||
			hasEthnic(Asian) ||
			hasEthnic(MiddleEastNorthAfrican) ||
			hasEthnic(MiddleEastSouthAmerican) ||
			hasEthnic(EasternEuropeanCentralAsian) {
			return SouthAmerican, nil
		}
		if person.ethnicPrimary != "" {
			return person.ethnicPrimary, nil
		}
		return person.ethnicSecondary, nil
	case 2:
		if hasEthnic(MiddleEastSouthAmerican) {
			return MiddleEastSouthAmerican, nil
		}
		return MiddleEastNorthAfrican, nil
	case 3, 6, 7, 8, 9:
		if person.ethnicValue == 7 {
			if person.ethnicPrimary == SouthAmericanMediterranean {
				return SouthAmericanMediterranean, nil
			}
			if person.ethnicPrimary == SouthAmerican {
				return SouthAmerican, nil
			}
		}
		return African, nil
	case 4:
		return MiddleEastSouthAmerican, nil
	case 5:
		return SouthEastAsian, nil
	case 10:
		if person.ethnicPrimary == SouthAmerican {
			return SouthAmerican, nil
		}
		return Asian, nil
	default:
		return "", errors.New("ethnic not found")
	}
}

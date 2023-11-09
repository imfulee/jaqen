package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ErrBadRTFFormat error = errors.New("bad RTF Format")

type Person struct {
	id     string
	ethnic string
}

type INationEthnicMapper interface {
	Map(string) (string, bool)
}

type RTF struct {
	nationEthnicMapper INationEthnicMapper
}

func (r RTF) getEthnic(nationality1, nationality2 string, ethnicValue int) (Ethnic, error) {
	ethnic1, ok := r.nationEthnicMapper.Map(nationality1)
	if !ok {
		return "", errors.New("ethnic not found")
	}

	ethnic2, _ := r.nationEthnicMapper.Map(nationality2)

	hasEthnic := func(ethnic string) bool {
		return ethnic1 == ethnic || ethnic2 == ethnic
	}

	switch ethnicValue {
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
		if ethnic1 != "" {
			return ethnic1, nil
		}
		return ethnic2, nil
	case 2:
		if hasEthnic(MiddleEastSouthAmerican) {
			return MiddleEastSouthAmerican, nil
		}
		return MiddleEastNorthAfrican, nil
	case 3, 6, 7, 8, 9:
		if ethnicValue == 7 {
			if ethnic1 == SouthAmericanMediterranean {
				return SouthAmericanMediterranean, nil
			}
			if ethnic1 == SouthAmerican {
				return SouthAmerican, nil
			}
		}
		return African, nil
	case 4:
		return MiddleEastSouthAmerican, nil
	case 5:
		return SouthEastAsian, nil
	case 10:
		if ethnic1 == SouthAmerican {
			return SouthAmerican, nil
		}
		return Asian, nil
	default:
		return "", errors.New("ethnic not found")
	}
}

func (r RTF) Parse(path string) ([]Person, error) {
	result := []Person{}

	rtfFile, rtfErr := os.Open(path)
	if rtfErr != nil {
		return nil, rtfErr
	}
	defer rtfFile.Close()

	UIDRegex, err := regexp.Compile("([0-9]){7,}")
	if err != nil {
		return nil, err
	}

	rtfScanner := bufio.NewScanner(rtfFile)
	for rtfScanner.Scan() {
		rtfLine := rtfScanner.Text()
		uidByte := UIDRegex.Find([]byte(rtfLine))
		if uidByte != nil {
			id := string(uidByte)

			rtfData := strings.Split(rtfLine, "|")
			if len(rtfData) < 8 {
				return nil, errors.Join(ErrBadRTFFormat, errors.New("not enough lines in RTF line"))
			}

			ethnicValue, ethniceValueErr := strconv.Atoi(strings.Trim(rtfData[7], " "))
			if ethniceValueErr != nil {
				return nil, ethniceValueErr
			}
			if ethnicValue > 10 {
				return nil, errors.Join(ErrBadRTFFormat, errors.New("ethnic value out of bounds"))
			}

			nationality1 := strings.Trim(rtfData[2], " ")
			nationality2 := strings.Trim(rtfData[3], " ")

			ethnic, error := r.getEthnic(nationality1, nationality2, ethnicValue)
			if error != nil {
				return nil, errors.Join(ErrBadRTFFormat,
					errors.Join(
						errors.New("could not find ethnic"),
						error,
					))
			}

			result = append(result, Person{
				id:     id,
				ethnic: ethnic,
			})
		}
	}

	if rtfScannerErr := rtfScanner.Err(); rtfScannerErr != nil {
		return nil, rtfScannerErr
	}

	return result, nil
}

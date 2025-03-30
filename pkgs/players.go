package mapper

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ErrBadRTFFormat string = "bad RTF Format:\n%w"

// change fileOpener to allow overriding in tests
var fileOpener = os.Open

func getEthnic(nationality1, nationality2 string, ethnicValue int) (Ethnic, error) {
	ethnic1, ok := NationEthnicMapping[nationality1]
	if !ok {
		return "", fmt.Errorf("ethnic not found for country initials: %s", nationality1)
	}

	ethnic2 := NationEthnicMapping[nationality2]

	hasEthnic := func(ethnic Ethnic) bool {
		return ethnic1 == ethnic || ethnic2 == ethnic
	}

	switch ethnicValue {
	case 0:
		if hasEthnic(Scandinavian) {
			return Scandinavian, nil
		}
		if hasEthnic(Caucasian) {
			return Caucasian, nil
		}
		return CentralEuropean, nil
	case 1:
		if hasEthnic(Scandinavian) ||
			hasEthnic(SouthEastAsian) ||
			hasEthnic(CentralEuropean) ||
			hasEthnic(Caucasian) ||
			hasEthnic(African) ||
			hasEthnic(Asian) ||
			hasEthnic(MiddleEastNorthAfrican) ||
			hasEthnic(MiddleEastSouthAsian) ||
			hasEthnic(EasternEuropeanCentralAsian) {
			return SouthAmerican, nil
		}
		if ethnic1 != "" {
			return ethnic1, nil
		}
		return ethnic2, nil // not sure how to test. guarded by statement above and ethnic not found for country error
	case 2:
		if hasEthnic(MiddleEastSouthAsian) {
			return MiddleEastSouthAsian, nil
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
		return MiddleEastSouthAsian, nil
	case 5:
		return SouthEastAsian, nil
	case 10:
		if ethnic1 == SouthAmerican {
			return SouthAmerican, nil
		}
		return Asian, nil
	default:
		return "", fmt.Errorf("ethnic value not found: %d", ethnicValue)
	}
}

func GetPlayers(rtfPath string) ([]Player, error) {
	players := make([]Player, 0)

	rtfFile, rtfErr := fileOpener(rtfPath)
	if rtfErr != nil {
		return nil, rtfErr
	}
	defer rtfFile.Close()

	UIDRegex := regexp.MustCompile("([0-9]){7,}")

	getEthnicErrors := make([]error, 0)

	rtfScanner := bufio.NewScanner(rtfFile)
	for rtfScanner.Scan() {
		rtfLine := rtfScanner.Text()
		uidByte := UIDRegex.Find([]byte(rtfLine))

		if uidByte != nil {
			id := string(uidByte)

			rtfData := strings.Split(rtfLine, "|")
			if len(rtfData) < 8 {
				return nil, fmt.Errorf(ErrBadRTFFormat, fmt.Errorf("not enough lines in RTF line: %s", rtfLine))
			}

			for rtfDataIndex := range rtfData {
				rtfData[rtfDataIndex] = strings.Trim(rtfData[rtfDataIndex], " ")
			}

			ethnicValue, ethniceValueErr := strconv.Atoi(rtfData[7])
			if ethniceValueErr != nil {
				return nil, ethniceValueErr
			}

			nationality1 := rtfData[2]
			nationality2 := rtfData[3]

			ethnic, err := getEthnic(nationality1, nationality2, ethnicValue)
			if err != nil {
				getEthnicErrors = append(getEthnicErrors, err)
				continue
			}

			players = append(players, Player{
				ID:     PlayerID(id),
				Ethnic: ethnic,
			})
		}
	}

	if len(getEthnicErrors) > 0 {
		return nil, fmt.Errorf(ErrBadRTFFormat, errors.Join(getEthnicErrors...))
	}

	if rtfScannerErr := rtfScanner.Err(); rtfScannerErr != nil {
		return nil, rtfScannerErr
	}

	return players, nil
}

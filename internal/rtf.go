package mapper

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ErrBadRTFFormat error = errors.New("bad RTF Format")

func getEthnic(nationality1, nationality2 string, ethnicValue int) (Ethnic, error) {
	ethnic1, ok := NationEthnicMapping[nationality1]
	if !ok {
		return "", errors.New("ethnic not found")
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

func GetPlayersBuilder(path string) func() ([]Player, error) {
	return func() ([]Player, error) {
		players := make([]Player, 0)

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

				for rtfDataIndex := range rtfData {
					rtfData[rtfDataIndex] = strings.Trim(rtfData[rtfDataIndex], " ")
				}

				ethnicValue, ethniceValueErr := strconv.Atoi(rtfData[7])
				if ethniceValueErr != nil {
					return nil, ethniceValueErr
				}
				if ethnicValue > 10 {
					return nil, errors.Join(ErrBadRTFFormat, errors.New("ethnic value out of bounds"))
				}

				nationality1 := rtfData[2]
				nationality2 := rtfData[3]

				ethnic, err := getEthnic(nationality1, nationality2, ethnicValue)
				if err != nil {
					return nil, errors.Join(ErrBadRTFFormat, err)
				}

				players = append(players, Player{
					id:     PlayerID(id),
					ethnic: ethnic,
				})
			}
		}

		if rtfScannerErr := rtfScanner.Err(); rtfScannerErr != nil {
			return nil, rtfScannerErr
		}

		return players, nil
	}
}

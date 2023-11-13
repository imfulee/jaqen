package mapper

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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
	NationEthnicMapper INationEthnicMapper
}

func (r RTF) getEthnic(nationality1, nationality2 string, ethnicValue int) (Ethnic, error) {
	ethnic1, ok := r.NationEthnicMapper.Map(nationality1)
	if !ok {
		return "", errors.New("ethnic not found")
	}

	ethnic2, _ := r.NationEthnicMapper.Map(nationality2)

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

type Record struct {
	From string `xml:"from,attr"`
	To   string `xml:"to,attr"`
}

type XMLStruct struct {
	XMLName xml.Name `xml:"record"`
	Boolean []struct {
		ID    string `xml:"id,attr"`
		Value string `xml:"value,attr"`
	} `xml:"boolean"`
	List struct {
		ID     string   `xml:"id,attr"`
		Record []Record `xml:"record"`
	} `xml:"list"`
}

type XML struct {
	data *XMLStruct
}

func (x *XML) ReadXML(xmlPath string) error {
	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		return errors.Join(errors.New("cannot open xml file"), err)
	}

	xmlBytes, err := io.ReadAll(xmlFile)
	if err != nil {
		return errors.Join(errors.New("cannot read all xml data"), err)
	}
	defer xmlFile.Close()

	if err := xml.Unmarshal(xmlBytes, &x.data); err != nil {
		return errors.Join(errors.New("cannot unmarshall xml file"), err)
	}

	return nil
}

func (x *XML) GetPreviousMappings() (map[string]PersonMap, error) {
	idRegex, err := regexp.Compile(`\d`)
	if err != nil {
		return nil, errors.Join(errors.New("cannot compile regex"), err)
	}

	mappings := make(map[string]PersonMap)
	for _, record := range x.data.List.Record {
		pm := PersonMap{
			FromPath: record.From,
			ToPath:   record.To,
		}

		idArray := idRegex.FindStringSubmatch(pm.FromPath)
		if idArray == nil {
			return nil, fmt.Errorf("cannot find id for %s", pm.FromPath)
		}
		id := idArray[0]

		mappings[id] = pm
	}

	return mappings, nil
}

func (x *XML) UpdateMappings(personMaps map[string]PersonMap) {
	x.data.List.Record = make([]Record, 0)

	for _, personMap := range personMaps {
		x.data.List.Record = append(x.data.List.Record, Record{
			From: personMap.FromPath,
			To:   personMap.ToPath,
		})
	}
}

func (x *XML) WriteXML(xmlPath string) error {
	xmlString, err := xml.MarshalIndent(&x.data, "  ", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(xmlPath, xmlString, 0777)
	if err != nil {
		return err
	}

	return nil
}

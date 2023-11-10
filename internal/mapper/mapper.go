package mapper

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

type PersonMap struct {
	fromPath string
	toPath   string
}

type IImages interface {
	Random(string) (string, error)
}

type Mapper struct {
	Preserve bool
	Images   IImages
	mappings map[string]PersonMap
	xmlData  *XMLStruct
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

func (m *Mapper) ReadXML(xmlPath string) error {
	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		return errors.Join(errors.New("cannot open xml file"), err)
	}
	xmlData, err := io.ReadAll(xmlFile)
	if err != nil {
		return errors.Join(errors.New("cannot read all xml data"), err)
	}
	defer xmlFile.Close()

	record := XMLStruct{}

	if err := xml.Unmarshal(xmlData, &record); err != nil {
		return errors.Join(errors.New("cannot unmarshall xml file"), err)
	}

	m.xmlData = &record

	idRegex, err := regexp.Compile(`\d`)
	if err != nil {
		return errors.Join(errors.New("cannot compile regex"), err)
	}

	for _, record := range record.List.Record {
		pm := PersonMap{
			fromPath: record.From,
			toPath:   record.To,
		}

		idArray := idRegex.FindStringSubmatch(pm.fromPath)
		if idArray == nil {
			return fmt.Errorf("cannot find id for %s", pm.fromPath)
		}
		id := idArray[0]

		m.mappings[id] = pm
	}

	return nil
}

func (m *Mapper) CreateMap(persons []Person) error {
	for _, person := range persons {
		if _, hasMapped := m.mappings[person.id]; hasMapped && m.Preserve {
			continue
		}

		fromImage, err := m.Images.Random(person.ethnic)
		if err != nil {
			return errors.Join(
				errors.New("could not get random image from ethnic"),
				err,
			)
		}

		m.mappings[person.id] = PersonMap{
			fromPath: PathForImage("", person.ethnic, fromImage),
			toPath:   PathForGame(person.id),
		}
	}

	m.xmlData.List.Record = make([]Record, 0)
	for _, personMap := range m.mappings {
		m.xmlData.List.Record = append(m.xmlData.List.Record, Record{
			From: personMap.fromPath,
			To:   personMap.toPath,
		})
	}

	return nil
}

func (m *Mapper) WriteXML(xmlPath string) error {
	xmlString, err := xml.MarshalIndent(&m.xmlData, "  ", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(xmlPath, xmlString, 0777)
	if err != nil {
		return err
	}

	return nil
}

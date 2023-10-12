package main

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"regexp"

	mapset "github.com/deckarep/golang-set"
)

type PersonMap struct {
	id        string
	imagePath string
}

type IImages interface {
	Random(string, bool) (string, error)
}

type Mapper struct {
	preserve bool
	mappings []PersonMap
	images   IImages
}

type XMLStruct struct {
	XMLName xml.Name `xml:"record"`
	Boolean []struct {
		ID    string `xml:"id,attr"`
		Value string `xml:"value,attr"`
	} `xml:"boolean"`
	List struct {
		ID     string `xml:"id,attr"`
		Record []struct {
			From string `xml:"from,attr"`
			To   string `xml:"to,attr"`
		} `xml:"record"`
	} `xml:"list"`
}

func (m *Mapper) ReadPreviousMap(xmlPath string) error {
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

	idRegex, err := regexp.Compile(`\d`)
	if err != nil {
		return errors.Join(errors.New("cannot compile regex"), err)
	}

	m.mappings = []PersonMap{}
	for _, record := range record.List.Record {
		m.mappings = append(m.mappings, PersonMap{
			id:        idRegex.FindStringSubmatch(record.To)[0],
			imagePath: record.From,
		})
	}

	return nil
}

func (m *Mapper) CreateMap(persons []Person) error {
	if m.preserve {
		m.mappings = []PersonMap{}
	}

	idSet := mapset.NewSet()
	for _, mapping := range m.mappings {
		idSet.Add(mapping.id)
	}

	for _, person := range persons {
		if idSet.Contains(person.uid) {
			continue
		}

		ethnic, ethnicError := person.GetEthnic()
		if ethnicError != nil {
			return errors.Join(errors.New("cannot get ethnic"), ethnicError)
		}

		randomImage, err := m.images.Random(ethnic, m.preserve)
		if err != nil {
			return errors.Join(errors.New("unable to get random image"), err)
		}

		m.mappings = append(m.mappings, PersonMap{
			id:        person.uid,
			imagePath: randomImage,
		})

		idSet.Add(person.uid)
	}

	return nil
}

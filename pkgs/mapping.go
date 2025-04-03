package mapper

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

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

type Mapping struct {
	instance   *XMLStruct
	idImageMap map[PlayerID]FilePath
	fmVersion  string
}

func convertToPathToPlayerID(toPath string, fmVersion string) PlayerID {
	var idRegex *regexp.Regexp
	switch fmVersion {
	case FMVersion2024:
		idRegex = regexp.MustCompile(`r-\d+`)
	default:
		idRegex = regexp.MustCompile(`\d+`)
	}

	playerID := PlayerID(idRegex.Find([]byte(toPath)))

	switch fmVersion {
	case FMVersion2024:
		playerID = PlayerID(strings.TrimPrefix(string(playerID), "r-"))
	default:
	}

	return playerID
}

func convertPlayerIDToToPath(id PlayerID) string {
	return fmt.Sprintf("graphics/pictures/person/%s/portrait", string(id))
}

func NewMapping(xmlPath string, fmVersion string) (*Mapping, error) {
	parser := &Mapping{
		instance:   nil,
		idImageMap: make(map[PlayerID]FilePath),
		fmVersion:  fmVersion,
	}

	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		return nil, errors.Join(errors.New("cannot open xml file"), err)
	}

	xmlBytes, err := io.ReadAll(xmlFile)
	if err != nil {
		return nil, errors.Join(errors.New("cannot read all xml data"), err)
	}
	defer xmlFile.Close()

	if err := xml.Unmarshal(xmlBytes, &parser.instance); err != nil {
		return nil, errors.Join(errors.New("cannot unmarshall xml file"), err)
	}

	for _, record := range parser.instance.List.Record {
		playerID := convertToPathToPlayerID(record.To, parser.fmVersion)
		filepath := FilePath(record.From)
		parser.idImageMap[playerID] = filepath
	}

	return parser, nil
}

func (m *Mapping) AssignedImages() []FilePath {
	return MapValues(m.idImageMap)
}

func (m *Mapping) Exist(id PlayerID) bool {
	_, ok := m.idImageMap[id]
	return ok
}

func (m *Mapping) MapToImage(id PlayerID, filepath FilePath) {
	m.idImageMap[id] = filepath
}

func (m *Mapping) Save() error {
	if m.instance == nil {
		return errors.New("unintialised instance")
	}

	m.instance.List.Record = make([]Record, 0)

	for id, filename := range m.idImageMap {
		var playerID PlayerID
		switch m.fmVersion {
		case FMVersion2024:
			playerID = PlayerID(fmt.Sprintf("r-%s", id))
		default:
			playerID = id
		}

		m.instance.List.Record = append(m.instance.List.Record, Record{From: string(filename), To: convertPlayerIDToToPath(playerID)})
	}

	return nil
}

func (m *Mapping) Write(xmlPath string) error {
	rtnXML, err := xml.MarshalIndent(m.instance, "", "\t")
	if err != nil {
		return err
	}

	xmlFile, err := os.Create(xmlPath)
	if err != nil {
		return err
	}

	if _, err := xmlFile.Write(rtnXML); err != nil {
		return err
	}

	return nil
}

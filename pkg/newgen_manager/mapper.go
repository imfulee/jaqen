package newgen_manager

import (
	"errors"
)

type PersonMap struct {
	FromPath string
	ToPath   string
}

type IImages interface {
	Random(string) (string, error)
}

// feels weird that I have to re insert the root path in mapper again
type Mapper struct {
	Preserve bool
	Images   IImages
	Mappings map[string]PersonMap
}

func (m *Mapper) Map(persons []Person) error {
	for _, person := range persons {
		if _, hasMapped := m.Mappings[person.id]; hasMapped && m.Preserve {
			continue
		}

		fromImage, err := m.Images.Random(person.ethnic)
		if err != nil {
			return errors.Join(
				errors.New("could not get random image from ethnic"),
				err,
			)
		}

		m.Mappings[person.id] = PersonMap{
			FromPath: PathForImage(person.ethnic, fromImage),
			ToPath:   PathForGame(person.id),
		}
	}

	return nil
}

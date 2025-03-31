package mapper

import (
	"errors"
	"fmt"
)

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func OverrideNationEthnicMapping(overrides map[string]string) error {
	overrideErrors := []error{}

	for nation, ethnic := range overrides {
		if !IsValidEthnic(ethnic) {
			overrideErrors = append(overrideErrors, fmt.Errorf(`ethnic value "%s" is not valid ethnic for "%s"`, ethnic, nation))
		} else {
			NationEthnicMapping[nation] = Ethnic(ethnic)
		}
	}

	if len(overrideErrors) > 0 {
		return errors.Join(overrideErrors...)
	}

	return nil
}

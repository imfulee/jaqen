package main

import (
	"encoding/json"
	"io"
	"os"
)

func NationToEthnicReader(ethnicJsonPath string) (map[string]string, error) {
	jsonFile, jsonFileErr := os.Open(ethnicJsonPath)
	if jsonFileErr != nil {
		return nil, jsonFileErr
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	nationToEthnicMap := map[string]string{}
	json.Unmarshal([]byte(byteValue), &nationToEthnicMap)

	return nationToEthnicMap, nil
}

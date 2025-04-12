package internal

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func ParseConfig(configBytes []byte) (JaqenConfig, error) {
	var config JaqenConfig

	err := toml.Unmarshal(configBytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func ReadConfig(filePath string) (JaqenConfig, error) {
	var config JaqenConfig

	if _, err := os.Stat(filePath); err != nil {
		return config, fmt.Errorf("could not find file: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return config, err
	}

	config, err = ParseConfig(bytes)

	return config, err
}

func MarshalConfig(config JaqenConfig) ([]byte, error) {
	bytes, err := toml.Marshal(config)
	if err != nil {
		return []byte{}, err
	}

	return bytes, err
}

func WriteConfig(config JaqenConfig, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	marshalledConfig, err := MarshalConfig(config)
	if err != nil {
		return err
	}

	if _, err := file.Write(marshalledConfig); err != nil {
		return err
	}

	return nil
}

func FindFmPaths() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fmPath := ""

	for fmVersion, fmSteamId := range SteamAppIdMap {
		possibleFmPaths := map[string]string{
			"arch":   fmt.Sprintf("%s/.local/shareSteam/steamapps/compatdata/%s/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager %s/", homePath, fmSteamId, fmVersion),
			"debian": fmt.Sprintf("%s/.steam/debian-installation/steamapps/compatdata/%s/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager %s", homePath, fmSteamId, fmVersion),
		}

		for _, possibleFmPath := range possibleFmPaths {
			if _, err := os.Stat(possibleFmPath); err != nil {
				fmPath = possibleFmPath
			}
		}
	}

	if fmPath == "" {
		return fmPath, errors.New("could not find steam football manager path")
	}

	return fmPath, nil
}

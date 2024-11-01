package internal

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

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
	err = toml.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

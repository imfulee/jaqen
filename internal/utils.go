package internal

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// openFileFunc opens files when reading configuration. Defined as a variable to allow overriding during tests.
var openFileFunc = func(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

func ReadConfig(filePath string) (JaqenConfig, error) {
	var config JaqenConfig

	if _, err := os.Stat(filePath); err != nil {
		return config, fmt.Errorf("could not find file: %w", err)
	}

	file, err := openFileFunc(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return config, fmt.Errorf("failed to read file %q: %w", filePath, err)
	}
	err = toml.Unmarshal(bytes, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal file %q: %w", filePath, err)
	}

	return config, nil
}

// tomlMarshal as a variable to allow overriding during tests.
var tomlMarshal = toml.Marshal

// fileCreatorFunc creates a new file for writing config. Defined as a variable to allow overriding during tests.
var (
	fileCreatorFunc = func(filename string) (io.WriteCloser, error) {
		return os.Create(filename)
	}
)

func WriteConfig(config JaqenConfig, filePath string) error {
	bytes, err := tomlMarshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	file, err := fileCreatorFunc(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", filePath, err)
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return fmt.Errorf("failed to write configuration to file %q: %w", filePath, err)
	}

	return nil
}

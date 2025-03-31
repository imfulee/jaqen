package internal

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func ptrBool(v bool) *bool {
	return &v
}

func ptrString(v string) *string {
	return &v
}

func TestReadConfig(t *testing.T) {
	testCases := []struct {
		name        string
		filePath    string
		expected    JaqenConfig
		expectError bool
	}{
		{
			name:     "Valid config file",
			filePath: "../test/testdata/valid_config.toml",
			expected: JaqenConfig{
				Preserve:        ptrBool(true),
				XMLPath:         ptrString("config.xml"),
				RTFPath:         ptrString("newgen.rtf"),
				IMGPath:         ptrString("images"),
				FMVersion:       ptrString("2024"),
				AllowDuplicate:  ptrBool(false),
				MappingOverride: &map[string]string{"USA": "African", "IND": "Caucasian"},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config, err := ReadConfig(tc.filePath)
			if tc.expectError && err == nil {
				t.Errorf("Expected an error, but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(config, tc.expected) {
				t.Errorf("Expected config: %+v, but got: %+v", tc.expected, config)
			}
		})
	}
}

func TestReadConfigInvalidTOML(t *testing.T) {
	filePath := "../test/testdata/invalid_config.toml"
	_, err := ReadConfig(filePath)
	var decodeErr *toml.DecodeError
	if !errors.As(err, &decodeErr) {
		t.Errorf("Expected DecodeError, got %v", err)
	}
}

func TestReadConfigFileDoesNotExist(t *testing.T) {
	filePath := "/path/to/nonexistent.toml"
	_, err := ReadConfig(filePath)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected ErrNotExist, got %v", err)
	}
}

func TestReadConfigFileOpenError(t *testing.T) {
	filePath := "../testdata/cannot_open_config.toml"
	_, err := ReadConfig(filePath)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Errorf("Expected a PathError, got: %v", err)
	}
}

type errorReadCloser struct{}

func (e *errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func (e *errorReadCloser) Close() error {
	return nil
}

// fakeOpenFile returns a fake io.ReadCloser that simulates a read error.
func fakeOpenFile(name string) (io.ReadCloser, error) {
	return &errorReadCloser{}, nil
}

func TestReadConfig_ReadAllError(t *testing.T) {
	origOpen := openFileFunc
	openFileFunc = fakeOpenFile
	defer func() { openFileFunc = origOpen }()

	// Create a temporary file so the os.Stat check passes.
	tmpFile := "test_dummy_config.toml"
	if err := os.WriteFile(tmpFile, []byte("key = \"value\""), 0644); err != nil {
		t.Fatalf("cannot write temporary file: %v", err)
	}
	defer os.Remove(tmpFile)

	_, err := ReadConfig(tmpFile)
	if err == nil {
		t.Fatal("expected error from ReadConfig, got nil")
	}
	if err.Error() != "failed to read file \"test_dummy_config.toml\": simulated read error" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWriteConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.toml")

	config := JaqenConfig{
		Preserve:        ptrBool(true),
		XMLPath:         ptrString("config.xml"),
		RTFPath:         ptrString("newgen.rtf"),
		IMGPath:         ptrString("images"),
		FMVersion:       ptrString("2024"),
		AllowDuplicate:  ptrBool(false),
		MappingOverride: &map[string]string{"USA": "African", "IND": "Caucasian"},
	}

	err := WriteConfig(config, configPath)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	writtenConfig, err := ReadConfig(configPath)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(writtenConfig, config) {
		t.Errorf("Expected written config to match original config")
	}
}

func TestWriteConfigMarshalError(t *testing.T) {
	origMarshal := tomlMarshal
	defer func() { tomlMarshal = origMarshal }()

	// Force the marshaler to fail.
	tomlMarshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("marshal error")
	}

	config := JaqenConfig{XMLPath: ptrString("config.xml")}
	err := WriteConfig(config, "dummy-path")
	if err == nil {
		t.Fatal("Expected an error from toml.Marshal, got nil")
	}
	if !strings.Contains(err.Error(), "marshal error") {
		t.Errorf("Expected marshal error, got: %v", err)
	}
}

func TestWriteConfigCreateError(t *testing.T) {
	origCreator := fileCreatorFunc
	defer func() { fileCreatorFunc = origCreator }()

	// Force fileCreator to fail.
	fileCreatorFunc = func(name string) (io.WriteCloser, error) {
		return nil, errors.New("create error")
	}

	config := JaqenConfig{XMLPath: ptrString("config.xml")}
	err := WriteConfig(config, "dummy-path")
	if err == nil {
		t.Fatal("Expected an error from os.Create, got nil")
	}
	if !strings.Contains(err.Error(), "create error") {
		t.Errorf("Expected create error, got: %v", err)
	}
}

// fakeWriteErrorFile simulates a file that is open (not closed) and returns an error when Write is called.
type fakeWriteErrorFile struct{}

func (f *fakeWriteErrorFile) Write(b []byte) (int, error) {
	return 0, errors.New("simulated write error")
}
func (f *fakeWriteErrorFile) Close() error {
	return nil
}

// fakeFileCreator returns our fake file that behaves as expected.
func fakeFileCreator(filename string) (io.WriteCloser, error) {
	return &fakeWriteErrorFile{}, nil
}

func TestWriteConfig_WriteError(t *testing.T) {
	origCreator := fileCreatorFunc
	fileCreatorFunc = fakeFileCreator
	defer func() { fileCreatorFunc = origCreator }()

	config := JaqenConfig{
		// populate fields if necessary
	}

	err := WriteConfig(config, "dummy_file.toml")
	if err == nil {
		t.Fatal("expected an error from WriteConfig, got nil")
	}
	if err.Error() != "failed to write configuration to file \"dummy_file.toml\": simulated write error" {
		t.Fatalf("got unexpected error: %v", err)
	}
}

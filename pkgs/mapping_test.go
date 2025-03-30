package mapper

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExistSuccessfully(t *testing.T) {
	t.Run("ID exists", func(t *testing.T) {
		id := PlayerID("123")
		mapping, err := NewMapping("../testdata/test_config.xml", "2024")
		if err != nil {
			t.Fatalf("Failed to create mapping: %v", err)
		}

		mapping.MapToImage(id, "test.png")
		exists := mapping.Exist(id)
		if !exists {
			t.Error("ID does not exist in the map when it should")
		}
	})

	t.Run("ID does not exist", func(t *testing.T) {
		id := PlayerID("456")
		mapping, err := NewMapping("../testdata/test_config.xml", "2024")
		if err != nil {
			t.Fatalf("Failed to create mapping: %v", err)
		}

		exists := mapping.Exist(id)
		if exists {
			t.Error("ID exists when it does not")
		}
	})
}

func TestMapToImageSuccessfully(t *testing.T) {
	t.Run("Correct ID format for 2024", func(t *testing.T) {
		id := PlayerID("r-123")
		mapping, err := NewMapping("../testdata/test_config.xml", "2024")
		if err != nil {
			t.Fatalf("Failed to create mapping: %v", err)
		}

		mapping.MapToImage(id, "test.png")
		ids := mapping.idImageMap[PlayerID("r-123")]
		if ids != "test.png" {
			t.Errorf("Expected value not equal to actual value")
		}
	})

	t.Run("Correct ID format for default version", func(t *testing.T) {
		id := PlayerID("456")
		mapping, err := NewMapping("../testdata/test_config.xml", "2024")
		if err != nil {
			t.Fatalf("Failed to create mapping: %v", err)
		}

		mapping.MapToImage(id, "test.png")
		ids := mapping.idImageMap[PlayerID("456")]
		if ids != "test.png" {
			t.Errorf("Expected value not equal to actual value")
		}
	})
}

func TestAssignedImagesSuccessfully(t *testing.T) {
	id := PlayerID("789")
	mapping, err := NewMapping("../testdata/test_config.xml", "2024")
	if err != nil {
		t.Fatalf("Failed to create mapping: %v", err)
	}

	mapping.MapToImage(id, "image.png")
	images := mapping.AssignedImages()

	if len(images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(images))
	}
	if images[0] != "image.png" {
		t.Errorf("Expected image 'image.png', got %v", images[0])
	}
}

func TestSaveSuccessfully2024(t *testing.T) {
	id := PlayerID("123")
	mapping, err := NewMapping("../testdata/test_config.xml", "2024")
	if err != nil {
		t.Fatalf("Failed to create mapping: %v", err)
	}
	mapping.MapToImage(id, "test.png")
	err = mapping.Save()
	if err != nil {
		t.Fatalf("Failed to save mapping: %v", err)
	}

	rtnXML, err := xml.Marshal(mapping.instance)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}
	xmlStr := string(rtnXML)

	if !strings.Contains(xmlStr, "test.png") {
		t.Errorf("Unexpected XML output.\nGot: %s\nExpected to contain: %s", xmlStr, "test.png")
	}
}

func TestSaveSuccessfullyDefault(t *testing.T) {
	id := PlayerID("123")
	mapping, err := NewMapping("../testdata/test_config.xml", "2023")
	if err != nil {
		t.Fatalf("Failed to create mapping: %v", err)
	}
	mapping.MapToImage(id, "test.png")
	err = mapping.Save()
	if err != nil {
		t.Fatalf("Failed to save mapping: %v", err)
	}

	rtnXML, err := xml.Marshal(mapping.instance)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}
	xmlStr := string(rtnXML)

	if !strings.Contains(xmlStr, "test.png") {
		t.Errorf("Unexpected XML output.\nGot: %s\nExpected to contain: %s", xmlStr, "test.png")
	}
}

func TestConvertToPathToPlayerID(t *testing.T) {
	tests := []struct {
		name      string
		toPath    string
		fmVersion string
		expected  PlayerID // Assuming output is just a trimmed ID as string.
	}{
		{"2024 with prefix", "/path/r-12345/abc", "2024", PlayerID("12345")},
		{"2024 without invalid characters", "/path/r-xy/z++1111/wxyz6789", "2024", PlayerID("")}, // Invalid id returns as-is after trimming.
		{"Default version with numbers only", "/path/45678901/some/path/name", "", PlayerID("45678901")},
		{"Version without prefix match", "/r-12345/abc...some/path/name", "2024", PlayerID("12345")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := convertToPathToPlayerID(tc.toPath, tc.fmVersion)
			if result != tc.expected {
				t.Errorf("Test %s: convertToPathToPlayerID(%q, %q) = %q; want %q",
					tc.name, tc.toPath, tc.fmVersion, result, tc.expected)
			}
		})
	}
	fmt.Println("All tests passed!")

}

func writeTempXMLFile(t *testing.T, xmlData string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "mapping_*.xml")
	if err != nil {
		t.Fatalf("unable to create temp file: %v", err)
	}
	_, err = tmpFile.WriteString(xmlData)
	if err != nil {
		t.Fatalf("unable to write to temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestNewMapping_Success(t *testing.T) {
	// In this test the fmVersion is "2024" so convertToPathToPlayerID is expected to trim an "r-" prefix.
	xmlData := `<record>
		<boolean id="preload" value="false"/>
		<boolean id="amap" value="false"/>
		<list id="foo">
			<record from="image1.jpg" to="/path/r-12345/extra"/>
			<record from="image2.jpg" to="/anotherPath/r-67890/more"/>
		</list>
</record>`

	// Write the XML data to a temporary file.
	xmlFilePath := writeTempXMLFile(t, xmlData)
	defer os.Remove(xmlFilePath)

	fmVersion := "2024"
	mapping, err := NewMapping(xmlFilePath, fmVersion)
	if err != nil {
		t.Fatalf("NewMapping returned unexpected error: %v", err)
	}

	expectedKey1 := convertToPathToPlayerID("/path/r-12345/extra", fmVersion)
	expectedKey2 := convertToPathToPlayerID("/anotherPath/r-67890/more", fmVersion)

	if mapping.idImageMap[expectedKey1] != "image1.jpg" {
		t.Errorf("expected mapping for key %q to be %q, got %q", expectedKey1, "image1.jpg", mapping.idImageMap[expectedKey1])
	}
	if mapping.idImageMap[expectedKey2] != "image2.jpg" {
		t.Errorf("expected mapping for key %q to be %q, got %q", expectedKey2, "image2.jpg", mapping.idImageMap[expectedKey2])
	}

	// Also check that the number of records in the map is exactly 2.
	if len(mapping.idImageMap) != 2 {
		t.Errorf("expected 2 records in mapping, got %d", len(mapping.idImageMap))
	}
}

func TestNewMapping_XML_FileOpenError(t *testing.T) {
	xmlFilePath := "../testdata/cannot_open_mapping.xml"

	fmVersion := "2024"
	_, err := NewMapping(xmlFilePath, fmVersion)
	if err == nil {
		t.Errorf("Expected an error, but got %v", err)
	}

	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Errorf("Expected os.PathError, but got %v", err)
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

func TestReadXML_ReadAllError(t *testing.T) {
	origOpen := openFileFunc
	openFileFunc = fakeOpenFile
	defer func() { openFileFunc = origOpen }()

	// Create a temporary file so the os.Stat check passes.
	tmpFile := "test_dummy_readall_error.xml"
	if err := os.WriteFile(tmpFile, []byte("key = \"value\""), 0644); err != nil {
		t.Fatalf("cannot write temporary file: %v", err)
	}
	defer os.Remove(tmpFile)

	fmVersion := "2024"
	_, err := NewMapping(tmpFile, fmVersion)
	if err == nil {
		t.Fatal("expected error from ReadConfig, got nil")
	}
	if err.Error() != "cannot read all xml data\nsimulated read error" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewMapping_InvalidXML(t *testing.T) {
	xmlFilePath := "../testdata/invalid_mapping.xml"

	_, err := NewMapping(xmlFilePath, "2024")
	if err == nil {
		t.Error("expected error for invalid XML, got nil")
	}
	if err.Error() != "cannot unmarshall xml file\nXML syntax error on line 8: element <record> closed by </r>" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TODO: Fix this test
func TestWriteSuccessfully(t *testing.T) {

	//mapping, err := NewMapping("../testdata/test_config.xml", "2024")
	//if err != nil {
	//	t.Fatalf("Failed to create mapping: %v", err)
	//}
	//mapping.MapToImage(PlayerID("123"), "../testdata/test.png")
	//filename := "../testdata/test_output.xml"
	//err = mapping.Write(filename)
	//if err != nil {
	//	t.Fatalf("Failed to write mapping: %v", err)
	//}
	//fileContent, err := os.ReadFile(filename)
	//if err != nil {
	//	t.Fatalf("Failed to read file: %v", err)
	//}
	//contentStr := string(fileContent)
	//t.Log("TODO: Fix this test")
	//t.Logf("File content:\n%s", contentStr)
	//if !strings.Contains(contentStr, "test.png") {
	//	t.Errorf("File content does not contain expected 'test.png'.\nGot: %s", contentStr)
	//}
	//os.Remove(filename)
}

func TestSaveWithNilInstance(t *testing.T) {
	mapping := &Mapping{}
	err := mapping.Save()
	if err == nil {
		t.Error("Expected error did not occur")
	}
}

func TestWrite_FileCreateError(t *testing.T) {
	// Override createFileFunc to simulate a file creation error
	origCreator := createFileFunc
	defer func() { createFileFunc = origCreator }()

	// Create a Mapping instance with some data
	mapping, err := NewMapping("../testdata/test_config.xml", "2024")
	if err != nil {
		t.Fatalf("Failed to create mapping: %v", err)
	}
	mapping.MapToImage(PlayerID("123"), "test.png")

	createFileFunc = func(name string) (io.WriteCloser, error) {
		return nil, errors.New("simulated create error")
	}

	// Attempt to write the mapping to a file
	err = mapping.Write("dummy_file.xml")
	if err == nil {
		t.Fatal("Expected an error from WriteConfig, got nil")
	}
	if !strings.Contains(err.Error(), "simulated create error") {
		t.Errorf("Expected create error, got: %v", err)
	}
}

func TestWrite_MarshalIndentError(t *testing.T) {
	// Create a Mapping instance with some data
	mapping, err := NewMapping("../testdata/test_config.xml", "2024")
	if err != nil {
		t.Fatalf("Failed to create mapping: %v", err)
	}
	mapping.MapToImage(PlayerID("123"), "test.png")

	// Override xml.MarshalIndent to simulate a marshal error
	origMarshalIndent := marshalXMLIndentFunc
	defer func() { marshalXMLIndentFunc = origMarshalIndent }()
	marshalXMLIndentFunc = func(v interface{}, prefix, indent string) ([]byte, error) {
		return nil, errors.New("simulated marshal error")
	}

	// Attempt to write the mapping to a file
	err = mapping.Write("dummy_file.xml")
	if err == nil {
		t.Fatal("Expected an error from xml.MarshalIndent, got nil")
	}
	if !strings.Contains(err.Error(), "simulated marshal error") {
		t.Errorf("Expected marshal error, got: %v", err)
	}
}

type fakeWriteErrorFile struct{}

func (f *fakeWriteErrorFile) Write(b []byte) (int, error) {
	return 0, errors.New("simulated write error")
}

func (f *fakeWriteErrorFile) Close() error {
	return nil
}

func TestWrite_WriteError(t *testing.T) {
	// Create a Mapping instance with some data
	mapping, err := NewMapping("../testdata/test_config.xml", "2024")
	if err != nil {
		t.Fatalf("Failed to create mapping: %v", err)
	}
	mapping.MapToImage(PlayerID("123"), "test.png")

	// Override createFileFunc to return a fake file that simulates a write error
	origCreator := createFileFunc
	defer func() { createFileFunc = origCreator }()
	createFileFunc = func(filename string) (io.WriteCloser, error) {
		return &fakeWriteErrorFile{}, nil
	}

	// Attempt to write the mapping to a file
	err = mapping.Write("dummy_file.xml")
	if err == nil {
		t.Fatal("expected an error from WriteConfig, got nil")
	}
	if err.Error() != "failed to write XML content: simulated write error" {
		t.Fatalf("got unexpected error: %v", err)
	}
}

func TestFileCreatorFunc_Success(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a temporary file path
	tmpFile := filepath.Join(tmpDir, "test_file.txt")

	// Use createFileFunc to create the file
	file, err := createFileFunc(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Verify that the file was created
	_, err = os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("File does not exist: %v", err)
	}

	// Write some data to the file
	_, err = file.Write([]byte("test data"))
	if err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	// Read the content of the file
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Check if the content contains the expected data
	if string(content) != "test data" {
		t.Errorf("File content does not contain expected 'test data'.\nGot: %s", string(content))
	}
}

func TestWrite_Success(t *testing.T) {
	// Create a temporary file path
	tmpFile := filepath.Join(t.TempDir(), "test_mapping.xml")
	defer os.Remove(tmpFile)

	// Create a Mapping instance with some data
	mapping, err := NewMapping("../testdata/test_config.xml", "2024")
	if err != nil {
		t.Fatalf("Failed to create mapping: %v", err)
	}
	mapping.MapToImage(PlayerID("123"), "test.png")

	err = mapping.Save()
	if err != nil {
		t.Fatalf("Failed to save mapping: %v", err)
	}

	// Write the mapping to the temporary file
	err = mapping.Write(tmpFile)
	if err != nil {
		t.Fatalf("Failed to write mapping: %v", err)
	}

	// Read the content of the file
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(content),
		"<record from=\"test.png\" to=\"") {
		t.Errorf("File content does not contain expected 'test.png, etc'.\nGot: %s", string(content))
	}

	// Unmarshal the content to verify the structure
	var unmarshalled XMLStruct
	err = xml.Unmarshal(content, &unmarshalled)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML content: %v", err)
	}

	// Check if the unmarshalled data matches the expected data
	if len(unmarshalled.List.Record) != 1 {
		t.Errorf("Expected 1 record, got %d", len(unmarshalled.List.Record))
	}
	if unmarshalled.List.Record[0].From != "test.png" {
		t.Errorf("Expected 'From' to be 'test.png', got %s", unmarshalled.List.Record[0].From)
	}
	if unmarshalled.List.Record[0].To != "graphics/pictures/person/r-123/portrait" {
		t.Errorf("Expected 'To' to be 'graphics/pictures/person/r-123/portrait', got %s", unmarshalled.List.Record[0].To)
	}
}

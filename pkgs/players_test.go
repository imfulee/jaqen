package mapper

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestGetEthnic_Parameterized(t *testing.T) {
	tests := []struct {
		name           string
		nationality1   string
		nationality2   string
		ethnicValue    int
		expectedEthnic Ethnic
	}{
		{
			name:           "Finnish Canadian with ethnic 0 is Scandinavian",
			nationality1:   "FIN",
			nationality2:   "CAN",
			ethnicValue:    0,
			expectedEthnic: Scandinavian,
		},
		{
			name:           "Bhutanese Canadian with ethnic 0 is Caucasian",
			nationality1:   "BHU",
			nationality2:   "CAN",
			ethnicValue:    0,
			expectedEthnic: Caucasian,
		},
		{
			name:           "Austrian with ethnic value 0 is Central European",
			nationality1:   "AUT",
			nationality2:   "",
			ethnicValue:    0,
			expectedEthnic: CentralEuropean,
		},
		{
			name:           "Albanian with ethnic 1 is YugoSlavGreek",
			nationality1:   "ALB",
			nationality2:   "",
			ethnicValue:    1,
			expectedEthnic: YugoslavGreek,
		},
		{
			name:           "Pakistan with ethnic 2 is MiddleEastSouthAsian",
			nationality1:   "PAK",
			nationality2:   "ITA",
			ethnicValue:    2,
			expectedEthnic: MiddleEastSouthAsian,
		},
		{
			name:           "English with ethnic 2 is MiddleEastNorthAfrican",
			nationality1:   "ENG",
			nationality2:   "ITA",
			ethnicValue:    2,
			expectedEthnic: MiddleEastNorthAfrican,
		},
		{
			name:           "English with ethnic 3 is African",
			nationality1:   "ENG",
			nationality2:   "",
			ethnicValue:    3,
			expectedEthnic: African,
		},
		{
			name:           "English with ethnic 6 is African",
			nationality1:   "ENG",
			nationality2:   "",
			ethnicValue:    6,
			expectedEthnic: African,
		},
		{
			name:           "English with ethnic 8 is African",
			nationality1:   "ENG",
			nationality2:   "",
			ethnicValue:    8,
			expectedEthnic: African,
		},
		{
			name:           "English with ethnic 9 is African",
			nationality1:   "ENG",
			nationality2:   "",
			ethnicValue:    9,
			expectedEthnic: African,
		},
		{
			name:           "Argentinian with ethnic 7 is SouthAmericanMediterranean",
			nationality1:   "ARG",
			nationality2:   "",
			ethnicValue:    7,
			expectedEthnic: SouthAmericanMediterranean,
		},
		{
			name:           "Colombian with ethnic 7 is SouthAmerican",
			nationality1:   "COL",
			nationality2:   "",
			ethnicValue:    7,
			expectedEthnic: SouthAmerican,
		},
		{
			name:           "Bangladeshi with ethnic 4 is MiddleEastSouthAsian",
			nationality1:   "BAN",
			nationality2:   "",
			ethnicValue:    4,
			expectedEthnic: MiddleEastSouthAsian,
		},
		{
			name:           "Cambodian with ethnic 5 is SouthEastAsian",
			nationality1:   "CAM",
			nationality2:   "",
			ethnicValue:    5,
			expectedEthnic: SouthEastAsian,
		},
		{
			name:           "Cambodian with ethnic 10 is Asian",
			nationality1:   "CAM",
			nationality2:   "",
			ethnicValue:    10,
			expectedEthnic: Asian,
		},
		{
			name:           "Colombian with ethnic 10 is SouthAmerican",
			nationality1:   "COL",
			nationality2:   "",
			ethnicValue:    10,
			expectedEthnic: SouthAmerican,
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			eth, err := getEthnic(tc.nationality1, tc.nationality2, tc.ethnicValue)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if eth != tc.expectedEthnic {
				t.Fatalf("expected %q, got %q", tc.expectedEthnic, eth)
			}
		})
	}
}

// Test getEthnic with a nationality not in NationEthnicMapping.
func TestGetEthnic_UnknownNationality1_XXX(t *testing.T) {
	_, err := getEthnic("XXX", "CAN", 0)
	if err == nil {
		t.Fatalf("expected error for unknown nationality, got nil")
	}
	if !strings.Contains(err.Error(), "ethnic not found") {
		t.Errorf("expected error containing 'ethnic not found', got %v", err)
	}
}

func TestGetEthnic_InvalidEthnicValue(t *testing.T) {
	_, err := getEthnic("CAN", "NZL", 999)
	if err == nil {
		t.Errorf("Expected error for invalid ethnic value, got nil")
	}
}

func TestGetPlayers_Success(t *testing.T) {
	// Create a temporary RTF file
	content := `| UID       | Nat       | 2nd Nat   | Name                       |           |           |           | 
| ---------------------------------------------------------------------------------------------------| 
| 2000133469| GER       | RSA       | Tebogo Maluleke            | 1         | 16        | 3         |
| ---------------------------------------------------------------------------------------------------|
| 2000133381| FRA       | MTQ       | Anthony Marlet             | 1         | 5         | 1         |
| ---------------------------------------------------------------------------------------------------|
`
	tmpFile, err := os.CreateTemp("", "test_rtf_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	fileOpener = os.Open
	players, err := GetPlayers(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(players) != 2 {
		t.Errorf("expected 2 players, got %d", len(players))
	}

	if players[0].ID != "2000133469" {
		t.Errorf("expected player id 2000133469, got %s", players[0].ID)
	}
}

// Test getEthnic with a nationality not in NationEthnicMapping.
func TestGetEthnic_UnknownNationality(t *testing.T) {
	_, err := getEthnic("XXX", "CAN", 0)
	if err == nil {
		t.Fatalf("expected error for unknown nationality, got nil")
	}
	if !strings.Contains(err.Error(), "ethnic not found") {
		t.Errorf("expected error containing 'ethnic not found', got %v", err)
	}
}

func TestGetPlayers_GetEthnicError(t *testing.T) {
	// In this test the first nationality is unknown.
	content := `| UID       | Nat       | 2nd Nat   | Name                       |           |           |           | 
| ---------------------------------------------------------------------------------------------------| 
| 2000133469| XXX       | RSA       | Tebogo Maluleke            | 1         | 16        | 3         |
| ---------------------------------------------------------------------------------------------------|
`
	tmpFile, err := os.CreateTemp("", "rtf_ethnicerror_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err = tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	_, err = GetPlayers(tmpFile.Name())
	if err == nil {
		t.Fatalf("expected error due to ethnicity lookup failure, got nil")
	}
	if !strings.Contains(err.Error(), "ethnic not found") {
		t.Errorf("expected error containing 'ethnic not found', got %v", err)
	}
}

// Test GetPlayers with a badly formatted line with too few fields.
func TestGetPlayers_BadFormat(t *testing.T) {
	// This content line produces only 3 fields.
	content := `| UID       | Nat       | 2nd Nat   | Name                       |           |           |           | 
| ---------------------------------------------------------------------------------------------------| 
| 2000133469| XXX       | RSA       | Tebogo Maluleke            |
| ---------------------------------------------------------------------------------------------------|
`

	tmpFile, err := os.CreateTemp("", "rtf_badformat_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	_, err = GetPlayers(tmpFile.Name())

	if err == nil {
		t.Fatalf("expected error due to bad RTF format, got nil")
	}

	if !strings.Contains(err.Error(), "not enough") {
		t.Errorf("expected error containing 'not enough', got %v", err)
	}
}

func TestGetPlayers_FileOpenError(t *testing.T) {
	origOpener := fileOpener
	defer func() { fileOpener = origOpener }()

	// Override fileOpener to return an error.
	fileOpener = func(path string) (*os.File, error) {
		return nil, errors.New("open file error")
	}

	_, err := GetPlayers("dummy.txt")
	if err == nil {
		t.Fatal("expected error from fileOpener, got nil")
	}
	if !strings.Contains(err.Error(), "open file error") {
		t.Errorf("expected error 'open file error', got %v", err)
	}
}

func TestGetPlayers_InvalidEthnicValue(t *testing.T) {
	// Create a temporary file with a line that has exactly 8 fields,
	// but field 8 (index 7) is non-numeric so that strconv.Atoi fails.
	// For example: "1234567|Player One|FIN|CAN|Data|Data|Data|NaN"
	content := "1234567|Player One|FIN|CAN|Data|Data|Data|NaN\n"
	tmpFile, err := os.CreateTemp("", "rtf_invalid_ethnic_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err = tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	_, err = GetPlayers(tmpFile.Name())
	if err == nil {
		t.Fatal("expected error due to invalid ethnic value, got nil")
	}
	if !strings.Contains(err.Error(), "invalid syntax") {
		t.Errorf("expected error containing 'invalid syntax', got %v", err)
	}
}

func TestGetPlayers_ScannerError(t *testing.T) {
	origOpener := fileOpener
	defer func() { fileOpener = origOpener }()

	// Create a temporary file that is closed immediately.
	fileOpener = func(path string) (*os.File, error) {
		tmp, err := os.CreateTemp("", "scanner_error_*.txt")
		if err != nil {
			return nil, err
		}
		tmp.Close() // return a closed file; subsequent reads should error.
		return tmp, nil
	}

	_, err := GetPlayers("dummy.txt")
	if err == nil {
		t.Fatal("expected error from scanner, got nil")
	}
	if !strings.Contains(err.Error(), "file already closed") && !strings.Contains(err.Error(), "bad file descriptor") {
		t.Errorf("expected error indicating closed file, got %v", err)
	}
}

package mapper

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func setup() {
	NationEthnicMapping = make(map[string]Ethnic)

	EthnicSet = mapset.NewSet[Ethnic]()
	EthnicSet.Add(African)
	EthnicSet.Add(Asian)
	EthnicSet.Add(Caucasian)
}

func TestOverrideNationEthnicMapping_ValidOverrides(t *testing.T) {
	setup()

	overrides := map[string]string{
		"USA": "African",
		"IND": "Caucasian",
	}

	err := OverrideNationEthnicMapping(overrides)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if NationEthnicMapping["USA"] != African || NationEthnicMapping["IND"] != Caucasian {
		t.Fatal("expected mappings to be updated, but they were not")
	}
}

func TestOverrideNationEthnicMapping_InvalidOverrides(t *testing.T) {
	setup()

	overrides := map[string]string{
		"USA": "Caucasian",
		"IND": "FakeEthnic", // Invalid ethnic input
	}

	err := OverrideNationEthnicMapping(overrides)
	if err == nil {
		t.Fatal("expected an error but got none")
	}

	expectedErrorMsg := `ethnic value "FakeEthnic" is not valid ethnic for "IND"`
	if err.Error() != expectedErrorMsg {
		t.Fatalf("expected error message to be %q, got %q", expectedErrorMsg, err.Error())
	}
}

func TestOverrideNationEthnicMapping_MixedValidAndInvalid(t *testing.T) {
	setup()

	overrides := map[string]string{
		"USA": "Caucasian",
		"IND": "FakeEthnic", // Invalid ethnic input
		"GBR": "Asian",
	}

	err := OverrideNationEthnicMapping(overrides)
	if err == nil {
		t.Fatal("expected an error but got none")
	}

	// Valid entries should still be added
	if NationEthnicMapping["USA"] != Caucasian || NationEthnicMapping["GBR"] != Asian {
		t.Fatal("valid mappings should have been updated, but they were not")
	}

	// Invalid entry should cause an error
	expectedErrorMsg := `ethnic value "FakeEthnic" is not valid ethnic for "IND"`
	if err.Error() != expectedErrorMsg {
		t.Fatalf("expected error message to be %q, got %q", expectedErrorMsg, err.Error())
	}
}

func TestOverrideNationEthnicMapping_NoOverrides(t *testing.T) {
	setup()

	overrides := map[string]string{} // empty map

	err := OverrideNationEthnicMapping(overrides)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

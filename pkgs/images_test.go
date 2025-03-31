package mapper

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock file creation for testing
func createMockFiles(t *testing.T, dir string, ethnicities []Ethnic, filesPerEthnicity int) {
	for _, ethnic := range ethnicities {
		ethnicDir := filepath.Join(dir, string(ethnic))
		if err := os.MkdirAll(ethnicDir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", ethnicDir, err)
		}
		for i := 0; i < filesPerEthnicity; i++ {
			filePath := filepath.Join(ethnicDir, fmt.Sprintf("image%d.png", i))
			if err := ioutil.WriteFile(filePath, []byte{}, 0644); err != nil {
				t.Fatalf("Failed to create file %s: %v", filePath, err)
			}
		}
	}
}

func mockEthnicities() []Ethnic {
	ethnicities := []Ethnic{
		African,
		Asian,
		Caucasian,
		CentralEuropean,
		EasternEuropeanCentralAsian,
		ItalianMediterranean,
		MiddleEastNorthAfrican,
		MiddleEastSouthAsian,
		SouthAmericanMediterranean,
		Scandinavian,
		SouthEastAsian,
		SouthAmerican,
		SpanishMediterranean,
		YugoslavGreek,
	}
	return ethnicities
}

func TestNewImagePool_Success(t *testing.T) {
	tmpDir := t.TempDir()

	ethnicities := mockEthnicities()
	filesPerEthnicity := 3
	createMockFiles(t, tmpDir, ethnicities, filesPerEthnicity)

	imagePool, err := NewImagePool(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create ImagePool: %v", err)
	}

	// Verify that the pool contains the correct number of images for each ethnicity
	for _, ethnic := range ethnicities {
		if len(imagePool.pool[ethnic]) != filesPerEthnicity {
			t.Errorf("Expected %d images for ethnicity %s, got %d", filesPerEthnicity, ethnic, len(imagePool.pool[ethnic]))
		}
	}
}

func TestNewImagePool_FolderError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a non-directory file to simulate an error
	nonDirFile := filepath.Join(tmpDir, "nonDirFile")
	if err := ioutil.WriteFile(nonDirFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create non-directory file: %v", err)
	}

	// Create a new ImagePool with a non-directory path
	imagePool, err := NewImagePool(nonDirFile)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "cannot get ethnic folder") {
		t.Errorf("Expected error containing 'cannot get ethnic folder', got %v", err)
	}
	if imagePool != nil {
		t.Error("Expected imagePool to be nil, got non-nil")
	}
}

func TestExcludeImages_NoExcludes(t *testing.T) {
	tmpDir := t.TempDir()

	ethnicities := mockEthnicities()
	filesPerEthnicity := 3
	createMockFiles(t, tmpDir, ethnicities, filesPerEthnicity)

	imagePool, err := NewImagePool(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create ImagePool: %v", err)
	}

	// Define no excludes
	excludes := []FilePath{}

	// Exclude images
	err = imagePool.ExcludeImages(excludes)
	if err != nil {
		t.Fatalf("Failed to exclude images: %v", err)
	}

	for _, ethnic := range ethnicities {
		if len(imagePool.pool[ethnic]) != filesPerEthnicity {
			t.Errorf("Expected %d images for ethnicity %s, got %d", filesPerEthnicity, ethnic, len(imagePool.pool[ethnic]))
		}
	}
}

func TestGetRandomImagePath_Success(t *testing.T) {
	tmpDir := t.TempDir()

	ethnicities := mockEthnicities()
	filesPerEthnicity := 3
	createMockFiles(t, tmpDir, ethnicities, filesPerEthnicity)

	imagePool, err := NewImagePool(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create ImagePool: %v", err)
	}

	// Get a random image path without removing it from the pool
	randImagePath, err := imagePool.GetRandomImagePath(Asian, false)
	if err != nil {
		t.Fatalf("Failed to get random image path: %v", err)
	}
	if randImagePath == "" {
		t.Errorf("Expected a random image path, got empty string")
	}

	if len(imagePool.pool[Asian]) != filesPerEthnicity {
		t.Errorf("Expected %d images for ethnicity %s, got %d", filesPerEthnicity, Asian, len(imagePool.pool[Asian]))
	}

	// Get a random image path and remove it from the pool
	randImagePath, err = imagePool.GetRandomImagePath(African, true)
	if err != nil {
		t.Fatalf("Failed to get random image path: %v", err)
	}
	if randImagePath == "" {
		t.Errorf("Expected a random image path, got empty string")
	}

	// Verify that the pool contains one less image
	if len(imagePool.pool[African]) != filesPerEthnicity-1 {
		t.Errorf("Expected %d images for ethnicity %s, got %d", filesPerEthnicity-1, African, len(imagePool.pool[African]))
	}
}

func TestGetRandomImagePath_NoImages(t *testing.T) {
	tmpDir := t.TempDir()

	ethnicities := mockEthnicities()
	filesPerEthnicity := 0 // No files
	createMockFiles(t, tmpDir, ethnicities, filesPerEthnicity)

	imagePool, err := NewImagePool(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create ImagePool: %v", err)
	}

	// Get a random image path for an ethnicity with no images
	_, err = imagePool.GetRandomImagePath(Asian, false)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "ran out of images for ethnicity") {
		t.Errorf("Expected error containing 'ran out of images for ethnicity', got %v", err)
	}
}

func TestGetRandomImagePath_SingleImage(t *testing.T) {
	tmpDir := t.TempDir()

	ethnicities := mockEthnicities()
	filesPerEthnicity := 1 // Single file
	createMockFiles(t, tmpDir, ethnicities, filesPerEthnicity)

	imagePool, err := NewImagePool(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create ImagePool: %v", err)
	}

	// Get a random image path for an ethnicity with a single image
	randImagePath, err := imagePool.GetRandomImagePath(Asian, false)
	if err != nil {
		t.Fatalf("Failed to get random image path: %v", err)
	}
	if randImagePath == "" {
		t.Errorf("Expected a random image path, got empty string")
	}

	// Verify that the pool still contains the correct number of images
	if len(imagePool.pool[Asian]) != filesPerEthnicity {
		t.Errorf("Expected %d images for ethnicity %s, got %d", filesPerEthnicity, Asian, len(imagePool.pool[Asian]))
	}

	// Remove the single image from the pool
	randImagePath, err = imagePool.GetRandomImagePath(African, true)
	if err != nil {
		t.Fatalf("Failed to get random image path: %v", err)
	}
	if randImagePath == "" {
		t.Errorf("Expected a random image path, got empty string")
	}

	// Verify that the pool contains no images
	if len(imagePool.pool[African]) != 0 {
		t.Errorf("Expected 0 images for ethnicity %s, got %d", African, len(imagePool.pool[African]))
	}

	// Attempt to get a random image path again, which should fail
	_, err = imagePool.GetRandomImagePath(African, false)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "ran out of images for ethnicity") {
		t.Errorf("Expected error containing 'ran out of images for ethnicity', got %v", err)
	}
}

func mockRandIntn(n int) int {
	return 0 // Always return the first index
}

func TestGetRandomImagePath_RandIntn(t *testing.T) {
	tmpDir := t.TempDir()

	ethnicities := mockEthnicities()
	filesPerEthnicity := 3
	createMockFiles(t, tmpDir, ethnicities, filesPerEthnicity)

	imagePool, err := NewImagePool(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create ImagePool: %v", err)
	}

	origRandIntn := rand.Intn
	defer func() { randIntnFunc = origRandIntn }()
	randIntnFunc = mockRandIntn

	randImagePath, err := imagePool.GetRandomImagePath(Asian, false)
	if err != nil {
		t.Fatalf("Failed to get random image path: %v", err)
	}
	if randImagePath != "image0" {
		t.Errorf("Expected 'image0', got %s", randImagePath)
	}

	// Verify that the pool still contains the correct number of images
	if len(imagePool.pool[Asian]) != filesPerEthnicity {
		t.Errorf("Expected %d images for ethnicity %s, got %d", filesPerEthnicity, Asian, len(imagePool.pool[Asian]))
	}
}

func TestExcludeImages_WithExcludes(t *testing.T) {
	tmpDir := t.TempDir()

	ethnicities := mockEthnicities()
	filesPerEthnicity := 3
	createMockFiles(t, tmpDir, ethnicities, filesPerEthnicity)

	imagePool, err := NewImagePool(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create ImagePool: %v", err)
	}

	// football manager requires filenames but not filename.png
	//TODO Should ExcludedImages be "Asian/image1" or "Asian/image1.png"
	excludes := []FilePath{
		"Asian/image1",
		"Central European/image2",
	}

	// Exclude images
	err = imagePool.ExcludeImages(excludes)
	if err != nil {
		t.Fatalf("Failed to exclude images: %v", err)
	}

	for ethnic, pool := range imagePool.pool {
		switch ethnic {
		case Asian:
			assert.Equal(t, filesPerEthnicity-1, len(pool))
			assert.NotContains(t, pool, "Asian/image1.png")
		case CentralEuropean:
			assert.Equal(t, filesPerEthnicity-1, len(pool))
			assert.NotContains(t, pool, "Central European/image2.png")
		default:
			assert.Equal(t, filesPerEthnicity, len(pool))
		}
	}
}

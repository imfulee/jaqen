package mapper

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"

	mapset "github.com/deckarep/golang-set/v2"
)

type Images struct {
	imagePool map[string][]string // images in folders
}

// `exclude` is a map of ethnic string to images
func (imgs *Images) Init(imageFolderPath string, perserve bool, exclude map[Ethnic][]string) error {
	if perserve && exclude == nil {
		return errors.New("if preseve mode is on, exclude map should be given")
	}

	if imageFolderPath == "" {
		return errors.New("no image folder root path")
	}

	ethnicities := [...]Ethnic{
		African,
		Asian,
		Caucasian,
		CentralEuropean,
		EasternEuropeanCentralAsian,
		ItalianMediterranean,
		MiddleEastNorthAfrican,
		MiddleEastSouthAmerican,
		SouthAmericanMediterranean,
		Scandinavian,
		SouthEastAsian,
		SouthAmerican,
		SpanishMediterranean,
		YugoslavGreek,
	}

	for ethnic := range exclude {
		found := false
		for _, allowedEthnic := range ethnicities {
			if ethnic == allowedEthnic {
				found = true
				break
			}
		}

		if !found {
			return errors.New("ethnic not allowed in exclude list")
		}
	}

	imgs.imagePool = make(map[string][]string)

	excludeSets := make(map[Ethnic]mapset.Set[string])
	for ethnic, excludedImages := range exclude {
		ethnicExcludedSet := mapset.NewSet[string]()

		for _, excludedImage := range excludedImages {
			ethnicExcludedSet.Add(excludedImage)
		}

		excludeSets[ethnic] = ethnicExcludedSet
	}

	imagePool := make(map[Ethnic]mapset.Set[string])
	for _, ethnic := range ethnicities {
		if _, ok := imagePool[ethnic]; !ok {
			imagePool[ethnic] = mapset.NewSet[string]()
		}

		ethnicImageFolderPath := path.Join(imageFolderPath, ethnic)
		ethnicImageFiles, error := os.ReadDir(ethnicImageFolderPath)
		if error != nil {
			return errors.Join(fmt.Errorf("cannot get ethnic folder %s", ethnicImageFolderPath), error)
		}

		for _, ethnicImageFile := range ethnicImageFiles {
			if ethnicImageFile.IsDir() {
				continue
			}

			filename := ethnicImageFile.Name()
			if filename == "" {
				continue
			}

			imagePool[ethnic].Add(filename)
		}

		currentPool := imagePool[ethnic]
		excludePool := excludeSets[ethnic]
		if perserve && excludePool != nil {
			imagePool[ethnic] = currentPool.Difference(excludePool)
		}

		imgs.imagePool[ethnic] = imagePool[ethnic].ToSlice()
	}

	return nil
}

func (imgs *Images) Random(ethnic string) (string, error) {
	if len(imgs.imagePool) == 0 {
		return "", errors.New("there are no images in image pool")
	}

	if _, ok := imgs.imagePool[ethnic]; !ok {
		return "", errors.New("no images for this ethnic")
	}

	chosenImage := ""

	// get images from ethnic provided
	images, hasEthnicInMap := imgs.imagePool[ethnic]
	if !hasEthnicInMap {
		return "", fmt.Errorf("no ethnic %s in image pool", ethnic)
	}

	imagePoolLength := len(images)
	if imagePoolLength == 0 {
		return "", fmt.Errorf("image pool of ethnic %s has no images", ethnic)
	}

	chosenImage = images[rand.Intn(imagePoolLength)]

	return chosenImage, nil
}

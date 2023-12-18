package newgen_manager

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
)

type Images struct {
	imagePool map[string][]string // images in folders
}

// `exclude` is a map of ethnic string to images
func (imgs *Images) Init(perserve bool, exclude map[Ethnic][]string) error {
	if perserve && exclude == nil {
		return errors.New("if preseve mode is on, exclude map should be given")
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
	ethnicErr := make(chan error)
	var wg sync.WaitGroup
	for _, ethnic := range ethnicities {
		wg.Add(1)
		go func(e string) {
			defer wg.Done()

			if _, ok := imagePool[e]; !ok {
				imagePool[e] = mapset.NewSet[string]()
			}

			ethnicImageFiles, err := os.ReadDir(e)
			if err != nil {
				ethnicErr <- errors.Join(fmt.Errorf("cannot get ethnic folder %s", e), err)
				return
			}

			for _, ethnicImageFile := range ethnicImageFiles {
				if ethnicImageFile.IsDir() {
					continue
				}

				filename := ethnicImageFile.Name()
				if filename == "" {
					continue
				}

				imagePool[e].Add(strings.TrimSuffix(filename, filepath.Ext(filename)))
			}

			currentPool := imagePool[e]
			excludePool := excludeSets[e]
			if perserve && excludePool != nil {
				imagePool[e] = currentPool.Difference(excludePool)
			}

			imgs.imagePool[e] = imagePool[e].ToSlice()
		}(ethnic)
	}

	go func() {
		wg.Wait()
		close(ethnicErr)
	}()

	if err := <-ethnicErr; err != nil {
		return err
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

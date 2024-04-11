package mapper

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"regexp"

	mapset "github.com/deckarep/golang-set/v2"
)

type ImagePool struct {
	pool map[Ethnic][]FilePath
}

func NewImagePool(imageRootPath string, excludes []FilePath) (*ImagePool, error) {
	pool := make(map[Ethnic][]FilePath)
	excludeSets := make(map[Ethnic]mapset.Set[FilePath])

	for _, ethnic := range Ethnicities {
		excludeSets[ethnic] = mapset.NewSet[FilePath]()
	}

	ethnicRegex := regexp.MustCompile(`^([a-zA-Z]+)`)
	filenameRegex := regexp.MustCompile(`\d+`)

	for _, filePath := range excludes {
		ethnic := Ethnic(ethnicRegex.Find([]byte(filePath)))
		filename := FilePath(filenameRegex.Find([]byte(filePath)))
		excludeSets[ethnic].Add(filename)
	}

	for _, ethnic := range Ethnicities {
		pool[ethnic] = make([]FilePath, 0)

		files, err := os.ReadDir(path.Join(imageRootPath, string(ethnic)))
		if err != nil {
			return nil, errors.Join(fmt.Errorf("cannot get ethnic folder %s", ethnic), err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			filename := file.Name()

			if excludeSet, hasSet := excludeSets[Ethnic(ethnic)]; hasSet {
				if excludeSet.Contains(FilePath(filename)) {
					continue
				}
			}

			pool[ethnic] = append(pool[ethnic], FilePath(filename))
		}
	}

	return &ImagePool{pool}, nil
}

func (images *ImagePool) Random(ethnic Ethnic) FilePath {
	length := len(images.pool[ethnic])
	if length == 0 {
		return ""
	}

	index := rand.Intn(length - 1)

	filename := images.pool[ethnic][index]

	images.pool[ethnic][index] = images.pool[ethnic][length-1]
	images.pool[ethnic] = images.pool[ethnic][:length-1]

	return filename
}

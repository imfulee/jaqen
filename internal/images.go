package mapper

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

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

	ethnictiesStrs := make([]string, len(Ethnicities))
	for i, ethnic := range Ethnicities {
		ethnictiesStrs[i] = string(ethnic)
	}
	ethnicRegexPattern := strings.Join(ethnictiesStrs, "|")
	ethnicRegexPattern = strings.ReplaceAll(ethnicRegexPattern, " ", `\s`)
	ethnicRegex := regexp.MustCompile(fmt.Sprintf(`\b(%s)\b`, ethnicRegexPattern))

	imageFilenameRegex := regexp.MustCompile(`[^\/]+$`)

	for _, filePath := range excludes {
		ethnic := Ethnic(ethnicRegex.FindString(string(filePath)))
		filename := FilePath(imageFilenameRegex.FindString(string(filePath)))
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

			// football manager requires filenames but not filename.png
			fullFilename := file.Name()
			filename := strings.TrimSuffix(filepath.Base(fullFilename), filepath.Ext(fullFilename))

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

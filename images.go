package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"slices"
)

type Images struct {
	unusedPool map[string][]string // images that are unused
	imagePool  map[string][]string // images in folders
}

func (imgs *Images) CreatePool(rootPath string, usedImages map[string][]string) error {
	// finding all the images in the rootPath
	rootFolder, rootFolderError := os.Open(rootPath)
	if rootFolderError != nil {
		return fmt.Errorf("cannot access root path: %w", rootFolderError)
	}
	defer rootFolder.Close()

	subFolders, subFoldersError := rootFolder.ReadDir(-1)
	if subFoldersError != nil {
		return fmt.Errorf("cannot read directory: %w", subFoldersError)
	}

	for _, subFolder := range subFolders {
		if !subFolder.IsDir() {
			continue
		}

		subFName := subFolder.Name()

		subfolderRoot, subfolderRootError := os.Open(path.Join(rootPath, subFName))
		if subfolderRootError != nil {
			return subfolderRootError
		}
		defer subfolderRoot.Close()

		subFolderFiles, subFolderFilesError := subfolderRoot.ReadDir(-1)
		if subFolderFilesError != nil {
			return subFolderFilesError
		}

		imgs.imagePool[subFName] = []string{}

		for _, subFolerFile := range subFolderFiles {
			if subFolerFile.IsDir() {
				continue
			}

			imgs.imagePool[subFName] = append(imgs.imagePool[subFName], subFolerFile.Name())
		}
	}

	if len(usedImages) > 0 {
		for ethnicInImagePool := range imgs.imagePool {
			usedImage, excludeHasSameEthnic := usedImages[ethnicInImagePool]
			if !excludeHasSameEthnic {
				return fmt.Errorf("exclude has bad format: ethnic %s is not found", ethnicInImagePool)
			}

			for _, image := range imgs.imagePool[ethnicInImagePool] {
				if !slices.Contains(usedImage, image) {
					if len(imgs.unusedPool[ethnicInImagePool]) == 0 {
						imgs.unusedPool[ethnicInImagePool] = []string{}
					}

					imgs.unusedPool[ethnicInImagePool] = append(imgs.unusedPool[ethnicInImagePool], image)
				}
			}
		}
	}

	return nil
}

func (imgs *Images) Random(ethnic string, excludeUsed bool) (string, error) {
	if len(imgs.imagePool) == 0 {
		return "", errors.New("there are no images in image pool")
	}

	var chosenImage string

	// get images from ethnic provided
	images, hasEthnicInMap := imgs.imagePool[ethnic]
	if !hasEthnicInMap {
		return "", fmt.Errorf("no ethnic %s in image pool", ethnic)
	}

	chosenImage = images[rand.Intn(len(images))]

	// if exclude used then exclude the used images
	if excludeUsed && len(imgs.unusedPool) > 0 {
		unusedImages, hasEthnicInMap := imgs.unusedPool[ethnic]
		if !hasEthnicInMap {
			return "", fmt.Errorf("no ethnic %s in unused image pool", ethnic)
		}

		chosenImage = unusedImages[rand.Intn(len(unusedImages))]
	}

	// remove chosen image from unused pool
	imgs.unusedPool[ethnic] = slices.DeleteFunc(imgs.unusedPool[ethnic], func(img string) bool {
		return img == chosenImage
	})

	return chosenImage, nil
}

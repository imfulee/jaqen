package mapper

import (
	"fmt"
	"path"
)

func PathForGame(id string) string {
	return fmt.Sprintf("graphics/pictures/person/%s/portrait", id)
}

func PathForImage(imgRoot, ethnic, imageFilename string) string {
	return path.Join(imgRoot, ethnic, imageFilename)
}

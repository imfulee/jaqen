package newgen_manager

import (
	"fmt"
	"path"
)

func PathForGame(id string) string {
	return fmt.Sprintf("graphics/pictures/person/%s/portrait", id)
}

func PathForImage(ethnic, imageFilename string) string {
	return path.Join(ethnic, imageFilename)
}

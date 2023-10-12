package main

import "fmt"

func PathForGame(id string) string {
	return fmt.Sprintf("graphics/pictures/person/%s/portrait", id)
}

func PathForImage(ethnic, id, imgRoot string) string {
	return fmt.Sprintf("%s/%s/%s", imgRoot, ethnic, id)
}

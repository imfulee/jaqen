package main

import (
	"fmt"
	mapper "jaqen/internal/mapper"
	"os"
	"strings"
)

func main() {
	Preserve := true
	ImageFolderRootPath := ""
	XMLPath := ""
	RTFPath := ""

	neMapper := mapper.NationEthnicMapper{}
	neMapper.Init("")
	rtfParser := mapper.RTF{NationEthnicMapper: &neMapper}
	fmPersons, err := rtfParser.Parse(RTFPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	xmlParser := mapper.XML{}
	err = xmlParser.ReadXML(XMLPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	previousMappings, err := xmlParser.GetPreviousMappings()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	excludeImages := make(map[string][]string)
	for _, mapping := range previousMappings {
		strs := strings.Split(mapping.FromPath, "/")
		if len(strs) != 3 {
			fmt.Println("error")
			os.Exit(0)
		}

		ethnic := strs[1]
		imageFilename := strs[2]

		if _, ok := excludeImages[ethnic]; !ok {
			excludeImages[ethnic] = make([]string, 0)
		}

		excludeImages[ethnic] = append(excludeImages[ethnic], imageFilename)
	}

	fmImages := mapper.Images{}
	err = fmImages.Init(ImageFolderRootPath, Preserve, excludeImages)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmMapper := mapper.Mapper{
		Preserve:            Preserve,
		Images:              &fmImages,
		ImageFolderRootPath: ImageFolderRootPath,
		Mappings:            previousMappings,
	}

	err = fmMapper.Map(fmPersons)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	xmlParser.UpdateMappings(fmMapper.Mappings)

	err = xmlParser.WriteXML(XMLPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

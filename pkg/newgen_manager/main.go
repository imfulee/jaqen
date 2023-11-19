package newgen_manager

import (
	"fmt"
	"os"
	"strings"
)

func Map(preserve bool, xmlPath, rtfPath string) {
	rtf := RTF{GetEthnicFromNation: NationMapToEthnic}
	persons, err := rtf.Parse(rtfPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	xml := &XML{}
	previousMappings, err := xml.Read(xmlPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	excludeImages := make(map[string][]string)
	for _, mapping := range previousMappings {
		strs := strings.Split(mapping.FromPath, "/")
		if len(strs) != 2 {
			fmt.Println("xml file bad format")
			os.Exit(0)
		}

		ethnic := strs[1]
		imageFilename := strs[2]

		if _, ok := excludeImages[ethnic]; !ok {
			excludeImages[ethnic] = make([]string, 0)
		}

		excludeImages[ethnic] = append(excludeImages[ethnic], imageFilename)
	}

	images := Images{}
	err = images.Init(preserve, excludeImages)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	mapper := Mapper{
		Preserve: preserve,
		Images:   &images,
		Mappings: previousMappings,
	}

	err = mapper.Map(persons)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = xml.Write(xmlPath, mapper.Mappings)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

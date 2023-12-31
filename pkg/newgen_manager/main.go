package newgen_manager

import (
	"fmt"
	"os"
	"path"
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
		ethnic, imageFilename := path.Split(mapping.FromPath)
		// ethnic would be ex. "YugoGreek/" so we should remove the trailing slash
		ethnic = strings.Trim(ethnic, "/")
		if ethnic == "" || imageFilename == "" {
			fmt.Println("xml file bad format")
			os.Exit(0)
		}

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

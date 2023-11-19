package newgen_manager

import (
	"fmt"
	"os"
	"strings"
)

func Map(preserve bool, xmlPath, rtfPath string) {
	Preserve := preserve
	XMLPath := xmlPath
	RTFPath := rtfPath

	rtf := RTF{GetEthnicFromNation: NationMapToEthnic}
	persons, err := rtf.Parse(RTFPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	xml := &XML{}
	previousMappings, err := xml.Read(XMLPath)
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

	fmImages := Images{}
	err = fmImages.Init(Preserve, excludeImages)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmMapper := Mapper{
		Preserve: Preserve,
		Images:   &fmImages,
		Mappings: previousMappings,
	}

	err = fmMapper.Map(persons)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = xml.Write(XMLPath, fmMapper.Mappings)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

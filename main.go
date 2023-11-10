package main

import (
	"fmt"
	mapper "jaqen/internal/mapper"
	"os"
)

func main() {
	rtfParser := mapper.RTF{NationEthnicMapper: &mapper.NationEthnicMapper{}}
	fmPersons, err := rtfParser.Parse("")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmImages := mapper.Images{}
	err = fmImages.Init("", false, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmMapper := mapper.Mapper{Images: &fmImages}
	err = fmMapper.ReadXML("")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = fmMapper.CreateMap(fmPersons)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = fmMapper.WriteXML("")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

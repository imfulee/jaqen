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
	fmImages.Init("", false, nil)
	fmMapper := mapper.Mapper{Images: &fmImages}
	fmMapper.ReadXML("")
	fmMapper.CreateMap(fmPersons)
	fmMapper.WriteXML("")
}

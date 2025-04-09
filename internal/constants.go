package internal

import mapper "jaqen/pkgs"

const (
	DefaultPreserve       = false
	DefaultXMLPath        = "./config.xml"
	DefaultRTFPath        = "./newgen.rtf"
	DefaultImagesPath     = "./"
	DefaultFMVersion      = mapper.FMVersion2024
	DefaultConfigPath     = "./jaqen.toml"
	DefaultAllowDuplicate = false
)

var SteamAppIdMap = map[string]string{
	mapper.FMVersion2020: "1100600",
	mapper.FMVersion2021: "1263850",
	mapper.FMVersion2022: "1569040",
	mapper.FMVersion2023: "1904540",
	mapper.FMVersion2024: "2252570",
}

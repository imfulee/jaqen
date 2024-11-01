package cmd

import (
	"fmt"
	internal "jaqen/internal"
	mapper "jaqen/pkgs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	preserve   bool
	xmlPath    string
	rtfPath    string
	imgDir     string
	fmVersion  string
	configPath string
)

const (
	flagkeys_preserve  = "preserve"
	flagkeys_xml       = "xml"
	flagkeys_rtf       = "rtf"
	flagkeys_img       = "img"
	flagkeys_fmversion = "version"
	flagkeys_config    = "config"
)

func mapFaces(cmd *cobra.Command, _ []string) {
	if _, err := os.Stat(configPath); err == nil {
		configFromFile, err := internal.ReadConfig(configPath)
		if err != nil {
			log.Fatalln(err)
		}

		if !cmd.Flags().Changed(flagkeys_preserve) && configFromFile.Preserve != nil {
			preserve = *configFromFile.Preserve
		}
		if !cmd.Flags().Changed(flagkeys_img) && configFromFile.IMGPath != nil {
			imgDir = *configFromFile.IMGPath
		}
		if !cmd.Flags().Changed(flagkeys_img) && configFromFile.IMGPath != nil {
			imgDir = *configFromFile.IMGPath
		}
		if !cmd.Flags().Changed(flagkeys_xml) && configFromFile.XMLPath != nil {
			xmlPath = *configFromFile.XMLPath
		}
		if !cmd.Flags().Changed(flagkeys_rtf) && configFromFile.RTFPath != nil {
			rtfPath = *configFromFile.RTFPath
		}
		if !cmd.Flags().Changed(flagkeys_fmversion) && configFromFile.FMVersion != nil {
			fmVersion = *configFromFile.FMVersion
		}

		err = mapper.OverrideNationEthnicMapping(*configFromFile.MappingOverride)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := os.Stat(imgDir); err != nil {
		log.Fatalln(fmt.Errorf("image directory could not be found: %w", err))
	}

	if _, err := os.Stat(xmlPath); err != nil {
		log.Fatalln(fmt.Errorf("xml file could not be found: %w", err))
	}

	if _, err := os.Stat(rtfPath); err != nil {
		log.Fatalln(fmt.Errorf("rtf file could not be found: %w", err))
	}

	mapping, err := mapper.NewMapping(xmlPath, fmVersion)
	if err != nil {
		log.Fatalln(err)
	}

	imagePool, err := mapper.NewImagePool(imgDir, mapping.AssignedImages())
	if err != nil {
		log.Fatalln(err)
	}

	players, err := mapper.GetPlayers(rtfPath)
	if err != nil {
		log.Fatalln(err)
	}

	imgDirPathAbs, err := filepath.Abs(imgDir)
	if err != nil {
		log.Fatalln(err)
	}
	xmlFilePathAbs, err := filepath.Abs(xmlPath)
	if err != nil {
		log.Fatalln(err)
	}

	rel := ""
	isXMLFileInsideImgDir := imgDirPathAbs == filepath.Dir(xmlFilePathAbs)
	if !isXMLFileInsideImgDir {
		rel, err = filepath.Rel(xmlFilePathAbs, imgDirPathAbs)
		if err != nil {
			log.Fatalln(err)
		}
	}
	rel = strings.TrimPrefix(rel, "./")

	for _, player := range players {
		if preserve && mapping.Exist(player.ID) {
			continue
		}

		imgFilename, err := imagePool.Random(player.Ethnic)
		if err != nil {
			log.Fatalln(err)
		}

		mapping.MapToImage(player.ID, mapper.FilePath(path.Join(rel, string(player.Ethnic), string(imgFilename))))
	}

	if err := mapping.Save(); err != nil {
		log.Fatalln(err)
	}

	if err := mapping.Write(xmlPath); err != nil {
		log.Fatalln(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "jaqen",
	Short: "Creates your mapping file for Football Manager regen images",
	Long:  `CLI that creates your mapping file for Football Manager regen images.`,
	Run:   mapFaces,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&preserve, flagkeys_preserve, "p", internal.DefaultPreserve, "Preserve previous settings")
	rootCmd.Flags().StringVarP(&xmlPath, flagkeys_xml, "x", internal.DefaultXMLPath, "Specify XML file path")
	rootCmd.Flags().StringVarP(&rtfPath, flagkeys_rtf, "r", internal.DefaultRTFPath, "Specify RTF file path")
	rootCmd.Flags().StringVarP(&imgDir, flagkeys_img, "i", internal.DefaultImagesPath, "Specify the image directory path")
	rootCmd.Flags().StringVarP(&fmVersion, flagkeys_fmversion, "v", internal.DefaultFMVersion, "Specify the football manager version")
	rootCmd.Flags().StringVarP(&configPath, flagkeys_config, "c", internal.DefaultConfigPath, "Specify the config file path")
}

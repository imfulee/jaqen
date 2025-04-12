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
	preserve       bool
	xmlPath        string
	rtfPath        string
	imgDir         string
	fmVersion      string
	configPath     string
	allowDuplicate bool
)

const (
	flagkeysPreserve  = "preserve"
	flagkeysXml       = "xml"
	flagkeysRtf       = "rtf"
	flagkeysImg       = "img"
	flagkeyFmVersion = "version"
	flagkeyConfig    = "config"
	flagkeyDuplicate = "allow_duplicate"
)

func mapFaces(cmd *cobra.Command, _ []string) {
	if _, err := os.Stat(configPath); err == nil {
		configFromFile, err := internal.ReadConfig(configPath)
		if err != nil {
			log.Fatalln(err)
		}

		if !cmd.Flags().Changed(flagkeysPreserve) && configFromFile.Preserve != nil {
			preserve = *configFromFile.Preserve
		}
		if !cmd.Flags().Changed(flagkeysImg) && configFromFile.IMGPath != nil {
			imgDir = *configFromFile.IMGPath
		}
		if !cmd.Flags().Changed(flagkeysImg) && configFromFile.IMGPath != nil {
			imgDir = *configFromFile.IMGPath
		}
		if !cmd.Flags().Changed(flagkeysXml) && configFromFile.XMLPath != nil {
			xmlPath = *configFromFile.XMLPath
		}
		if !cmd.Flags().Changed(flagkeysRtf) && configFromFile.RTFPath != nil {
			rtfPath = *configFromFile.RTFPath
		}
		if !cmd.Flags().Changed(flagkeyFmVersion) && configFromFile.FMVersion != nil {
			fmVersion = *configFromFile.FMVersion
		}
		if !cmd.Flags().Changed(flagkeyDuplicate) && configFromFile.AllowDuplicate != nil {
			allowDuplicate = *configFromFile.AllowDuplicate
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

	imagePool, err := mapper.NewImagePool(imgDir)
	if err != nil {
		log.Fatalln(err)
	}

	if !allowDuplicate {
		if err := imagePool.ExcludeImages(mapping.AssignedImages()); err != nil {
			log.Fatalln(err)
		}
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

		imgFilename, err := imagePool.GetRandomImagePath(player.Ethnic, !allowDuplicate)
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
	rootCmd.Flags().BoolVarP(&preserve, flagkeysPreserve, "p", internal.DefaultPreserve, "Preserve previous settings")
	rootCmd.Flags().StringVarP(&xmlPath, flagkeysXml, "x", internal.DefaultXMLPath, "Specify XML file path")
	rootCmd.Flags().StringVarP(&rtfPath, flagkeysRtf, "r", internal.DefaultRTFPath, "Specify RTF file path")
	rootCmd.Flags().StringVarP(&imgDir, flagkeysImg, "i", internal.DefaultImagesPath, "Specify the image directory path")
	rootCmd.Flags().StringVarP(&fmVersion, flagkeyFmVersion, "v", internal.DefaultFMVersion, "Specify the football manager version")
	rootCmd.Flags().StringVarP(&configPath, flagkeyConfig, "c", internal.DefaultConfigPath, "Specify the config file path")
	rootCmd.Flags().BoolVarP(&allowDuplicate, flagkeyDuplicate, "d", internal.DefaultAllowDuplicate, "Allow duplicate images")
}

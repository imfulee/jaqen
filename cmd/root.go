package cmd

import (
	"fmt"
	mapper "jaqen/internal"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	preserve  bool
	xmlPath   string
	rtfPath   string
	imgDir    string
	fmVersion string
)

var rootCmd = &cobra.Command{
	Use:   "jaqen",
	Short: "Creates your mapping file for Football Manager regen images",
	Long:  `CLI that creates your mapping file for Football Manager regen images.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(imgDir); err != nil {
			log.Fatalln(fmt.Errorf("image directory could not be found: %w", err))
		}

		if _, err := os.Stat(xmlPath); err != nil {
			log.Fatalln(fmt.Errorf("xml file could not be found: %w", err))
		}

		if _, err := os.Stat(rtfPath); err != nil {
			log.Fatalln(fmt.Errorf("rtf file could not be found: %w", err))
		}

		getPlayers := mapper.GetPlayersBuilder(rtfPath)

		mapping, err := mapper.NewMapping(xmlPath, fmVersion)
		if err != nil {
			log.Fatalln(err)
		}

		imagePool, err := mapper.NewImagePool(imgDir, mapping.AssignedImages())
		if err != nil {
			log.Fatalln(err)
		}

		players, err := getPlayers()
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
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&preserve, "preserve", "p", false, "Preserve previous settings")
	rootCmd.Flags().StringVar(&xmlPath, "xml", "./config.xml", "Specify XML file path")
	rootCmd.Flags().StringVar(&rtfPath, "rtf", "./newgen.rtf", "Specify RTF file path")
	rootCmd.Flags().StringVar(&imgDir, "img", "./", "Specify the image directory path")
	rootCmd.Flags().StringVar(&fmVersion, "ver", "2024", "Specify the football manager version")
}

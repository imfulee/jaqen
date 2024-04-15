package cmd

import (
	"errors"
	mapper "jaqen/internal"
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
			panic(errors.Join(errors.New("image directory could not be found"), err))
		}

		if _, err := os.Stat(xmlPath); err != nil {
			panic(errors.Join(errors.New("xml file could not be found"), err))
		}

		if _, err := os.Stat(rtfPath); err != nil {
			panic(errors.Join(errors.New("rtf file could not be found"), err))
		}

		getPlayers := mapper.GetPlayersBuilder(rtfPath)

		mapping, err := mapper.NewMapping(xmlPath, fmVersion)
		if err != nil {
			panic(err)
		}

		imagePool, err := mapper.NewImagePool(imgDir, mapping.AssignedImages())
		if err != nil {
			panic(err)
		}

		players, err := getPlayers()
		if err != nil {
			panic(err)
		}

		imgDirPathAbs, err := filepath.Abs(imgDir)
		if err != nil {
			panic(err)
		}
		xmlFilePathAbs, err := filepath.Abs(xmlPath)
		if err != nil {
			panic(err)
		}

		rel := ""
		isXMLFileInsideImgDir := imgDirPathAbs == filepath.Dir(xmlFilePathAbs)
		if !isXMLFileInsideImgDir {
			rel, err = filepath.Rel(xmlFilePathAbs, imgDirPathAbs)
			if err != nil {
				panic(err)
			}
		}
		rel = strings.TrimPrefix(rel, "./")

		for _, player := range players {
			if preserve && mapping.Exist(player.ID) {
				continue
			}

			imgFilename := imagePool.Random(player.Ethnic)
			mapping.MapToImage(player.ID, mapper.FilePath(path.Join(rel, string(player.Ethnic), string(imgFilename))))
		}

		if err := mapping.Save(); err != nil {
			panic(err)
		}

		mapping.Write(xmlPath)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&preserve, "preserve", "p", false, "Preserve previous settings")
	rootCmd.Flags().StringVar(&xmlPath, "xml", "./config.xml", "Specify XML file path")
	rootCmd.Flags().StringVar(&rtfPath, "rtf", "./newgen.rtf", "Specify RTF file path")
	rootCmd.Flags().StringVar(&imgDir, "img", "./", "Specify the image directory path")
	rootCmd.Flags().StringVar(&fmVersion, "ver", "2024", "Specify the football manager version")
}

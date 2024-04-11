package cmd

import (
	mapper "jaqen/internal"
	"os"
	"path"

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

		for _, player := range players {
			if preserve && mapping.Exist(player.ID) {
				continue
			}

			imgFilename := imagePool.Random(player.Ethnic)
			mapping.MapToImage(player.ID, mapper.FilePath(path.Join(string(player.Ethnic), string(imgFilename))))
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
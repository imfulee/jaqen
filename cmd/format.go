package cmd

import (
	"errors"
	"log"
	"os"
	"sort"

	internal "jaqen/internal"

	"github.com/spf13/cobra"
)

func formatConfig(cmd *cobra.Command, args []string) {
	configPath := "./jaqen.toml"
	if len(args) == 1 {
		configPath = args[0]
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalln(errors.New("config file not found"))
	}

	config, err := internal.ReadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	if config.MappingOverride != nil {
		nations := make([]string, 0, len(*config.MappingOverride))
		for nation := range *config.MappingOverride {
			nations = append(nations, nation)
		}
		sort.Strings(nations)

		mappingOverride := make(map[string]string)
		for _, nation := range nations {
			mappingOverride[nation] = (*config.MappingOverride)[nation]
		}
		*config.MappingOverride = mappingOverride
	}

	if err = internal.WriteConfig(config, configPath); err != nil {
		log.Fatalln(err)
	}
}

var formatCmd = &cobra.Command{
	Use:   "format /path/to/config/file",
	Short: "Formats config file",
	Long:  "Formats config file specified. Defaults to ./jaqen.toml",
	Run:   formatConfig,
}

func init() {
	rootCmd.AddCommand(formatCmd)
}

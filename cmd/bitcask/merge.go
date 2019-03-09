package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prologic/bitcask"
)

var mergeCmd = &cobra.Command{
	Use:     "merge",
	Aliases: []string{"clean", "compact", "defrag"},
	Short:   "Merges the Datafiles in the Database",
	Long: `This merges all non-active Datafiles in the Database and
compacts the data stored on disk. Old values are removed as well as deleted
keys.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("path")
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.WithError(err).Error("error parsing force flag")
			os.Exit(1)
		}

		os.Exit(merge(path, force))
	},
}

func init() {
	RootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().BoolP(
		"force", "f", false,
		"Force a re-merge even if .hint files exist",
	)
}

func merge(path string, force bool) int {
	err := bitcask.Merge(path, force)
	if err != nil {
		log.WithError(err).Error("error merging database")
		return 1
	}

	return 0
}

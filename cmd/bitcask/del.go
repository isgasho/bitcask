package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prologic/bitcask"
)

var delCmd = &cobra.Command{
	Use:     "del <key>",
	Aliases: []string{"delete", "remove", "rm"},
	Short:   "Delete a key and its value",
	Long:    `This deletes a key and its value`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("path")

		key := args[0]

		os.Exit(del(path, key))
	},
}

func init() {
	RootCmd.AddCommand(delCmd)
}

func del(path, key string) int {
	db, err := bitcask.Open(path)
	if err != nil {
		log.WithError(err).Error("error opening database")
		return 1
	}

	err = db.Delete(key)
	if err != nil {
		log.WithError(err).Error("error deleting key")
		return 1
	}

	return 0
}

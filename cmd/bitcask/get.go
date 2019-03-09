package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prologic/bitcask"
)

var getCmd = &cobra.Command{
	Use:     "get <key>",
	Aliases: []string{"view"},
	Short:   "Get a new Key and display its Value",
	Long:    `This retrieves a key and display its value`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("path")

		key := args[0]

		os.Exit(get(path, key))
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}

func get(path, key string) int {
	db, err := bitcask.Open(path)
	if err != nil {
		log.WithError(err).Error("error opening database")
		return 1
	}

	value, err := db.Get(key)
	if err != nil {
		log.WithError(err).Error("error reading key")
		return 1
	}

	fmt.Printf("%s\n", string(value))
	log.WithField("key", key).WithField("value", value).Debug("key/value")

	return 0
}

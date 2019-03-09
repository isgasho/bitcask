package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prologic/bitcask"
)

var setCmd = &cobra.Command{
	Use:     "set <key> [<value>]",
	Aliases: []string{"add"},
	Short:   "Add/Set a new Key/Value pair",
	Long: `This adds or sets a new key/value pair.

If the value is not specified as an argument it is read from standard input.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("path")

		key := args[0]

		var value io.Reader
		if len(args) > 1 {
			value = bytes.NewBufferString(args[1])
		} else {
			value = os.Stdin
		}

		os.Exit(set(path, key, value))
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}

func set(path, key string, value io.Reader) int {
	db, err := bitcask.Open(path)
	if err != nil {
		log.WithError(err).Error("error opening database")
		return 1
	}

	data, err := ioutil.ReadAll(value)
	if err != nil {
		log.WithError(err).Error("error writing key")
		return 1
	}

	err = db.Put(key, data)
	if err != nil {
		log.WithError(err).Error("error writing key")
		return 1
	}

	return 0
}

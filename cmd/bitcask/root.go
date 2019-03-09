package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prologic/bitcask"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "bitcask",
	Version: bitcask.FullVersion(),
	Short:   "Command-line tools for bitcask",
	Long: `This is the command-line tool to interact with a bitcask database.

This lets you get, set and delete key/value pairs as well as perform merge
(or compaction) operations. This tool serves as an example implementation
however is also intended to be useful in shell scripts.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// set logging level
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	},
}

// Execute adds all child commands to the root command
// and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolP(
		"debug", "d", false,
		"Enable debug logging",
	)

	RootCmd.PersistentFlags().StringP(
		"path", "p", "/tmp/bitcask",
		"Path to Bitcask database",
	)

	viper.BindPFlag("path", RootCmd.PersistentFlags().Lookup("path"))
	viper.SetDefault("path", "/tmp/bitcask")

	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.SetDefault("debug", false)
}

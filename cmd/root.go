package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"watchmen/config"
	"watchmen/log"
)

var (
	asciiArt = `
██     ██  █████  ████████  ██████ ██   ██ ███    ███ ███████ ███    ██ 
██     ██ ██   ██    ██    ██      ██   ██ ████  ████ ██      ████   ██ 
██  █  ██ ███████    ██    ██      ███████ ██ ████ ██ █████   ██ ██  ██ 
██ ███ ██ ██   ██    ██    ██      ██   ██ ██  ██  ██ ██      ██  ██ ██ 
 ███ ███  ██   ██    ██     ██████ ██   ██ ██      ██ ███████ ██   ████ 
`
)

var rootCMD = &cobra.Command{
	Use:   "watchmen",
	Short: "watchmen monitor",
}

var (
	configFilePath string
)

func init() {
	cobra.OnInitialize(configure)

	rootCMD.PersistentFlags().StringVarP(&configFilePath, "config", "c", "config.yml", "config file")

	rootCMD.AddCommand(serveCMD)
	rootCMD.AddCommand(databaseCMD)
}

func configure() {
	cfg := config.Init(configFilePath)
	log.SetupLogger(config.C.Logger)

	time.Local = cfg.Server.Location
	fmt.Print(asciiArt)
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Printf("failed to execute root command: %s\n", err.Error())
		os.Exit(1)
	}
}

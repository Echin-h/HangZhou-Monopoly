package cmd

import (
	"github.com/Echin-h/HangZhou-Monopoly/cmd/config"
	"github.com/Echin-h/HangZhou-Monopoly/cmd/create"
	"github.com/Echin-h/HangZhou-Monopoly/cmd/server"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "app",
	Short:        "app",
	SilenceUsage: true,
	Long:         `app`,
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(create.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

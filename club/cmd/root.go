package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	logFile string
	isDev   bool
)

var RootCmd = &cobra.Command{
	Use:   "club",
	Short: "Application club",
	Long:  `Application club`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&logFile, "log", "l", "", "Log file")
	RootCmd.PersistentFlags().BoolVar(&isDev, "dev", false, "Development mode")
}

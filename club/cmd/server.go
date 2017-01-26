package cmd

import (
	"fmt"

	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/manifest"
	"github.com/shumkovdenis/club/server"
	"github.com/spf13/cobra"
)

var cfgFile string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Start server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := manifest.Read(); err != nil {
			return fmt.Errorf("read manifest failed: %s", err)
		}

		if err := config.Read(cfgFile); err != nil {
			return fmt.Errorf("read config failed: %s", err)
		}

		if err := server.Start(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&cfgFile, "config", config.File, "config file")
}

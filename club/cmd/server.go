package cmd

import (
	"github.com/shumkovdenis/club"
	"github.com/shumkovdenis/club/manifest"
	"github.com/spf13/cobra"
)

var config string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Start server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := manifest.ReadConfig(config); err != nil {
			return err
		}
		if err := club.StartServer(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().StringVar(&config, "config", "", "config file")
}

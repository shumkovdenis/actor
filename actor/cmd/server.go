package cmd

import (
	"github.com/shumkovdenis/actor"
	"github.com/shumkovdenis/actor/manifest"
	"github.com/spf13/cobra"
)

var config string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Long:  "Run server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := manifest.ReadConfig(config); err != nil {
			return err
		}
		if err := actor.StartServer(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().StringVar(&config, "config", "", "config file")
}

package cmd

import (
	"github.com/shumkovdenis/club"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/manifest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	mnfDef = "manifest.json"
	cfgDef = "config.toml"
)

var mnf = viper.New()
var cfg = viper.New()
var cfgFile string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `Start server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := readManifest(); err != nil {
			return err
		}

		if err := readConfig(); err != nil {
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

	serverCmd.Flags().StringVar(&cfgFile, "config", cfgDef, "config file")
}

func readManifest() error {
	mnf.SetConfigType("json")
	mnf.SetConfigFile(mnfDef)

	if err := mnf.ReadInConfig(); err != nil {
		return err
	}

	if err := mnf.Unmarshal(manifest.Get()); err != nil {
		return err
	}

	return nil
}

func readConfig() error {
	cfg.SetConfigType("toml")
	cfg.SetConfigFile(cfgFile)

	if err := cfg.ReadInConfig(); err != nil {
		return err
	}

	if err := cfg.Unmarshal(config.Get()); err != nil {
		return err
	}

	return nil
}

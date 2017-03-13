package cmd

import (
	"github.com/shumkovdenis/club"
	"github.com/spf13/cobra"
)

var dataFile string

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start client",
	Long:  `Start client.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := club.StartClient(dataFile); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVar(&dataFile, "data", "", "data file")
}

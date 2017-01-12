package cmd

import (
	"github.com/shumkovdenis/club"
	"github.com/spf13/cobra"
)

var dataFile string

// clientCmd represents the client command
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

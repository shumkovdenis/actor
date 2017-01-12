package cmd

import (
	"github.com/shumkovdenis/club/packer"
	"github.com/spf13/cobra"
)

var (
	srcFile string
	trgPath string
)

var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Unpack update",
	Long:  `Unpack update.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := packer.Unpack(srcFile, trgPath); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(unpackCmd)

	unpackCmd.Flags().StringVarP(&srcFile, "source", "s", "", "Update file")
	unpackCmd.Flags().StringVarP(&trgPath, "target", "t", "", "Target path")
}

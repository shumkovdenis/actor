package cmd

import (
	"github.com/shumkovdenis/club/packer"
	"github.com/spf13/cobra"
)

var (
	oldPath string
	newPath string
	dstFile string
)

var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack update",
	Long:  `Pack update.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := packer.Pack(oldPath, newPath, dstFile); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(packCmd)

	packCmd.Flags().StringVarP(&oldPath, "old", "o", "", "Old version path")
	packCmd.Flags().StringVarP(&newPath, "new", "n", "", "New version path")
	packCmd.Flags().StringVarP(&dstFile, "dest", "d", "", "Out update file")
}

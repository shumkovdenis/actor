package cmd

import (
	"os/exec"
	"runtime"

	"github.com/shumkovdenis/club/packer"
	"github.com/spf13/cobra"
)

var (
	updFile string
	appPath string
	appFile string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update",
	Long:  `Update.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := packer.Unpack(updFile, appPath); err != nil {
			return err
		}

		var c *exec.Cmd
		if runtime.GOOS == "darwin" {
			c = exec.Command("open", appFile)
		} else {
			c = exec.Command(appFile)
		}

		if err := c.Start(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	unpackCmd.Flags().StringVarP(&updFile, "source", "s", "", "Update file")
	unpackCmd.Flags().StringVarP(&appPath, "target", "t", "", "Application path")
	unpackCmd.Flags().StringVarP(&appFile, "app", "a", "", "Application file")
}

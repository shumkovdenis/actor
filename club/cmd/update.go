package cmd

import (
	"os/exec"
	"runtime"

	"github.com/shumkovdenis/club/logger"
	"github.com/shumkovdenis/club/packer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
		logger.InitProduction(logFile)

		logger.L().Info("start update")

		if err := packer.Unpack(updFile, appPath); err != nil {
			logger.L().Error("unpack failed",
				zap.Error(err),
			)
			return err
		}

		logger.L().Info("unpack complete")

		var c *exec.Cmd
		if runtime.GOOS == "darwin" {
			c = exec.Command("open", appFile)
		} else {
			c = exec.Command(appFile)
		}

		logger.L().Info("start app")

		if err := c.Start(); err != nil {
			logger.L().Error("start app failed",
				zap.Error(err),
			)
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updFile, "source", "s", "", "Update file")
	updateCmd.Flags().StringVarP(&appPath, "target", "t", "", "Application path")
	updateCmd.Flags().StringVarP(&appFile, "app", "a", "", "Application file")
}

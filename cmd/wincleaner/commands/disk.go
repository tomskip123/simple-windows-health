package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewDiskCommand returns the cobra command for 'disk'
func NewDiskCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "disk",
		Short: "Run Disk Cleanup utility",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Disk Cleanup", cleaner.RunDiskCleanup, 0)
		},
	}
} 
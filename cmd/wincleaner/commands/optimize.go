package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewOptimizeCommand returns the cobra command for 'optimize'
func NewOptimizeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "optimize",
		Short: "Run Disk Optimization (defrag for HDDs, TRIM for SSDs)",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Disk Optimization", cleaner.RunDiskOptimization, 0)
		},
	}
} 
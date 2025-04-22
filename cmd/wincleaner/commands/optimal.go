package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewOptimalCommand returns the cobra command for 'optimal'
func NewOptimalCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "optimal",
		Short: "Apply optimal Windows settings (e.g., disable Fast Boot)",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Set Optimal Windows Settings", func() error { return cleaner.SetOptimalWindowsSettings(core.Verbose) }, 0)
		},
	}
} 
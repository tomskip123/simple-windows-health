package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewResetNetCommand returns the cobra command for 'resetnet'
func NewResetNetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "resetnet",
		Short: "Reset Windows network configuration",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Reset Network Configuration", cleaner.ResetNetworkConfig, 0)
		},
	}
} 
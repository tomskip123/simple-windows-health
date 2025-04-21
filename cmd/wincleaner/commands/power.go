package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewPowerCommand returns the cobra command for 'power'
func NewPowerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "power",
		Short: "Optimize power configuration settings",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Optimize Power Configuration", cleaner.OptimizePowerConfig, 0)
		},
	}
} 
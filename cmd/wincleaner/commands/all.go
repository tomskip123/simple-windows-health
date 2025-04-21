package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
)

// NewAllCommand returns the cobra command for 'all'
func NewAllCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Run all cleaning operations",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunAllOperations(cmd.Context())
		},
	}
} 
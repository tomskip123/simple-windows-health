package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewRecycleCommand returns the cobra command for 'recycle'
func NewRecycleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "recycle",
		Short: "Empty Recycle Bin",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Empty Recycle Bin", cleaner.EmptyRecycleBin, 0)
		},
	}
} 
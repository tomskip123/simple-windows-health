package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewTempCommand returns the cobra command for 'temp'
func NewTempCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "temp",
		Short: "Clean temporary files",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Temporary Files Cleaning", cleaner.CleanTempFiles, 0)
		},
	}
} 
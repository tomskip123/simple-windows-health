package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
	"time"
)

// NewSFCCommand returns the cobra command for 'sfc'
func NewSFCCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sfc",
		Short: "Run System File Checker",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "System File Checker", func() error { return cleaner.RunSystemFileChecker(core.Verbose) }, 1000*time.Second)
		},
	}
} 
package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewMemcheckCommand returns the cobra command for 'memcheck'
func NewMemcheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "memcheck",
		Short: "Run Windows Memory Diagnostic tool",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Windows Memory Diagnostic", func() error { return cleaner.RunMemoryDiagnostic(core.Verbose) }, 0)
		},
	}
} 
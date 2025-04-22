package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewPrefetchCommand returns the cobra command for 'prefetch'
func NewPrefetchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "prefetch",
		Short: "Clean Windows prefetch directory",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Clean Prefetch Cache", func() error { return cleaner.CleanPrefetch(core.Verbose) }, 0)
		},
	}
} 
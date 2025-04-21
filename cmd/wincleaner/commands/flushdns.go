package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewFlushDNSCommand returns the cobra command for 'flushdns'
func NewFlushDNSCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "flushdns",
		Short: "Flush DNS resolver cache",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Flush DNS Cache", cleaner.FlushDNSCache, 0)
		},
	}
} 
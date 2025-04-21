package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewEventsCommand returns the cobra command for 'events'
func NewEventsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "events",
		Short: "Clear Windows event logs",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Event Logs Clearing", cleaner.ClearEventLogs, 0)
		},
	}
} 
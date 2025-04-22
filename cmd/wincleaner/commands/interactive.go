package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/pkg/interactive"
)

// NewInteractiveCommand returns the cobra command for 'interactive'
func NewInteractiveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "interactive",
		Short: "Launch the interactive menu interface",
		Run: func(cmd *cobra.Command, args []string) {
			interactive.RunInteractiveMode()
		},
	}
} 
package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// NewStatusCommand returns the cobra command for 'status'
func NewStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Display system status information",
		Run: func(cmd *cobra.Command, args []string) {
			core.Logger.Info("Retrieving system status...")
			fmt.Println("Retrieving system status...")
			status, err := cleaner.GetSystemStatus()
			if err != nil {
				fmt.Printf("Error retrieving system status: %v\n", err)
				core.Logger.Errorf("Error retrieving system status: %v", err)
			} else {
				core.DisplaySystemStatus(status)
				core.Logger.Info("System status displayed successfully.")
			}
		},
	}
} 
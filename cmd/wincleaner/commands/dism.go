package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
	"time"
)

// NewDismCommand returns the cobra command for 'dism'
func NewDismCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dism",
		Short: "Run DISM to repair Windows image",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "DISM Windows Image Repair", func() error { return cleaner.RunDISM(core.Verbose) }, 1000*time.Second)
		},
	}
} 
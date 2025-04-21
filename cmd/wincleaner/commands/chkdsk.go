package commands

import (
	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
	"time"
)

// NewChkdskCommand returns the cobra command for 'chkdsk'
func NewChkdskCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "chkdsk",
		Short: "Run Check Disk utility",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunOperation(cmd.Context(), "Check Disk", cleaner.RunCheckDisk, 90*time.Second)
		},
	}
} 
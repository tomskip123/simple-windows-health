package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/windows_health/cmd/wincleaner/commands"
	"github.com/user/windows_health/cmd/wincleaner/core"
)

const version = "1.0.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "wincleaner",
		Short: "Windows Health Cleaner - A utility for system maintenance",
		Long:  `Windows Health Cleaner is a comprehensive Windows system maintenance utility that helps keep your Windows system in optimal condition.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			core.LoadConfig()
			core.SetupLogger()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(core.Config.DefaultOps) > 0 {
				for _, op := range core.Config.DefaultOps {
					for _, c := range cmd.Commands() {
						if c.Name() == op {
							c.Run(c, []string{})
						}
					}
				}
			} else {
				cmd.Help()
			}
		},
	}

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("Windows Health Cleaner version {{.Version}}\n")
	rootCmd.PersistentFlags().StringVar(&core.ConfigFile, "config", "", "Path to config YAML file")
	rootCmd.PersistentFlags().BoolVarP(&core.Verbose, "verbose", "v", false, "Enable verbose logging to console")

	rootCmd.AddCommand(
		commands.NewDiskCommand(),
		commands.NewTempCommand(),
		commands.NewEventsCommand(),
		commands.NewSFCCommand(),
		commands.NewDismCommand(),
		commands.NewRecycleCommand(),
		commands.NewOptimizeCommand(),
		commands.NewChkdskCommand(),
		commands.NewFlushDNSCommand(),
		commands.NewMemcheckCommand(),
		commands.NewPrefetchCommand(),
		commands.NewResetNetCommand(),
		commands.NewAllCommand(),
		commands.NewStatusCommand(),
		commands.NewOptimalCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

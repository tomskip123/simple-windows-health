package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/user/windows_health/cmd/wincleaner/core"
	"github.com/user/windows_health/pkg/cleaner"
)

// MenuOption represents a menu option in the interactive interface
type MenuOption struct {
	Name        string
	Description string
	Action      func() error
}

// RunInteractiveMode starts an interactive console-based interface
func RunInteractiveMode() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Clear the screen and show the menu
		fmt.Print("\033[H\033[2J") // ANSI escape sequence to clear screen

		fmt.Println("======================================")
		fmt.Println("   Windows Health Cleaner Utility    ")
		fmt.Println("======================================")

		if !cleaner.IsAdmin() {
			fmt.Println("\n⚠️  WARNING: Not running with administrator privileges.")
			fmt.Println("    Some operations may not work correctly.")
			fmt.Println("    Choose option 0 to restart with admin rights.")
		}

		// Display menu options
		options := getMenuOptions()
		for i, option := range options {
			fmt.Printf("%d. %s - %s\n", i, option.Name, option.Description)
		}

		fmt.Println("\nq. Exit program")
		fmt.Print("\nEnter your choice: ")

		// Read user input
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Handle quit option
		if input == "q" || input == "Q" {
			fmt.Println("Exiting program. Goodbye!")
			return
		}

		// Parse the choice
		var choice int
		fmt.Sscanf(input, "%d", &choice)

		// Validate choice
		if choice < 0 || choice >= len(options) {
			fmt.Println("Invalid choice. Please try again.")
			time.Sleep(2 * time.Second)
			continue
		}

		// Execute the chosen action
		fmt.Printf("\nRunning %s...\n", options[choice].Name)
		err := options[choice].Action()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("%s completed successfully.\n", options[choice].Name)
		}

		fmt.Print("\nPress Enter to continue...")
		reader.ReadString('\n')
	}
}

// getMenuOptions returns the list of menu options
func getMenuOptions() []MenuOption {
	options := []MenuOption{
		{
			Name:        "Restart with Admin Rights",
			Description: "Restart the application with administrator privileges",
			Action:      cleaner.RunAsAdmin,
		},
		{
			Name:        "System Status",
			Description: "Display detailed system status information",
			Action: func() error {
				status, err := cleaner.GetSystemStatus()
				if err != nil {
					return err
				}

				fmt.Println("\n=== System Status Information ===")
				fmt.Printf("Windows Version: %s\n", status.WindowsVersion)
				fmt.Printf("Last Boot Time: %s\n", status.LastBootTime)

				fmt.Println("\nDisk Space Information:")
				fmt.Println("------------------------")
				for drive, info := range status.DiskSpace {
					fmt.Printf("Drive %s:\n", drive)
					fmt.Printf("  Total Size: %s\n", info.TotalSize)
					fmt.Printf("  Free Space: %s\n", info.FreeSpace)
					fmt.Printf("  Used Space: %s (%s)\n", info.UsedSpace, info.UsedPercent)
					fmt.Println()
				}

				return nil
			},
		},
	}

	// Add a visual separator for optimal settings
	options = append(options, MenuOption{
		Name:        "--- Optimal Settings ---",
		Description: "",
		Action:      func() error { return nil },
	})
	options = append(options, MenuOption{
		Name:        "Apply Optimal Windows Settings",
		Description: "Apply recommended settings (e.g., disables Fast Boot)",
		Action:      func() error { return cleaner.SetOptimalWindowsSettings(core.Verbose) },
	})

	options = append(options,
		MenuOption{Name: "Disk Cleanup", Description: "Run Windows Disk Cleanup utility", Action: func() error { return cleaner.RunDiskCleanup(core.Verbose) }},
		MenuOption{Name: "Clean Temporary Files", Description: "Remove temporary files from Windows directories", Action: func() error { return cleaner.CleanTempFiles(core.Verbose) }},
		MenuOption{Name: "Clear Event Logs", Description: "Clear Windows event logs", Action: func() error { return cleaner.ClearEventLogs(core.Verbose) }},
		MenuOption{Name: "System File Checker", Description: "Run SFC to scan and repair Windows system files", Action: func() error { return cleaner.RunSystemFileChecker(core.Verbose) }},
		MenuOption{Name: "DISM Repair", Description: "Run DISM to repair the Windows image", Action: func() error { return cleaner.RunDISM(core.Verbose) }},
		MenuOption{Name: "Empty Recycle Bin", Description: "Empty the Windows Recycle Bin", Action: func() error { return cleaner.EmptyRecycleBin(core.Verbose) }},
		MenuOption{Name: "Disk Optimization", Description: "Run disk optimization (defrag for HDDs, TRIM for SSDs)", Action: func() error { return cleaner.RunDiskOptimization(core.Verbose) }},
		MenuOption{Name: "Check Disk", Description: "Run CHKDSK to scan and repair disk errors", Action: func() error { return cleaner.RunCheckDisk(core.Verbose) }},
		MenuOption{Name: "Flush DNS Cache", Description: "Clear Windows DNS resolver cache", Action: func() error { return cleaner.FlushDNSCache(core.Verbose) }},
		MenuOption{Name: "Memory Diagnostic", Description: "Run Windows Memory Diagnostic tool", Action: func() error { return cleaner.RunMemoryDiagnostic(core.Verbose) }},
		MenuOption{Name: "Clean Prefetch Cache", Description: "Clean Windows prefetch directory", Action: func() error { return cleaner.CleanPrefetch(core.Verbose) }},
		MenuOption{Name: "Reset Network", Description: "Reset Windows network configuration", Action: func() error { return cleaner.ResetNetworkConfig(core.Verbose) }},
		MenuOption{
			Name:        "Run All Cleaning Operations",
			Description: "Execute all cleaning operations sequentially",
			Action: func() error {
				fmt.Println("Running all cleaning operations...")

				operations := []struct {
					name   string
					action func() error
				}{
					{"Disk Cleanup", func() error { return cleaner.RunDiskCleanup(core.Verbose) }},
					{"Temporary Files Cleaning", func() error { return cleaner.CleanTempFiles(core.Verbose) }},
					{"Event Logs Clearing", func() error { return cleaner.ClearEventLogs(core.Verbose) }},
					{"System File Checker", func() error { return cleaner.RunSystemFileChecker(core.Verbose) }},
					{"DISM Windows Image Repair", func() error { return cleaner.RunDISM(core.Verbose) }},
					{"Empty Recycle Bin", func() error { return cleaner.EmptyRecycleBin(core.Verbose) }},
					{"Disk Optimization", func() error { return cleaner.RunDiskOptimization(core.Verbose) }},
					{"Check Disk", func() error { return cleaner.RunCheckDisk(core.Verbose) }},
					{"Flush DNS Cache", func() error { return cleaner.FlushDNSCache(core.Verbose) }},
					{"Clean Prefetch Cache", func() error { return cleaner.CleanPrefetch(core.Verbose) }},
				}

				for _, op := range operations {
					fmt.Printf("\nRunning %s...\n", op.name)
					err := op.action()
					if err != nil {
						fmt.Printf("Error running %s: %v\n", op.name, err)
					} else {
						fmt.Printf("%s completed successfully.\n", op.name)
					}
					// Small pause between operations
					time.Sleep(1 * time.Second)
				}

				return nil
			},
		},
	)

	return options
} 
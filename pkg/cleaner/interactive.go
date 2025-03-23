package cleaner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
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

		if !IsAdmin() {
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
	return []MenuOption{
		{
			Name:        "Restart with Admin Rights",
			Description: "Restart the application with administrator privileges",
			Action:      RunAsAdmin,
		},
		{
			Name:        "System Status",
			Description: "Display detailed system status information",
			Action: func() error {
				status, err := GetSystemStatus()
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
		{
			Name:        "Disk Drive Status",
			Description: "Check disk drive health using WMIC",
			Action:      CheckDiskDriveStatus,
		},
		{
			Name:        "Disk Cleanup",
			Description: "Run Windows Disk Cleanup utility",
			Action:      RunDiskCleanup,
		},
		{
			Name:        "Clean Temporary Files",
			Description: "Remove temporary files from Windows directories",
			Action:      CleanTempFiles,
		},
		{
			Name:        "Clear Event Logs",
			Description: "Clear Windows event logs",
			Action:      ClearEventLogs,
		},
		{
			Name:        "System File Checker",
			Description: "Run SFC to scan and repair Windows system files",
			Action:      RunSystemFileChecker,
		},
		{
			Name:        "DISM Repair",
			Description: "Run DISM to repair the Windows image",
			Action:      RunDISM,
		},
		{
			Name:        "Empty Recycle Bin",
			Description: "Empty the Windows Recycle Bin",
			Action:      EmptyRecycleBin,
		},
		{
			Name:        "Disk Optimization",
			Description: "Run disk optimization (defrag for HDDs, TRIM for SSDs)",
			Action:      RunDiskOptimization,
		},
		{
			Name:        "Check Disk",
			Description: "Run CHKDSK to scan and repair disk errors",
			Action:      RunCheckDisk,
		},
		{
			Name:        "Flush DNS Cache",
			Description: "Clear Windows DNS resolver cache",
			Action:      FlushDNSCache,
		},
		{
			Name:        "Memory Diagnostic",
			Description: "Run Windows Memory Diagnostic tool",
			Action:      RunMemoryDiagnostic,
		},
		{
			Name:        "Clean Prefetch Cache",
			Description: "Clean Windows prefetch directory",
			Action:      CleanPrefetch,
		},
		{
			Name:        "Reset Network Configuration",
			Description: "Reset Windows network configuration",
			Action:      ResetNetworkConfig,
		},
		{
			Name:        "Optimize Startup Items",
			Description: "Analyze and manage startup programs with performance impact assessment",
			Action:      OptimizeStartup,
		},
		{
			Name:        "Run All Maintenance Tasks",
			Description: "Execute all cleaning and optimization operations",
			Action: func() error {
				fmt.Println("Running all maintenance tasks...")
				
				tasks := []struct {
					name   string
					action func() error
				}{
					{"Disk Cleanup", RunDiskCleanup},
					{"Clean Temporary Files", CleanTempFiles},
					{"Clear Event Logs", ClearEventLogs},
					{"Empty Recycle Bin", EmptyRecycleBin},
					{"Flush DNS Cache", FlushDNSCache},
					{"Clean Prefetch Cache", CleanPrefetch},
					{"Optimize Power Configuration", OptimizePowerConfig},
					{"Disk Drive Status", CheckDiskDriveStatus},
				}
				
				// Only add network reset if we have admin privileges
				if IsAdmin() {
					tasks = append(tasks, struct {
						name   string
						action func() error
					}{"Reset Network Configuration", ResetNetworkConfig})
					
					// Add other admin-only tasks
					adminTasks := []struct {
						name   string
						action func() error
					}{
						{"System File Checker", RunSystemFileChecker},
						{"DISM Windows Image Repair", RunDISM},
						{"Disk Optimization", RunDiskOptimization},
						{"Check Disk", RunCheckDisk},
						{"Windows Memory Diagnostic", RunMemoryDiagnostic},
					}
					
					tasks = append(tasks, adminTasks...)
				} else {
					fmt.Println("\nNote: Some tasks require administrator privileges and will be skipped.")
					fmt.Println("Restart the application with admin rights to run all tasks.")
				}
				
				// Run all tasks and keep track of results
				successCount := 0
				failedCount := 0
				var failedTasks []string
				
				for _, task := range tasks {
					fmt.Printf("\nRunning %s...\n", task.name)
					err := task.action()
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						failedCount++
						failedTasks = append(failedTasks, task.name)
					} else {
						fmt.Printf("%s completed successfully.\n", task.name)
						successCount++
					}
					time.Sleep(500 * time.Millisecond) // Small delay between tasks
				}
				
				// Print summary
				fmt.Println("\n========== Maintenance Summary ==========")
				fmt.Printf("Total tasks attempted: %d\n", successCount+failedCount)
				fmt.Printf("Successfully completed: %d\n", successCount)
				
				if failedCount > 0 {
					fmt.Printf("Failed tasks: %d\n", failedCount)
					for i, task := range failedTasks {
						fmt.Printf("  %d. %s\n", i+1, task)
					}
					fmt.Println("\nSome tasks failed. You might need to run them individually or with admin privileges.")
				} else {
					fmt.Println("All maintenance tasks completed successfully!")
				}
				
				fmt.Println("\nPress Enter to continue...")
				bufio.NewReader(os.Stdin).ReadString('\n')
				return nil
			},
		},
	}
}

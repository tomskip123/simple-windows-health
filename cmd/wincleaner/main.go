package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/windows_health/pkg/cleaner"
)

const (
	version = "1.0.0"
)

func main() {
	// Define command-line flags
	diskCleanup := flag.Bool("disk", false, "Run Disk Cleanup utility")
	tempFiles := flag.Bool("temp", false, "Clean temporary files")
	eventLogs := flag.Bool("events", false, "Clear Windows event logs")
	sfc := flag.Bool("sfc", false, "Run System File Checker")
	dism := flag.Bool("dism", false, "Run DISM to repair Windows image")
	recycle := flag.Bool("recycle", false, "Empty Recycle Bin")
	all := flag.Bool("all", false, "Run all cleaning operations")
	status := flag.Bool("status", false, "Display system status information")
	versionFlag := flag.Bool("version", false, "Display version information")

	// Override default usage
	flag.Usage = usage
	flag.Parse()

	// Print version information if requested
	if *versionFlag {
		fmt.Printf("Windows Health Cleaner version %s\n", version)
		return
	}

	// Check if at least one flag is provided
	if !*diskCleanup && !*tempFiles && !*eventLogs && !*sfc && !*dism && !*recycle && !*all && !*status {
		usage()
		os.Exit(1)
	}

	// Check for administrator privileges for certain operations
	if (*all || *sfc || *dism) && !isAdmin() {
		fmt.Println("Warning: Some operations (SFC, DISM) require administrator privileges.")
		fmt.Println("Please run this program as administrator for full functionality.")
	}

	// Show status if requested
	if *status {
		fmt.Println("Retrieving system status...")
		sysStatus, err := cleaner.GetSystemStatus()
		if err != nil {
			fmt.Printf("Error retrieving system status: %v\n", err)
		} else {
			displaySystemStatus(sysStatus)
		}
	}

	// Run appropriate cleaning operations
	if *all || *diskCleanup {
		runOperation("Disk Cleanup", cleaner.RunDiskCleanup)
	}

	if *all || *tempFiles {
		runOperation("Temporary Files Cleaning", cleaner.CleanTempFiles)
	}

	if *all || *eventLogs {
		runOperation("Event Logs Clearing", cleaner.ClearEventLogs)
	}

	if *all || *sfc {
		runOperation("System File Checker", cleaner.RunSystemFileChecker)
	}

	if *all || *dism {
		runOperation("DISM Windows Image Repair", cleaner.RunDISM)
	}

	if *all || *recycle {
		runOperation("Empty Recycle Bin", cleaner.EmptyRecycleBin)
	}

	if *all || *diskCleanup || *tempFiles || *eventLogs || *sfc || *dism || *recycle {
		fmt.Println("Cleaning operations completed.")
	}
}

// usage prints a custom usage message
func usage() {
	fmt.Println("Windows Health Cleaner - A utility for system maintenance")
	fmt.Printf("Version: %s\n\n", version)
	fmt.Println("Usage: wincleaner [options]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  wincleaner -disk -temp       # Run Disk Cleanup and clean temporary files")
	fmt.Println("  wincleaner -status           # Display system status information")
	fmt.Println("  wincleaner -all              # Run all cleaning operations")
}

// displaySystemStatus formats and prints the system status information
func displaySystemStatus(status *cleaner.SystemStatus) {
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
}

// runOperation runs a cleaning operation and handles errors
func runOperation(name string, operation func() error) {
	fmt.Printf("Running %s...\n", name)
	err := operation()
	if err != nil {
		fmt.Printf("Error running %s: %v\n", name, err)
	} else {
		fmt.Printf("%s completed successfully.\n", name)
	}
}

// isAdmin checks if the program is running with administrator privileges
func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

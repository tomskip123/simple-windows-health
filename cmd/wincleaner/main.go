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
	diskOptimization := flag.Bool("optimize", false, "Run Disk Optimization (defrag for HDDs, TRIM for SSDs)")
	checkDisk := flag.Bool("chkdsk", false, "Run Check Disk utility")
	flushDNS := flag.Bool("flushdns", false, "Flush DNS resolver cache")
	memoryDiag := flag.Bool("memcheck", false, "Run Windows Memory Diagnostic tool")
	prefetch := flag.Bool("prefetch", false, "Clean Windows prefetch directory")
	powerCfg := flag.Bool("power", false, "Optimize power configuration settings")
	resetNet := flag.Bool("resetnet", false, "Reset Windows network configuration")
	all := flag.Bool("all", false, "Run all cleaning operations")
	status := flag.Bool("status", false, "Display system status information")
	versionFlag := flag.Bool("version", false, "Display version information")
	interactive := flag.Bool("interactive", false, "Launch interactive console mode")
	requireAdmin := flag.Bool("admin", false, "Request administrator privileges")

	// Override default usage
	flag.Usage = usage
	flag.Parse()

	// Print version information if requested
	if *versionFlag {
		fmt.Printf("Windows Health Cleaner version %s\n", version)
		return
	}

	// Check if admin privileges are explicitly requested
	if *requireAdmin && !cleaner.IsAdmin() {
		fmt.Println("Requesting administrator privileges...")
		err := cleaner.RunAsAdmin()
		if err != nil {
			fmt.Printf("Error requesting admin privileges: %v\n", err)
			os.Exit(1)
		}
		// The program will exit in RunAsAdmin if successful
	}

	// Launch interactive mode if requested
	if *interactive {
		cleaner.RunInteractiveMode()
		return
	}

	// Check if at least one flag is provided
	if !*diskCleanup && !*tempFiles && !*eventLogs && !*sfc && !*dism && !*recycle &&
		!*diskOptimization && !*checkDisk && !*flushDNS && !*memoryDiag && !*prefetch &&
		!*powerCfg && !*resetNet && !*all && !*status {
		// If no flags, default to interactive mode
		cleaner.RunInteractiveMode()
		return
	}

	// Check for administrator privileges for certain operations
	if (*all || *sfc || *dism || *diskOptimization || *checkDisk || *memoryDiag) && !cleaner.IsAdmin() {
		fmt.Println("Warning: Some operations (SFC, DISM, Disk Optimization, Check Disk, Memory Diagnostic) require administrator privileges.")
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

	if *all || *diskOptimization {
		runOperation("Disk Optimization", cleaner.RunDiskOptimization)
	}

	if *all || *checkDisk {
		runOperation("Check Disk", cleaner.RunCheckDisk)
	}

	if *all || *flushDNS {
		runOperation("Flush DNS Cache", cleaner.FlushDNSCache)
	}

	if *all || *memoryDiag {
		runOperation("Windows Memory Diagnostic", cleaner.RunMemoryDiagnostic)
	}

	if *all || *prefetch {
		runOperation("Clean Prefetch Cache", cleaner.CleanPrefetch)
	}

	if *all || *powerCfg {
		runOperation("Optimize Power Configuration", cleaner.OptimizePowerConfig)
	}

	if *all || *resetNet {
		runOperation("Reset Network Configuration", cleaner.ResetNetworkConfig)
	}

	if *all || *diskCleanup || *tempFiles || *eventLogs || *sfc || *dism || *recycle ||
		*diskOptimization || *checkDisk || *flushDNS || *memoryDiag || *prefetch ||
		*powerCfg || *resetNet {
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
	fmt.Println("  wincleaner -interactive      # Launch interactive console mode")
	fmt.Println("  wincleaner -admin            # Request administrator privileges")
	fmt.Println("  wincleaner -disk -temp       # Run Disk Cleanup and clean temporary files")
	fmt.Println("  wincleaner -optimize         # Run Disk Optimization (defrag for HDDs, TRIM for SSDs)")
	fmt.Println("  wincleaner -chkdsk           # Run Check Disk utility")
	fmt.Println("  wincleaner -prefetch         # Clean Windows prefetch directory")
	fmt.Println("  wincleaner -power            # Optimize power configuration settings")
	fmt.Println("  wincleaner -resetnet         # Reset Windows network configuration")
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

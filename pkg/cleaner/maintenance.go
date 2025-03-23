package cleaner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunDiskOptimization runs appropriate optimization based on drive type (defrag for HDDs, TRIM for SSDs)
func RunDiskOptimization() error {
	// First check if drive is SSD or HDD using PowerShell
	checkCmd := exec.Command("powershell", "-Command",
		"Get-PhysicalDisk | Select-Object DeviceId, MediaType | ConvertTo-Json")
	output, err := checkCmd.Output()
	if err != nil {
		// If we cannot determine drive types, use the safest option (automatic)
		cmd := exec.Command("defrag", "/C", "/O", "/U", "/V")
		return cmd.Run()
	}

	// Parse output to determine drive types
	// Simple check - if any SSD is found, use /O which automatically selects the correct optimization
	if strings.Contains(string(output), "SSD") {
		// Use /O which will automatically select proper optimization method based on media type
		cmd := exec.Command("defrag", "/C", "/O", "/U", "/V")
		return cmd.Run()
	} else {
		// Traditional defrag for HDDs
		cmd := exec.Command("defrag", "/C", "/D", "/U", "/V")
		return cmd.Run()
	}
}

// RunCheckDisk schedules the Windows Check Disk utility to run on next reboot
func RunCheckDisk() error {
	// Schedule CHKDSK to run on next boot for system drive
	systemDrive := os.Getenv("SYSTEMDRIVE")
	if systemDrive == "" {
		systemDrive = "C:"
	}
	
	// Run with exactly the parameters requested: /f (fix errors), /r (recover bad sectors), /x (force dismount)
	cmd := exec.Command("chkdsk", systemDrive, "/f", "/r", "/x")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("error scheduling disk check: %w - %s", err, string(output))
	}
	
	// Output a helpful message to the user
	fmt.Printf("CheckDisk has been scheduled for the next system reboot.\n")
	fmt.Printf("The check will be performed on %s drive when you restart your computer.\n", systemDrive)
	
	return nil
}

// CheckDiskDriveStatus uses WMIC to display disk drive model and status
func CheckDiskDriveStatus() error {
	fmt.Println("Checking disk drive status...")
	
	// Run the exact WMIC command requested
	cmd := exec.Command("wmic", "diskdrive", "get", "model,status")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("error checking disk drive status: %w", err)
	}
	
	fmt.Println("\nDisk Drive Status:")
	fmt.Println("------------------")
	fmt.Println(string(output))
	
	return nil
}

// FlushDNSCache flushes the Windows DNS resolver cache
func FlushDNSCache() error {
	cmd := exec.Command("ipconfig", "/flushdns")
	return cmd.Run()
}

// RunMemoryDiagnostic runs the Windows Memory Diagnostic tool
func RunMemoryDiagnostic() error {
	cmd := exec.Command("mdsched")
	return cmd.Run()
}

// OptimizePowerConfig optimizes Windows power settings
func OptimizePowerConfig() error {
	// First check if we're on a laptop by checking for a battery
	checkBatteryCmd := exec.Command("powershell", "-Command", 
		"(Get-WmiObject -Class Win32_Battery) -ne $null")
	batteryOutput, _ := checkBatteryCmd.Output()
	isLaptop := strings.TrimSpace(string(batteryOutput)) == "True"

	var powerScheme string
	if isLaptop {
		// For laptops, use balanced scheme
		powerScheme = "SCHEME_BALANCED"
	} else {
		// For desktops, use high performance
		powerScheme = "SCHEME_MIN"
	}

	// Set the appropriate power scheme
	setSchemeCmd := exec.Command("powercfg", "/setactive", powerScheme)
	err := setSchemeCmd.Run()
	if err != nil {
		return err
	}

	// Set display timeout
	displayTimeoutCmd := exec.Command("powercfg", "/change", "monitor-timeout-ac", "15")
	err = displayTimeoutCmd.Run()
	if err != nil {
		return err
	}

	// Set sleep timeout
	sleepTimeoutCmd := exec.Command("powercfg", "/change", "standby-timeout-ac", "30")
	return sleepTimeoutCmd.Run()
}

// CleanPrefetch cleans the Windows prefetch directory
func CleanPrefetch() error {
	// Using PowerShell to clean prefetch directory with proper error handling
	cmd := exec.Command("powershell", "-Command",
		"Remove-Item -Path \"$env:SystemRoot\\Prefetch\\*\" -Force -ErrorAction SilentlyContinue")
	return cmd.Run()
}

// ResetNetworkConfig resets Windows network configuration
func ResetNetworkConfig() error {
	// Check for admin rights first
	if !IsAdmin() {
		return fmt.Errorf("administrator privileges required for network reset")
	}

	fmt.Println("Resetting network configuration (this may temporarily disconnect your network)...")
	
	// Run commands with proper error handling
	commands := []struct {
		desc string
		args []string
	}{
		{"Resetting Winsock catalog", []string{"netsh", "winsock", "reset"}},
		{"Resetting TCP/IP stack", []string{"netsh", "int", "ip", "reset"}},
		{"Resetting Windows Firewall", []string{"netsh", "advfirewall", "reset"}},
		{"Flushing DNS cache", []string{"ipconfig", "/flushdns"}},
		{"Releasing IP address", []string{"ipconfig", "/release"}},
		{"Renewing IP address", []string{"ipconfig", "/renew"}},
	}

	var errorMsgs []string
	
	for _, command := range commands {
		fmt.Printf("  %s...\n", command.desc)
		cmd := exec.Command(command.args[0], command.args[1:]...)
		output, err := cmd.CombinedOutput()
		
		// Some commands might fail but we want to continue with others
		if err != nil {
			errMsg := fmt.Sprintf("  - %s failed: %v - %s", command.desc, err, string(output))
			errorMsgs = append(errorMsgs, errMsg)
		}
	}
	
	// If we encountered errors, return them but only after running all commands
	if len(errorMsgs) > 0 {
		return fmt.Errorf("some network reset commands failed:\n%s", strings.Join(errorMsgs, "\n"))
	}
	
	fmt.Println("Network configuration has been reset successfully.")
	return nil
}

// OptimizeStartup lists startup items and provides info on system performance impact
func OptimizeStartup() error {
	// Use WMIC to get startup items
	cmd := exec.Command("powershell", "-Command", 
		"Get-CimInstance Win32_StartupCommand | Select-Object Name, Command, Location | Format-Table -AutoSize")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("error getting startup items: %w", err)
	}
	
	fmt.Println("=== Current Startup Items ===")
	fmt.Println(string(output))
	fmt.Println("To disable startup items, you can use Task Manager > Startup tab")
	fmt.Println("Or run: msconfig > Startup tab")
	
	// Get startup impact assessment
	impactCmd := exec.Command("powershell", "-Command", 
		"Get-StartupApp | Select-Object Name, Command, StartupType, StartupImpact | Format-Table -AutoSize")
	impactOutput, err := impactCmd.CombinedOutput()
	
	if err == nil {
		fmt.Println("\n=== Startup Impact Assessment ===")
		fmt.Println(string(impactOutput))
	}
	
	return nil
}

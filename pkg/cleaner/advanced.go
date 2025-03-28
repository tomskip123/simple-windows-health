package cleaner

import (
	"fmt"
	"os/exec"
)

// ClearEventLogs clears Windows event logs using the wevtutil command
func ClearEventLogs() error {
	// Get a list of all event logs
	cmd := exec.Command("wevtutil", "el")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// Clear each event log
	logs := splitLines(string(output))
	for _, logName := range logs {
		if logName == "" {
			continue
		}
		clearCmd := exec.Command("wevtutil", "cl", logName)
		clearCmd.Run() // Ignore errors, as some logs might be protected
	}

	return nil
}

// RunSystemFileChecker runs the Windows System File Checker to repair system files
func RunSystemFileChecker() error {
	// Use exactly the command specified: sfc /scannow
	fmt.Println("Running System File Checker (sfc /scannow)...")
	fmt.Println("This may take some time. Please wait...")
	
	cmd := exec.Command("sfc", "/scannow")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("SFC scan failed: %w - %s", err, string(output))
	}
	
	fmt.Println(string(output))
	return nil
}

// RunDISM runs the Deployment Image Servicing and Management tool to repair Windows image
func RunDISM() error {
	// Use exactly the command specified: DISM /Online /Cleanup-Image /RestoreHealth
	fmt.Println("Running DISM Image Repair (DISM /Online /Cleanup-Image /RestoreHealth)...")
	fmt.Println("This may take some time. Please wait...")
	
	cmd := exec.Command("DISM", "/Online", "/Cleanup-Image", "/RestoreHealth")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("DISM repair failed: %w - %s", err, string(output))
	}
	
	fmt.Println(string(output))
	return nil
}

// EmptyRecycleBin empties the Windows Recycle Bin
func EmptyRecycleBin() error {
	// Use only the official PowerShell cmdlet to empty the recycle bin
	// This is the Microsoft-approved way to do this operation
	cmd := exec.Command("powershell", "-Command", "Clear-RecycleBin", "-Force", "-ErrorAction", "SilentlyContinue")
	err := cmd.Run()
	
	if err != nil {
		// Instead of using shell commands that directly manipulate protected folders
		// which may trigger antivirus, just report the error
		return fmt.Errorf("failed to empty recycle bin: %w", err)
	}
	
	return nil
}

// Helper function to split command output into lines
func splitLines(s string) []string {
	var lines []string
	var line string
	for _, c := range s {
		if c == '\n' || c == '\r' {
			if line != "" {
				lines = append(lines, line)
				line = ""
			}
		} else {
			line += string(c)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

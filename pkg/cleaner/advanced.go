package cleaner

import (
	"os/exec"
	"fmt"
)

// ClearEventLogs clears Windows event logs using the wevtutil command
func ClearEventLogs(verbose bool) error {
	cmd := exec.Command("wevtutil", "el")
	if verbose {
		fmt.Println("[VERBOSE] Running command: wevtutil el")
	}
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
		if verbose {
			fmt.Printf("[VERBOSE] Running command: wevtutil cl %s\n", logName)
		}
		clearCmd.Run() // Ignore errors, as some logs might be protected
	}

	return nil
}

// RunSystemFileChecker runs the Windows System File Checker to repair system files
func RunSystemFileChecker(verbose bool) error {
	if verbose {
		fmt.Println("[VERBOSE] Running command: sfc /scannow")
	}
	cmd := exec.Command("sfc", "/scannow")
	return cmd.Run()
}

// RunDISM runs the Deployment Image Servicing and Management tool to repair Windows image
func RunDISM(verbose bool) error {
	if verbose {
		fmt.Println("[VERBOSE] Running command: DISM /Online /Cleanup-Image /RestoreHealth")
	}
	cmd := exec.Command("DISM", "/Online", "/Cleanup-Image", "/RestoreHealth")
	return cmd.Run()
}

// EmptyRecycleBin empties the Windows Recycle Bin
func EmptyRecycleBin(verbose bool) error {
	if verbose {
		fmt.Println("[VERBOSE] Running command: powershell -Command Clear-RecycleBin -Force")
	}
	// Using PowerShell to clear recycle bin
	cmd := exec.Command("powershell", "-Command", "Clear-RecycleBin", "-Force")
	return cmd.Run()
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

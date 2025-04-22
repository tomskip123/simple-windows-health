package cleaner

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunDiskOptimization runs appropriate optimization based on drive type (defrag for HDDs, TRIM for SSDs)
func RunDiskOptimization(verbose bool) error {
	checkCmd := exec.Command("powershell", "-Command", "Get-PhysicalDisk | Select-Object DeviceId, MediaType | ConvertTo-Json")
	if verbose {
		fmt.Println("[VERBOSE] Running command: powershell -Command Get-PhysicalDisk | Select-Object DeviceId, MediaType | ConvertTo-Json")
	}
	output, err := checkCmd.Output()
	if err != nil {
		return err
	}

	// Parse output to determine drive types
	// Simple check - if any SSD is found, use /O which automatically selects the correct optimization
	if strings.Contains(string(output), "SSD") {
		if verbose {
			fmt.Println("[VERBOSE] Running command: defrag /C /O /U /V")
		}
		// Use /O which will automatically select proper optimization method based on media type
		cmd := exec.Command("defrag", "/C", "/O", "/U", "/V")
		return cmd.Run()
	} else {
		if verbose {
			fmt.Println("[VERBOSE] Running command: defrag /C /D /U /V")
		}
		// Traditional defrag for HDDs
		cmd := exec.Command("defrag", "/C", "/D", "/U", "/V")
		return cmd.Run()
	}
}

// RunCheckDisk runs the Windows Check Disk utility
func RunCheckDisk(verbose bool) error {
	if verbose {
		fmt.Println("[VERBOSE] Running command: chkdsk /f /r /c")
	}
	// Schedule CHKDSK to run on next boot since it requires exclusive access
	cmd := exec.Command("chkdsk", "/f", "/r", "/c")
	return cmd.Run()
}

// FlushDNSCache flushes the Windows DNS resolver cache
func FlushDNSCache(verbose bool) error {
	if verbose {
		fmt.Println("[VERBOSE] Running command: ipconfig /flushdns")
	}
	cmd := exec.Command("ipconfig", "/flushdns")
	return cmd.Run()
}

// RunMemoryDiagnostic runs the Windows Memory Diagnostic tool
func RunMemoryDiagnostic(verbose bool) error {
	if verbose {
		fmt.Println("[VERBOSE] Running command: mdsched")
	}
	cmd := exec.Command("mdsched")
	return cmd.Run()
}

// OptimizePowerConfig optimizes Windows power settings
func OptimizePowerConfig(verbose bool) error {
	if verbose {
		fmt.Println("[VERBOSE] Running command: powercfg /setactive SCHEME_BALANCED")
	}
	// Reset power scheme to balanced
	cmd := exec.Command("powercfg", "/setactive", "SCHEME_BALANCED")
	return cmd.Run()
}

// CleanPrefetch cleans the Windows prefetch directory
func CleanPrefetch(verbose bool) error {
	if verbose {
		fmt.Println("[VERBOSE] Running command: powershell -Command Remove-Item -Path $env:SystemRoot\\Prefetch\\* -Force -ErrorAction SilentlyContinue")
	}
	// Using PowerShell to clean prefetch directory with proper error handling
	cmd := exec.Command("powershell", "-Command",
		"Remove-Item -Path \"$env:SystemRoot\\Prefetch\\*\" -Force -ErrorAction SilentlyContinue")
	return cmd.Run()
}

// ResetNetworkConfig resets Windows network configuration
func ResetNetworkConfig(verbose bool) error {
	var errors []string

	if verbose {
		fmt.Println("[VERBOSE] Running command: netsh winsock reset")
	}
	cmd := exec.Command("netsh", "winsock", "reset")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if verbose {
			errors = append(errors, fmt.Sprintf("netsh winsock reset failed: %v\nOutput: %s", err, string(out)))
		} else {
			errors = append(errors, "Partial success: 'netsh winsock reset' failed. You can try running this command manually in an elevated command prompt: netsh winsock reset")
		}
	}

	if verbose {
		fmt.Println("[VERBOSE] Running command: netsh int ip reset")
	}
	cmd = exec.Command("netsh", "int", "ip", "reset")
	out, err = cmd.CombinedOutput()
	if err != nil {
		if verbose {
			errors = append(errors, fmt.Sprintf("netsh int ip reset failed: %v\nOutput: %s", err, string(out)))
		} else {
			errors = append(errors, "Partial success: 'netsh int ip reset' failed. You can try running this command manually in an elevated command prompt: netsh int ip reset")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

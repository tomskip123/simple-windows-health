package cleaner

import (
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
		return err
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

// RunCheckDisk runs the Windows Check Disk utility
func RunCheckDisk() error {
	// Schedule CHKDSK to run on next boot since it requires exclusive access
	cmd := exec.Command("chkdsk", "/f", "/r", "/c")
	return cmd.Run()
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
	// Reset power scheme to balanced
	cmd := exec.Command("powercfg", "/setactive", "SCHEME_BALANCED")
	return cmd.Run()
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
	cmd := exec.Command("netsh", "winsock", "reset")
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("netsh", "int", "ip", "reset")
	return cmd.Run()
}

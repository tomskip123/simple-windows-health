package cleaner

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// SystemStatus represents the overall system status information
type SystemStatus struct {
	DiskSpace      map[string]DiskInfo
	WindowsVersion string
	LastBootTime   string
}

// DiskInfo contains information about a disk drive
type DiskInfo struct {
	TotalSize   string
	FreeSpace   string
	UsedSpace   string
	UsedPercent string
}

// GetSystemStatus retrieves the current system status
func GetSystemStatus() (*SystemStatus, error) {
	status := &SystemStatus{
		DiskSpace: make(map[string]DiskInfo),
	}

	// Get disk space information
	if err := getDiskSpace(status); err != nil {
		return nil, fmt.Errorf("failed to get disk space: %w", err)
	}

	// Get Windows version
	if ver, err := getWindowsVersion(); err == nil {
		status.WindowsVersion = ver
	}

	// Get last boot time
	if bootTime, err := getLastBootTime(); err == nil {
		status.LastBootTime = bootTime
	}

	return status, nil
}

// getDiskSpace retrieves disk space information for all drives
func getDiskSpace(status *SystemStatus) error {
	cmd := exec.Command("powershell", "-Command",
		"Get-WmiObject -Class Win32_LogicalDisk | Select-Object DeviceID, Size, FreeSpace | ConvertTo-Csv -NoTypeInformation")

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\r\n")
	if len(lines) < 2 {
		return fmt.Errorf("unexpected output format")
	}

	// Skip header line
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		// Parse CSV line
		parts := strings.Split(line, ",")
		if len(parts) >= 3 {
			drive := strings.Trim(parts[0], "\"")
			sizeStr := strings.Trim(parts[1], "\"")
			freeStr := strings.Trim(parts[2], "\"")

			// Skip drives with empty size
			if sizeStr == "" {
				continue
			}

			// Parse size values
			size, err1 := strconv.ParseInt(sizeStr, 10, 64)
			free, err2 := strconv.ParseInt(freeStr, 10, 64)

			if err1 == nil && err2 == nil && size > 0 {
				used := size - free
				usedPercent := float64(used) / float64(size) * 100

				status.DiskSpace[drive] = DiskInfo{
					TotalSize:   formatBytes(float64(size)),
					FreeSpace:   formatBytes(float64(free)),
					UsedSpace:   formatBytes(float64(used)),
					UsedPercent: fmt.Sprintf("%.1f%%", usedPercent),
				}
			}
		}
	}

	return nil
}

// getWindowsVersion retrieves the Windows version
func getWindowsVersion() (string, error) {
	cmd := exec.Command("powershell", "-Command", "(Get-WmiObject -class Win32_OperatingSystem).Caption")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getLastBootTime retrieves the last system boot time
func getLastBootTime() (string, error) {
	cmd := exec.Command("powershell", "-Command",
		"(Get-CimInstance -ClassName Win32_OperatingSystem).LastBootUpTime.ToString('yyyy-MM-dd HH:mm:ss')")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// formatBytes formats bytes to a human-readable string
func formatBytes(bytes float64) string {
	const unit = 1024.0
	if bytes < unit {
		return fmt.Sprintf("%.0f B", bytes)
	}

	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", bytes/div, "KMGTPE"[exp])
}

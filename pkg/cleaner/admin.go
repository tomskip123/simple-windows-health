package cleaner

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// IsAdmin checks if the current process is running with administrator privileges
func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

// RunAsAdmin restarts the current application with elevated privileges using UAC
func RunAsAdmin() error {
	// Get the path to the current executable
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// Prepare the runas command to request elevation
	cmd := exec.Command("powershell.exe", "-Command", "Start-Process", "-FilePath", exe, "-ArgumentList", strings.Join(os.Args[1:], " "), "-Verb", "RunAs")

	// Set up the command to create no window (HideWindow not available on this platform)
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	// Run the elevation command
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Exit the current non-elevated process
	os.Exit(0)
	return nil // This will not execute, but is needed for the return signature
}

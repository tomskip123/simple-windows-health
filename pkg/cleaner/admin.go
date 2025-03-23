package cleaner

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// IsAdmin checks if the current process is running with administrator privileges
func IsAdmin() bool {
	// Use a more reliable method to check for admin privileges
	cmd := exec.Command("powershell", "-Command", 
		"([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	return strings.TrimSpace(string(output)) == "True"
}

// RunAsAdmin restarts the current application with elevated privileges using UAC
func RunAsAdmin() error {
	// Get the path to the current executable
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// Build arguments string, preserving quotes for arguments with spaces
	args := make([]string, 0, len(os.Args[1:]))
	for _, arg := range os.Args[1:] {
		if strings.Contains(arg, " ") {
			args = append(args, "\""+arg+"\"")
		} else {
			args = append(args, arg)
		}
	}

	// Prepare the runas command to request elevation
	cmd := exec.Command("powershell.exe", "-Command", "Start-Process", 
		"-FilePath", exe, 
		"-ArgumentList", strings.Join(args, " "), 
		"-Verb", "RunAs")

	// Set up the command to create no window
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	// Run the elevation command
	err = cmd.Start()
	if err != nil {
		return err
	}

	// Exit the current non-elevated process
	os.Exit(0)
	return nil // This will not execute, but is needed for the return signature
}

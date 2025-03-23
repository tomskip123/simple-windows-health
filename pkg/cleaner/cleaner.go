package cleaner

import (
	"os"
	"os/exec"
	"path/filepath"
)

// RunDiskCleanup executes the Windows built-in Disk Cleanup utility (cleanmgr.exe)
func RunDiskCleanup() error {
	// Using sageset and sagerun with a specific registry key (102)
	// First, set up the configuration with sageset
	setupCmd := exec.Command("cleanmgr", "/sageset:102")
	err := setupCmd.Run()
	if err != nil {
		return err
	}

	// Then run the cleanup with the saved settings
	cmd := exec.Command("cleanmgr", "/sagerun:102")
	return cmd.Run()
}

// CleanTempFiles removes files from Windows temporary directories
func CleanTempFiles() error {
	// Get the Windows temp directory
	tempDir := os.Getenv("TEMP")
	if tempDir == "" {
		// Fallback to typical location if environment variable is not set
		tempDir = filepath.Join(os.Getenv("SYSTEMDRIVE"), "Windows", "Temp")
	}

	// Also clean user temp directory
	userTempDir := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Temp")

	// Clean the system temp directory
	if err := cleanDirectory(tempDir); err != nil {
		return err
	}

	// Clean the user temp directory
	return cleanDirectory(userTempDir)
}

// cleanDirectory removes files from the specified directory
// It skips files that are in use and returns no error in that case
func cleanDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		// Skip removal if it's a directory
		info, err := entry.Info()
		if err != nil {
			// Just log and continue if we can't get file info
			continue
		}

		// Remove files only, not directories
		if !info.IsDir() {
			// Attempt to remove the file, ignore errors for files in use
			os.Remove(path)
		}
	}

	return nil
}

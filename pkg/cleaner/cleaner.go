package cleaner

import (
	"os"
	"os/exec"
	"path/filepath"
	"fmt"
)

// RunDiskCleanup executes the Windows built-in Disk Cleanup utility (cleanmgr.exe)
func RunDiskCleanup(verbose bool) error {
	if verbose {
		fmt.Println("[VERBOSE] Running command: cleanmgr /sageset:102")
	}
	// Using sageset and sagerun with a specific registry key (102)
	// First, set up the configuration with sageset
	setupCmd := exec.Command("cleanmgr", "/sageset:102")
	err := setupCmd.Run()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Println("[VERBOSE] Running command: cleanmgr /sagerun:102")
	}
	// Then run the cleanup with the saved settings
	cmd := exec.Command("cleanmgr", "/sagerun:102")
	return cmd.Run()
}

// CleanTempFiles removes files from Windows temporary directories
func CleanTempFiles(verbose bool) error {
	// Get the Windows temp directory
	tempDir := os.Getenv("TEMP")
	if tempDir == "" {
		// Fallback to typical location if environment variable is not set
		tempDir = filepath.Join(os.Getenv("SYSTEMDRIVE"), "Windows", "Temp")
	}

	// Also clean user temp directory
	userTempDir := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Temp")

	if verbose {
		fmt.Printf("[VERBOSE] Cleaning system temp directory: %s\n", tempDir)
	}
	// Clean the system temp directory
	if err := cleanDirectory(tempDir, verbose); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("[VERBOSE] Cleaning user temp directory: %s\n", userTempDir)
	}
	// Clean the user temp directory
	return cleanDirectory(userTempDir, verbose)
}

// cleanDirectory removes files from the specified directory
// It skips files that are in use and returns no error in that case
func cleanDirectory(dir string, verbose bool) error {
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
			if verbose {
				fmt.Printf("[VERBOSE] Could not get info for %s: %v\n", path, err)
			}
			continue
		}

		// Remove files only, not directories
		if !info.IsDir() {
			if verbose {
				fmt.Printf("[VERBOSE] Removing file: %s\n", path)
			}
			// Attempt to remove the file, ignore errors for files in use
			os.Remove(path)
		}
	}

	return nil
}

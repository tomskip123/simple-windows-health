package cleaner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	var errors []string

	// Get the Windows temp directory
	tempDir := os.Getenv("TEMP")
	if tempDir == "" {
		// Fallback to typical location if environment variable is not set
		tempDir = filepath.Join(os.Getenv("SYSTEMDRIVE"), "Windows", "Temp")
	}

	// User temp directory
	userTempDir := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Temp")
	
	// Additional temporary folders
	winTempDir := filepath.Join(os.Getenv("SYSTEMDRIVE"), "Windows", "Temp")
	prefetchDir := filepath.Join(os.Getenv("SYSTEMDRIVE"), "Windows", "Prefetch")
	
	// Windows Update cleanup folder
	softwareDistDir := filepath.Join(os.Getenv("SYSTEMDRIVE"), "Windows", "SoftwareDistribution", "Download")

	// Clean each directory
	dirs := []string{tempDir, userTempDir, winTempDir, prefetchDir, softwareDistDir}
	
	for _, dir := range dirs {
		err := cleanDirectory(dir)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error cleaning %s: %v", dir, err))
		}
	}

	// If we have errors, combine them
	if len(errors) > 0 {
		return fmt.Errorf("errors while cleaning temporary files: %s", strings.Join(errors, "; "))
	}

	return nil
}

// cleanDirectory removes files from the specified directory
// It skips files that are in use and returns no error in that case
func cleanDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("could not read directory %s: %w", dir, err)
	}

	var failedCount int
	var totalCount int
	var permissionDenied int

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		// Get file info
		info, err := entry.Info()
		if err != nil {
			// Just continue if we can't get file info
			continue
		}

		if info.IsDir() {
			// Recursively clean subdirectories
			cleanDirectory(path)
		} else {
			totalCount++
			// Attempt to remove the file, ignore errors for files in use
			err := os.Remove(path)
			if err != nil {
				if os.IsPermission(err) {
					permissionDenied++
				} else {
					failedCount++
				}
			}
		}
	}

	// Only print summary information instead of individual files
	if failedCount > 0 {
		fmt.Printf("%d/%d files in %s could not be removed (in use by other processes)\n", 
			failedCount, totalCount, dir)
	}
	if permissionDenied > 0 {
		fmt.Printf("%d/%d files in %s could not be removed due to permission denied\n", 
			permissionDenied, totalCount, dir)
	}

	return nil
}

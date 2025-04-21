package core

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/user/windows_health/pkg/cleaner"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

// ConfigData holds settings from the YAML config file
// default_ops: list of command names to run by default
// log_file: path to the log file
// timeout: global timeout for operations
// timeouts: per-operation timeout overrides
// json_output: toggle JSON output mode for supported commands
type ConfigData struct {
	DefaultOps []string                 `yaml:"default_ops"`
	LogFile    string                   `yaml:"log_file"`
	Timeout    time.Duration            `yaml:"timeout"`
	Timeouts   map[string]time.Duration `yaml:"timeouts"`
	JSONOutput bool                     `yaml:"json_output"`
}

var (
	// ConfigFile is the path to the YAML config; set via --config
	ConfigFile string

	// Config is populated by LoadConfig
	Config ConfigData

	// Logger is the global log target; set by SetupLogger
	Logger *logrus.Logger
)

// LoadConfig reads the YAML config (if present) into Config
func LoadConfig() {
	if ConfigFile == "" {
		ConfigFile = "wincleaner.yaml"
	}
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		// No config file; silent fallback to defaults
		return
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config file: %v\n", err)
	}
}

// SetupLogger initializes the structured logger with logrus and lumberjack for log rotation
func SetupLogger() {
	logPath := Config.LogFile
	if logPath == "" {
		logPath = "wincleaner.log"
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})
	logger.SetLevel(logrus.InfoLevel)
	Logger = logger
}

// Replace existing RunOperation with a unified context-aware runner supporting optional timeout
func RunOperation(ctx context.Context, name string, operation func() error, timeout time.Duration) {
	// apply config overrides for operation timeouts
	if t, ok := Config.Timeouts[name]; ok {
		timeout = t
	} else if timeout == 0 && Config.Timeout > 0 {
		timeout = Config.Timeout
	}
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	start := time.Now()
	fmt.Printf("Running %s...\n", name)
	Logger.Infof("Running %s...", name)
	done := make(chan error, 1)
	go func() {
		done <- operation()
	}()

	if timeout > 0 {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case err := <-done:
				if err != nil {
					fmt.Printf("Error running %s: %v\n", name, err)
					Logger.Errorf("Error running %s: %v", name, err)
				} else {
					fmt.Printf("%s completed successfully.\n", name)
					Logger.Infof("%s completed successfully.", name)
				}
				return
			case <-ctx.Done():
				fmt.Printf("\nOperation %s canceled: %v\n", name, ctx.Err())
				Logger.Infof("Operation %s canceled: %v", name, ctx.Err())
				return
			case <-ticker.C:
				elapsed := time.Since(start).Truncate(time.Second)
				fmt.Printf("%s: %v elapsed...\r", name, elapsed)
			}
		}
	}

	// No timeout: simple execution
	err := <-done
	if err != nil {
		fmt.Printf("Error running %s: %v\n", name, err)
		Logger.Errorf("Error running %s: %v", name, err)
	} else {
		fmt.Printf("%s completed successfully.\n", name)
		Logger.Infof("%s completed successfully.", name)
	}
}

// Simplify RunAllOperations to use unified RunOperation
func RunAllOperations(ctx context.Context) {
	fmt.Println("Running all cleaning operations...")
	Logger.Info("Running all cleaning operations...")
	ops := []struct {
		name    string
		op      func() error
		timeout time.Duration
	}{
		{"Disk Cleanup", cleaner.RunDiskCleanup, 0},
		{"Temporary Files Cleaning", cleaner.CleanTempFiles, 0},
		{"Event Logs Clearing", cleaner.ClearEventLogs, 0},
		{"System File Checker", cleaner.RunSystemFileChecker, 120 * time.Second},
		{"DISM Windows Image Repair", cleaner.RunDISM, 180 * time.Second},
		{"Empty Recycle Bin", cleaner.EmptyRecycleBin, 0},
		{"Disk Optimization", cleaner.RunDiskOptimization, 0},
		{"Check Disk", cleaner.RunCheckDisk, 90 * time.Second},
		{"Flush DNS Cache", cleaner.FlushDNSCache, 0},
		{"Windows Memory Diagnostic", cleaner.RunMemoryDiagnostic, 0},
		{"Clean Prefetch Cache", cleaner.CleanPrefetch, 0},
		{"Optimize Power Configuration", cleaner.OptimizePowerConfig, 0},
		{"Reset Network Configuration", cleaner.ResetNetworkConfig, 0},
	}
	for _, op := range ops {
		RunOperation(ctx, op.name, op.op, op.timeout)
	}
	fmt.Println("All cleaning operations completed.")
	Logger.Info("All cleaning operations completed.")
}

// DisplaySystemStatus formats and prints the system status info
func DisplaySystemStatus(status *cleaner.SystemStatus) {
	// JSON output mode
	if Config.JSONOutput {
		fmt.Println("\n=== System Status Information ===")
		fmt.Printf("Windows Version: %s\n", status.WindowsVersion)
		fmt.Printf("Last Boot Time: %s\n", status.LastBootTime)
		fmt.Println("\nDisk Space Information:")
		fmt.Println("------------------------")
		for drive, info := range status.DiskSpace {
			fmt.Printf("Drive %s:\n", drive)
			fmt.Printf("  Total Size: %s\n", info.TotalSize)
			fmt.Printf("  Free Space: %s\n", info.FreeSpace)
			fmt.Printf("  Used Space: %s (%s)\n", info.UsedSpace, info.UsedPercent)
			fmt.Println()
		}
	}
}

# Windows Health Cleaner

A comprehensive Windows system maintenance utility that helps keep your Windows system in optimal condition.

## Features

- **Disk Cleanup**: Run Windows built-in disk cleanup utility
- **Temporary Files Cleaning**: Remove temporary files from Windows directories
- **Event Logs Clearing**: Clear Windows event logs
- **System File Checker**: Run SFC to scan and repair Windows system files
- **DISM Repair**: Run DISM to repair Windows image
- **Empty Recycle Bin**: Clear the Windows Recycle Bin
- **Disk Optimization**: Smart optimization based on drive type (defrag for HDDs, TRIM for SSDs)
- **Check Disk**: Run CHKDSK to scan and repair disk errors
- **Flush DNS Cache**: Clear Windows DNS resolver cache
- **Memory Diagnostic**: Run Windows Memory Diagnostic tool
- **Clean Prefetch Cache**: Clean Windows prefetch directory
- **Power Configuration**: Optimize power settings
- **Network Reset**: Reset Windows network configuration

## Usage

```bash
wincleaner [flags] <command> [arguments]
```

### Flags

- `-h, --help`: Show help information
- `--config`: Path to YAML config file; supports advanced settings (default_ops, log_file, timeout, timeouts, json_output)
- `--version`: Display version information

### Configuration File

You can provide a `wincleaner.yaml` configuration file to customize behavior. It supports the following fields:

```yaml
default_ops:
  - disk
  - temp
log_file: wincleaner.log
timeout: 120s
timeouts:
  "Disk Cleanup": 30s
  "Check Disk": 60s
json_output: true
```

Field descriptions:
- `default_ops`: List of operation names to run by default when no subcommand is provided.
- `log_file`: Path to write structured logs.
- `timeout`: Global timeout (Go duration) applied to all operations if no per-operation override is set.
- `timeouts`: Map of individual operation names to Go duration strings to override the global timeout.
- `json_output`: Enable JSON output mode for commands that support it.

### Commands

- `disk`: Run Disk Cleanup utility
- `temp`: Clean temporary files
- `events`: Clear Windows event logs
- `sfc`: Run System File Checker
- `dism`: Run DISM to repair Windows image
- `recycle`: Empty Recycle Bin
- `optimize`: Run Disk Optimization (defrag for HDDs, TRIM for SSDs)
- `chkdsk`: Run Check Disk utility
- `flushdns`: Flush DNS resolver cache
- `memcheck`: Run Windows Memory Diagnostic tool
- `prefetch`: Clean Windows prefetch directory
- `power`: Optimize power configuration settings
- `resetnet`: Reset Windows network configuration
- `all`: Run all cleaning operations
- `status`: Display system status information
- `interactive`: Launch interactive console mode
- `admin`: Request administrator privileges

### Examples

```bash
wincleaner interactive      # Launch interactive console mode
wincleaner admin            # Request administrator privileges
wincleaner disk temp        # Run Disk Cleanup and clean temporary files
wincleaner optimize         # Run Disk Optimization
wincleaner chkdsk           # Run Check Disk utility
wincleaner status           # Display system status information
wincleaner all              # Run all cleaning operations
```

## Installation

1. Download the latest release from the [Releases](https://github.com/yourusername/windows_health/releases) page
2. Extract the ZIP file to a location of your choice
3. Run `wincleaner.exe` (use `-admin` flag or right-click and select "Run as administrator" for full functionality)

## Building from Source

```bash
go mod download
go build -o bin/wincleaner.exe ./cmd/wincleaner
```

You can also build using the provided scripts:

```bash
build.bat
build_release.bat
```

## Requirements

- Windows 10 or later
- Administrator privileges for some operations

## License

MIT License 
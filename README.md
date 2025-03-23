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

```
wincleaner [options]
```

### Options

- `-disk`: Run Disk Cleanup utility
- `-temp`: Clean temporary files
- `-events`: Clear Windows event logs
- `-sfc`: Run System File Checker
- `-dism`: Run DISM to repair Windows image
- `-recycle`: Empty Recycle Bin
- `-optimize`: Run Disk Optimization (defrag for HDDs, TRIM for SSDs)
- `-chkdsk`: Run Check Disk utility
- `-flushdns`: Flush DNS resolver cache
- `-memcheck`: Run Windows Memory Diagnostic tool
- `-prefetch`: Clean Windows prefetch directory
- `-power`: Optimize power configuration settings
- `-resetnet`: Reset Windows network configuration
- `-all`: Run all cleaning operations
- `-status`: Display system status information
- `-version`: Display version information
- `-interactive`: Launch interactive console mode
- `-admin`: Request administrator privileges

### Examples

```
wincleaner -interactive      # Launch interactive console mode
wincleaner -admin            # Request administrator privileges
wincleaner -disk -temp       # Run Disk Cleanup and clean temporary files
wincleaner -optimize         # Run Disk Optimization
wincleaner -chkdsk           # Run Check Disk utility
wincleaner -status           # Display system status information
wincleaner -all              # Run all cleaning operations
```

## Installation

1. Download the latest release from the [Releases](https://github.com/yourusername/windows_health/releases) page
2. Extract the ZIP file to a location of your choice
3. Run `wincleaner.exe` (use `-admin` flag or right-click and select "Run as administrator" for full functionality)

## Building from Source

```
go build -o bin/wincleaner.exe cmd/wincleaner/main.go
```

## Requirements

- Windows 10 or later
- Administrator privileges for some operations

## License

MIT License 
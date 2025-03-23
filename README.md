# Windows Health Cleaner

A Go utility for running various Windows cleaning and maintenance operations.

## Features

- Disk Cleanup - Runs the built-in Windows Disk Cleanup utility
- Temporary Files Cleanup - Removes files from Windows temp directories
- Event Log Clearing - Clears Windows event logs
- System File Checker - Runs SFC to repair Windows system files
- DISM Image Repair - Runs DISM to repair the Windows system image
- Recycle Bin Emptying - Empties the Windows Recycle Bin
- System Status - Displays disk space and system information

## Requirements

- Windows 10/11
- Go 1.16 or later (for building from source)

## Installation

### From Source

1. Clone the repository:
```
git clone https://github.com/user/windows_health.git
cd windows_health
```

2. Build the application:
```
go build -o wincleaner.exe ./cmd/wincleaner
```

## Usage

Run the application with one or more flags to specify which operations to perform:

```
wincleaner.exe -disk -temp
```

### Available Options

- `-disk`: Run Disk Cleanup utility
- `-temp`: Clean temporary files
- `-events`: Clear Windows event logs
- `-sfc`: Run System File Checker (requires admin rights)
- `-dism`: Run DISM to repair Windows image (requires admin rights)
- `-recycle`: Empty Recycle Bin
- `-status`: Display system status information
- `-all`: Run all cleaning operations
- `-version`: Display version information

### Administrator Privileges

Some operations (SFC, DISM) require administrator privileges. To run these operations, launch the command prompt or PowerShell as administrator before running the program.

## Examples

```
# Run Disk Cleanup and clean temporary files
wincleaner.exe -disk -temp

# Display system status information
wincleaner.exe -status

# Run all cleaning operations
wincleaner.exe -all

# Show version information
wincleaner.exe -version
```

## Building for Release

To build a release version:

```
go build -ldflags="-s -w" -o wincleaner.exe ./cmd/wincleaner
```

## License

MIT 
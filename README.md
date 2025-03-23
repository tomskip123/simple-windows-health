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
- Interactive Mode - Console-based menu interface for easy navigation
- Admin Elevation - Ability to request and run with administrator privileges

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

2. Build the application using one of the build scripts:
```
build.bat        # Full build with manifest (requires Windows SDK)
build_simple.bat # Simple build without manifest embedding
```

Or build manually:
```
go build -o wincleaner.exe ./cmd/wincleaner
```

## Usage

### Interactive Mode

The easiest way to use the application is in interactive mode, which provides a menu-driven interface:

```
wincleaner.exe -interactive
```

If no command-line flags are provided, interactive mode launches by default.

### Command-Line Options

Run the application with one or more flags to specify which operations to perform:

```
wincleaner.exe -disk -temp
```

### Available Options

- `-interactive`: Launch interactive console mode
- `-admin`: Request administrator privileges
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

Some operations (SFC, DISM) require administrator privileges. You can:

1. Right-click the executable and select "Run as administrator" 
2. Use the `-admin` flag to request elevation: `wincleaner.exe -admin`
3. Use the properly built executable with manifest which will automatically prompt for elevation
4. In interactive mode, select the "Restart with Admin Rights" option

## Examples

```
# Launch interactive mode
wincleaner.exe -interactive

# Request admin privileges and run all cleaning operations
wincleaner.exe -admin -all

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

To include the manifest for UAC elevation, use the provided build script:
```
build.bat
```

## License

MIT 
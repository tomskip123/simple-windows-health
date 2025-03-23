# Windows Health Check Script
# Run this script as Administrator for full functionality

# Function to check if script is running as Administrator
function Test-Admin {
    $currentUser = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    $isAdmin = $currentUser.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    return $isAdmin
}

# Function to format size in GB
function Format-Size {
    param([long]$size)
    return [math]::Round($size/1GB, 2)
}

# Check if running as administrator
$isAdmin = Test-Admin
if (-not $isAdmin) {
    Write-Host "WARNING: This script is not running with administrator privileges." -ForegroundColor Yellow
    Write-Host "Some checks will be limited or unavailable." -ForegroundColor Yellow
    Write-Host "Please consider running this script as administrator for complete diagnostics." -ForegroundColor Yellow
    Write-Host ""
}

# Create results directory if it doesn't exist
$resultsDir = "$env:USERPROFILE\WindowsHealthResults"
$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$resultsFile = "$resultsDir\HealthCheck_$timestamp.txt"

if (-not (Test-Path $resultsDir)) {
    New-Item -ItemType Directory -Path $resultsDir | Out-Null
}

# Start transcript to capture all output
Start-Transcript -Path $resultsFile -Append

# Function to display section headers
function Write-SectionHeader {
    param([string]$title)
    Write-Host "`n============= $title =============" -ForegroundColor Cyan
}

# System Information
Write-SectionHeader "System Information"
Get-ComputerInfo | Select-Object WindowsProductName, WindowsVersion, OsHardwareAbstractionLayer, OsArchitecture, CsManufacturer, CsModel, CsProcessors, CsNumberOfLogicalProcessors, CsNumberOfProcessors, CsTotalPhysicalMemory

# Boot Time
Write-SectionHeader "System Boot Time"
$os = Get-CimInstance Win32_OperatingSystem
$bootTime = $os.LastBootUpTime
$uptime = (Get-Date) - $bootTime
Write-Host "Last Boot Time: $bootTime"
Write-Host "System Uptime: $($uptime.Days) days, $($uptime.Hours) hours, $($uptime.Minutes) minutes"

# CPU Usage
Write-SectionHeader "Current CPU Usage (Top 10 Processes)"
Get-Process | Sort-Object -Property CPU -Descending | Select-Object -First 10 | Format-Table -Property ProcessName, ID, CPU, @{Name="MemoryMB";Expression={[math]::Round($_.WorkingSet / 1MB, 2)}}, @{Name="PrivateMemoryMB";Expression={[math]::Round($_.PrivateMemorySize / 1MB, 2)}} -AutoSize

# Memory Usage
Write-SectionHeader "Memory Status"
$os = Get-CimInstance Win32_OperatingSystem
$totalRAM = $os.TotalVisibleMemorySize / 1MB
$freeRAM = $os.FreePhysicalMemory / 1MB
$usedRAM = $totalRAM - $freeRAM
$ramPercentUsed = [math]::Round(($usedRAM / $totalRAM) * 100, 2)

Write-Host "Total RAM: $([math]::Round($totalRAM, 2)) GB"
Write-Host "Used RAM: $([math]::Round($usedRAM, 2)) GB ($ramPercentUsed%)"
Write-Host "Free RAM: $([math]::Round($freeRAM, 2)) GB"

# Disk Space
Write-SectionHeader "Disk Space"
Get-CimInstance -ClassName Win32_LogicalDisk -Filter "DriveType=3" | Format-Table -Property DeviceID, @{Name='Size(GB)';Expression={Format-Size $_.Size}}, @{Name='FreeSpace(GB)';Expression={Format-Size $_.FreeSpace}}, @{Name='%Free';Expression={[math]::Round(($_.FreeSpace/$_.Size)*100, 2)}} -AutoSize

# Disk Health
Write-SectionHeader "Physical Disk Health"
if ($isAdmin) {
    Get-PhysicalDisk | Format-Table DeviceId, FriendlyName, MediaType, Size, HealthStatus, OperationalStatus -AutoSize
} else {
    Write-Host "Disk health check requires administrator privileges" -ForegroundColor Yellow
}

# Recent System Errors
Write-SectionHeader "Recent System Errors (Last 10)"
try {
    Get-WinEvent -LogName System -MaxEvents 10 | Where-Object {$_.LevelDisplayName -eq "Error"} | 
    Format-Table TimeCreated, Id, LevelDisplayName, @{Name="Message";Expression={$_.Message.Substring(0, [Math]::Min(80, $_.Message.Length)) + "..."}} -AutoSize
} catch {
    Write-Host "Unable to retrieve system errors: $_" -ForegroundColor Yellow
}

# Recent Application Errors
Write-SectionHeader "Recent Application Errors (Last 10)"
try {
    Get-WinEvent -LogName Application -MaxEvents 10 | Where-Object {$_.LevelDisplayName -eq "Error"} | 
    Format-Table TimeCreated, Id, LevelDisplayName, @{Name="Message";Expression={$_.Message.Substring(0, [Math]::Min(80, $_.Message.Length)) + "..."}} -AutoSize
} catch {
    Write-Host "Unable to retrieve application errors: $_" -ForegroundColor Yellow
}

# Windows Update Status
Write-SectionHeader "Windows Update Status"
Write-Host "Recent Windows Updates:"
Get-HotFix | Sort-Object -Property InstalledOn -Descending | Select-Object -First 5 | Format-Table -Property Description, HotFixID, InstalledOn -AutoSize

# Services Status
Write-SectionHeader "Stopped Services (Automatic Start Type)"
Get-Service | Where-Object {$_.Status -eq "Stopped" -and $_.StartType -eq "Automatic"} | Format-Table -Property Name, DisplayName -AutoSize

# Startup Items
Write-SectionHeader "Startup Programs"
Get-CimInstance Win32_StartupCommand | Format-Table -Property Name, Command, Location, User -AutoSize

# Network Connectivity Test
Write-SectionHeader "Network Connectivity"
$testConnection = Test-Connection -ComputerName 8.8.8.8 -Count 4 -ErrorAction SilentlyContinue
if ($testConnection) {
    $averageTime = ($testConnection | Measure-Object -Property ResponseTime -Average).Average
    Write-Host "Internet Connection: ACTIVE"
    Write-Host "Average Response Time: $([math]::Round($averageTime, 2)) ms"
} else {
    Write-Host "Internet Connection: FAILED or UNSTABLE" -ForegroundColor Red
}

# System File Check (requires admin)
Write-SectionHeader "System File Check"
if ($isAdmin) {
    Write-Host "Running System File Check (SFC)..."
    $sfc = Start-Process -FilePath "sfc.exe" -ArgumentList "/scannow" -Wait -PassThru -NoNewWindow
    if ($sfc.ExitCode -eq 0) {
        Write-Host "SFC completed successfully" -ForegroundColor Green
    } else {
        Write-Host "SFC encountered issues (Exit code: $($sfc.ExitCode))" -ForegroundColor Yellow
    }
} else {
    Write-Host "System File Check requires administrator privileges" -ForegroundColor Yellow
}

# Check DISM Health (requires admin)
Write-SectionHeader "DISM Health Check"
if ($isAdmin) {
    Write-Host "Checking Windows component store health..."
    $dism = Start-Process -FilePath "DISM.exe" -ArgumentList "/Online /Cleanup-Image /CheckHealth" -Wait -PassThru -NoNewWindow
    if ($dism.ExitCode -eq 0) {
        Write-Host "DISM health check completed successfully" -ForegroundColor Green
    } else {
        Write-Host "DISM health check encountered issues (Exit code: $($dism.ExitCode))" -ForegroundColor Yellow
    }
} else {
    Write-Host "DISM health check requires administrator privileges" -ForegroundColor Yellow
}

# Driver Issues
Write-SectionHeader "Driver Issues"
if ($isAdmin) {
    $driverIssues = Get-CimInstance Win32_PNPEntity | Where-Object {$_.ConfigManagerErrorCode -ne 0}
    if ($driverIssues) {
        $driverIssues | Format-Table -Property Caption, DeviceID, ConfigManagerErrorCode -AutoSize
    } else {
        Write-Host "No driver issues detected" -ForegroundColor Green
    }
} else {
    Write-Host "Driver check requires administrator privileges" -ForegroundColor Yellow
}

# CPU and Memory Performance (last 5 minutes)
Write-SectionHeader "Performance Metrics (Last 5 Minutes)"
try {
    $cpuCounter = Get-Counter -Counter "\Processor(_Total)\% Processor Time" -SampleInterval 2 -MaxSamples 5
    $cpuAvg = ($cpuCounter.CounterSamples.CookedValue | Measure-Object -Average).Average
    $memCounter = Get-Counter -Counter "\Memory\% Committed Bytes In Use" -SampleInterval 2 -MaxSamples 5
    $memAvg = ($memCounter.CounterSamples.CookedValue | Measure-Object -Average).Average
    
    Write-Host "Average CPU Usage: $([math]::Round($cpuAvg, 2))%"
    Write-Host "Average Memory Usage: $([math]::Round($memAvg, 2))%"
} catch {
    Write-Host "Unable to get performance metrics: $_" -ForegroundColor Yellow
}

# Summary
Write-SectionHeader "Health Check Summary"
Write-Host "Health check completed on $(Get-Date)" -ForegroundColor Green
Write-Host "Results saved to: $resultsFile" -ForegroundColor Green

# Stop transcript
Stop-Transcript

# Open results folder
Invoke-Item $resultsDir 
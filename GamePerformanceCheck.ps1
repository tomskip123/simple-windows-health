# Windows Game Performance Diagnostic Script
# Run this script as Administrator for full functionality

# Function to check if script is running as Administrator
function Test-Admin {
    $currentUser = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    $isAdmin = $currentUser.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    return $isAdmin
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
$resultsFile = "$resultsDir\GamePerformanceCheck_$timestamp.txt"

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
Get-ComputerInfo | Select-Object WindowsProductName, WindowsVersion, OsHardwareAbstractionLayer, OsArchitecture

# GPU Information
Write-SectionHeader "GPU Information"
try {
    $gpuInfo = Get-CimInstance -ClassName Win32_VideoController | Select-Object Name, DriverVersion, CurrentHorizontalResolution, CurrentVerticalResolution, CurrentRefreshRate, AdapterRAM
    $gpuInfo | Format-Table -AutoSize
} catch {
    Write-Host "Unable to retrieve GPU information: $_" -ForegroundColor Yellow
}

# Temperature Monitoring (if OpenHardwareMonitor is installed)
Write-SectionHeader "Temperature Check"
try {
    $ohwPath = "C:\Program Files\OpenHardwareMonitor\OpenHardwareMonitor.exe"
    if (Test-Path $ohwPath) {
        Write-Host "OpenHardwareMonitor found. Please check temperatures manually using OpenHardwareMonitor."
    } else {
        Write-Host "OpenHardwareMonitor not found. Consider installing it to monitor system temperatures." -ForegroundColor Yellow
        Write-Host "Download from: https://openhardwaremonitor.org/downloads/" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Error checking for temperature monitoring software: $_" -ForegroundColor Yellow
}

# CPU Performance and Throttling
Write-SectionHeader "CPU Throttling Check"
try {
    $powerInfo = Get-CimInstance -Namespace root\wmi -ClassName MSAcpi_ThermalZoneTemperature -ErrorAction SilentlyContinue
    if ($powerInfo) {
        $kelvin = $powerInfo.CurrentTemperature / 10
        $celsius = $kelvin - 273.15
        Write-Host "Current CPU temperature: $([math]::Round($celsius, 2))°C"
    } else {
        Write-Host "Unable to retrieve temperature information" -ForegroundColor Yellow
    }
    
    # Check for throttling in event logs
    $throttlingEvents = Get-WinEvent -LogName System -MaxEvents 100 -ErrorAction SilentlyContinue | 
        Where-Object { $_.Id -eq 37 -and $_.ProviderName -eq "Microsoft-Windows-Kernel-Processor-Power" }
    
    if ($throttlingEvents -and $throttlingEvents.Count -gt 0) {
        Write-Host "WARNING: Found $($throttlingEvents.Count) CPU throttling events in recent logs" -ForegroundColor Red
        $throttlingEvents | Format-Table TimeCreated, Id, LevelDisplayName -AutoSize
    } else {
        Write-Host "No CPU throttling events detected in recent logs" -ForegroundColor Green
    }
} catch {
    Write-Host "Unable to check for CPU throttling: $_" -ForegroundColor Yellow
}

# Power Plan Check
Write-SectionHeader "Power Plan"
try {
    $powerPlan = powercfg /GetActiveScheme
    Write-Host $powerPlan
    
    # Check if high performance plan is active
    if ($powerPlan -notlike "*High performance*") {
        Write-Host "TIP: Consider switching to the High performance power plan for gaming" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Unable to retrieve power plan: $_" -ForegroundColor Yellow
}

# Game Mode Status
Write-SectionHeader "Windows Game Mode"
try {
    $gameMode = Get-ItemProperty -Path "HKCU:\Software\Microsoft\GameBar" -Name "AllowAutoGameMode" -ErrorAction SilentlyContinue
    if ($gameMode -and $gameMode.AllowAutoGameMode -eq 1) {
        Write-Host "Windows Game Mode is ENABLED"
    } else {
        Write-Host "Windows Game Mode may be DISABLED" -ForegroundColor Yellow
        Write-Host "Consider enabling Game Mode in Windows Settings > Gaming > Game Mode"
    }
} catch {
    Write-Host "Unable to determine Game Mode status: $_" -ForegroundColor Yellow
}

# Graphics Driver Check
Write-SectionHeader "Graphics Driver Check"
try {
    $graphicsDevices = Get-CimInstance Win32_VideoController
    foreach ($device in $graphicsDevices) {
        $driverDate = [DateTime]::ParseExact($device.DriverDate.Split('.')[0], 'yyyyMMdd', $null)
        $driverAge = (Get-Date) - $driverDate
        
        Write-Host "Graphics Device: $($device.Name)"
        Write-Host "Driver Version: $($device.DriverVersion)"
        Write-Host "Driver Date: $($driverDate.ToShortDateString()) (Age: $($driverAge.Days) days)"
        
        if ($driverAge.Days -gt 180) {
            Write-Host "WARNING: Graphics driver is older than 6 months. Consider updating." -ForegroundColor Yellow
        }
    }
} catch {
    Write-Host "Unable to check graphics driver age: $_" -ForegroundColor Yellow
}

# Background Process Analysis
Write-SectionHeader "High Resource Background Processes"
Get-Process | Where-Object { $_.CPU -gt 5 -or $_.WorkingSet -gt 100MB } | 
    Sort-Object -Property CPU -Descending | 
    Select-Object -First 15 | 
    Format-Table -Property ProcessName, Id, CPU, @{Name="MemoryMB";Expression={[math]::Round($_.WorkingSet / 1MB, 2)}}, @{Name="PriorityClass";Expression={$_.PriorityClass}} -AutoSize

# DPC Latency Check
Write-SectionHeader "DPC Latency Check"
Write-Host "High DPC latency can cause stuttering in games."
Write-Host "Consider running LatencyMon to check for DPC latency issues:"
Write-Host "https://www.resplendence.com/latencymon" -ForegroundColor Yellow

# DirectX Diagnostic
Write-SectionHeader "DirectX Diagnostic"
if ($isAdmin) {
    try {
        $dxdiag = Start-Process -FilePath "dxdiag.exe" -ArgumentList "/t", "$env:TEMP\dxdiag.txt" -Wait -PassThru -NoNewWindow
        if (Test-Path "$env:TEMP\dxdiag.txt") {
            Write-Host "DirectX Diagnostic completed successfully"
            Write-Host "Checking for DirectX errors..."
            
            $dxContent = Get-Content "$env:TEMP\dxdiag.txt"
            $errors = $dxContent | Select-String -Pattern "Error:|WARNING:|Failed"
            
            if ($errors) {
                Write-Host "Found potential DirectX issues:" -ForegroundColor Yellow
                $errors | ForEach-Object { Write-Host $_ -ForegroundColor Yellow }
            } else {
                Write-Host "No DirectX errors detected" -ForegroundColor Green
            }
            
            # Copy the dxdiag output to results
            Copy-Item "$env:TEMP\dxdiag.txt" -Destination "$resultsDir\dxdiag_$timestamp.txt"
            Write-Host "Full DirectX diagnostic saved to: $resultsDir\dxdiag_$timestamp.txt"
        }
    } catch {
        Write-Host "Unable to run DirectX diagnostic: $_" -ForegroundColor Yellow
    }
} else {
    Write-Host "DirectX diagnostic requires administrator privileges" -ForegroundColor Yellow
}

# Game-Specific Recommendations
Write-SectionHeader "Common Game Freeze Fixes"
Write-Host "1. Verify game file integrity through the game launcher (Steam, Epic, etc.)"
Write-Host "2. Check if your antivirus is interfering with the game executable"
Write-Host "3. Disable overlays (Discord, Steam, NVIDIA GeForce Experience)"
Write-Host "4. Try running with specialized tools like ISLC (Intelligent Standby List Cleaner)"
Write-Host "5. For specific games, check PCGamingWiki for known issues and fixes:"
Write-Host "   https://www.pcgamingwiki.com/" -ForegroundColor Yellow

# RAM Check
Write-SectionHeader "Memory Diagnostics"
Write-Host "Consider running Windows Memory Diagnostic tool to check for RAM issues:"
Write-Host "Run 'mdsched.exe' and restart your computer to perform a memory check"

# Summary and Recommendations
Write-SectionHeader "Summary and Recommendations"
Write-Host "Game performance diagnostic completed on $(Get-Date)" -ForegroundColor Green
Write-Host "Results saved to: $resultsFile" -ForegroundColor Green
Write-Host ""
Write-Host "NEXT STEPS:" -ForegroundColor Cyan
Write-Host "1. Review any WARNING messages above"
Write-Host "2. Make sure your system drivers are up to date"
Write-Host "3. Check temperatures during gameplay - overheating can cause freezes"
Write-Host "4. Consider running games with performance monitoring tools like MSI Afterburner"
Write-Host "5. If problems persist, consider running a full Windows Memory Diagnostic"

# Stop transcript
Stop-Transcript

# Open results folder
Invoke-Item $resultsDir 
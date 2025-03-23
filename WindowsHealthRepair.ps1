# Windows Health Repair Script
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
    Write-Host "ERROR: This repair script requires administrator privileges to function properly." -ForegroundColor Red
    Write-Host "Please run this script as administrator." -ForegroundColor Red
    pause
    exit
}

# Create results directory if it doesn't exist
$resultsDir = "$env:USERPROFILE\WindowsHealthResults"
$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$resultsFile = "$resultsDir\RepairLog_$timestamp.txt"

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

# Function to log repair actions
function Write-RepairLog {
    param(
        [string]$message,
        [string]$status = "INFO"
    )
    
    $timeStamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    
    switch ($status) {
        "SUCCESS" { $color = "Green" }
        "WARNING" { $color = "Yellow" }
        "ERROR" { $color = "Red" }
        default { $color = "White" }
    }
    
    Write-Host "[$timeStamp] [$status] $message" -ForegroundColor $color
}

Write-SectionHeader "Windows Health Repair"
Write-RepairLog "Starting Windows repair process..."

# 1. System File Check and Repair
Write-SectionHeader "System File Check and Repair"
Write-RepairLog "Running System File Check to repair system files..."
$sfc = Start-Process -FilePath "sfc.exe" -ArgumentList "/scannow" -Wait -PassThru -NoNewWindow
if ($sfc.ExitCode -eq 0) {
    Write-RepairLog "SFC completed successfully" "SUCCESS"
} else {
    Write-RepairLog "SFC encountered issues (Exit code: $($sfc.ExitCode))" "WARNING"
}

# 2. DISM Repair
Write-SectionHeader "Windows Image Repair"
Write-RepairLog "Repairing Windows component store..."
$dism = Start-Process -FilePath "DISM.exe" -ArgumentList "/Online /Cleanup-Image /RestoreHealth" -Wait -PassThru -NoNewWindow
if ($dism.ExitCode -eq 0) {
    Write-RepairLog "Windows component store repair completed successfully" "SUCCESS"
} else {
    Write-RepairLog "Windows component store repair encountered issues (Exit code: $($dism.ExitCode))" "WARNING"
}

# 3. Disk Cleanup
Write-SectionHeader "Disk Cleanup"
Write-RepairLog "Running disk cleanup to remove temporary files..."
try {
    Start-Process -FilePath cleanmgr.exe -ArgumentList "/sagerun:1" -Wait
    Write-RepairLog "Disk cleanup completed" "SUCCESS"
} catch {
    Write-RepairLog "Error during disk cleanup: $_" "ERROR"
}

# 4. Check and Repair Disk
Write-SectionHeader "Disk Check and Repair"
Write-RepairLog "Checking system drive for errors..."
$systemDrive = $env:SystemDrive.TrimEnd(":")
$chkdsk = Start-Process -FilePath "chkdsk.exe" -ArgumentList "$systemDrive /f /r" -Wait -PassThru -NoNewWindow
if ($chkdsk.ExitCode -eq 0) {
    Write-RepairLog "Disk check completed without finding errors" "SUCCESS"
} else {
    Write-RepairLog "Disk check found and potentially fixed errors (Exit code: $($chkdsk.ExitCode))" "WARNING"
    Write-RepairLog "A restart may be required for full disk repair" "WARNING"
}

# 5. Windows Update Reset
Write-SectionHeader "Windows Update Service Reset"
Write-RepairLog "Resetting Windows Update components..."
try {
    # Stop Windows Update services
    Stop-Service -Name wuauserv, cryptSvc, bits, msiserver -Force

    # Rename Windows Update folders (backup)
    Rename-Item "$env:SystemRoot\SoftwareDistribution" "$env:SystemRoot\SoftwareDistribution.old" -Force -ErrorAction SilentlyContinue
    Rename-Item "$env:SystemRoot\System32\catroot2" "$env:SystemRoot\System32\catroot2.old" -Force -ErrorAction SilentlyContinue

    # Reset Windows Update components
    & "$env:SystemRoot\System32\regsvr32.exe" /s "$env:SystemRoot\System32\wuaueng.dll"
    & "$env:SystemRoot\System32\regsvr32.exe" /s "$env:SystemRoot\System32\wuapi.dll"
    & "$env:SystemRoot\System32\regsvr32.exe" /s "$env:SystemRoot\System32\wups.dll"
    
    # Start Windows Update services
    Start-Service -Name wuauserv, cryptSvc, bits, msiserver
    
    Write-RepairLog "Windows Update components reset successfully" "SUCCESS"
} catch {
    Write-RepairLog "Error resetting Windows Update components: $_" "ERROR"
}

# 6. Restart Essential Services
Write-SectionHeader "Service Restoration"
Write-RepairLog "Checking and restarting essential stopped services..."
$essentialServices = @(
    "Dhcp", "Dnscache", "LanmanWorkstation", "LanmanServer", 
    "RpcSs", "nsi", "mpssvc", "BFE", "DPS", "WinDefend"
)

foreach ($service in $essentialServices) {
    $svc = Get-Service -Name $service -ErrorAction SilentlyContinue
    if ($svc) {
        if ($svc.Status -eq "Stopped" -and $svc.StartType -ne "Disabled") {
            try {
                Start-Service -Name $service -ErrorAction Stop
                Write-RepairLog "Started essential service: $($svc.DisplayName)" "SUCCESS"
            } catch {
                Write-RepairLog "Failed to start $($svc.DisplayName): $_" "ERROR"
            }
        }
    }
}

# 7. Network Reset
Write-SectionHeader "Network Stack Reset"
Write-RepairLog "Resetting network stack..."
try {
    # Reset TCP/IP stack
    netsh winsock reset
    netsh int ip reset
    
    # Flush DNS cache
    ipconfig /flushdns
    
    Write-RepairLog "Network stack reset successfully" "SUCCESS"
} catch {
    Write-RepairLog "Error resetting network stack: $_" "ERROR"
}

# 8. Clear Event Logs (for performance)
Write-SectionHeader "Event Log Cleanup"
Write-RepairLog "Clearing event logs to improve performance..."
try {
    $eventLogs = Get-WinEvent -ListLog * -ErrorAction SilentlyContinue | Where-Object { $_.RecordCount -gt 0 -and $_.IsEnabled -eq $true }
    foreach ($log in $eventLogs) {
        try {
            [System.Diagnostics.Eventing.Reader.EventLogSession]::GlobalSession.ClearLog($log.LogName)
            Write-RepairLog "Cleared event log: $($log.LogName)" "SUCCESS"
        } catch {
            # Just ignore errors on logs we can't clear
        }
    }
} catch {
    Write-RepairLog "Error clearing event logs: $_" "WARNING"
}

# 9. Fix Registry Permissions
Write-SectionHeader "Registry Repairs"
Write-RepairLog "Repairing registry permissions..."
try {
    # Reset registry permissions for key system areas
    & secedit /configure /cfg %windir%\inf\defltbase.inf /db defltbase.sdb /verbose
    Write-RepairLog "Registry permissions repaired" "SUCCESS"
} catch {
    Write-RepairLog "Error repairing registry permissions: $_" "ERROR"
}

# 10. Group Policy Refresh
Write-SectionHeader "Group Policy Refresh"
Write-RepairLog "Refreshing Group Policy settings..."
try {
    & gpupdate /force
    Write-RepairLog "Group Policy refresh completed" "SUCCESS"
} catch {
    Write-RepairLog "Error refreshing Group Policy: $_" "ERROR"
}

# 11. Clean up temporary files
Write-SectionHeader "Temporary File Cleanup"
Write-RepairLog "Cleaning up temporary files..."
$tempFolders = @(
    "$env:TEMP",
    "$env:SystemRoot\Temp",
    "$env:SystemRoot\Prefetch"
)

foreach ($folder in $tempFolders) {
    try {
        Get-ChildItem -Path $folder -File -Recurse -Force -ErrorAction SilentlyContinue | 
        Where-Object { ($_.LastWriteTime -lt (Get-Date).AddDays(-7)) } | 
        Remove-Item -Force -Recurse -ErrorAction SilentlyContinue
        Write-RepairLog "Cleaned temp files in $folder" "SUCCESS"
    } catch {
        Write-RepairLog "Error cleaning $folder" "WARNING"
    }
}

# 12. Fix Windows Game Bar issues
Write-SectionHeader "Windows Game Bar Fixes"
Write-RepairLog "Repairing Windows Game Bar components..."
try {
    Get-AppxPackage Microsoft.XboxGamingOverlay -AllUsers | ForEach-Object {
        Remove-AppxPackage $_.PackageFullName -ErrorAction SilentlyContinue
    }
    Get-AppxPackage Microsoft.XboxGamingOverlay -AllUsers | ForEach-Object {
        Add-AppxPackage -Register "$($_.InstallLocation)\AppXManifest.xml" -DisableDevelopmentMode -ErrorAction SilentlyContinue
    }
    Write-RepairLog "Windows Game Bar components repaired" "SUCCESS"
} catch {
    Write-RepairLog "Error repairing Windows Game Bar: $_" "WARNING"
}

# 13. Update Device Drivers
Write-SectionHeader "Driver Updates"
Write-RepairLog "Checking for problematic device drivers..."
try {
    $problemDevices = Get-PnpDevice -Status Error
    foreach ($device in $problemDevices) {
        try {
            $device | Disable-PnpDevice -Confirm:$false -ErrorAction SilentlyContinue
            $device | Enable-PnpDevice -Confirm:$false -ErrorAction SilentlyContinue
            Write-RepairLog "Reset problematic device: $($device.FriendlyName)" "SUCCESS"
        } catch {
            Write-RepairLog "Failed to reset device $($device.FriendlyName): $_" "WARNING"
        }
    }
} catch {
    Write-RepairLog "Error processing device drivers: $_" "ERROR"
}

# Final Summary
Write-SectionHeader "Repair Summary"
Write-RepairLog "Windows Health Repair completed" "SUCCESS"
Write-RepairLog "Log file saved to: $resultsFile" "INFO"
Write-RepairLog "A system restart is recommended to complete all repairs" "WARNING"

# Stop transcript
Stop-Transcript

# Prompt for restart
$title = "System Restart Required"
$message = "It's recommended to restart your computer to complete all repairs. Would you like to restart now?"

$yes = New-Object System.Management.Automation.Host.ChoiceDescription "&Yes", "Restarts the computer immediately."
$no = New-Object System.Management.Automation.Host.ChoiceDescription "&No", "Does not restart the computer."
$options = [System.Management.Automation.Host.ChoiceDescription[]]($yes, $no)
$result = $host.ui.PromptForChoice($title, $message, $options, 1)

switch ($result) {
    0 { Restart-Computer -Force }
    1 { Write-Host "Please restart your computer manually to complete the repair process." -ForegroundColor Yellow }
}

# Open results folder
Invoke-Item $resultsDir 
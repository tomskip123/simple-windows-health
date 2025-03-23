=============================================
    WINDOWS HEALTH CHECK AND REPAIR TOOLS
=============================================

These scripts will help diagnose and fix common Windows issues that may cause system instability.

--------------------------
INCLUDED FILES:
--------------------------

1. WindowsHealthCheck.ps1 - Diagnostic script
2. RunHealthCheck.bat - Launcher for the diagnostic script
3. WindowsHealthRepair.ps1 - Automatic repair script
4. RunRepair.bat - Launcher for the repair script

--------------------------
RECOMMENDED USAGE:
--------------------------

STEP 1: Run the diagnostic tool first
--------------------------------------
Double-click "RunHealthCheck.bat" and approve the admin prompt.
This script will analyze your system and generate a detailed report without changing anything.
All results are saved to a folder in your user profile: %USERPROFILE%\WindowsHealthResults

STEP 2: Review the diagnostic results
--------------------------------------
Open the results folder and look at the most recent log file.
The diagnostic report will identify issues that may be causing system instability.

STEP 3: Run the repair tool
--------------------------------------
IMPORTANT: Close all applications and save your work before running this tool!
Double-click "RunRepair.bat" and approve the admin prompt.
This script will automatically fix many common Windows issues, including:

- Repairing corrupted system files
- Fixing Windows component issues
- Cleaning up temporary files
- Checking and repairing disk errors
- Resetting Windows Update components
- Restarting essential services
- Resetting network settings
- Clearing event logs
- Repairing registry issues
- Fixing driver problems
- Repairing Windows Game Bar issues (which showed errors in your system)

The repair process may take 10-30 minutes depending on your system.
Your computer may restart when repairs are complete.

--------------------------
IMPORTANT NOTES:
--------------------------

1. Both scripts require administrator privileges to function properly.

2. The repair script will make changes to your system, including:
   - Stopping and starting Windows services
   - Renaming system folders (with backups)
   - Running system cleanup utilities
   - Repairing system components
   - Potentially restarting your computer

3. All actions are logged to text files in: %USERPROFILE%\WindowsHealthResults

4. For best results, run both scripts in Safe Mode.
   
5. If you continue to experience issues after running these tools, consider:
   - Checking hardware diagnostics
   - Updating device drivers manually
   - Performing a Windows Reset as a last resort

--------------------------
SYSTEM REQUIREMENTS:
--------------------------

- Windows 10 or 11
- Administrative privileges
- PowerShell 5.1 or later 
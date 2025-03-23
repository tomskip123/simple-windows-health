@echo off
echo Windows Health Repair
echo ==========================================
echo This will perform AUTOMATIC REPAIRS on your Windows system.
echo.
echo WARNING: This script will modify system settings and restart services.
echo Please save your work before continuing.
echo.
echo IMPORTANT: Administrator privileges are REQUIRED for this repair tool.
echo.
powershell -Command "Start-Process PowerShell -ArgumentList '-ExecutionPolicy Bypass -File \"%~dp0WindowsHealthRepair.ps1\"' -Verb RunAs"
echo.
echo If a UAC prompt appeared and you approved it, the repair process is now running.
echo.
echo Your computer may restart automatically after repairs are complete.
pause 
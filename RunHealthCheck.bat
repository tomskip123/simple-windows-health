@echo off
echo Windows Health Check
echo ==========================================
echo This will run a comprehensive system health check.
echo Results will be saved to your user profile folder.
echo.
echo IMPORTANT: Administrator privileges are required for full functionality.
echo.
powershell -Command "Start-Process PowerShell -ArgumentList '-ExecutionPolicy Bypass -File \"%~dp0WindowsHealthCheck.ps1\"' -Verb RunAs"
echo.
echo If a UAC prompt appeared and you approved it, the health check is now running.
echo Please wait for the process to complete.
pause 
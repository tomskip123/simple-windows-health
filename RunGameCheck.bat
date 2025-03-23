@echo off
title Windows Game Performance Diagnostic
echo Windows Game Performance Diagnostic
echo Checking for administrator privileges...

NET SESSION >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo This script requires administrator privileges.
    echo Please right-click on this batch file and select "Run as administrator".
    pause
    exit
)

echo Running Game Performance Diagnostic with administrator privileges...
PowerShell -NoProfile -ExecutionPolicy Bypass -File "%~dp0GamePerformanceCheck.ps1"

echo.
echo Diagnostic complete. Check the results window for details.
pause 
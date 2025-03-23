@echo off
echo Building Windows Health Cleaner (simple build)...

set GOARCH=amd64
set GOOS=windows

REM Ensure directories exist
if not exist bin mkdir bin

REM Build the executable
go build -o bin/wincleaner.exe cmd/wincleaner/main.go

REM Check if the build was successful
if %ERRORLEVEL% NEQ 0 (
    echo Build failed.
    exit /b %ERRORLEVEL%
)

echo Build completed. Executable is in the bin directory.
echo Use 'bin\wincleaner.exe -interactive' to launch in interactive mode.
echo Note: This build does not embed the manifest file.
echo To request admin rights automatically, use 'bin\wincleaner.exe -admin' instead. 
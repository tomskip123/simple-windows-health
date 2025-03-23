@echo off
echo Building Windows Health Cleaner Release...

set VERSION=1.0.0
set GOARCH=amd64
set GOOS=windows

REM Ensure directories exist
if not exist bin mkdir bin
if not exist release mkdir release

REM Build the executable with version info
go build -ldflags="-s -w" -o bin/wincleaner.exe cmd/wincleaner/main.go

REM Create release zip file
powershell Compress-Archive -Path bin/wincleaner.exe, README.md, LICENSE -DestinationPath release/wincleaner_v%VERSION%.zip -Force

echo Release build completed. Release file is in the release directory:
echo release/wincleaner_v%VERSION%.zip 
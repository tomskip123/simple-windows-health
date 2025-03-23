@echo off
echo Building Windows Health Cleaner...

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

REM Include the manifest file using the Windows Resource Compiler (rc.exe)
echo ^
#define RT_MANIFEST 24 ^
1 RT_MANIFEST "cmd\wincleaner\wincleaner.exe.manifest" > bin\wincleaner.rc

REM Check if the Windows SDK is installed and set up
where /q rc.exe
if %ERRORLEVEL% NEQ 0 (
    echo Windows Resource Compiler (rc.exe) not found.
    echo You need to install the Windows SDK and add it to your PATH.
    echo Manifest will not be embedded.
    goto end
)

REM Compile the resource file
rc.exe bin\wincleaner.rc

REM Check if resource compilation was successful
if %ERRORLEVEL% NEQ 0 (
    echo Resource compilation failed.
    exit /b %ERRORLEVEL%
)

REM Link the resource to the executable
cvtres.exe /machine:x64 /out:bin\wincleaner.res.obj bin\wincleaner.res
link.exe /subsystem:windows /entry:mainCRTStartup /out:bin\wincleaner_with_manifest.exe bin\wincleaner.exe bin\wincleaner.res.obj

REM Check if linking was successful
if %ERRORLEVEL% NEQ 0 (
    echo Linking resources failed.
    echo Continuing with non-manifested executable.
    goto copy_exe
)

:copy_exe
copy /y bin\wincleaner_with_manifest.exe bin\wincleaner.exe >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Warning: Could not update executable with manifest.
    echo The program will still work but may not automatically request admin rights.
) else (
    echo Manifest successfully embedded into executable.
)

:end
echo Build completed. Executable is in the bin directory.
echo Use 'bin\wincleaner.exe -interactive' to launch in interactive mode. 
@echo off
echo Trying to build SilentCast with basic hotkey support...

cd source

echo Installing alternative hotkey library...
go get github.com/moutend/go-hook

echo Building...
go build -ldflags "-s -w" -o ..\silentcast-basic.exe .\cmd\silentcast

if %ERRORLEVEL% EQU 0 (
    echo Build successful!
    echo Run silentcast-basic.exe to test
) else (
    echo Build failed. Using AutoHotkey is recommended.
)

pause
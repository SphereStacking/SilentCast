@echo off
echo Building SilentCast for Windows...

REM 依存関係のダウンロード
echo Downloading dependencies...
go mod download

REM ビルド（フルバージョン）
echo Building full version...
go build -ldflags "-H windowsgui -s -w" -o silentcast-gui.exe ./cmd/silentcast

REM ビルド（コンソール版）
echo Building console version...
go build -ldflags "-s -w" -o silentcast.exe ./cmd/silentcast

echo Build complete!
echo.
echo Files created:
echo - silentcast.exe (console version)
echo - silentcast-gui.exe (GUI version - no console window)
pause
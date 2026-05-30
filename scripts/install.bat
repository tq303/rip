@echo off
setlocal

set ARCH=amd64
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" set ARCH=arm64

set URL=https://github.com/tq303/rip/releases/latest/download/rip-windows-%ARCH%.exe
set DEST=%SystemRoot%\System32\rip.exe

echo Installing rip for windows/%ARCH%...
curl -sL "%URL%" -o "%DEST%"
echo Installed to %DEST%

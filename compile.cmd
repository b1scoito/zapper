@echo off
:: https://go.dev/doc/install/source#environment

set /p TARGET_OS=Target OS: 

set GOOS=%TARGET_OS%

set /p TARGET_ARCH=Target Architecture: 
set GOARCH=%TARGET_ARCH%

:: Strip debug info
go build -ldflags "-s -w" plow/zapper

:: Run garble?

pause
exit /b 0
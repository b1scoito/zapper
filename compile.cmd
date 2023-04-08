@echo off

:: https://go.dev/doc/install/source#environment
set GOOS=windows
set GOARCH=amd64

:: Strip debug info and run garble
garble -literals -tiny -seed=random build -ldflags "-s -w" plow/zapper

exit /b 0
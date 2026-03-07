@echo off
set GOOS=windows
set GOARCH=amd64

echo Compiling System Information Collector (main.go) to system_info.exe...
go build -o system_info.exe main.go

echo Compiling XAMPP Collector (xampp_collector.go) to xampp_collector.exe...
go build -o xampp_collector.exe xampp_collector.go

echo Compiling Patch Collector (patch_collector.go) to patch_collector.exe...
go build -o patch_collector.exe patch_collector.go

echo.
if %ERRORLEVEL% EQU 0 (
    echo Build Successful!
    dir /B *.exe
) else (
    echo Build Failed!
)
pause

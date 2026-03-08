@echo off
set GOOS=windows
set GOARCH=amd64

echo Compiling System Information Collector (main.go) to system_info.exe...
go build -o system_info.exe main.go

echo Compiling XAMPP Collector (xampp_collector.go) to xampp_collector.exe...
go build -o xampp_collector.exe xampp_collector.go

echo Compiling Patch Collector (patch_collector.go) to patch_collector.exe...
go build -o patch_collector.exe patch_collector.go

echo Compiling Firewall Collector (firewall_collector.go) to firewall_collector.exe...
go build -o firewall_collector.exe firewall_collector.go

echo Compiling Archiver (archiver.go) to archiver.exe...
go build -o archiver.exe archiver.go

echo Compiling Unarchiver (unarchiver.go) to unarchiver.exe...
go build -o unarchiver.exe unarchiver.go

echo Compiling Database Manager (db_manager.go) to db_manager.exe...
go build -o db_manager.exe db_manager.go

echo Compiling Dev Tool Collector (dev_tool_collector.go) to dev_tool_collector.exe...
go build -o dev_tool_collector.exe dev_tool_collector.go

echo.
if %ERRORLEVEL% EQU 0 (
    echo Build Successful!
    dir /B *.exe
) else (
    echo Build Failed!
)
pause

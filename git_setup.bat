@echo off
set /p username="Enter your GitHub username: "
set repo_name=win-autit-2026

echo.
echo Checking for existing remote...
git remote remove origin 2>nul

echo Setting up remote for https://github.com/%username%/%repo_name%.git...
git remote add origin https://github.com/%username%/%repo_name%.git

echo.
echo Force pushing your local code to GitHub (this will overwrite remote files)...
git push -u origin main --force

echo.
if %ERRORLEVEL% EQU 0 (
    echo Successfully pushed to GitHub!
) else (
    echo.
    echo [ERROR] Failed to push. 
    echo Please ensure you have created the repository '%repo_name%' on GitHub first.
)
pause

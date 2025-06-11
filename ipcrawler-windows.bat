@echo off
setlocal enabledelayedexpansion

echo ========================================
echo       ipcrawler Windows Launcher
echo ========================================
echo.

REM Check if Docker is installed and running
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker is not installed or not in PATH
    echo.
    echo Please install Docker Desktop for Windows:
    echo https://docs.docker.com/desktop/install/windows/
    echo.
    echo Then restart this script
    echo.
    pause
    exit /b 1
)

REM Check if Docker daemon is running
docker ps >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker Desktop is not running
    echo.
    echo Please start Docker Desktop and try again
    echo.
    pause
    exit /b 1
)

echo ✅ Docker is ready!
docker --version
echo.

REM Check if ipcrawler image exists
for /f %%i in ('docker images -q ipcrawler 2^>nul') do set "IMAGE_ID=%%i"

if "%IMAGE_ID%"=="" (
    echo ℹ️ ipcrawler Docker image not found
    echo.
    echo 🔨 Building ipcrawler Docker image (this may take a few minutes)...
    echo.
    docker build -t ipcrawler .
    if errorlevel 1 (
        echo ❌ Docker build failed
        pause
        exit /b 1
    )
    echo ✅ Docker image built successfully!
    echo.
) else (
    echo ✅ ipcrawler Docker image found
    echo 🚀 Image ready! Starting Docker terminal...
    echo.
)

REM Create results directory if it doesn't exist
if not exist "results" mkdir results

echo 🚀 Starting ipcrawler Docker terminal...
echo.
echo 📋 Available commands once inside:
echo   • ipcrawler --help          (Show help)
echo   • ipcrawler 127.0.0.1       (Test scan)
echo   • ipcrawler target.com      (Scan target)
echo   • ls /scans                 (View results)
echo   • exit                      (Leave container)
echo.
echo 💾 Results will be saved to: %cd%\results\
echo.

REM Run the container interactively with proper volume mounting
docker run -it --rm ^
    -v "%cd%\results:/scans" ^
    -w /opt/ipcrawler ^
    --name ipcrawler-session ^
    ipcrawler bash

echo.
echo 👋 ipcrawler session ended
echo 📁 Check your results in: %cd%\results\
echo.
pause 
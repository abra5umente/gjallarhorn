@echo off
REM Build script for Gjallarhorn on Windows
REM This script builds both frontend and backend for production

echo 🔨 Building Gjallarhorn for production...

REM Check if required tools are installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Go is not installed. Please install Go 1.21 or later.
    exit /b 1
)

where npm >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ❌ npm is not installed. Please install Node.js and npm.
    exit /b 1
)

REM Install frontend dependencies
echo 📦 Installing frontend dependencies...
npm install
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Failed to install frontend dependencies
    exit /b 1
)

REM Build frontend
echo 🎨 Building frontend...
npm run build
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Failed to build frontend
    exit /b 1
)

REM Build backend with embedded frontend
echo 🔧 Building Go backend...
go build -o gjallarhorn.exe .
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Failed to build backend
    exit /b 1
)

echo ✅ Build complete!
echo    Binary: gjallarhorn.exe
echo    Run with: gjallarhorn.exe

@echo off
REM Development script for Gjallarhorn on Windows
REM This script starts both the Go backend and Vite frontend in development mode

echo 🚀 Starting Gjallarhorn development environment...

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

REM Install frontend dependencies if node_modules doesn't exist
if not exist "node_modules" (
    echo 📦 Installing frontend dependencies...
    npm install
    if %ERRORLEVEL% NEQ 0 (
        echo ❌ Failed to install frontend dependencies
        exit /b 1
    )
)

REM Start backend in background
echo 🔧 Starting Go backend on port 8080...
start /B go run .
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Failed to start backend
    exit /b 1
)

REM Wait a moment for backend to start
timeout /t 2 /nobreak >nul

REM Start frontend
echo 🎨 Starting Vite frontend on port 5173...
npm run dev
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Failed to start frontend
    exit /b 1
)

echo ✅ Development environment started!
echo    Backend: http://localhost:8080
echo    Frontend: http://localhost:5173
echo    Press Ctrl+C to stop all servers

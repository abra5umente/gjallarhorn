#!/bin/bash

# Development script for Gjallarhorn
# This script starts both the Go backend and Vite frontend in development mode

set -e

echo "🚀 Starting Gjallarhorn development environment..."

# Check if required tools are installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install Node.js and npm."
    exit 1
fi

# Install frontend dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    npm install
fi

# Start backend in background
echo "🔧 Starting Go backend on port 8080..."
go run . &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Start frontend
echo "🎨 Starting Vite frontend on port 5173..."
npm run dev &
FRONTEND_PID=$!

# Function to cleanup processes on exit
cleanup() {
    echo "🛑 Shutting down development servers..."
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    exit 0
}

# Trap Ctrl+C
trap cleanup SIGINT

echo "✅ Development environment started!"
echo "   Backend: http://localhost:8080"
echo "   Frontend: http://localhost:5173"
echo "   Press Ctrl+C to stop all servers"

# Wait for processes
wait

#!/bin/bash

# Build script for Gjallarhorn
# This script builds both frontend and backend for production

set -e

echo "🔨 Building Gjallarhorn for production..."

# Check if required tools are installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install Node.js and npm."
    exit 1
fi

# Install frontend dependencies
echo "📦 Installing frontend dependencies..."
npm install

# Build frontend
echo "🎨 Building frontend..."
npm run build

# Build backend with embedded frontend
echo "🔧 Building Go backend..."
CGO_ENABLED=0 go build -a -installsuffix cgo -o gjallarhorn .

echo "✅ Build complete!"
echo "   Binary: ./gjallarhorn"
echo "   Run with: ./gjallarhorn"

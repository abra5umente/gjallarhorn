#!/bin/bash

# Docker build script for Gjallarhorn
# This script builds a multi-architecture Docker image

set -e

echo "🐳 Building Gjallarhorn Docker image..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker."
    exit 1
fi

# Check if docker buildx is available
if ! docker buildx version &> /dev/null; then
    echo "❌ Docker Buildx is not available. Please enable Docker Buildx."
    exit 1
fi

# Create buildx builder if it doesn't exist
if ! docker buildx ls | grep -q "gjallarhorn-builder"; then
    echo "🔧 Creating buildx builder..."
    docker buildx create --name gjallarhorn-builder --use
fi

# Use the builder
docker buildx use gjallarhorn-builder

# Build multi-architecture image
echo "🏗️  Building multi-architecture image..."
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --tag gjallarhorn:latest \
    --tag gjallarhorn:$(date +%Y%m%d-%H%M%S) \
    --push \
    .

echo "✅ Docker image built and pushed successfully!"
echo "   Image: gjallarhorn:latest"
echo "   Platforms: linux/amd64, linux/arm64"

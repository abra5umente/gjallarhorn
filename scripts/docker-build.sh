#!/bin/bash

# Docker build script for Gjallarhorn
# This script builds Docker images for different platforms

set -e

echo "🐳 Building Gjallarhorn Docker image..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker."
    exit 1
fi

# Parse arguments
PLATFORM="${1:-local}"
PUSH=false

case "$PLATFORM" in
    local)
        echo "📦 Building for local architecture..."
        docker build -t gjallarhorn:latest .
        echo "✅ Built gjallarhorn:latest for local architecture"
        ;;
    pi|arm64)
        echo "🍓 Building for Raspberry Pi (ARM64)..."

        # Set up QEMU for cross-platform builds
        echo "🔧 Setting up QEMU for cross-platform builds..."
        docker run --rm --privileged multiarch/qemu-user-static --reset -p yes 2>/dev/null || true

        # Create/use buildx builder
        if ! docker buildx ls | grep -q "gjallarhorn-builder"; then
            echo "🔧 Creating buildx builder..."
            docker buildx create --name gjallarhorn-builder --use
        else
            docker buildx use gjallarhorn-builder
        fi

        # Build for ARM64
        docker buildx build \
            --platform linux/arm64 \
            -t gjallarhorn:latest \
            -t gjallarhorn:arm64 \
            --load \
            .

        echo "✅ Built gjallarhorn:latest for ARM64"
        echo "   To export: docker save gjallarhorn:latest | gzip > gjallarhorn-arm64.tar.gz"
        ;;
    push)
        if [ -z "$2" ]; then
            echo "❌ Usage: ./docker-build.sh push <registry/image:tag>"
            echo "   Example: ./docker-build.sh push alexschladetsch/gjallarhorn:latest"
            exit 1
        fi
        IMAGE="$2"
        echo "🚀 Building and pushing multi-arch image to $IMAGE..."

        # Set up QEMU for cross-platform builds
        echo "🔧 Setting up QEMU for cross-platform builds..."
        docker run --rm --privileged multiarch/qemu-user-static --reset -p yes 2>/dev/null || true

        # Create/use buildx builder
        if ! docker buildx ls | grep -q "gjallarhorn-builder"; then
            echo "🔧 Creating buildx builder..."
            docker buildx create --name gjallarhorn-builder --use
        else
            docker buildx use gjallarhorn-builder
        fi

        # Build and push multi-arch
        docker buildx build \
            --platform linux/amd64,linux/arm64 \
            -t "$IMAGE" \
            --push \
            .

        echo "✅ Built and pushed $IMAGE"
        echo "   Platforms: linux/amd64, linux/arm64"
        ;;
    *)
        echo "Usage: ./docker-build.sh [local|pi|push <image>]"
        echo ""
        echo "Options:"
        echo "  local         Build for local architecture (default)"
        echo "  pi, arm64     Build for Raspberry Pi (ARM64)"
        echo "  push <image>  Build multi-arch and push to registry"
        echo ""
        echo "Examples:"
        echo "  ./docker-build.sh                                    # Build for local"
        echo "  ./docker-build.sh pi                                 # Build for Pi"
        echo "  ./docker-build.sh push alexschladetsch/gjallarhorn   # Push to Docker Hub"
        exit 1
        ;;
esac

# Build script for different platforms
#!/bin/bash

echo "🔨 Building Go CRUD API for production..."

# Create build directory
mkdir -p build

# Build for Linux (most common for servers)
echo "Building for Linux..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o build/go-crud-api-linux ./cmd/api

# Build for Windows
echo "Building for Windows..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o build/go-crud-api-windows.exe ./cmd/api

# Build for macOS
echo "Building for macOS..."
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o build/go-crud-api-macos ./cmd/api

echo "✅ Build completed! Binaries are in ./build/ directory"
echo "📦 Binary sizes:"
ls -lh build/

#!/bin/bash

# IPCrawler Development Runner
# This script is for development only - production uses compiled binaries

set -e

echo "🚀 IPCrawler Development Mode"
echo "============================="

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "❌ Error: main.go not found. Please run from the project root directory."
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed. Please install Go first."
    exit 1
fi

echo "📦 Checking dependencies..."
go mod tidy

echo "🔨 Running IPCrawler..."
echo ""

# Run the application with all arguments passed through
go run . "$@"

#!/bin/bash

set -e

BUILD_TYPE=${BUILD_TYPE:-local}

echo "Building raspberryTemperatureSensor..."

if [ "$BUILD_TYPE" = "cross" ]; then
    echo "Cross-compiling for Raspberry Pi 3 (ARM7)..."
    echo "This requires arm-linux-gnueabihf-gcc installed"
    CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -o build/telemetry-server ./cmd/server
else
    echo "Building for current platform..."
    go build -o build/telemetry-server ./cmd/server
fi

echo "Build complete! Binary: build/telemetry-server"

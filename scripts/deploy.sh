#!/bin/bash

set -e

PI_USER=${PI_USER:-pi}
PI_HOST=${PI_HOST:-raspberrypi.local}
PI_PATH=${PI_PATH:-/opt/telemetry-service}
SERVICE_NAME=${SERVICE_NAME:-telemetry-service}

echo "Deploying telemetry service to Raspberry Pi..."

echo "Creating remote directory..."
ssh ${PI_USER}@${PI_HOST} "mkdir -p ${PI_PATH}/{data,static}"

echo "Copying source files..."
ssh ${PI_USER}@${PI_HOST} "mkdir -p ${PI_PATH}/src"
scp -r cmd internal pkg static go.mod go.sum ${PI_USER}@${PI_HOST}:${PI_PATH}/src/

echo "Copying service file..."
scp configs/telemetry-service.service ${PI_USER}@${PI_HOST}:/tmp/

echo "Building on remote Pi..."
ssh ${PI_USER}@${PI_HOST} "
    cd ${PI_PATH}/src
    go mod download 2>/dev/null || true
    go build -o ../telemetry-server ./cmd/server
    cd ${PI_PATH}
    cp -r src/static/* static/
    chmod +x telemetry-server
"

echo "Installing systemd service..."
ssh ${PI_USER}@${PI_HOST} "
    sudo mv /tmp/telemetry-service.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl enable ${SERVICE_NAME}
    sudo systemctl restart ${SERVICE_NAME}
"

echo "Checking service status..."
ssh ${PI_USER}@${PI_HOST} "sudo systemctl status ${SERVICE_NAME} --no-pager"

echo ""
echo "Deploy complete!"
echo "Dashboard: http://${PI_HOST}:8080"

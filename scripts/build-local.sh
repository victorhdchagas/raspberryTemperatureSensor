#!/bin/bash

set -e

APP_NAME=mysensoringo
APP_PATH=/opt/${APP_NAME}
GIT_REPO=$(git config --get remote.origin.url)
CURRENT_DIR=$(pwd)

echo "Building mysensoringo for local deployment..."

echo "Pulling latest changes..."
git pull

echo "Building binary..."
go build -o ${APP_NAME} ./cmd/server

echo "Creating target directory..."
sudo mkdir -p ${APP_PATH}

echo "Copying files..."
sudo cp ${APP_NAME} ${APP_PATH}/
sudo cp -r static ${APP_PATH}/

echo "Setting permissions..."
sudo chmod +x ${APP_PATH}/${APP_NAME}

echo "Creating systemd service file..."
cat <<EOF | sudo tee /etc/systemd/system/${APP_NAME}.service > /dev/null
[Unit]
Description=MySensorIngo - DHT11 Sensor Monitor
After=network.target

[Service]
Type=simple
User=pi
WorkingDirectory=${APP_PATH}
ExecStart=${APP_PATH}/${APP_NAME}
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

echo "Reloading systemd..."
sudo systemctl daemon-reload

echo "Enabling service..."
sudo systemctl enable ${APP_NAME}

echo "Restarting service..."
sudo systemctl restart ${APP_NAME}

echo ""
echo "Build and deploy complete!"
echo "Binary: ${APP_PATH}/${APP_NAME}"
echo "Dashboard: http://$(hostname -I | awk '{print $1}'):8080"
echo ""
echo "Service status:"
sudo systemctl status ${APP_NAME} --no-pager

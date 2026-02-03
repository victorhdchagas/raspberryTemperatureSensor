#!/bin/bash

set -e

export PATH=/usr/local/go/bin:$PATH

APP_NAME=mysensoringo
APP_PATH=/opt/${APP_NAME}

echo "=== Build com GPIO 17 (Solução para erro pin 4) ==="


echo "Pulling latest changes..."
git pull

echo "Checking dependencies..."
if ! command -v go &> /dev/null; then
    echo "❌ Erro: Go não está instalado."
    exit 1
fi

if ! python3 -c "import board; import adafruit_dht" &> /dev/null; then
    echo "❌ Erro: Dependências Python não encontradas."
    echo "   Por favor, execute: pip3 install Adafruit-Blinka adafruit-circuitpython-dht --break-system-packages"
    exit 1
fi

echo "Building binary..."
go build -o ${APP_NAME} ./cmd/server

echo "Creating target directory..."
sudo mkdir -p ${APP_PATH}

echo "Stopping service if running..."
sudo systemctl stop ${APP_NAME} 2>/dev/null || true

echo "Copying files..."
sudo cp ${APP_NAME} ${APP_PATH}/
sudo cp -r static ${APP_PATH}/
sudo cp -r scripts ${APP_PATH}/

echo "Setting permissions..."
sudo chmod +x ${APP_PATH}/${APP_NAME}

echo "Creating systemd service file..."
cat <<EOF | sudo tee /etc/systemd/system/${APP_NAME}.service > /dev/null
[Unit]
Description=MySensorIngo - DHT11 Sensor Monitor (GPIO 17)
After=network.target

[Service]
Type=simple
User=root
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

echo "Starting service..."
sudo systemctl start ${APP_NAME}

echo ""
echo "=== Build e Deploy Complete! ==="
echo ""
echo "⚠️ IMPORTANTE:"
echo "   Mova o fio DATA do DHT11:"
echo "   - De: Pin 7 (GPIO 4)"
echo "   - Para: Pin 11 (GPIO 17)"
echo ""
echo "Binary: ${APP_PATH}/${APP_NAME}"
echo "Dashboard: http://$(hostname -I | awk '{print $1}'):8080"
echo ""
echo "Verificando serviço..."
sleep 3
sudo systemctl status ${APP_NAME} --no-pager

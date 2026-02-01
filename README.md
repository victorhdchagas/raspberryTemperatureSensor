# Telemetry Service - Quick Start

## Quick Deploy (Recommended)

The easiest way to deploy is to compile directly on the Raspberry Pi:

```bash
# From your development machine
bash scripts/deploy.sh
```

This will:
1. Copy source files to the Pi
2. Build the binary on the Pi (requires Go installed)
3. Deploy and start the service

## Manual Deployment

If you prefer to build on the Pi manually:

```bash
# SSH into the Pi
ssh pi@raspberrypi.local

# Clone the project
git clone <your-repo-url> /opt/telemetry-service
cd /opt/telemetry-service

# Install dependencies and build
sudo apt update
sudo apt install -y git gcc build-essential
go build -o telemetry-server ./cmd/server

# Setup service
sudo cp configs/telemetry-service.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable telemetry-service
sudo systemctl start telemetry-service
```

## Access Dashboard

Open your browser and navigate to:
```
http://raspberrypi.local:8080
```

## Cross-Compilation Limitations

This project uses `d2r2/go-dht` which requires CGO for GPIO access. Cross-compilation is not straightforward because:
- Requires `arm-linux-gnueabihf-gcc` compiler
- May need additional cross-compilation setup

**Recommendation:** Build directly on the Raspberry Pi for simplicity.
# raspberryTemperatureSensor

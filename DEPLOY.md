# Raspberry Pi Setup

## Prerequisites

1. **Install Go on Raspberry Pi:**
   ```bash
   wget https://go.dev/dl/go1.21.5.linux-armv6l.tar.gz
   sudo tar -C /usr/local -xzf go1.21.5.linux-armv6l.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   ```

2. **Install system dependencies:**
   ```bash
   sudo apt update
   sudo apt install -y git gcc build-essential
   ```

3. **Clone and build on Pi:**
   ```bash
   cd /opt
   sudo git clone <your-repo-url> telemetry-service
   cd telemetry-service
   go mod download
   go build -o telemetry-server ./cmd/server
   ```

## Installation

1. **Create systemd service file:**
   ```bash
   ssh pi@raspberrypi.local
   sudo mkdir -p /opt/telemetry-service/{data,static}
   ```

2. **Copy service file:**
   ```bash
   exit
   scp configs/telemetry-service.service pi@raspberrypi.local:/tmp/
   ssh pi@raspberrypi.local "sudo mv /tmp/telemetry-service.service /etc/systemd/system/"
   ```

3. **Copy binary and files:**
   ```bash
   ssh pi@raspberrypi.local
   cd /opt/telemetry-service
   # Copy your built binary here
   cp ~/telemetry-server /opt/telemetry-service/
   cp -r static/* /opt/telemetry-service/static/
   sudo chmod +x /opt/telemetry-service/telemetry-server
   ```

4. **Enable and start service:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable telemetry-service
   sudo systemctl start telemetry-service
   ```

5. **Check service status:**
   ```bash
   sudo systemctl status telemetry-service
   ```

6. **View logs:**
   ```bash
   sudo journalctl -u telemetry-service -f
   ```

## Alternative: Cross-compile from x86

If you want to cross-compile from x86 to ARM, you need:

1. **Install cross-compiler:**
   ```bash
   sudo apt install gcc-arm-linux-gnueabihf
   ```

2. **Build with CGO:**
   ```bash
   export BUILD_TYPE=cross
   source configs/deploy.env
   bash scripts/build.sh
   ```

Note: This may require additional setup for CGO cross-compilation.

## Accessing the Dashboard

Once the service is running, access the dashboard at:
```
http://raspberrypi.local:8080
```

Or use the Pi's IP address:
```
http://<pi-ip>:8080
```

## GPIO Pin Configuration

The DHT11 sensor should be connected to GPIO pin 4 (BCM numbering).
If you need to change the pin, update the config in `internal/config/config.go`.

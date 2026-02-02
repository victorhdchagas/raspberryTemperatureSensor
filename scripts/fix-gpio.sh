#!/bin/bash

set -e

echo "=== Solução de Problemas GPIO para MySensorIngo ==="
echo ""
echo "O erro 'failed to export pin 4' ocorre porque o GPIO 4"
echo "conflita com o sistema 1-Wire do Raspberry Pi."
echo ""

echo "=== Testando pinos GPIO disponíveis ==="
echo ""

# Criar script de teste Go
cat > /tmp/test_gpio.go <<'EOF'
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/d2r2/go-dht"
)

func testPin(pin int) bool {
	_, _, _, err := dht.ReadDHTxxWithRetry(dht.DHT11, pin, false, 1)
	if err != nil {
		fmt.Printf("  Pin %d: ❌ ERRO - %v\n", pin, err)
		return false
	}
	fmt.Printf("  Pin %d: ✅ FUNCIONA!\n", pin)
	return true
}

func main() {
	pinsToTest := []int{4, 17, 18, 27, 22}

	fmt.Println("Testando pinos GPIO...")
	for _, pin := range pinsToTest {
		testPin(pin)
		time.Sleep(500 * time.Millisecond)
	}
}
EOF

cd /tmp
go mod init test_gpio 2>/dev/null
go get github.com/d2r2/go-dht 2>/dev/null
go run test_gpio.go

echo ""
echo "=== Recomendação ==="
echo ""
echo "Baseado nos resultados acima, escolha um pino que FUNCIONA."
echo "Recomendação: GPIO 17 (physical pin 11) ou GPIO 18 (physical pin 12)"
echo ""

# Ler pino atual
CURRENT_PIN="4"
if [ -f /opt/mysensoringo/internal/config/config.go ]; then
	CURRENT_PIN=$(grep -o 'GPIO.*"[^"]*"' /opt/mysensoringo/internal/config/config.go | head -1 | sed 's/.*GPIO.*"\([^"]*\)".*/\1/')
fi

echo "=== Configuração Atual ==="
echo "Pino GPIO atual: $CURRENT_PIN"
echo ""

read -p "Deseja mudar para outro pino? (17/18/27/22/n) " choice

case $choice in
	17|18|27|22)
		echo ""
		echo "Mudando pino para GPIO $choice..."
		cd /opt/mysensoringo
		sed -i "s/GPIO.*\"[0-9]*\"/GPIO: \"$choice\"/" internal/config/config.go
		echo "Pino mudado para GPIO $choice"
		echo ""
		echo "Rebuildando..."
		bash scripts/build-local.sh
		;;
	n)
		echo ""
		echo "=== Opção Manual ==="
		echo "1. Desabilitar 1-Wire:"
		echo "   sudo nano /boot/config.txt"
		echo "   # Adicionar ou modificar:"
		echo "   dtoverlay=w1-gpio-pullup,gpiopin=18"
		echo "   # Isso move 1-Wire para GPIO 18, liberando GPIO 4"
		echo "   sudo reboot"
		;;
	*)
		echo "Opção inválida. Tente novamente."
		exit 1
		;;
esac

echo ""
echo "Limpeza completa..."
rm -f /tmp/test_gpio.go

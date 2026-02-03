# Deploy - MySensorIngo

## Visão Geral

Este documento descreve o processo de deploy do sistema de telemetria térmica no Raspberry Pi (ssh alias: raspberry3).

## Pré-requisitos

- Git configurado no ambiente local
- SSH configurado com alias `raspberry3` apontando para o Raspberry Pi
- Repositório clonado no Raspberry Pi em `~/raspberryTemperatureSensor` ou similar
- Script de deploy: `scripts/build-gpio17.sh` (usa GPIO 17 para evitar conflito com 1-Wire)

## Processo de Deploy

### Passo 1: Fazer Push das Alterações para o GitHub

No diretório local do projeto:

```bash
# Adicionar todas as alterações
git add .

# Commit das alterações
git commit -m "mensagem do commit"

# Push para o GitHub
git push
```

### Passo 2: Conectar ao Servidor Raspberry Pi

```bash
ssh raspberry3
```

### Passo 3: Navegar para o Diretório do Projeto

O diretório padrão deve ser `raspberryTemperatureSensor` no home do usuário.

```bash
cd raspberryTemperatureSensor
```

**Verificação:** Se o diretório não existir, o deploy não pode continuar. É necessário que o repositório já tenha sido clonado anteriormente.

### Passo 4: Atualizar o Repositório

```bash
git pull
```

Isso baixa todas as alterações que foram feitas push no Passo 1.

### Passo 5: Executar o Script de Deploy

O script `scripts/build-gpio17.sh` realiza as seguintes operações:

1. Verifica se o repositório está atualizado (`git pull`)
2. Verifica se o Go está instalado
3. Verifica as dependências Python (adafruit-dht)
4. Compila o binário do projeto
5. Para o serviço existente (se estiver rodando)
6. Copia o binário compilado e arquivos estáticos para `/opt/mysensoringo`
7. Configura o serviço systemd
8. Inicia o serviço

**Execução:**

```bash
bash scripts/build-gpio17.sh
```

### Passo 6: Verificar Status do Serviço (Opcional)

Após o script finalizar, verifique se o serviço está rodando corretamente:

```bash
sudo systemctl status mysensoringo
```

O dashboard estará disponível em: `http://<ip-do-raspberry>:8080`

## Detalhes do Script de Deploy (scripts/build-gpio17.sh)

### Variáveis do Script

- `APP_NAME=mysensoringo`: Nome da aplicação e do serviço systemd
- `APP_PATH=/opt/mysensoringo`: Diretório de instalação do binário

### Operações Realizadas

1. **Pull automático** no início do script (redundante com o Passo 4, mas seguro)
2. **Verificação de dependências:**
   - Go instalado
   - Bibliotecas Python: `adafruit-blinka` e `adafruit-circuitpython-dht`
3. **Build:** Compila o binário com `go build -o mysensoringo ./cmd/server`
4. **Parada do serviço:** `sudo systemctl stop mysensoringo` (ignora erros se serviço não existir)
5. **Cópia de arquivos:**
   - Binário compilado → `/opt/mysensoringo/mysensoringo`
   - Diretório `static/` → `/opt/mysensoringo/static/`
   - Diretório `scripts/` → `/opt/mysensoringo/scripts/`
6. **Permissões:** `chmod +x` no binário
7. **Criação do serviço systemd:** Cria arquivo em `/etc/systemd/system/mysensoringo.service`
8. **Reload do systemd:** `sudo systemctl daemon-reload`
9. **Habilitação do serviço:** `sudo systemctl enable mysensoringo`
10. **Inicialização do serviço:** `sudo systemctl start mysensoringo`
11. **Verificação:** Mostra status do serviço após 3 segundos

### Configuração GPIO Importante

O script está configurado para usar **GPIO 17** (Pin 11 físico) em vez de GPIO 4 (Pin 7 físico).

**Motivo:** GPIO 4 conflita com o sistema 1-Wire do Raspberry Pi.

**Conexão do DHT11:**
```
DHT11       Raspberry Pi
---------   ---------------
VCC   -----> Pin 1  (3.3V)
GND   -----> Pin 6  (GND)
DATA  -----> Pin 11 (GPIO 17)
```

### Serviço Systemd

O serviço criado roda como `root` para evitar problemas de permissão de GPIO:

```ini
[Unit]
Description=MySensorIngo - DHT11 Sensor Monitor (GPIO 17)
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/mysensoringo
ExecStart=/opt/mysensoringo/mysensoringo
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

## Troubleshooting

### Erro: "Go não está instalado"

Execute no Raspberry Pi:
```bash
wget https://go.dev/dl/go1.21.5.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-armv6l.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

### Erro: "Dependências Python não encontradas"

Execute no Raspberry Pi:
```bash
pip3 install Adafruit-Blinka adafruit-circuitpython-dht --break-system-packages
```

### Serviço não inicia

Verificar logs:
```bash
sudo journalctl -u mysensoringo -f
```

**IMPORTANTE:** Se houver qualquer erro durante o deploy que não seja solucionado pelas instruções acima, NÃO tente resolver automaticamente. Exija intervenção humana.

### Problemas de GPIO

Se o sensor não estiver funcionando, verificar:
1. O fio DATA está conectado no Pin 11 (GPIO 17)
2. O serviço está rodando como root
3. Não há conflito com outro dispositivo

## Logs e Monitoramento

- **Status do serviço:** `sudo systemctl status mysensoringo`
- **Logs em tempo real:** `sudo journalctl -u mysensoringo -f`
- **Logs completos:** `sudo journalctl -u mysensoringo`

## Acesso ao Dashboard

Após o deploy bem-sucedido, acesse:
```
http://<ip-do-raspberry>:8080
```

Ou use o hostname:
```
http://raspberry3:8080
```

## Notas Importantes

1. **Intervenção Humana:** Se o script falhar por qualquer motivo não documentado acima, NÃO tente resolver. Exija intervenção humana.
2. **Serviço como Root:** O serviço roda como `root` para acessar o GPIO. Não alterar o usuário no systemd sem permissões adequadas.
3. **GPIO 17:** Sempre usar GPIO 17 para evitar conflitos. O script já está configurado para isso.
4. **Backup:** O script sobrescreve arquivos em `/opt/mysensoringo` sem fazer backup. Se necessário, faça backup manual antes do deploy.

## Referências

- Diretório do projeto: `~/raspberryTemperatureSensor` no Raspberry Pi
- Diretório de instalação: `/opt/mysensoringo`
- Script de deploy: `scripts/build-gpio17.sh`
- Nome do serviço: `mysensoringo`
- Porta HTTP: 8080

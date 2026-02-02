# MySensorIngo - Deploy Local

## Deploy no Raspberry Pi (Execução Direta)

Execute este script diretamente no Raspberry Pi:

```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

O script irá:
1. ✅ Fazer `git pull` para atualizar o código
2. ✅ Build do binário com nome `mysensoringo`
3. ✅ Copiar para `/opt/mysensoringo/`
4. ✅ Criar e configurar serviço systemd
5. ✅ Reiniciar o serviço automaticamente

## Comportamento do Sensor

**Importante:** O sensor faz a **primeira leitura imediatamente** ao iniciar o serviço!

- 🟢 **Execução inicial:** Captura imediatamente na inicialização
- ⏱️ **Leituras subsequentes:** A cada 50 minutos (configurável)

## Estrutura Final

```
/opt/mysensoringo/
├── mysensoringo          # Binário executável
├── static/               # Arquivos do dashboard web
├── data/                 # Banco de dados SQLite (criado automaticamente)
└── logs/                # Logs do systemd (systemctl status)
```

## Acesso ao Dashboard

```
http://<ip-do-raspberry>:8080
```

Ou se estiver no próprio Pi:
```
http://localhost:8080
```

## Gerenciamento do Serviço

```bash
# Ver status
sudo systemctl status mysensoringo

# Ver logs em tempo real
sudo journalctl -u mysensoringo -f

# Reiniciar manualmente
sudo systemctl restart mysensoringo

# Parar
sudo systemctl stop mysensoringo

# Iniciar
sudo systemctl start mysensoringo
```

## Primeiro Setup (Instalação)

Se for a primeira vez que está instalando:

```bash
# Criar diretório e clonar repositório
sudo mkdir -p /opt/mysensoringo
cd /opt/mysensoringo
git clone <seu-repositorio> .
git config user.email "pi@raspberrypi"
git config user.name "Pi"

# Instalar Go (se não tiver instalado)
wget https://go.dev/dl/go1.21.5.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-armv6l.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Instalar dependências
sudo apt update
sudo apt install -y git gcc build-essential

# Build e deploy
bash scripts/build-local.sh
```

## Atualizações Futuras

Sempre que quiser atualizar:

```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

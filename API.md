# System Monitor API REST

Uma API REST robusta para monitorar recursos do sistema em tempo real, incluindo CPU e memória.

## 🚀 Instalação e Execução

### Pré-requisitos
- Go 1.25.6 ou superior

### Build
```bash
go build -o monitor
```

### Executar
```bash
./monitor
```

O servidor iniciará na porta `8080` e exibirá os endpoints disponíveis.

## 📍 Endpoints

### 1. Health Check
**Endpoint:** `GET /health`

Verifica se o servidor está ativo e saudável.

**Resposta:**
```json
{
  "status": "healthy",
  "timestamp": "2026-02-05T23:20:49-03:00"
}
```

---

### 2. Stats Completo
**Endpoint:** `GET /stats`

Retorna estatísticas completas de CPU e memória.

**Resposta:**
```json
{
  "timestamp": "2026-02-05T23:20:56-03:00",
  "memory": {
    "total_kb": 7758044,
    "available_kb": 6167312,
    "used_kb": 1590732,
    "usage_percentage": 20.50
  },
  "cpu": {
    "user": 14494,
    "nice": 0,
    "system": 14239,
    "idle": 868949,
    "total": 897682,
    "usage_percentage": 3.20
  }
}
```

---

### 3. Stats de Memória
**Endpoint:** `GET /memory`

Retorna apenas as estatísticas de memória.

**Resposta:**
```json
{
  "timestamp": "2026-02-05T23:21:01-03:00",
  "total_kb": 7758044,
  "available_kb": 6161204,
  "used_kb": 1596840,
  "usage_percentage": 20.58
}
```

**Campos:**
- `total_kb`: Memória total disponível em kilobytes
- `available_kb`: Memória disponível em kilobytes
- `used_kb`: Memória em uso em kilobytes
- `usage_percentage`: Percentual de memória em uso

---

### 4. Stats de CPU
**Endpoint:** `GET /cpu`

Retorna apenas as estatísticas de CPU.

**Resposta:**
```json
{
  "timestamp": "2026-02-05T23:21:01-03:00",
  "user": 14546,
  "nice": 0,
  "system": 14302,
  "idle": 872831,
  "total": 901679,
  "usage_percentage": 3.20
}
```

**Campos:**
- `user`: Tempo em modo usuário
- `nice`: Tempo com nice
- `system`: Tempo em modo kernel
- `idle`: Tempo ocioso
- `total`: Tempo total
- `usage_percentage`: Percentual de CPU em uso

---

## 📊 Exemplos de Uso

### Com curl
```bash
# Health check
curl http://localhost:8080/health

# Todos os stats
curl http://localhost:8080/stats

# Apenas memória
curl http://localhost:8080/memory

# Apenas CPU
curl http://localhost:8080/cpu
```

### Com JavaScript/Fetch
```javascript
async function getStats() {
  const response = await fetch('http://localhost:8080/stats');
  const data = await response.json();
  console.log(data);
}

getStats();
```

### Com Python
```python
import requests

response = requests.get('http://localhost:8080/stats')
data = response.json()
print(data)
```

---

## 🏗️ Estrutura do Projeto

```
system-monitor/
├── main.go          # Ponto de entrada - servidor HTTP
├── api.go           # Handlers dos endpoints REST
├── cpu.go           # Lógica de coleta de stats de CPU
├── memory.go        # Lógica de coleta de stats de memória
├── go.mod           # Módulo Go
├── API.md           # Esta documentação
└── README.MD        # Documentação original
```

---

## 🔧 Configuração

### Alterar Porta
Edite o arquivo `main.go` e modifique a variável `port`:

```go
port := ":9000"  // Mudar de 8080 para 9000
```

---

## 📝 Notas

- Todos os timestamps estão em formato ISO 8601 (RFC3339)
- Memória é reportada em kilobytes (KB)
- CPU é calculada em percentual baseado em jiffy do Linux
- A API é thread-safe e pode receber múltiplas requisições simultâneas

---

## 📦 Dependências

Nenhuma dependência externa. Utiliza apenas a biblioteca padrão do Go.

---

## 📄 Licença

Consulte o arquivo LICENSE para mais informações.

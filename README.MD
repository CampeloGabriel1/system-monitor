# Sentinel Core - System Monitor

O **Sentinel Core** é uma ferramenta de monitoramento de baixo nível para sistemas Linux. Ele extrai métricas de CPU e Memória diretamente do Kernel, processa os dados de forma concorrente utilizando Goroutines e os persiste em um banco de dados relacional para análise histórica.

## Tecnologias
* **Linguagem:** Go (Golang) 1.25.6
* **Banco de Dados:** PostgreSQL 15
* **ORM:** GORM
* **Containerização:** Docker & Docker Compose

## Funcionalidades
* **Extração de Dados:** Leitura direta dos arquivos `/proc/stat` e `/proc/meminfo` do Linux.
* **Alta Performance:** Processamento paralelo de métricas com canais e goroutines para reduzir a latência da API.
* **Persistência:** Registro automático de telemetria no banco de dados a cada requisição de status.
* **Resiliência:** Lógica de *retry* integrada para garantir a conexão com o banco de dados durante a inicialização.

## Como rodar com Docker

Certifique-se de ter o **Docker** e o **Docker Compose** instalados em seu ambiente Linux ou WSL2.

1. **Clone o repositório:**
   ```bash
   git clone [https://https://github.com/CampeloGabriel1/system-monitor]
   cd sentinel-core

2. **Infraestrutura**
   O seguinte comando irá buildar a imagem em Go e subir o banco de dados PostgreSQL:
   ```bash
   docker compose up--build -d

3. **Acesso a API**
   Com o seguintes endpoints:

   GET, /, Discovery endpoint com guia de uso da API.
   GET, /health, Liveness probe para verificar saúde do sistema.
   GET, /stats, Coleta CPU/RAM atual e persiste os dados no banco.
   GET, /history, Retorna as últimas 20 métricas armazenadas no banco.
   GET, /memory, Retorna exclusivamente métricas de memória.
   GET, /cpu, Retorna exclusivamente métricas de CPU.   

   Projeto desenvolvido para fins de portfólio, focando em sistemas distribuídos, observabilidade e conteinerização.
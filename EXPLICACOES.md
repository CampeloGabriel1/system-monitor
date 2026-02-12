# Explicações para iniciantes – System Health Monitor

Este documento explica cada sugestão feita na análise do projeto, em linguagem simples, e o que foi implementado.

---

## 1. Concorrência (goroutines) – “não escolher entre um e outro”

### O que é
No endpoint **GET /stats** precisamos de duas coisas: estatísticas de **memória** e de **CPU**. Cada uma lê um arquivo em disco (`/proc/meminfo` e `/proc/stat`).

**Antes (sequencial):**
1. Buscar memória → esperar terminar  
2. Só depois buscar CPU → esperar terminar  
3. Montar a resposta  

O tempo total é **tempo(memória) + tempo(CPU)**.

**Depois (concorrente):**
1. Iniciar **ao mesmo tempo** a busca de memória e a busca de CPU (duas “tarefas” em paralelo)  
2. Esperar as duas terminarem  
3. Montar a resposta  

O tempo total é praticamente o **maior** dos dois (quem terminar por último), não a soma. Em um monitor de saúde, isso deixa o /stats mais rápido e você não precisa “escolher” priorizar só memória ou só CPU – as duas vêm juntas, mais rápido.

### Como fazemos em Go
- **Goroutine:** uma função que roda “em paralelo” com o resto do programa. Criamos com a palavra-chave `go`: `go funcao()`.
- **Canal (channel):** um canal é como um tubo por onde uma goroutine pode enviar um valor e outra (ou o código principal) pode receber. Usamos para “esperar” o resultado das goroutines.
- No código: duas goroutines rodam ao mesmo tempo (uma chama `GetMemoryStats()`, outra `GetCPUStats()`); cada uma envia o resultado (ou erro) por um canal; o handler só monta a resposta quando tiver os dois resultados.

### Por que é bom para o projeto
- Resposta do `/stats` mais rápida.  
- Você continua expondo memória e CPU juntos, sem ter que escolher um ou outro.  
- É um padrão muito comum em Go para “fazer várias coisas ao mesmo tempo e juntar os resultados”.

---

## 2. Tratar o erro do `json.Encode`

### O que é
`json.NewEncoder(w).Encode(response)` escreve o JSON na conexão HTTP (no `ResponseWriter`). Essa escrita **pode falhar**, por exemplo se o cliente fechou a conexão antes de receber a resposta. Em Go, funções que podem falhar retornam um `error` como último valor de retorno.

### O que estava acontecendo
O código não verificava esse retorno. Se `Encode` falhasse, o erro era ignorado e o cliente poderia receber uma resposta incompleta ou vazia sem que o servidor registrasse o problema.

### O que fazemos agora
Guardamos o retorno em `err` e, se for diferente de `nil`, registramos no log (e, se ainda não tivermos escrito nada útil no `w`, podemos responder com status 500). Assim você passa a **saber** quando a escrita da resposta falhou, o que ajuda a debugar e monitorar.

### Por que é importante
- Boas práticas em Go: sempre tratar erro quando a função retorna.  
- Em produção, logs de falha de escrita ajudam a ver clientes desconectando ou problemas de rede.

---

## 3. Porta por variável de ambiente (PORT)

### O que é
A porta em que o servidor escuta (ex.: 8080) estava “fixa” no código. Em ambientes diferentes (seu PC, servidor, container, cloud) às vezes precisamos usar outra porta sem alterar o código.

### O que fazemos
- Usamos `os.Getenv("PORT")`: lê a variável de ambiente `PORT` do sistema operacional.  
- Se `PORT` estiver definida (ex.: `PORT=3000`), usamos ela.  
- Se não estiver definida, usamos um valor padrão (ex.: 8080).  

Assim você pode rodar assim no terminal:
```bash
PORT=3000 ./monitor
```
e o servidor sobe na porta 3000, sem recompilar.

### Por que é útil
- Padrão em deploy (Heroku, Cloud Run, Kubernetes etc. costumam definir `PORT`).  
- Facilita rodar vários serviços na mesma máquina, cada um em uma porta diferente.

---

## 4. Graceful shutdown (não implementado agora)

### O que é
Quando você encerra o programa (Ctrl+C ou `kill`), o Go por padrão pode parar na hora. **Graceful shutdown** é: ao receber o sinal de parada, o servidor deixa de aceitar **novas** conexões mas espera as requisições **em andamento** terminarem antes de fechar.

### Por que não fizemos ainda
Envolve usar `http.Server` com `Shutdown(ctx)`, sinais do OS (`SIGTERM`, `SIGINT`) e um pouco mais de código. Para um monitor simples e para você estar no início em Go, não é obrigatório no primeiro passo. Vale como próximo passo quando se sentir confortável com o básico.

---

## 5. Uso de CPU com duas leituras (não implementado)

### O que é
O Linux em `/proc/stat` dá **contadores que só aumentam** (tempo total gasto em user, system, idle, etc.). Com **uma** leitura, temos só um “instantâneo”; o percentual que calculamos hoje é mais “instantâneo” e pode variar bastante.

O jeito clássico para um percentual mais estável é:
1. Ler `/proc/stat`  
2. Esperar 1 segundo  
3. Ler de novo  
4. Calcular: (diferença de “usado”) / (diferença de “total”) × 100  

Assim o valor representa “uso médio naquele segundo”.

### Por que não mudamos agora
Já temos concorrência para **não escolher entre memória e CPU** e responder rápido. A “média em 1 segundo” exigiria esperar 1s em cada request de CPU e mudaria um pouco a API (ou teríamos dois endpoints: instantâneo vs média). Para um primeiro estágio do monitor, o valor instantâneo é aceitável; você pode adicionar depois um endpoint como `GET /cpu?interval=1s` que usa duas leituras.

---

## Resumo do que foi implementado

| Melhoria              | Onde     | Benefício principal                          |
|-----------------------|----------|---------------------------------------------|
| Concorrência no /stats| api.go   | Memória e CPU em paralelo; resposta mais rápida |
| Tratar erro do Encode | api.go   | Log e tratamento quando a escrita da resposta falha |
| PORT por env          | main.go  | Configurar porta sem mudar código; bom para deploy |

Se quiser, o próximo passo pode ser: adicionar um endpoint de CPU com média em 1s ou implementar graceful shutdown.

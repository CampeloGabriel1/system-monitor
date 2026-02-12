# Go System Monitor

Um monitor de recursos do sistema em tempo real escrito em **Go**, projetado para rodar nativamente em ambientes **Linux**. Este projeto demonstra a integração entre a linguagem Go e o Kernel Linux através da leitura direta do sistema de arquivos `/proc`.



## Objetivo do Projeto
Desenvolver uma ferramenta CLI (Command Line Interface) leve e eficiente para monitorar a saúde do sistema (CPU e RAM), aplicando boas práticas de arquitetura de software e conceitos de engenharia de sistemas.

## Tecnologias e Conceitos Aplicados

- **Linguagem Go:** Uso de Structs, Métodos, Tipagem Forte e Concorrência (Time/Sleep).
- **Integração com Kernel Linux:** Parsing de arquivos do sistema `/proc/meminfo` e `/proc/stat`.
- **Arquitetura Modular:** Separação de responsabilidades em múltiplos arquivos (`cpu.go`, `memory.go`, `main.go`).
- **UX de Terminal:** Implementação de cores ANSI e interface dinâmica com loop de atualização em tempo real.
- **Git/GitHub:** Versionamento semântico e organização de repositório.

## Funcionalidades
- **Monitoramento de RAM:** Cálculo de memória total, disponível e percentual de uso.
- **Monitoramento de CPU:** Cálculo de carga de processamento via deltas de ticks do processador.
- **Interface Visual:** Alertas coloridos baseados no nível de utilização (Verde, Amarelo e Vermelho).
- **Real-time:** Atualização contínua a cada segundo.

## Como rodar o projeto

1. Certifique-se de estar em um ambiente Linux (ou WSL).
2. Clone o repositório:
   ```bash
   git clone [https://github.com/CampeloGabriel1/system-monitor](https://github.com/CampeloGabriel1/system-monitor)
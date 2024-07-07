# Terminal Chat Server em Go
Este é um projeto experimental de um servidor de chat em terminal desenvolvido em Go, utilizando websockets para comunicação em tempo real.

## Sobre o Projeto
O `terminal_chat_server_golang` foi criado para oferecer uma solução simples e eficiente de chat em ambiente de terminal, explorando os conceitos de comunicação em tempo real utilizando websockets em Go.

## Funcionalidades Principais
- **Comunicação em Tempo Real:** Utiliza websockets para estabelecer conexões bidirecionais entre o servidor e os clientes, permitindo o envio e recebimento de mensagens instantâneas.

- **Gerenciamento de Salas:** Suporta a criação dinâmica de salas de chat, onde os usuários podem se conectar e interagir separadamente.
  
- **Concorrência:** Implementa goroutines para suportar operações concorrentes de forma segura e eficiente, garantindo que várias conexões de clientes possam ser tratadas simultaneamente sem conflitos.
  
- **Broadcast de Mensagens:** As mensagens enviadas por um cliente são distribuídas automaticamente para todos os outros clientes na mesma sala, utilizando channels para garantir a entrega assíncrona e eficiente.

## Objetivos
- **Exploração de Websockets:** Implementação de interação bidirecional instantânea entre clientes e servidor usando a biblioteca `github.com/gorilla/websocket`.
- **Aprendizado de Go:** Foco no aprimoramento e entendimento dos conceitos fundamentais da linguagem Go através do desenvolvimento prático.

## Instalação e Uso
Para iniciar o servidor:

1. Clone o repositório:

```bash
git clone https://github.com/8january/terminal_chat_server_golang.git
cd terminal_chat_server_golang
```

2. Execute o servidor:

```bash
go run main.go
```

O servidor estará disponível em http://localhost:8080 por padrão.

### Download do Executável
Baixe o executável da versão v0.0.1-experimental [aqui](https://github.com/8january/terminal_chat_server_golang/releases/tag/v0.0.1-experimental).

### Contribuição
Este projeto é experimental e não está planejado para desenvolvimento futuro. Contribuições não são esperadas neste momento.

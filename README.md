# Homelab Message Queue

A lightweight, concurrent message broker written in Go for the homelab. This project implements a simple publish-subscribe message queue system using TCP connections, allowing multiple clients to subscribe and receive messages in real-time.

## Features

- **TCP-based Message Broker**: Listens on port 8080 for incoming connections
- **Publish-Subscribe Pattern**: Multiple subscribers can connect and receive messages published by any client
- **Concurrent Connection Handling**: Uses Go goroutines for non-blocking, concurrent client connections
- **Graceful Shutdown**: Handles OS signals (SIGINT, SIGTERM) for clean server termination
- **Efficient Message Broadcasting**: Thread-safe message distribution to all connected subscribers

## Architecture

The project uses a channel-based architecture to manage the message broker:

- **BrokerLoop**: Central goroutine that owns the subscriber list and processes all commands sequentially
- **HandleConnectionV2**: Per-connection handler that manages individual client connections
- **BrokerCommand**: Unified command structure for add subscriber, remove subscriber, and publish operations

This design ensures thread-safety without explicit locks by using a single goroutine to manage the subscriber state.

## Project Structure

```
.
├── main.go                    # Entry point
├── go.mod                     # Go module definition
├── tcp_server/
│   └── tcp_server.go         # TCP server bootstrap and connection loop
├── message_broker/
│   ├── message_broker_v2.go  # Broker loop and connection handler
│   └── message_broker.go     # Legacy broker implementation
└── test-script-pws/
    └── Hello-Server.ps1      # PowerShell test script
```

## Requirements

- Go 1.26.1 or later

## Installation

1. Clone the repository:
```bash
git clone https://github.com/ThanhTNV/homelab-messageq.git
cd homelab-messageq
```

2. Install dependencies (if any):
```bash
go mod download
```

## Usage

### Starting the Server

```bash
go run main.go
```

You should see:
```
Server listening on port 8080...
```

### Connecting Clients

You can connect to the server using any TCP client (telnet, netcat, custom applications, etc.):

```bash
# Example using netcat (nc)
nc localhost 8080
```

Or using telnet:
```bash
telnet localhost 8080
```

### Publishing Messages

Any text sent by a connected client is broadcast to all other subscribers:

```bash
# Terminal 1 (Subscriber 1)
nc localhost 8080
# (waiting to receive messages)

# Terminal 2 (Subscriber 2)
nc localhost 8080
# (waiting to receive messages)

# Terminal 3 (Publisher)
nc localhost 8080
Hello from Publisher!
# Terminals 1 and 2 will receive: Hello from Publisher!
```

## How It Works

1. **Server Bootstrap**: The TCP server starts and listens on port 8080
2. **Broker Initialization**: A BrokerLoop goroutine is spawned to manage subscriptions and message broadcasting
3. **Client Connection**: When a client connects, a new HandleConnectionV2 goroutine is created
4. **Add Subscriber**: The connection sends an "add" command to the broker loop
5. **Message Publishing**: When a client sends data, it's published to all subscribers
6. **Graceful Shutdown**: Sending SIGINT (Ctrl+C) cleanly closes all connections and exits

## Design Patterns

- **Channel-based Concurrency**: Uses Go channels for safe communication between goroutines
- **Command Pattern**: Encapsulates broker operations as BrokerCommand structs
- **Single Writer Principle**: Only the BrokerLoop goroutine modifies the subscriber list

## Future Enhancements

- Topic-based subscriptions (publish to specific topics)
- Message persistence
- Authentication and authorization
- Performance metrics and monitoring
- Support for different protocols (WebSocket, gRPC)
- Configuration file support
- Unit and integration tests

## License

MIT (or your preferred license)

## Author

ThanhTNV

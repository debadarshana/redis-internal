# Redis Internal - TCP Echo Server

A simple TCP echo server implementation in Go that mimics basic Redis-like functionality.

## Features

- 🚀 TCP socket server with configurable host and port
- 🔧 Command-line flag support 
- 👥 Client connection handling with concurrent client counting
- 🔄 Command echo functionality
- 🛡️ Graceful client disconnection handling
- 📡 Single-threaded server (handles one client at a time)

## Project Structure

```
RedisInternal/
├── main.go                      # Entry point and configuration
├── go.mod                       # Go module file
├── server/
│   ├── tcp_echo_server.go       # TCP server implementation
│   └── socket_read_write.go     # Socket I/O utilities
└── README.md                    # This file
```

## Usage

### Running the Server

```bash
# Default configuration (host: 0.0.0.0, port: 7379)
go run main.go

# Custom host and port
go run main.go --host=localhost --port=8080
```

### Available Flags

- `--host`: Server host address (default: "0.0.0.0")
- `--port`: Server port number (default: 7379)

### Testing the Server

Once the server is running, you can test it using various clients:

#### Using curl
```bash
curl localhost:7379
```

#### Using telnet
```bash
telnet localhost 7379
# Type your commands and press Enter
# Type 'quit' or Ctrl+C to disconnect
```

#### Using netcat
```bash
echo "PING" | nc localhost 7379
```

## How It Works

1. **Server Start**: The server listens on the specified host and port
2. **Client Connection**: When a client connects, the server accepts the connection
3. **Command Processing**: The server reads commands from the client
4. **Echo Response**: Each command is echoed back to the client
5. **Disconnection Handling**: When a client disconnects (EOF), the server gracefully closes the connection and waits for new clients

## Example Session

```bash
# Terminal 1: Start the server
$ go run main.go
Starting the NiniDB server...
Listening on 0.0.0.0:7379
Accept connection: 127.0.0.1:54321 concurrent client: 1
command received: PING
client Disconnected 127.0.0.1:54321
Closing the Current connection and ready to accept new client

# Terminal 2: Connect as client
$ telnet localhost 7379
PING
PING
^C
```

## Building

To build the project:

```bash
go build -o redis-internal main.go
./redis-internal --host=0.0.0.0 --port=7379
```

## Requirements

- Go 1.16 or higher
- Network access for TCP connections

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).

## Author

Built with ❤️ using Go

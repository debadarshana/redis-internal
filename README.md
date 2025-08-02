# Redis Internal - TCP Echo Server

A TCP server implementation in Go that handles Redis RESP (Redis Serialization Protocol) format, providing basic echo functionality similar to Redis.

## Features

- 🚀 **TCP Socket Server**: Handles multiple client connections sequentially
- 📡 **RESP Protocol Support**: Parses Redis Serialization Protocol format (Simple Strings)
- 🔧 **Command-line Configuration**: Host and port configuration via flags
- 👥 **Client Connection Management**: Tracks concurrent clients and handles graceful disconnections
- 🔄 **Echo Functionality**: Returns received commands back to clients
- 🛡️ **Error Handling**: Proper connection cleanup and protocol error management
- � **Debug Output**: ASCII code debugging for protocol analysis

## Project Structure

```
redis-internal/
├── main.go                     # Entry point and CLI configuration
├── server/
│   ├── tcp_echo_server.go     # TCP server implementation with connection handling
│   ├── socket_read_write.go   # Socket I/O utilities and command processing
│   └── RESP.go                # Redis RESP protocol parser (Simple Strings support)
├── go.mod                     # Go module definition
├── .gitignore                 # Git ignore rules
└── README.md                  # Project documentation
```

## Installation & Setup

```bash
# Clone the repository
git clone https://github.com/debadarshana/redis-internal.git
cd redis-internal

# Initialize Go module dependencies
go mod tidy

# Build the project
go build -o redis-internal

# Or run directly
go run main.go
```

## Usage

### Starting the Server

```bash
# Default configuration (host: 0.0.0.0, port: 7379)
go run main.go

# Custom host and port
go run main.go --host=localhost --port=8080

# Help and available flags
go run main.go --help
```

### Connecting to the Server

#### Using Redis CLI (Recommended)
```bash
redis-cli -h localhost -p 7379
# Try commands like: PING, HELLO, SET key value
```

#### Using Telnet
```bash
telnet localhost 7379
# Type: +OK
# Press Enter
```

#### Using Netcat
```bash
echo "+PING" | nc localhost 7379
```

#### Using Curl
```bash
curl localhost:7379
```

## RESP Protocol Implementation

Currently supports **Simple Strings only** in RESP format. This is a basic implementation that will be expanded gradually.

### Currently Supported
- **Simple Strings**: `+OK\r\n`, `+PONG\r\n`, `+Hello World\r\n`

### Coming Soon
- Bulk Strings (`$length\r\ndata\r\n`)
- Arrays (`*count\r\n...`)
- Integers (`:number\r\n`)
- Errors (`-message\r\n`)

### Example RESP Data Flow
```
Client sends: +PING\r\n
Server receives: ASCII codes [43, 80, 73, 78, 71, 13, 10]
Server parses: "PING"
Server echoes: "PING"
```

### Debug Mode
The server outputs ASCII codes for received data to help with protocol debugging:
```
ASCII: 80  # P
ASCII: 73  # I  
ASCII: 78  # N
ASCII: 71  # G
```

**Note**: This is an early-stage implementation. RESP protocol support will be expanded incrementally as the project develops.

## Server Architecture

### Connection Flow
1. **Listen**: Server binds to specified host:port
2. **Accept**: Accepts incoming TCP connections (one at a time)
3. **Read**: Receives RESP-formatted data from client
4. **Parse**: Uses RESP parser to extract command string
5. **Echo**: Sends the parsed command back to client
6. **Loop**: Continues reading commands until client disconnects
7. **Cleanup**: Handles EOF and decrements client counter

### Technical Implementation
- **Single-threaded**: Processes one client at a time (no concurrent connections)
- **Blocking I/O**: Uses synchronous socket operations
- **RESP Parsing**: Custom parser for Redis protocol compatibility
- **Memory Management**: Fixed 1KB buffer for command reading
- **Error Recovery**: Handles malformed data and network errors

## Configuration Options

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | `0.0.0.0` | Host address for the server to bind to |
| `--port` | `7379` | Port number for the server (Redis default) |

## Example Session

```bash
# Terminal 1: Start server
$ go run main.go
Starting the NiniDB server...
Listening on 0.0.0.0:7379
Accept connection: 127.0.0.1:54321 concurrent client: 1
ASCII: 80
ASCII: 73
ASCII: 78
ASCII: 71
command received: PING
client Disconnected 127.0.0.1:54321
Closing the Current connection and ready to accept new client

# Terminal 2: Connect with redis-cli
$ redis-cli -h localhost -p 7379
localhost:7379> +PING
"PING"
localhost:7379> +HELLO WORLD
"HELLO WORLD"
```

## Development

### Prerequisites
- **Go**: Version 1.16 or higher
- **Network**: TCP connection capability
- **OS**: Linux, macOS, or Windows

### Code Architecture

#### `main.go`
- Command-line flag parsing
- Server configuration setup
- Application entry point

#### `server/tcp_echo_server.go`
- TCP listener creation and management
- Client connection acceptance
- Connection lifecycle management
- Concurrent client counting

#### `server/socket_read_write.go`
- Socket read/write operations
- Command parsing and response handling
- Buffer management (1KB)

#### `server/RESP.go`
- RESP protocol parser implementation (Simple Strings only)
- `readSimpleString()` function for parsing `+string\r\n` format
- ASCII debugging output for development
- Basic protocol error handling
- **Note**: Will be expanded to support full RESP specification gradually

### Building and Testing

```bash
# Clean build
go clean -cache
go build -v .

# Run server
./redis-internal --host=0.0.0.0 --port=7379

# Test with different clients
redis-cli -h localhost -p 7379
telnet localhost 7379
echo "+TEST" | nc localhost 7379
```

## Limitations & Current Status

- **Early Development**: This is a basic implementation, features will be added incrementally
- **Single Client**: Only handles one client at a time
- **Limited RESP**: Only supports Simple Strings (`+string\r\n`) - other RESP types coming soon
- **No Persistence**: Data is not stored (echo server only)
- **Basic Commands**: No Redis command implementation yet (SET, GET, etc.)
- **Simple Error Handling**: Basic error recovery

**Project Status**: 🚧 **Work in Progress** - This server is in early development and will be enhanced step by step.

## Future Roadmap

- [ ] **Concurrent Connections**: Multi-client support with goroutines
- [ ] **Full RESP Support**: Arrays, Bulk Strings, Integers, Errors
- [ ] **Redis Commands**: Implement SET, GET, DEL, PING, etc.
- [ ] **Data Persistence**: In-memory and disk storage
- [ ] **Configuration File**: YAML/JSON configuration support
- [ ] **Logging**: Structured logging with levels
- [ ] **Metrics**: Connection and command metrics
- [ ] **Tests**: Comprehensive test suite

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/resp-arrays`)
3. Commit your changes (`git commit -m 'Add RESP array support'`)
4. Push to the branch (`git push origin feature/resp-arrays`)
5. Open a Pull Request

### Development Guidelines
- Follow Go conventions and `gofmt`
- Add comments for complex logic
- Handle errors gracefully
- Test with multiple Redis clients

## License

This project is open source and available under the [MIT License](LICENSE).

## Author

**Debadarsh Naparida**
- 📧 Email: debadarshnaparida@yahoo.com
- 🐙 GitHub: [@debadarshana](https://github.com/debadarshana)
- 🔗 Repository: [redis-internal](https://github.com/debadarshana/redis-internal)

---

**Built with ❤️ using Go**

*This project serves as an educational implementation to understand Redis internals, network programming, and protocol parsing in Go.*

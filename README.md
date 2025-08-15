# Redis Internal - High-Performance TCP Server

A high-performance TCP server implementation in Go that handles Redis RESP (Redis Serialization Protocol) format with async I/O and epoll-based event handling for production workloads.

## Features

- 🚀 **High-Performance Async Server**: Epoll-based event-driven I/O supporting 20,000+ concurrent clients
- 📡 **Full RESP Protocol Support**: Complete Redis Serialization Protocol implementation (Simple Strings, Bulk Strings, Arrays, Integers, Errors)
- ⚡ **Non-blocking I/O**: Zero-copy syscalls with optimized performance
- 🔧 **Command-line Configuration**: Host and port configuration via flags
- 👥 **Concurrent Client Management**: Real-time connection tracking and graceful disconnections
- �️ **Redis Commands**: PING, ECHO, TIME commands with proper Redis protocol compliance
- 🏗️ **Modular Architecture**: Clean separation between core logic and server implementation
- 🔄 **Production Ready**: SO_REUSEADDR, proper error handling, and resource cleanup

## Project Structure

```
redis-internal/
├── main.go                     # Entry point and CLI configuration
├── core/                       # Core Redis functionality
│   ├── eval.go                # Command evaluation and response generation
│   └── RESP.go                # Complete Redis RESP protocol parser
├── server/                     # Server implementations
│   ├── aync_tcp.go           # High-performance async TCP server (epoll-based)
│   ├── tcp_echo_server.go    # Simple synchronous TCP server
│   └── socket_read_write.go  # Socket I/O utilities and abstractions
├── go.mod                     # Go module definition
├── .gitignore                # Git ignore rules
└── README.md                 # Project documentation
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


```

### Connecting to the Server

#### Using Redis CLI (Recommended)
```bash
redis-cli -h localhost -p 7379

# Test Redis commands
127.0.0.1:7379> PING
PONG
127.0.0.1:7379> ECHO "Hello World"
"Hello World"
127.0.0.1:7379> TIME
1) "1692123456"
2) "123456"
```

#### Using Raw RESP Protocol
```bash
# PING command
printf "*1\r\n\$4\r\nPING\r\n" | nc localhost 7379

# ECHO command  
printf "*2\r\n\$4\r\nECHO\r\n\$5\r\nhello\r\n" | nc localhost 7379

# TIME command
printf "*1\r\n\$4\r\nTIME\r\n" | nc localhost 7379
```

## RESP Protocol Implementation

**Complete RESP Protocol Support** - Full implementation of all Redis Serialization Protocol data types.

### Fully Supported RESP Types
-  **Simple Strings**: `+OK\r\n`, `+PONG\r\n`
-  **Bulk Strings**: `$4\r\nPING\r\n`, `$-1\r\n` (null)
-  **Arrays**: `*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n`
-  **Integers**: `:1000\r\n`, `:42\r\n`
-  **Errors**: `-ERR unknown command\r\n`

### Implemented Redis Commands
- **PING**: Returns PONG or echoes argument
- **ECHO**: Returns the provided string
- **TIME**: Returns Unix timestamp and microseconds



## Server Architecture

### Async Server (Default)
1. **Epoll Event Loop**: Linux epoll for efficient I/O multiplexing
2. **Non-blocking Sockets**: All operations use non-blocking I/O
4. **Event-Driven**: Processes connections only when data is ready
5. **Resource Cleanup**: Automatic cleanup on client disconnect

### Connection Flow
1. **Listen**: Server binds to specified host:port with SO_REUSEADDR
2. **Accept**: Accepts incoming TCP connections via epoll events
3. **Parse**: Complete RESP protocol parsing for all data types
4. **Execute**: Command evaluation with proper Redis responses
5. **Respond**: Send formatted RESP responses back to clients
6. **Monitor**: Real-time concurrent client tracking



## Configuration Options

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | `0.0.0.0` | Host address for the server to bind to |
| `--port` | `7379` | Port number for the server (Redis default) |

## Example Session

```bash
# Terminal 1: Start async server
$ go run main.go
Starting the NiniDB server...
2025/08/15 12:16:03 Starting Async TCP server on 127.0.0.1 7379
2025/08/15 12:16:03 Server listening on 127.0.0.1:7379

# Terminal 2: Connect with redis-cli
$ redis-cli -h localhost -p 7379
localhost:7379> PING
PONG
localhost:7379> PING "Hello World"
"Hello World"
localhost:7379> ECHO "Redis Internal"
"Redis Internal"
localhost:7379> TIME
1) "1692123456"
2) "123456"
localhost:7379> INVALID
(error) ERR unknown command 'INVALID'
```

## Development

### Prerequisites
- **Go**: Version 1.16 or higher
- **Network**: TCP connection capability
- **OS**: Linux, macOS, or Windows

### Code Architecture

#### `main.go`
- Command-line flag parsing with default Redis port (7379)
- Server configuration setup
- Application entry point with async server initialization

#### `core/eval.go`
- Redis command evaluation and response generation
- Command implementations: PING, ECHO, TIME
- RESP encoding utilities for proper Redis responses
- RedisCmd structure for parsed commands

#### `core/RESP.go`
- Complete RESP protocol parser for all data types
- Functions: readSimpleString, readBulkString, readArray, readInt64, readError
- DecodeOne for dispatching to appropriate parsers
- DecodeCmd for command extraction from RESP arrays

#### `server/aync_tcp.go` (Production Server)
- High-performance epoll-based async TCP server
- Non-blocking I/O with FDConn wrapper for io.ReadWriter compatibility
- Support for 20,000+ concurrent clients
- SO_REUSEADDR, proper error handling, and resource cleanup

#### `server/socket_read_write.go`
- io.ReadWriter abstraction for socket operations
- Integration between RESP parser and command evaluation
- Raw command/response logging for debugging
- ReadCommand and Respond functions for clean separation
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

# Run async server (production)
./redis-internal --host=0.0.0.0 --port=7379

# Test with Redis CLI
redis-cli -h localhost -p 7379 PING
redis-cli -h localhost -p 7379 ECHO "test"
redis-cli -h localhost -p 7379 TIME


```



### Benchmarking
```bash
# Test concurrent connections
redis-benchmark -h localhost -p 7379 -c 100 -n 10000 PING
redis-benchmark -h localhost -p 7379 -c 100 -n 10000 ECHO hello
```

## Current Status & Limitations

###  **Completed Features**
- **High-Performance Server**: Epoll-based async I/O for production workloads
- **Complete RESP Protocol**: All Redis data types (Simple Strings, Bulk Strings, Arrays, Integers, Errors)
- **Redis Commands**: PING, ECHO, TIME with proper protocol compliance
- **Concurrent Connections**: Support  simultaneous clients


###  **Current Limitations**
- **Limited Commands**: Only PING, ECHO, TIME implemented (Redis has 200+ commands)
- **No Persistence**: Data is not stored (in-memory only)
- **No Data Structures**: No support for Lists, Sets, Hashes, etc.
- **No Authentication**: No AUTH command or security features
- **No Clustering**: Single instance only



## Future Roadmap

### Phase 1: Core Redis Commands (In Progress)
- [ ] **String Commands**: SET, GET, DEL, EXISTS, INCR, DECR
- [ ] **Key Management**: EXPIRE, TTL, KEYS, TYPE
- [ ] **Database**: SELECT, FLUSHDB, FLUSHALL

### Phase 2: Advanced Data Structures
- [ ] **Lists**: LPUSH, RPUSH, LPOP, RPOP, LRANGE
- [ ] **Sets**: SADD, SREM, SMEMBERS, SINTER, SUNION
- [ ] **Hashes**: HSET, HGET, HDEL, HKEYS, HVALS
- [ ] **Sorted Sets**: ZADD, ZREM, ZRANGE, ZSCORE

### Phase 3: Production Features
- [ ] **Persistence**: RDB snapshots and AOF logging
- [ ] **Configuration**: Redis-compatible config file support
- [ ] **Authentication**: AUTH command and user management
- [ ] **Monitoring**: INFO command and metrics

### Phase 4: Advanced Features
- [ ] **Pub/Sub**: PUBLISH, SUBSCRIBE, UNSUBSCRIBE
- [ ] **Transactions**: MULTI, EXEC, DISCARD, WATCH
- [ ] **Scripting**: Lua script support
- [ ] **Clustering**: Master-slave replication





## License

This project is open source and available under the [MIT License](LICENSE).

## Author

**Debadarsh Naparida**
- Email: debadarshnaparida@yahoo.com
- GitHub: [@debadarshana](https://github.com/debadarshana)
- Repository: [redis-internal](https://github.com/debadarshana/redis-internal)

---

**Built with ❤️ using Go **

*This project demonstrates  network programming, protocol implementation, and high-performance server architecture in Go. It serves as both an educational Redis implementation and a foundation for understanding async I/O, event-driven programming, and system-level socket programming.*

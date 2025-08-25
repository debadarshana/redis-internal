# Redis Internal - High-Performance TCP Server

A async TCP server implementation in Go that handles Redis RESP (Redis Serialization Protocol) format with async I/O and epoll-based event handling for production workloads.

## Features

- **High-Performance Async Server**: Epoll-based event-driven architecture supporting 20,000+ concurrent connections
- **Complete RESP Protocol Support**: Full Redis Serialization Protocol implementation
- **Redis Command Compatibility**: Core Redis commands (PING, ECHO, TIME, SET, GET, TTL, DEL, EXPIRE)
- **Automatic Key Expiration**: Background auto-deletion of expired keys using Redis-compatible sampling algorithm
- **Memory Management**: Configurable key limits with automatic eviction when limits are reached
- **Flexible Configuration**: JSON config file with command-line overrides for all settings
- **Key Eviction Strategies**: Multiple eviction policies (simple-first, future: LRU, random)
- **Production Ready**: Comprehensive validation, error handling, and logging
- **Non-blocking I/O**: Efficient network operations with proper error handling



## Project Structure

```
redis-internal/
├── main.go                     # Entry point and CLI configuration
├── config.json                 # Configuration file with server settings
├── config/                     # Configuration management
│   └── config.go              # Config loading, validation, and CLI flag handling
├── core/                       # Core Redis functionality
│   ├── eval.go                # Command evaluation and response generation
│   ├── eviction.go            # Key eviction strategies and memory management
│   ├── expire.go              # Auto-deletion and key expiration management
│   ├── store.go               # In-memory key-value store with expiration
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

## Configuration

### Configuration File (`config.json`)

Redis Internal uses a JSON configuration file for server settings. The default `config.json` contains:

```json
{
  "host": "0.0.0.0",
  "port": 7379,
  "keysLimit": 5,
  "evictionStrategy": "simple-first",
  "autoDeleteFrequency": "1s",
  "maxClients": 20000,
  "logLevel": "info"
}
```

### Configuration Options

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `host` | string | `"0.0.0.0"` | Host address to bind server (0.0.0.0 for all interfaces) |
| `port` | int | `7379` | Port number for the server (Redis standard) |
| `keysLimit` | int | `1000` | Maximum number of keys before eviction is triggered |
| `evictionStrategy` | string | `"simple-first"` | Strategy for key eviction (`simple-first`, `lru`, `random`) |
| `autoDeleteFrequency` | string | `"1s"` | How often to run auto-deletion of expired keys |
| `maxClients` | int | `20000` | Maximum number of concurrent client connections |
| `logLevel` | string | `"info"` | Logging level (`debug`, `info`, `warn`, `error`) |

### Command Line Overrides

Any configuration setting can be overridden via command line flags:

```bash
# Override individual settings
./redis-internal --host=127.0.0.1 --port=8080
./redis-internal --keys-limit=100 --eviction=simple-first
./redis-internal --max-clients=50000 --log-level=debug

# Use different config file
./redis-internal --config=production-config.json

# Combine config file with overrides
./redis-internal --config=base-config.json --keys-limit=500
```

### Memory Management & Eviction

#### Key Limit Enforcement
- When `keysLimit` is reached, the eviction strategy is triggered
- Only applies when adding **new** keys (updates to existing keys don't trigger eviction)
- Eviction removes one key before adding the new key

#### Eviction Strategies
- **`simple-first`**: Removes the first key encountered in the map iteration
- **`lru`**: *(Future)* Least Recently Used eviction
- **`random`**: *(Future)* Random key eviction

#### Example Eviction Behavior
```bash
# Set key limit to 3
./redis-internal --keys-limit=3

# In redis-cli:
SET key1 "value1"  # OK - 1 key
SET key2 "value2"  # OK - 2 keys  
SET key3 "value3"  # OK - 3 keys (limit reached)
SET key4 "value4"  # OK - evicts key1, stores key4 (still 3 keys)
```

## Usage

### Starting the Server

```bash
# Use default config.json settings
./redis-internal

# Override specific settings
./redis-internal --host=127.0.0.1 --port=8080 --keys-limit=100

# Use custom config file
./redis-internal --config=my-config.json

# Development mode with debug logging
./redis-internal --log-level=debug --keys-limit=10

# Production mode with high limits
./redis-internal --keys-limit=10000 --max-clients=50000
```

### Server Output
```bash
$ ./redis-internal --keys-limit=5
Starting the Redis Internal server...
Loaded configuration from config.json
=== Redis Internal Configuration ===
Host: 0.0.0.0
Port: 7379
Keys Limit: 5
Eviction Strategy: simple-first
Auto Delete Frequency: 1s
Max Clients: 20000
Log Level: info
===================================
2025/08/25 21:43:10 Starting Async TCP server on 0.0.0.0:7379
2025/08/25 21:43:10 Configuration: MaxClients=20000, KeysLimit=5, EvictionStrategy=simple-first
2025/08/25 21:43:11 Deleted the expired Keys. total keys  0
```


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

# SET command
printf "*3\r\n\$3\r\nSET\r\n\$3\r\nkey\r\n\$5\r\nvalue\r\n" | nc localhost 7379

# SET command with expiration
printf "*5\r\n\$3\r\nSET\r\n\$3\r\nkey\r\n\$5\r\nvalue\r\n\$2\r\nEX\r\n\$2\r\n10\r\n" | nc localhost 7379

# GET command
printf "*2\r\n\$3\r\nGET\r\n\$3\r\nkey\r\n" | nc localhost 7379

# TTL command
printf "*2\r\n\$3\r\nTTL\r\n\$3\r\nkey\r\n" | nc localhost 7379

# DEL command (single key)
printf "*2\r\n\$3\r\nDEL\r\n\$3\r\nkey\r\n" | nc localhost 7379

# DEL command (multiple keys)
printf "*4\r\n\$3\r\nDEL\r\n\$4\r\nkey1\r\n\$4\r\nkey2\r\n\$4\r\nkey3\r\n" | nc localhost 7379

# EXPIRE command
printf "*3\r\n\$6\r\nEXPIRE\r\n\$3\r\nkey\r\n\$2\r\n60\r\n" | nc localhost 7379
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
- **SET**: Store key-value pairs with optional expiration (EX parameter)
- **GET**: Retrieve values by key, returns nil if key doesn't exist or expired
- **TTL**: Get time-to-live for keys in seconds (-1 for no expiry, -2 for non-existent)
- **DEL**: Delete one or more keys, returns number of keys deleted
- **EXPIRE**: Set expiration time for a key in seconds, returns 1 if successful, 0 if key doesn't exist



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

## Automatic Key Expiration

### Redis-Compatible Auto-Deletion
The server implements **automatic background expiration** using the same algorithm as Redis:

1. **Sampling-Based Approach**: Every second, a sample of 20 keys with expiration is tested
2. **Adaptive Deletion**: If more than 25% of sampled keys are expired, continue sampling and deleting
3. **Efficient Cleanup**: Process stops when less than 25% of sampled keys are expired
4. **Event-Loop Integration**: Auto-deletion runs in the main event loop with 1-second timeout
5. **Zero Blocking**: Uses epoll timeout to ensure deletion runs even when server is idle

### Implementation Details
- **Frequency**: Runs every 1 second (configurable via `cronFreq` variable)
- **Sample Size**: 20 keys per iteration (Redis-standard approach)
- **Threshold**: 25% expired keys trigger additional cleanup cycles
- **Integration**: Built into epoll event loop with 1000ms timeout
- **Performance**: Non-blocking operation that doesn't affect client request handling

### Benefits
- ✅ **Memory Efficient**: Automatic cleanup prevents memory leaks from expired keys
- ✅ **Redis Compatible**: Uses the same expiration algorithm as Redis
- ✅ **Performance Optimized**: Sampling approach scales well with large datasets
- ✅ **Always Active**: Runs continuously even when no clients are connected
- ✅ **Non-Intrusive**: Doesn't block client operations or degrade performance
- ✅ **Configurable**: Auto-deletion frequency can be adjusted via `cronFreq` variable
- ✅ **Observable**: Logs cleanup statistics for monitoring and debugging

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
2025/08/23 21:20:47 Starting Async TCP server on 0.0.0.0 7379
Deleted the expired Keys. total keys 0
Deleted the expired Keys. total keys 0
# Auto-deletion runs every second in background

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
localhost:7379> SET mykey "Hello Redis"
OK
localhost:7379> GET mykey
"Hello Redis"
localhost:7379> SET tempkey "expires" EX 10
OK
localhost:7379> TTL tempkey
(integer) 7
localhost:7379> GET tempkey
"expires"
localhost:7379> TTL tempkey
(integer) 4
localhost:7379> SET another_key "test value"
OK
localhost:7379> EXPIRE another_key 30
(integer) 1
localhost:7379> TTL another_key
(integer) 28
localhost:7379> SET key1 "value1"
OK
localhost:7379> SET key2 "value2"
OK
localhost:7379> DEL key1
(integer) 1
localhost:7379> DEL key1 key2 nonexistent
(integer) 1
localhost:7379> GET key1
(nil)
localhost:7379> INVALID
(error) ERR unknown command 'INVALID'
```

## Development

### Prerequisites
- **Go**: Version 1.16 or higher
- **Network**: TCP connection capability
- **OS**: Linux

### Code Architecture

#### `main.go`
- Configuration loading and validation with `config.LoadConfig()`
- Store initialization with memory limits and eviction strategy
- Application entry point with async server initialization
- Error handling and graceful shutdown

#### `config/config.go`
- JSON configuration file parsing and validation
- Command-line flag definitions and override handling
- Configuration validation (port ranges, limits, strategies)
- Default configuration values and help text
- Support for custom config file paths

#### `core/eval.go`
- Redis command evaluation and response generation
- Command implementations: PING, ECHO, TIME, SET, GET, TTL, DEL, EXPIRE
- RESP encoding utilities for proper Redis responses
- RedisCmd structure for parsed commands

#### `core/eviction.go`
- Key eviction strategies for memory management
- `Evict()` function with strategy selection (simple-first, future: LRU, random)
- `evictFirst()` implementation for simple-first strategy
- Debug logging for eviction monitoring and troubleshooting

#### `core/expire.go`
- Automatic key expiration and cleanup functionality
- Redis-compatible sampling algorithm for efficient memory management
- `DeleteExpireKeys()` function for background cleanup (called every second)
- `expireSample()` function implementing 20-key sampling with 25% threshold
- Integration with the main event loop for non-blocking operation

#### `core/store.go`
- In-memory key-value store with expiration and eviction support
- Configuration-aware memory limit enforcement
- `Put()` function with automatic eviction when limits exceeded
- Expiration timestamp management and cleanup
- Debug logging for store operations and key management

#### `core/RESP.go`
- Complete RESP protocol parser for all data types
- Functions: readSimpleString, readBulkString, readArray, readInt64, readError
- DecodeOne for dispatching to appropriate parsers
- DecodeCmd for command extraction from RESP arrays

#### `server/aync_tcp.go` (Production Server)
- High-performance epoll-based async TCP server with auto-deletion integration
- Non-blocking I/O with FDConn wrapper for io.ReadWriter compatibility
- Epoll timeout (1000ms) to ensure auto-deletion runs every second
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
redis-benchmark -h localhost -p 7379 -c 100 -n 10000 SET test value
redis-benchmark -h localhost -p 7379 -c 100 -n 10000 GET test
redis-benchmark -h localhost -p 7379 -c 100 -n 10000 DEL test
```

## Current Status & Limitations

###  **Completed Features**
- **High-Performance Server**: Epoll-based async I/O for production workloads
- **Complete RESP Protocol**: All Redis data types (Simple Strings, Bulk Strings, Arrays, Integers, Errors)
- **Redis Commands**: PING, ECHO, TIME, SET, GET, TTL, DEL, EXPIRE with proper protocol compliance
- **Flexible Configuration**: JSON config file with command-line overrides for all settings
- **Memory Management**: Configurable key limits with automatic eviction strategies
- **Key Eviction**: Simple-first eviction strategy with future support for LRU and random
- **Automatic Key Expiration**: Redis-compatible background auto-deletion using sampling algorithm
- **Configuration Validation**: Comprehensive validation of all configuration parameters
- **Event-Loop Integration**: Auto-deletion runs every second within the main epoll event loop
- **Debug Logging**: Detailed logging for troubleshooting eviction and store operations
- **Production Ready**: Error handling, graceful configuration, and robust architecture
- **Concurrent Connections**: Support for simultaneous clients with configurable limits


###  **Current Limitations**
- **Limited Commands**: Only 8 basic commands implemented (Redis has 200+ commands)
- **No Persistence**: Data is not stored (in-memory only)
- **No Data Structures**: No support for Lists, Sets, Hashes, etc.
- **No Authentication**: No AUTH command or security features
- **No Clustering**: Single instance only



## Future Roadmap

### Phase 1: Core Redis Commands (In Progress)
- [x] **String Commands**: SET, GET (completed)
- [x] **Key Management**: TTL, DEL, EXPIRE (completed)
- [x] **Memory Management**: Key limits and eviction strategies (completed)
- [x] **Configuration System**: JSON config with CLI overrides (completed)
- [ ] **String Commands**: EXISTS, INCR, DECR
- [ ] **Key Management**: KEYS, TYPE
- [ ] **Database**: SELECT, FLUSHDB, FLUSHALL

### Phase 2: Advanced Data Structures
- [ ] **Lists**: LPUSH, RPUSH, LPOP, RPOP, LRANGE
- [ ] **Sets**: SADD, SREM, SMEMBERS, SINTER, SUNION
- [ ] **Hashes**: HSET, HGET, HDEL, HKEYS, HVALS
- [ ] **Sorted Sets**: ZADD, ZREM, ZRANGE, ZSCORE

### Phase 3: Production Features
- [ ] **Persistence**: RDB snapshots and AOF logging
- [ ] **Advanced Eviction**: LRU, LFU, and TTL-based eviction strategies
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

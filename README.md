# CPU Noise Test Tool

A simple CLI tool designed to generate CPU and network noise to interfere with WireGuard-like operations. The tool supports two test modes: cryptographic load testing and UDP noise generation.

## Architecture

### Project Structure

```
cpu-noise-test/
├── main.go          # Entry point, flag parsing, signal handling, orchestration
├── config.go        # Configuration structures and validation
├── crypto.go        # CryptoLoad test implementation
├── udp.go           # UDPNoise test implementation
├── Makefile         # Build configuration for Linux amd64
└── go.mod           # Go module definition
```

### Component Overview

#### 1. Main Entry Point (`main.go`)

The main component orchestrates the entire application:

- **Flag Parsing**: Parses command-line arguments using Go's `flag` package
- **Signal Handling**: Sets up graceful shutdown handlers for SIGINT (Ctrl+C) and SIGTERM
- **Context Management**: Creates a cancellable context that propagates shutdown signals to all goroutines
- **Mode Selection**: Routes execution to the appropriate test mode based on user input
- **Error Handling**: Validates configuration and provides clear error messages

#### 2. Configuration (`config.go`)

Defines and validates configuration structures:

- **CryptoConfig**: Configuration for crypto load tests
  - `DataSize`: Size of data buffer for SHA256 operations
  - `Workers`: Number of concurrent worker goroutines
  
- **UDPConfig**: Configuration for UDP noise tests
  - `Target`: Target address and port (e.g., "127.0.0.1:1")
  - `Rate`: Packets per second (0 = unlimited, rate is divided per worker)
  - `Workers`: Number of concurrent worker goroutines

Both configurations include validation functions to ensure parameters are within acceptable ranges.

#### 3. Crypto Load Test (`crypto.go`)

Implements CPU-intensive cryptographic operations:

- **Purpose**: Simulates WireGuard's ChaCha20-Poly1305 encryption workload by performing SHA256 operations
- **Implementation**: 
  - Spawns multiple worker goroutines (configurable)
  - Each worker continuously performs SHA256 hashing on a data buffer
  - Uses vector instruction sets (AVX/NEON) similar to WireGuard's encryption
- **Cancellation**: Responds to context cancellation for graceful shutdown

#### 4. UDP Noise Test (`udp.go`)

Generates network stack interference:

- **Purpose**: Creates UDP traffic to stress the network stack, similar to WireGuard's UDP dependency
- **Implementation**:
  - Spawns multiple worker goroutines (configurable)
  - Each worker creates its own UDP connection to the target address
  - Sends packets at a configurable rate (or maximum speed if rate is 0)
  - When rate is specified, it's divided per worker (total rate = rate × workers)
  - Uses `time.Ticker` for rate limiting when specified
- **Cancellation**: Responds to context cancellation for graceful shutdown

### Design Patterns

#### Graceful Shutdown

The application uses a context-based cancellation pattern:

1. Signal handler receives SIGINT/SIGTERM
2. Context cancellation is triggered
3. All goroutines check `ctx.Done()` in their loops
4. Workers exit cleanly when context is cancelled
5. `sync.WaitGroup` ensures all goroutines complete before program exit

#### Concurrency Model

- **Crypto Mode**: Uses `sync.WaitGroup` to manage multiple worker goroutines
- **UDP Mode**: Uses `sync.WaitGroup` to manage multiple worker goroutines, each with its own UDP connection
- Both modes respect context cancellation for coordinated shutdown

### Data Flow

```
User Input (Flags)
    ↓
main.go: Parse & Validate
    ↓
main.go: Setup Signal Handler
    ↓
main.go: Create Cancellable Context
    ↓
main.go: Select Mode
    ├──→ crypto.go: RunCryptoLoad()
    │       └──→ Spawn N workers → SHA256 loops
    │
    └──→ udp.go: RunUDPNoise()
            └──→ Spawn N workers → Each worker: UDP send loop (rate-limited or unlimited)
    
Signal Received (Ctrl+C)
    ↓
Context Cancelled
    ↓
All Goroutines Exit
    ↓
WaitGroup.Wait() Completes
    ↓
Program Exits Gracefully
```

## Usage

### Building

Build for Linux amd64 using the Makefile:

```bash
make build
```

This will create the binary at `bin/cpu-noise-test`.

To clean build artifacts:

```bash
make clean
```

### Command-Line Options

#### Required Flags

- `--mode`: Test mode selection
  - `crypto`: Run cryptographic load test
  - `udp`: Run UDP noise test

#### Crypto Mode Options

- `--data-size`: Data size in bytes per SHA256 operation (default: 1024)
- `--workers`: Number of concurrent worker goroutines (default: 1)

#### UDP Mode Options

- `--target`: Target address:port (default: "127.0.0.1:1")
- `--rate`: Send rate in packets per second per worker (default: 0 = unlimited)
- `--workers`: Number of concurrent worker goroutines (default: 1)

### Examples

#### Crypto Load Test

Run with default settings (1 worker, 1024 bytes):
```bash
./bin/cpu-noise-test --mode crypto
```

Run with 4 workers and 2048 byte data blocks:
```bash
./bin/cpu-noise-test --mode crypto --workers 4 --data-size 2048
```

#### UDP Noise Test

Run with default settings (1 worker, unlimited rate to 127.0.0.1:1):
```bash
./bin/cpu-noise-test --mode udp
```

Run with rate limiting (1000 packets per second per worker):
```bash
./bin/cpu-noise-test --mode udp --target 127.0.0.1:1 --rate 1000
```

Run with multiple workers (5 workers, unlimited rate):
```bash
./bin/cpu-noise-test --mode udp --workers 5
```

Run with multiple workers and rate limiting (3 workers, 30 pps total = 10 pps per worker):
```bash
./bin/cpu-noise-test --mode udp --workers 3 --rate 30
```

Run targeting a specific address and port with multiple workers:
```bash
./bin/cpu-noise-test --mode udp --target 192.168.1.100:51820 --rate 500 --workers 4
```

### Graceful Shutdown

The tool supports graceful shutdown via Ctrl+C or SIGTERM:

1. Press `Ctrl+C` or send `SIGTERM` signal
2. The tool will display: "Received signal: <signal>. Shutting down gracefully..."
3. All worker goroutines will complete their current operations
4. The program exits cleanly

This ensures no resource leaks and allows the system to clean up properly.

### Help

View all available options:

```bash
./bin/cpu-noise-test -h
```

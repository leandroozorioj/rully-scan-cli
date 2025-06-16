# Rully CLI 🚀

A multi-platform network tool developed in Go Lang for IP/DNS verification and port scanning.

## ✨ Features

- **IP and DNS Verification**: DNS resolution, reverse DNS, CNAME and MX records
- **Port Scanning**: TCP port scanner for remote hosts
- **Local Analysis**: Check open ports on local system
- **Multi-threaded**: Multiple thread support for fast scanning
- **Multi-platform**: Windows, Linux and macOS

## 🛠️ Installation

### Prerequisites
- Go 1.21 or higher
- Make (optional, for automated builds)

### Compilation

1. Clone the project:
```bash
git git clone https://github.com/leandroozorioj/rully-scan-cli.git
cd rully-scan-cli
```

2. Install dependencies:
```bash
go mod download
```

3. Build for development:
```bash
go build -o rully-scan main.go
```

### Multi-platform Build

Use the Makefile to generate binaries for all platforms:

```bash
# Build for all platforms
make all

# Specific build
make windows  # Windows
make linux    # Linux
make darwin   # macOS
```

Binaries will be generated in the `build/` folder.

## 📖 Usage

### IP and DNS Verification

Check DNS information for a host:

```bash
# Check DNS information
./rully ip google.com
./rully ip github.com

# Example output:
[DNS] IP addresses found:
  172.217.29.14 (IPv4)
  2800:3f0:4001:81b::200e (IPv6)

[CNAME] www.google.com.
[MX RECORDS]
  smtp.google.com. (priority: 10)
```

### Port Scanning

Scan TCP ports of a remote host:

```bash
# Scan ports 1-1000
./rully-scan scan google.com

# Scan specific ports
./rully-scan scan google.com -p 80,443,8080

# Scan custom range
./rully-scan scan 192.168.1.1 -p 1-500

# With more threads and custom timeout
./rully-scan scan example.com -p 1-1000 -T 200 -t 5
```

### Local Analysis

Check open ports on local system:

```bash
# Check open local ports
./rully-scan local

# Example output:
[OPEN LOCAL PORTS]
  80/tcp    HTTP
  443/tcp   HTTPS
  3306/tcp  MySQL
```

## 🔧 Parameters

### Global Flags

- `--timeout, -t`: Timeout in seconds (default: 3)
- `--threads, -T`: Number of threads (default: 100)
- `--verbose, -v`: Verbose mode
- `--help, -h`: Show help

### `scan` Command

- `--ports, -p`: Port range (default: 1-1000)
  - Examples: `1-1000`, `80,443,8080`, `22`

## 📋 Usage Examples

```bash
# Basic DNS verification
./rully-scan ip example.com

# Quick scan of common ports
./rully-scan scan example.com -p 21,22,23,25,53,80,110,443

# Complete scan with verbose
./rully-scan scan 192.168.1.1 -p 1-65535 -T 500 -t 2 -v

# Check local services
./rully-scan local

# Specific scan for web services
./rully-scan scan myserver.com -p 80,443,8080,8443
```

## 🏗️ Project Structure

```
rully-scan-cli/
├── main.go          # Main code
├── go.mod           # Go dependencies
├── Makefile         # Build scripts
├── README.md        # Documentation
└── build/           # Compiled binaries
    ├── rully-windows-amd64.exe
    ├── rully-linux-amd64
    └── rully-darwin-amd64
```

## 🚀 Advanced Features

### Performance
- Multi-threaded scanning with concurrency control
- Configurable timeout per connection
- Optimized for slow and fast networks

### Compatibility
- Full IPv4 and IPv6 support
- Automatic service detection by port
- Terminal colors for better visualization

### Security
- No logs or sensitive data storage
- Safe timeouts to prevent hangs
- Strict parameter validation

## 🔮 Upcoming Versions

- [ ] Results export (JSON, CSV, XML)
- [ ] OS and service version detection
- [ ] UDP scanner
- [ ] Threat intelligence API integration
- [ ] Optional web interface
- [ ] Predefined scan profiles

## 📄 License

MIT License - see LICENSE file for details.

---

**Rully Scan CLI v1.0.2** - Multi-platform Network Tool in Go Lang
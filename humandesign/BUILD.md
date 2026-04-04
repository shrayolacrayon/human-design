# Build Instructions

This guide covers how to build and run the Human Design Calculator.

## Prerequisites

- **Go 1.21 or later**
- **Swiss Ephemeris library** (for planetary calculations)
- **Docker** (optional, but recommended for easiest setup)

## Quick Start (Docker - Recommended)

The easiest way to build and run everything:

```bash
cd humandesign

# Build and run the web server
docker-compose up --build

# Or build the CLI in Docker
docker build -f Dockerfile.test -t humandesign-dev .
docker run humandesign-dev go build -o humandesign-cli ./cmd/cli/main.go
```

## Building the Web Server

### Option 1: Docker

```bash
cd humandesign

# Build
docker build -t humandesign .

# Run
docker run -p 8080:8080 humandesign

# Access at http://localhost:8080
```

### Option 2: Local Build (Ubuntu/Debian)

```bash
# Install Swiss Ephemeris
sudo apt-get update
sudo apt-get install libswe-dev

# Build
cd humandesign
go build -o humandesign-server ./cmd/server/main.go

# Run
./humandesign-server
```

### Option 3: Local Build (macOS)

```bash
# Install Swiss Ephemeris from source
cd /tmp
curl -L https://www.astro.com/ftp/swisseph/swe_unix_src_2.10.03.tar.gz -o swe.tar.gz
tar xzf swe.tar.gz
cd src
make

# Install headers and libraries
sudo mkdir -p /usr/local/include /usr/local/lib
sudo cp *.h /usr/local/include/
sudo cp libswe.a /usr/local/lib/

# Create dynamic library
gcc -shared -o libswe.dylib -fPIC swedate.o swehouse.o swejpl.o swemmoon.o swemplan.o swepcalc.o sweph.o swepdate.o swephlib.o swecl.o swehel.o
sudo cp libswe.dylib /usr/local/lib/

# Build the app
cd /path/to/humandesign
CGO_ENABLED=1 go build -o humandesign-server ./cmd/server/main.go

# Run
./humandesign-server
```

## Building the CLI Tool

### Option 1: Docker

```bash
cd humandesign

# Build in Docker
docker build -f Dockerfile.test -t humandesign-dev .
docker run -v $(pwd):/app humandesign-dev sh -c "go build -o humandesign-cli ./cmd/cli/main.go"

# The binary will be in the current directory
./humandesign-cli help
```

### Option 2: Local Build

```bash
cd humandesign

# Install dependencies (if not already installed)
go mod download

# Build
CGO_ENABLED=1 go build -o humandesign-cli ./cmd/cli/main.go

# Optional: Install globally
sudo cp humandesign-cli /usr/local/bin/

# Use it
humandesign-cli help
```

## Build Flags

### For production (optimized binary):

```bash
CGO_ENABLED=1 go build -ldflags="-s -w" -o humandesign-cli ./cmd/cli/main.go
```

- `-ldflags="-s -w"` - Strip debug info for smaller binary

### For development (with debugging):

```bash
CGO_ENABLED=1 go build -gcflags="all=-N -l" -o humandesign-cli ./cmd/cli/main.go
```

- `-gcflags="all=-N -l"` - Disable optimizations for debugging

## Cross-Compilation

Note: Cross-compilation with CGO is complex due to Swiss Ephemeris dependency.

### For Linux (from macOS):

```bash
# Requires cross-compilation toolchain
# Not recommended - use Docker instead
```

### Recommended approach for different platforms:

Use Docker multi-stage builds:

```dockerfile
FROM golang:1.22-bookworm AS builder
RUN apt-get update && apt-get install -y libswe-dev
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 go build -o humandesign-cli ./cmd/cli/main.go

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y libswe2.0
COPY --from=builder /app/humandesign-cli /usr/local/bin/
```

## Verifying the Build

### Test the web server:

```bash
# Start server
./humandesign-server

# In another terminal
curl http://localhost:8080
# Should return HTML homepage
```

### Test the CLI:

```bash
# Show version
./humandesign-cli version

# Calculate a reading
./humandesign-cli calculate \
  -name "Test" \
  -date "1990-06-15T14:30:00Z" \
  -lat 40.7128 \
  -lon -74.0060 \
  -location "New York"
```

## Running Tests

```bash
cd humandesign

# All tests (requires Swiss Ephemeris)
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test -v ./internal/calculator

# In Docker
docker-compose -f docker-compose.test.yml up test
```

## Development Workflow

### 1. Make changes to code

```bash
vim internal/calculator/calculator.go
```

### 2. Run tests

```bash
go test ./internal/calculator -v
```

### 3. Build

```bash
go build -o humandesign-cli ./cmd/cli/main.go
```

### 4. Test manually

```bash
./humandesign-cli save -name "Test" -date "1990-01-01T12:00:00Z" -lat 0 -lon 0 -location "Test"
```

## Troubleshooting

### "swephexp.h: No such file or directory"

**Problem:** Swiss Ephemeris library not installed

**Solution:**
- Ubuntu/Debian: `sudo apt-get install libswe-dev`
- macOS: Build from source (see Option 3 above)
- Or use Docker: `docker build -t humandesign .`

### "undefined reference to swe_calc"

**Problem:** Swiss Ephemeris library not linked

**Solution:**
```bash
# Ensure CGO is enabled
export CGO_ENABLED=1

# Check library path
export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH  # Linux
export DYLD_LIBRARY_PATH=/usr/local/lib:$DYLD_LIBRARY_PATH  # macOS
```

### "permission denied" when running binary

**Solution:**
```bash
chmod +x humandesign-cli
```

### Build is slow

**Problem:** CGO compilation is slower than pure Go

**Solution:**
- Use caching: `go build -cache`
- Build in Docker (caches layers)
- Use `go install` for development

## Binary Locations

After building:

```
humandesign/
├── humandesign-server   # Web server binary
├── humandesign-cli      # CLI binary
└── data/
    └── people/          # Storage directory (created on first use)
```

## Clean Build

```bash
# Remove built binaries
rm -f humandesign-server humandesign-cli

# Clean Go cache
go clean -cache

# Clean modules
go clean -modcache

# Start fresh
go mod download
go build ./...
```

## CI/CD

### GitHub Actions Example

```yaml
name: Build
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Install dependencies
        run: sudo apt-get install -y libswe-dev
      - name: Build server
        run: CGO_ENABLED=1 go build -o humandesign-server ./cmd/server/main.go
      - name: Build CLI
        run: CGO_ENABLED=1 go build -o humandesign-cli ./cmd/cli/main.go
      - name: Test
        run: go test ./...
```

## Performance Tips

1. **Use `-ldflags` for smaller binaries**
2. **Enable Go module caching** in CI/CD
3. **Use Docker layer caching** for faster builds
4. **Consider static linking** for distribution

## Distribution

### Single Binary

The CLI is a single binary with no runtime dependencies except Swiss Ephemeris:

```bash
# Linux
ldd humandesign-cli
# Should show libswe dependency

# macOS
otool -L humandesign-cli
```

### Docker Image

Smallest image:

```bash
docker build -t humandesign:latest .
docker images humandesign:latest
# Should be ~100MB
```

## Next Steps

After building:

1. Read [README.md](README.md) for usage
2. Read [TESTING.md](TESTING.md) for testing guide
3. Read [STORAGE.md](STORAGE.md) for storage system
4. Run `./humandesign-cli help` for CLI commands

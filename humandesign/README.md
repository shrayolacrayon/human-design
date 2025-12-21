# Human Design Calculator

A Go web application that generates Human Design readings with body graph visualizations based on birth data.

## Features

- **Swiss Ephemeris Integration**: Accurate planetary calculations using the Swiss Ephemeris library
- **Birth Data Input**: Enter date, time, and location of birth
- **Human Design Calculations**: Determines Type, Authority, Profile, Strategy, and more
- **Body Graph Visualization**: SVG-based visualization of defined/undefined centers and channels
- **REST API**: JSON endpoints for integration with other applications
- **CLI Tool**: Command-line interface for batch calculations and testing
- **CSV Integration**: Read birth data from CSV files for batch processing
- **Comprehensive Testing**: Unit tests, integration tests, and validation framework

## Quick Start with Docker (Recommended)

The easiest way to run this application is with Docker - no dependencies to install!

```bash
# Clone/extract the project
cd humandesign

# Option 1: Docker Compose (easiest)
docker-compose up

# Option 2: Docker directly
docker build -t humandesign .
docker run -p 8080:8080 humandesign
```

Then open http://localhost:8080 in your browser.

## Manual Installation (Alternative)

If you prefer to run without Docker:

### Ubuntu/Debian
```bash
sudo apt-get install libswe-dev
go run cmd/server/main.go
```

### macOS (build Swiss Ephemeris from source)
```bash
# Download and build Swiss Ephemeris
cd /tmp
curl -L https://www.astro.com/ftp/swisseph/swe_unix_src_2.10.03.tar.gz -o swe.tar.gz
tar xzf swe.tar.gz
cd src
make

# Install
sudo mkdir -p /usr/local/include /usr/local/lib
sudo cp *.h /usr/local/include/
sudo cp libswe.a /usr/local/lib/

# Create dynamic library
gcc -shared -o libswe.dylib -fPIC swedate.o swehouse.o swejpl.o swemmoon.o swemplan.o swepcalc.o sweph.o swepdate.o swephlib.o swecl.o swehel.o
sudo cp libswe.dylib /usr/local/lib/

# Run the app
cd /path/to/humandesign
go run cmd/server/main.go
```

## Project Structure

```
humandesign/
├── Dockerfile              # Docker build configuration
├── docker-compose.yml      # Docker Compose for easy startup
├── cmd/
│   ├── server/
│   │   └── main.go         # Web server entry point
│   └── cli/
│       └── main.go         # CLI tool entry point
├── internal/
│   ├── calculator/
│   │   ├── types.go        # Data structures and types
│   │   ├── gates.go        # Gate and channel definitions (36 channels)
│   │   ├── calculator.go   # Main calculation logic
│   │   ├── gates_test.go   # Unit tests for gates
│   │   └── calculator_test.go # Unit tests for calculator
│   ├── csvreader/
│   │   ├── csvreader.go    # CSV reading/writing functionality
│   │   └── csvreader_test.go # CSV reader tests
│   ├── ephemeris/
│   │   └── ephemeris.go    # Swiss Ephemeris integration via cgo
│   ├── bodygraph/
│   │   └── generator.go    # SVG body graph generation
│   └── handlers/
│       └── handlers.go     # HTTP request handlers
├── testdata/
│   ├── birth_data.csv      # Test cases with famous people
│   └── gate_validation.csv # Gate and channel validation tests
├── integration_test.go     # Integration tests
├── go.mod
├── README.md
└── TESTING.md             # Comprehensive testing guide
```

## API Endpoints

### GET /
The main web interface with a form to enter birth data.

### POST /api/reading
Generate a Human Design reading (returns HTML).

**Request Body:**
```json
{
  "datetime": "1990-06-15T14:30:00Z",
  "latitude": 40.7128,
  "longitude": -74.0060,
  "location": "New York, NY"
}
```

### POST /api/reading/json
Generate a Human Design reading (returns JSON).

**Response includes:**
- Type (Generator, Manifestor, Projector, Reflector, Manifesting Generator)
- Authority (Emotional, Sacral, Splenic, Ego, Self, Environmental, Lunar)
- Profile (1-6 combinations)
- Strategy
- Signature and Not-Self Theme
- Centers (defined/undefined)
- Channels (defined)
- Gates (Personality and Design)
- Incarnation Cross

## CLI Tool

### Building the CLI

```bash
cd humandesign/cmd/cli
go build -o humandesign-cli
```

### CLI Commands

```bash
# Calculate reading from individual parameters
humandesign-cli calculate \
  -name "John Doe" \
  -date "1990-06-15T14:30:00Z" \
  -lat 40.7128 \
  -lon -74.0060 \
  -location "New York, NY"

# Calculate from CSV file
humandesign-cli calculate -csv testdata/birth_data.csv -output json

# Validate test cases
humandesign-cli validate -csv testdata/birth_data.csv -verbose

# Show help
humandesign-cli help
```

### CSV Format

```csv
name,datetime,latitude,longitude,location,expected_type,expected_authority,expected_profile_conscious,expected_profile_unconscious,expected_strategy
Steve Jobs,1955-02-24T19:15:00Z,37.7749,-122.4194,"San Francisco, CA",Manifestor,Splenic,5,1,Inform Before Acting
```

See `TESTING.md` for complete CLI documentation and CSV format details.

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -v -run Integration

# Validate test data
humandesign-cli validate -csv testdata/birth_data.csv
```

For comprehensive testing documentation, see [TESTING.md](TESTING.md).

## Understanding Human Design

### Types
- **Generator**: Life force energy, wait to respond (defined Sacral)
- **Manifesting Generator**: Multi-passionate, wait to respond then inform
- **Projector**: Guides and managers, wait for invitations (no defined Sacral, no motor to Throat)
- **Manifestor**: Initiators, inform before acting (motor to Throat, no Sacral)
- **Reflector**: Mirrors, wait a lunar cycle (no defined centers)

### The 9 Centers
| Center | Function | Motor? |
|--------|----------|--------|
| Head | Inspiration, mental pressure | No |
| Ajna | Conceptualization | No |
| Throat | Communication, manifestation | No |
| G Center | Identity, love, direction | No |
| Heart/Ego | Willpower | Yes |
| Sacral | Life force, work energy | Yes |
| Solar Plexus | Emotions | Yes |
| Spleen | Intuition, health | No |
| Root | Adrenaline, drive | Yes |

### Gates and Channels
- **64 Gates**: Correspond to I Ching hexagrams, each spans 5.625° of the zodiac
- **6 Lines**: Each gate has 6 lines (0.9375° each)
- **36 Channels**: Connect two centers when both gates are activated
- **Personality (Black)**: Conscious - calculated at birth time
- **Design (Red)**: Unconscious - calculated when Sun was 88° earlier

## License

MIT License - Feel free to use and modify as needed.

# Human Design Calculator

A Go web application that generates Human Design readings with body graph visualizations based on birth data.

## Features

- **Swiss Ephemeris Integration**: Accurate planetary calculations using the Swiss Ephemeris library
- **Birth Data Input**: Enter date, time, and location of birth
- **Human Design Calculations**: Determines Type, Authority, Profile, Strategy, and more
- **Body Graph Visualization**: SVG-based visualization of defined/undefined centers and channels
- **REST API**: JSON endpoints for integration with other applications

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
│   └── server/
│       └── main.go         # Application entry point
├── internal/
│   ├── calculator/
│   │   ├── types.go        # Data structures and types
│   │   ├── gates.go        # Gate and channel definitions (36 channels)
│   │   └── calculator.go   # Main calculation logic
│   ├── ephemeris/
│   │   └── ephemeris.go    # Swiss Ephemeris integration via cgo
│   ├── bodygraph/
│   │   └── generator.go    # SVG body graph generation
│   └── handlers/
│       └── handlers.go     # HTTP request handlers
├── go.mod
└── README.md
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

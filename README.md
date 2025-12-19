# Human Design Calculator

A Go web application that generates Human Design readings with body graph visualizations based on birth data.

## Features

- **Birth Data Input**: Enter date, time, and location of birth
- **Human Design Calculations**: Determines Type, Authority, Profile, Strategy, and more
- **Body Graph Visualization**: SVG-based visualization of defined/undefined centers and channels
- **REST API**: JSON endpoints for integration with other applications

## Project Structure

```
humandesign/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── calculator/
│   │   ├── types.go         # Data structures and types
│   │   ├── gates.go         # Gate and channel definitions
│   │   └── calculator.go    # Main calculation logic
│   ├── ephemeris/
│   │   └── ephemeris.go     # Planetary position calculations
│   ├── bodygraph/
│   │   └── generator.go     # SVG body graph generation
│   └── handlers/
│       └── handlers.go      # HTTP request handlers
├── go.mod
└── README.md
```

## Quick Start

```bash
# Navigate to the project directory
cd humandesign

# Run the application
go run cmd/server/main.go

# The server starts at http://localhost:8080
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

## Understanding Human Design

### Types
- **Generator**: Life force energy, wait to respond
- **Manifesting Generator**: Multi-passionate, wait to respond then inform
- **Projector**: Guides and managers, wait for invitations
- **Manifestor**: Initiators, inform before acting
- **Reflector**: Mirrors, wait a lunar cycle

### Centers
The 9 centers in the body graph:
1. **Head** - Inspiration and mental pressure
2. **Ajna** - Conceptualization and mental processing
3. **Throat** - Communication and manifestation
4. **G Center** - Identity, love, and direction
5. **Heart/Ego** - Willpower and ego
6. **Sacral** - Life force and work energy
7. **Solar Plexus** - Emotions and feelings
8. **Spleen** - Intuition, health, and survival
9. **Root** - Adrenaline and drive

### Gates and Channels
- **Gates**: 64 gates corresponding to the I Ching hexagrams
- **Channels**: Connections between centers formed when both gates of a pair are activated
- **Personality (Black)**: Conscious traits from birth time
- **Design (Red)**: Unconscious traits from ~88 days before birth

## Important Notes

### Ephemeris Accuracy
⚠️ **The ephemeris calculations in this project are simplified approximations.**

For production use, you should integrate a proper astronomical library such as:
- [Swiss Ephemeris](https://www.astro.com/swisseph/) - The gold standard
- [github.com/mshafiee/swephgo](https://github.com/mshafiee/swephgo) - Go bindings for Swiss Ephemeris
- An external ephemeris API service

### Timezone Handling
Birth times should ideally be converted to UTC accounting for the timezone of the birth location. The current implementation expects UTC times.

## Extending the Application

### Adding More Accurate Calculations

1. **Integrate Swiss Ephemeris**:
   ```go
   // Replace internal/ephemeris with proper Swiss Ephemeris bindings
   import "github.com/mshafiee/swephgo"
   ```

2. **Add Timezone Support**:
   ```go
   // Use a timezone database to convert local birth time to UTC
   import "github.com/zsefvlol/timezonemapper"
   ```

3. **Enhance Gate Mapping**:
   The gate-to-longitude mapping should use the precise Human Design wheel positions.

### Adding More Features

- **Composite Charts**: Compare two charts for relationships
- **Transit Charts**: Current planetary positions overlaid on birth chart
- **Line Descriptions**: Detailed descriptions for each gate.line
- **Cross Descriptions**: Full incarnation cross interpretations

## License

MIT License - Feel free to use and modify as needed.

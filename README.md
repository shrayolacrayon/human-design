# Cosmic Blueprint

A Go web application for Human Design, Western Astrology, and Astrocartography — all calculated from birth data.

## Features

- **Human Design**: Type, Authority, Profile, Strategy, body graph visualization
- **Astrology**: Natal chart with planetary placements, houses, aspects, elements
- **Astrocartography**: Planetary line influence by location
- **People Database**: Save birth data for multiple people and generate any chart in one click
- **City Search**: Searchable dropdown of ~170 cities — no manual lat/long entry needed

## Running the App

### Option 1: Docker (recommended)

Requires [Docker](https://www.docker.com/) or [Colima](https://github.com/abiosoft/colima) (macOS).

```bash
# macOS: start the Docker daemon first
colima start

# Build and run
cd humandesign
docker compose up --build
```

Then open **http://localhost:8080** in your browser.

To stop: `Ctrl+C`, then `docker compose down`.

### Option 2: Native (macOS)

Requires Go 1.21+ and the Swiss Ephemeris C library.

```bash
cd humandesign

# Install libswe (one-time setup)
bash install-macos-deps.sh

# Run
go run ./cmd/server/
```

Then open **http://localhost:8080**.

---

## Pages

| URL | Description |
|-----|-------------|
| `/` | Human Design chart |
| `/astrology` | Western natal astrology chart |
| `/astrocartography` | Planetary line map by location |
| `/people` | Saved people database |

## People Database

Birth data is stored in `humandesign/data/people.json`. You can commit this file to keep a shared database across machines. Edit it directly or use the `/people` UI to add/remove entries.

## Project Structure

```
humandesign/
├── cmd/server/main.go           # Entry point
├── data/people.json             # People database (committable)
├── internal/
│   ├── astrology/               # Natal chart calculations
│   ├── astrocartography/        # Planetary line calculations
│   ├── bodygraph/               # SVG body graph generator
│   ├── calculator/              # Human Design engine
│   ├── cities/                  # City coordinates for dropdown
│   ├── database/                # JSON file storage
│   ├── ephemeris/               # Swiss Ephemeris C bindings
│   └── handlers/                # HTTP handlers and HTML views
├── docker-compose.yml
├── Dockerfile
└── install-macos-deps.sh        # macOS libswe installer
```

## API Endpoints

### POST /api/reading
Human Design reading (returns HTML).

### POST /api/reading/json
Human Design reading (returns JSON).

### POST /api/astrology
Natal astrology chart (returns HTML).

### POST /api/astrocartography
Astrocartography report (returns HTML).

### GET /api/people
List all saved people (JSON).

### POST /api/people
Add a person (JSON body).

### DELETE /api/people/{id}
Delete a person.

### GET /api/cities
List all available cities with coordinates (JSON).

**Chart request body:**
```json
{
  "datetime": "1990-06-15T14:30:00Z",
  "latitude": 40.7128,
  "longitude": -74.0060,
  "location": "New York, US"
}
```

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

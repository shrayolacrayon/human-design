# Human Design Testing Guide

This guide explains how to run tests and validate Human Design calculations using CSV test data.

## Overview

The testing framework includes:
- **Unit Tests**: Test individual components (gates, channels, calculator logic)
- **Integration Tests**: Test full calculations against CSV test data
- **CLI Tool**: Command-line interface for running calculations and validations

## Test Data Structure

### CSV Test Files

Test data is stored in CSV files in the `testdata/` directory:

#### `birth_data.csv`
Contains test cases for famous people with expected Human Design types:

```csv
name,datetime,latitude,longitude,location,expected_type,expected_authority,expected_profile_conscious,expected_profile_unconscious,expected_strategy
Steve Jobs,1955-02-24T19:15:00Z,37.7749,-122.4194,"San Francisco, CA",Manifestor,Splenic,5,1,Inform Before Acting
```

**Columns:**
- `name`: Person's name (for reference)
- `datetime`: Birth date/time in RFC3339 format (YYYY-MM-DDTHH:MM:SSZ)
- `latitude`: Birth latitude (decimal degrees)
- `longitude`: Birth longitude (decimal degrees)
- `location`: Birth location name
- `expected_type`: Expected Human Design type
- `expected_authority`: Expected authority
- `expected_profile_conscious`: Expected conscious profile number (1-6)
- `expected_profile_unconscious`: Expected unconscious profile number (1-6)
- `expected_strategy`: Expected strategy

#### `gate_validation.csv`
Contains test cases for validating specific gates and channels:

```csv
name,datetime,latitude,longitude,location,expected_gates,expected_channels,notes
Gate 1 Test,1990-01-01T12:00:00Z,0.0,0.0,"Test Location","1,8","G-Throat Inspiration",Testing gate 1 activation
```

**Additional Columns:**
- `expected_gates`: Comma-separated list of gate numbers that should be activated
- `expected_channels`: Comma-separated list of channel names that should be defined
- `notes`: Additional notes about the test case

## Running Tests

### Unit Tests

Run all unit tests:
```bash
cd humandesign
go test ./...
```

Run tests for a specific package:
```bash
go test ./internal/calculator
go test ./internal/csvreader
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run a specific test:
```bash
go test -v -run TestAllGatesExist ./internal/calculator
```

### Integration Tests

Integration tests read from CSV files and validate calculations:

```bash
go test -v ./... -run Integration
```

### Test Coverage

Generate test coverage report:
```bash
go test -cover ./...
```

Generate detailed coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## CLI Tool

### Building the CLI

```bash
cd humandesign/cmd/cli
go build -o humandesign-cli
```

Or install it globally:
```bash
cd humandesign/cmd/cli
go install
```

### CLI Commands

#### Calculate Command

Calculate Human Design reading from individual parameters:
```bash
humandesign-cli calculate \
  -name "John Doe" \
  -date "1990-06-15T14:30:00Z" \
  -lat 40.7128 \
  -lon -74.0060 \
  -location "New York, NY"
```

Calculate from CSV file:
```bash
humandesign-cli calculate -csv testdata/birth_data.csv
```

Output as JSON:
```bash
humandesign-cli calculate -csv testdata/birth_data.csv -output json
```

Output as CSV:
```bash
humandesign-cli calculate -csv testdata/birth_data.csv -output csv
```

#### Validate Command

Validate calculations against expected results in CSV:
```bash
humandesign-cli validate -csv testdata/birth_data.csv
```

Verbose validation output:
```bash
humandesign-cli validate -csv testdata/birth_data.csv -verbose
```

JSON validation output:
```bash
humandesign-cli validate -csv testdata/birth_data.csv -output json
```

#### Version Command

Show CLI version:
```bash
humandesign-cli version
```

### CLI Output Examples

#### Text Output (default)
```
Human Design Reading
Birth Date: 1990-06-15 14:30:00 UTC
Location: New York, NY (40.7128, -74.0060)

Type: Generator
Strategy: Wait to Respond
Authority: Sacral
Profile: 1/3 - Investigator/Martyr
...
```

#### Validation Output
```
✓ PASS Steve Jobs
✗ FAIL Oprah Winfrey
  - Type mismatch: expected Generator, got Manifesting Generator

--------------------------------------------------
Total: 5 | Passed: 4 | Failed: 1
Success Rate: 80.0%
```

## Creating Test Cases

### Adding New Test Data

1. Create or edit a CSV file in `testdata/`:
```csv
name,datetime,latitude,longitude,location,expected_type,expected_authority,expected_profile_conscious,expected_profile_unconscious,expected_strategy
Your Name,1985-03-15T10:30:00Z,51.5074,-0.1278,"London, UK",Projector,Splenic,2,4,Wait for the Invitation
```

2. Run validation to check:
```bash
humandesign-cli validate -csv testdata/your_file.csv -verbose
```

### Date/Time Format

**Important**: All dates must be in RFC3339 format with timezone:
- UTC: `1990-06-15T14:30:00Z`
- With timezone: `1990-06-15T14:30:00-05:00`

Convert local time to UTC before adding to CSV.

### Coordinates

- Latitude: -90 (South Pole) to +90 (North Pole)
- Longitude: -180 (West) to +180 (East)

## Testing Best Practices

### 1. Test Data Quality

- Use accurate birth times (preferably from birth certificates)
- Verify coordinates for birth locations
- Cross-reference expected results with reliable Human Design software

### 2. Understanding Test Results

**Note**: This implementation uses simplified ephemeris calculations. Results may differ slightly from professional Human Design software that uses Swiss Ephemeris.

Test warnings (not failures) are normal for:
- Gate positions (approximate calculations)
- Some profile variations

Test failures should be investigated if:
- Type is consistently wrong
- Strategy doesn't match type
- Authority hierarchy is violated

### 3. Adding Comprehensive Tests

For thorough testing, include cases for:
- All 5 types (Generator, Manifesting Generator, Projector, Manifestor, Reflector)
- All 7 authorities
- All 12 profiles
- Edge cases (Reflectors with no definition, etc.)

## Troubleshooting

### Tests Failing

1. **Check ephemeris accuracy**:
   - This implementation uses approximate calculations
   - Consider integrating Swiss Ephemeris for production use

2. **Verify test data**:
   - Ensure dates are in RFC3339 format
   - Check coordinates are correct
   - Verify expected results against reliable sources

3. **Check CSV format**:
   - Ensure proper quoting for locations with commas
   - Use consistent column ordering
   - Include all required columns

### CSV Reading Errors

If you get CSV parsing errors:
- Check for unquoted commas in text fields
- Ensure consistent column count across all rows
- Verify file encoding is UTF-8

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.21'
      - run: go test -v ./...
      - run: cd cmd/cli && go build
      - run: ./cmd/cli/humandesign-cli validate -csv testdata/birth_data.csv
```

## Additional Resources

- [Human Design System](https://www.jovianarchive.com/)
- [Swiss Ephemeris](https://www.astro.com/swisseph/)
- [Go Testing Documentation](https://golang.org/pkg/testing/)

## Contributing Tests

When contributing test data:
1. Ensure accuracy of birth information
2. Verify expected results with multiple sources
3. Add meaningful test names and notes
4. Include edge cases when possible
5. Document any special circumstances

## License

Test data and testing framework are part of the Human Design project and share the same MIT license.

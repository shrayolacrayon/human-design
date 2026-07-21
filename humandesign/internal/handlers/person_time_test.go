package handlers

import (
	"testing"

	"humandesign/internal/calculator"
	"humandesign/internal/database"
)

// A person's birth_time is local wall-clock time at the birth location, so it
// must be interpreted in that location's timezone before being converted to
// UTC for the ephemeris. Previously it was parsed directly as UTC, which
// shifted every non-UTC birth and produced the wrong gates/profile/type.
func TestPersonBirthUTC_UsesLocationTimezone(t *testing.T) {
	tests := []struct {
		name    string
		person  database.Person
		wantUTC string // RFC3339 in UTC
	}{
		{
			name: "explicit timezone (Hong Kong, UTC+8)",
			person: database.Person{
				BirthDate: "1985-01-19", BirthTime: "11:00",
				Timezone: "Asia/Hong_Kong",
				Location: "Hong Kong, HK", Latitude: 22.3193, Longitude: 114.1694,
			},
			wantUTC: "1985-01-19T03:00:00Z",
		},
		{
			name: "timezone resolved from location label",
			person: database.Person{
				BirthDate: "1985-01-19", BirthTime: "11:00",
				Location: "Hong Kong, HK", Latitude: 22.3193, Longitude: 114.1694,
			},
			wantUTC: "1985-01-19T03:00:00Z",
		},
		{
			name: "timezone resolved from coordinates when label unknown",
			person: database.Person{
				BirthDate: "1985-01-19", BirthTime: "11:00",
				Location: "Somewhere near Hong Kong", Latitude: 22.30, Longitude: 114.20,
			},
			wantUTC: "1985-01-19T03:00:00Z",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := personBirthUTC(&tc.person)
			if err != nil {
				t.Fatalf("personBirthUTC: %v", err)
			}
			if got.UTC().Format("2006-01-02T15:04:05Z") != tc.wantUTC {
				t.Errorf("got UTC %s, want %s", got.UTC().Format("2006-01-02T15:04:05Z"), tc.wantUTC)
			}
		})
	}
}

// End-to-end guard: the Hong Kong birth at 11:00 local must read as a 3/5,
// not the 4/6 produced by the old parse-as-UTC bug.
func TestPersonBirthUTC_HongKongProfileIs35(t *testing.T) {
	person := &database.Person{
		BirthDate: "1985-01-19", BirthTime: "11:00",
		Location: "Hong Kong, HK", Latitude: 22.3193, Longitude: 114.1694,
	}

	dt, err := personBirthUTC(person)
	if err != nil {
		t.Fatalf("personBirthUTC: %v", err)
	}

	reading, err := calculator.NewCalculator().Calculate(calculator.BirthData{
		DateTime: dt, Latitude: person.Latitude, Longitude: person.Longitude, Location: person.Location,
	})
	if err != nil {
		t.Fatalf("calculate: %v", err)
	}

	if got := reading.Profile.Conscious; got != 3 {
		t.Errorf("conscious line = %d, want 3", got)
	}
	if got := reading.Profile.Unconscious; got != 5 {
		t.Errorf("unconscious line = %d, want 5", got)
	}
	if reading.Type != calculator.TypeManifestingGenerator {
		t.Errorf("type = %q, want %q", reading.Type, calculator.TypeManifestingGenerator)
	}
}

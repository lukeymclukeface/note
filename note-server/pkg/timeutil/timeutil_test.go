package timeutil

import (
	"testing"
	"time"
)

func TestGetCurrentTimestamp(t *testing.T) {
	timestamp := GetCurrentTimestamp()
	
	// Should be able to parse the returned timestamp
	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Errorf("GetCurrentTimestamp() returned invalid RFC3339 timestamp: %v", err)
	}
	
	// Should be recent (within last second)
	parsed, _ := time.Parse(time.RFC3339, timestamp)
	now := time.Now().UTC()
	diff := now.Sub(parsed)
	if diff > time.Second {
		t.Errorf("GetCurrentTimestamp() returned timestamp too old: %v", diff)
	}
}

func TestFormatTimestamp(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	
	result := FormatTimestamp(testTime)
	expected := "2023-12-25T15:30:45Z"
	
	if result != expected {
		t.Errorf("FormatTimestamp() = %v, want %v", result, expected)
	}
}

func TestFormatTimestampWithTimezone(t *testing.T) {
	// Test with a timezone - should convert to UTC
	loc, _ := time.LoadLocation("America/New_York")
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, loc)
	
	result := FormatTimestamp(testTime)
	
	// Should be in UTC format
	parsed, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Errorf("FormatTimestamp() returned invalid RFC3339 timestamp: %v", err)
	}
	
	if parsed.Location() != time.UTC {
		t.Errorf("FormatTimestamp() did not convert to UTC: %v", parsed.Location())
	}
}

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "valid RFC3339 timestamp",
			input:     "2023-12-25T15:30:45Z",
			wantError: false,
		},
		{
			name:      "valid RFC3339 timestamp with timezone",
			input:     "2023-12-25T15:30:45-05:00",
			wantError: false,
		},
		{
			name:      "invalid timestamp format",
			input:     "2023-12-25 15:30:45",
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "invalid date",
			input:     "2023-13-25T15:30:45Z",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTimestamp(tt.input)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseTimestamp() expected error for input %v, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseTimestamp() unexpected error: %v", err)
				}
				
				// Verify round-trip conversion
				formatted := FormatTimestamp(result)
				parsed2, err := ParseTimestamp(formatted)
				if err != nil {
					t.Errorf("ParseTimestamp() round-trip failed: %v", err)
				}
				
				if !result.Equal(parsed2) {
					t.Errorf("ParseTimestamp() round-trip mismatch: %v != %v", result, parsed2)
				}
			}
		})
	}
}

func TestTimestampConsistency(t *testing.T) {
	// Test that current timestamp can be parsed and formatted consistently
	timestamp1 := GetCurrentTimestamp()
	parsed, err := ParseTimestamp(timestamp1)
	if err != nil {
		t.Errorf("Failed to parse current timestamp: %v", err)
	}
	
	timestamp2 := FormatTimestamp(parsed)
	if timestamp1 != timestamp2 {
		t.Errorf("Timestamp consistency failed: %v != %v", timestamp1, timestamp2)
	}
}

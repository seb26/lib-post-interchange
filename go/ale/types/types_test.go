package types

import (
	"encoding/json"
	"testing"
)

func TestBaseField(t *testing.T) {
	field := BaseField{
		Key:   "test_key",
		Value: "test_value",
	}

	if field.GetKey() != "test_key" {
		t.Errorf("GetKey() = %v, want %v", field.GetKey(), "test_key")
	}

	if field.GetValue() != "test_value" {
		t.Errorf("GetValue() = %v, want %v", field.GetValue(), "test_value")
	}
}

func TestFieldTypes(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		key   string
		value string
	}{
		{
			name:  "FieldDelimiter",
			field: FieldDelimiter{BaseField{Key: "Field", Value: "TAB"}},
			key:   "Field",
			value: "TAB",
		},
		{
			name:  "VideoFormat",
			field: VideoFormat{BaseField{Key: "Video", Value: "1080"}},
			key:   "Video",
			value: "1080",
		},
		{
			name:  "AudioFormat",
			field: AudioFormat{BaseField{Key: "Audio", Value: "48khz"}},
			key:   "Audio",
			value: "48khz",
		},
		{
			name:  "FrameRate",
			field: FrameRate{BaseField{Key: "FPS", Value: "23.98"}},
			key:   "FPS",
			value: "23.98",
		},
		{
			name:  "FilmFormat",
			field: FilmFormat{BaseField{Key: "Film", Value: "35mm"}},
			key:   "Film",
			value: "35mm",
		},
		{
			name:  "Tape",
			field: Tape{BaseField{Key: "Tape", Value: "A001"}},
			key:   "Tape",
			value: "A001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field.GetKey() != tt.key {
				t.Errorf("GetKey() = %v, want %v", tt.field.GetKey(), tt.key)
			}
			if tt.field.GetValue() != tt.value {
				t.Errorf("GetValue() = %v, want %v", tt.field.GetValue(), tt.value)
			}
		})
	}
}

func TestColumn(t *testing.T) {
	col := Column{
		Name:  "Scene",
		Order: 1,
	}

	if col.Name != "Scene" {
		t.Errorf("Name = %v, want %v", col.Name, "Scene")
	}

	if col.Order != 1 {
		t.Errorf("Order = %v, want %v", col.Order, 1)
	}
}

func TestValue(t *testing.T) {
	col := Column{Name: "Scene", Order: 1}

	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{
			name:     "StringValue",
			value:    StringValue{Column: col, Value: "1A"},
			expected: "1A",
		},
		{
			name:     "IntValue",
			value:    IntValue{Column: col, Value: 42},
			expected: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value.String() != tt.expected {
				t.Errorf("String() = %v, want %v", tt.value.String(), tt.expected)
			}
		})
	}
}

func TestObjectJSON(t *testing.T) {
	obj := Object{
		HeaderFields: []Field{
			FieldDelimiter{BaseField{Key: "Field", Value: "TAB"}},
			VideoFormat{BaseField{Key: "Video", Value: "1080"}},
		},
		FieldDelimiter: FieldDelimiter{BaseField{Key: "Field", Value: "TAB"}},
		VideoFormat:    VideoFormat{BaseField{Key: "Video", Value: "1080"}},
		AudioFormat:    AudioFormat{BaseField{Key: "Audio", Value: "48khz"}},
		FPS:            FrameRate{BaseField{Key: "FPS", Value: "23.98"}},
		FilmFormat:     FilmFormat{BaseField{Key: "Film", Value: "35mm"}},
		Tape:           Tape{BaseField{Key: "Tape", Value: "A001"}},
		Columns: []Column{
			{Name: "Scene", Order: 0},
			{Name: "Take", Order: 1},
		},
		Rows: []Row{
			{
				Columns: []Column{{Name: "Scene", Order: 0}, {Name: "Take", Order: 1}},
				ValueMap: map[Column]Value{
					{Name: "Scene", Order: 0}: StringValue{Column: Column{Name: "Scene", Order: 0}, Value: "1"},
					{Name: "Take", Order: 1}:  StringValue{Column: Column{Name: "Take", Order: 1}, Value: "1"},
				},
				Order: 0,
			},
		},
	}

	data, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("Failed to marshal Object: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify essential fields are present
	expectedFields := []string{
		"header_fields",
		"field_delimiter",
		"video_format",
		"audio_format",
		"fps",
		"film_format",
		"tape",
		"columns",
		"data",
	}

	for _, field := range expectedFields {
		if _, ok := unmarshaled[field]; !ok {
			t.Errorf("Missing field in JSON output: %s", field)
		}
	}
}

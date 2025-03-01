package libale

import (
	"testing"
)

func TestNew(t *testing.T) {
	handler := New()
	if handler == nil {
		t.Error("Expected non-nil Handler")
	}
}

func TestHandler_Read(t *testing.T) {
	handler := New()
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "valid ale data",
			input: `Heading
FIELD_DELIM	TABS

Column
Name	Scene	Take

Data
A001	1	1
`,
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := handler.Read(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && obj == nil {
				t.Error("Expected non-nil Object for valid input")
			}
		})
	}
}

func TestHandler_ReadFile(t *testing.T) {
	handler := New()

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		_, err := handler.ReadFile("testdata/nonexistent.ale")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	// Test a single valid file to verify handler behavior
	t.Run("valid file", func(t *testing.T) {
		// Using a known good file from samples
		obj, err := handler.ReadFile("../../samples/ALE/A001R1AA_AVID.ale")
		if err != nil {
			t.Errorf("Handler.ReadFile() error = %v", err)
			return
		}
		if obj == nil {
			t.Error("Expected non-nil Object for valid file")
			return
		}

		// Basic validation that handler processed the file
		if obj.FieldDelimiter.GetValue() == "" {
			t.Error("Expected non-empty field delimiter")
		}
		if len(obj.Columns) == 0 {
			t.Error("Expected at least one column")
		}
		if len(obj.Rows) == 0 {
			t.Error("Expected at least one row")
		}
	})
}

package ale

import (
	"strings"
	"testing"

	"lib-post-interchange/libale/types"
)

func TestWriteMatchesSample(t *testing.T) {
	// Read the sample file
	ale, err := ReadFile("../../../samples/ALE/A001R1AA_AVID.ale")
	if err != nil {
		t.Fatalf("Failed to read sample file: %v", err)
	}

	// Write the ALE object back to string
	output, err := Write(ale)
	if err != nil {
		t.Fatalf("Failed to write ALE object: %v", err)
	}

	// Read the output back as an ALE object to compare structures
	outputAle, err := Read(output)
	if err != nil {
		t.Fatalf("Failed to read written output: %v", err)
	}

	// Compare the two ALE objects
	compareALEObjects(t, ale, outputAle)
}

func compareALEObjects(t *testing.T, expected, actual *types.Object) {
	t.Helper()

	// Compare header fields
	if len(expected.HeaderFields) != len(actual.HeaderFields) {
		t.Errorf("Header fields length mismatch: expected %d, got %d", len(expected.HeaderFields), len(actual.HeaderFields))
		return
	}

	// Create maps for easier comparison
	expectedHeaders := make(map[string]string)
	actualHeaders := make(map[string]string)

	for _, field := range expected.HeaderFields {
		expectedHeaders[field.GetKey()] = field.GetValue()
	}
	for _, field := range actual.HeaderFields {
		actualHeaders[field.GetKey()] = field.GetValue()
	}

	for key, expectedValue := range expectedHeaders {
		if actualValue, ok := actualHeaders[key]; !ok {
			t.Errorf("Missing header field %s", key)
		} else if expectedValue != actualValue {
			t.Errorf("Header field %s value mismatch: expected %q, got %q", key, expectedValue, actualValue)
		}
	}

	// Compare columns
	if len(expected.Columns) != len(actual.Columns) {
		t.Errorf("Columns length mismatch: expected %d, got %d", len(expected.Columns), len(actual.Columns))
		return
	}

	for i, expectedCol := range expected.Columns {
		actualCol := actual.Columns[i]
		if expectedCol.Name != actualCol.Name {
			t.Errorf("Column %d name mismatch: expected %q, got %q", i, expectedCol.Name, actualCol.Name)
		}
		if expectedCol.Order != actualCol.Order {
			t.Errorf("Column %d order mismatch: expected %d, got %d", i, expectedCol.Order, actualCol.Order)
		}
	}

	// Compare rows
	if len(expected.Rows) != len(actual.Rows) {
		t.Errorf("Rows length mismatch: expected %d, got %d", len(expected.Rows), len(actual.Rows))
		return
	}

	for i, expectedRow := range expected.Rows {
		actualRow := actual.Rows[i]

		// Compare values for each column
		for _, col := range expected.Columns {
			expectedVal, expectedOk := expectedRow.ValueMap[col]
			actualVal, actualOk := actualRow.ValueMap[col]

			if !expectedOk && !actualOk {
				continue
			}
			if expectedOk != actualOk {
				t.Errorf("Row %d, column %q presence mismatch: expected %v, got %v", i, col.Name, expectedOk, actualOk)
				continue
			}
			if expectedVal.String() != actualVal.String() {
				t.Errorf("Row %d, column %q value mismatch: expected %q, got %q", i, col.Name, expectedVal.String(), actualVal.String())
			}
		}
	}
}

func TestWriteNilObject(t *testing.T) {
	_, err := Write(nil)
	if err == nil {
		t.Error("Expected error when writing nil object, got nil")
	}
	if !strings.Contains(err.Error(), "cannot write nil ALE object") {
		t.Errorf("Expected error message to contain 'cannot write nil ALE object', got %q", err.Error())
	}
}

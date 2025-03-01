package ale

import (
	"fmt"
	"os"
	"strings"

	"lib-post-interchange/libale/format"
	"lib-post-interchange/libale/types"
)

// WriteFile writes an ALE object to a file at the specified path.
func WriteFile(filepath string, ale *types.Object) error {
	data, err := Write(ale)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, []byte(data), 0644)
}

// Write converts an ALE object to its string representation in ALE format.
func Write(ale *types.Object) (string, error) {
	if ale == nil {
		return "", fmt.Errorf("cannot write nil ALE object")
	}

	var builder strings.Builder

	// Write Heading section
	builder.WriteString(format.Heading + "\n")

	// Write header fields
	for _, field := range ale.HeaderFields {
		builder.WriteString(fmt.Sprintf("%s\t%s\n", field.GetKey(), field.GetValue()))
	}
	builder.WriteString("\n")

	// Write Column section
	builder.WriteString(format.Column + "\n")

	// Write column names
	columnNames := make([]string, len(ale.Columns))
	for _, col := range ale.Columns {
		columnNames[col.Order] = col.Name
	}
	builder.WriteString(strings.Join(columnNames, "\t") + "\n\n")

	// Write Data section
	builder.WriteString(format.Data + "\n")

	// Write rows
	for _, row := range ale.Rows {
		values := make([]string, len(ale.Columns))
		for _, col := range ale.Columns {
			if val, ok := row.ValueMap[col]; ok {
				values[col.Order] = val.String()
			}
		}
		builder.WriteString(strings.Join(values, "\t") + "\n")
	}

	return builder.String(), nil
}

// Package types provides core type definitions for the ALE format.
package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"lib-post-interchange/libale/validate"
)

// Field defines the common behavior for all ALE field types.
type Field interface {
	GetKey() string
	GetValue() string
}

// BaseField represents a base type for header fields.
type BaseField struct {
	Key   string
	Value string
}

func (f BaseField) GetKey() string   { return f.Key }
func (f BaseField) GetValue() string { return f.Value }

// FieldDelimiter represents the field delimiter value in the header.
type FieldDelimiter struct{ BaseField }

// VideoFormat represents the video format value in the header.
type VideoFormat struct{ BaseField }

// AudioFormat represents the audio format value in the header.
type AudioFormat struct{ BaseField }

// FrameRate represents the framerate value in the header.
type FrameRate struct{ BaseField }

// FilmFormat represents the film format value in the header.
type FilmFormat struct{ BaseField }

// Tape represents the tape name value in the header.
type Tape struct{ BaseField }

// Column represents a column in the ALE data table.
type Column struct {
	Name  string
	Order int
}

// Row represents a row in the ALE data table.
type Row struct {
	Columns  []Column
	ValueMap map[Column]Value
	Order    int
}

// Value represents a value in the ALE data table.
type Value interface {
	String() string
}

// StringValue represents a string value.
type StringValue struct {
	Column Column
	Value  string
}

func (v StringValue) String() string { return v.Value }

// IntValue represents an integer value.
type IntValue struct {
	Column Column
	Value  int
}

func (v IntValue) String() string { return fmt.Sprintf("%d", v.Value) }

// Object represents a structured Avid Log Exchange file.
type Object struct {
	HeaderFields   []Field
	FieldDelimiter FieldDelimiter
	VideoFormat    VideoFormat
	AudioFormat    AudioFormat
	FPS            FrameRate
	FilmFormat     FilmFormat
	Tape           Tape
	Columns        []Column
	Rows           []Row
}

// MarshalJSON implements the json.Marshaler interface.
func (o Object) MarshalJSON() ([]byte, error) {
	// Validate required fields
	if o.Columns == nil {
		return nil, validate.ErrValueNilColumns
	}
	if o.Rows == nil {
		return nil, validate.ErrValueNilRows
	}

	type jsonObject struct {
		HeaderFields   []map[string]string `json:"header_fields,omitempty"`
		FieldDelimiter string              `json:"field_delimiter,omitempty"`
		VideoFormat    string              `json:"video_format,omitempty"`
		AudioFormat    string              `json:"audio_format,omitempty"`
		FPS            string              `json:"fps,omitempty"`
		FilmFormat     string              `json:"film_format,omitempty"`
		Tape           string              `json:"tape,omitempty"`
		Columns        []string            `json:"columns,omitempty"`
		Data           []map[string]string `json:"data,omitempty"`
	}

	// Convert header fields to map
	headerFields := make([]map[string]string, 0, len(o.HeaderFields))
	for _, field := range o.HeaderFields {
		if field == nil {
			continue // Skip nil fields
		}
		headerFields = append(headerFields, map[string]string{
			"key":   field.GetKey(),
			"value": field.GetValue(),
		})
	}

	// Convert columns to string slice
	columns := make([]string, len(o.Columns))
	for i, col := range o.Columns {
		if col.Name == "" {
			return nil, validate.ErrValueEmptyColumnName
		}
		columns[i] = col.Name
	}

	// Convert rows to map slice
	data := make([]map[string]string, len(o.Rows))
	for i, row := range o.Rows {
		if row.ValueMap == nil {
			return nil, validate.ErrValueNilRowMap
		}
		rowData := make(map[string]string)
		for col, val := range row.ValueMap {
			if val == nil {
				continue // Skip nil values
			}
			rowData[col.Name] = val.String()
		}
		data[i] = rowData
	}

	// Create JSON object
	obj := jsonObject{
		HeaderFields: headerFields,
		Columns:      columns,
		Data:         data,
	}

	// Add optional fields only if they have valid values
	if v := o.FieldDelimiter.GetValue(); v != "" {
		obj.FieldDelimiter = v
	}
	if v := o.VideoFormat.GetValue(); v != "" {
		obj.VideoFormat = v
	}
	if v := o.AudioFormat.GetValue(); v != "" {
		obj.AudioFormat = v
	}
	if v := o.FPS.GetValue(); v != "" {
		obj.FPS = v
	}
	if v := o.FilmFormat.GetValue(); v != "" {
		obj.FilmFormat = v
	}
	if v := o.Tape.GetValue(); v != "" {
		obj.Tape = v
	}

	return json.Marshal(obj)
}

// String returns a string representation of the Object.
func (o Object) String() string {
	// Format columns and rows for display, limiting output length
	var columnsDisplay, rowsDisplay string

	if len(o.Columns) > 0 {
		columnNames := make([]string, 0, len(o.Columns))
		for _, col := range o.Columns {
			columnNames = append(columnNames, col.Name)
		}
		if len(columnNames) > 3 {
			columnsDisplay = fmt.Sprintf("%s, %s, %s, ...", columnNames[0], columnNames[1], columnNames[2])
		} else {
			columnsDisplay = strings.Join(columnNames, ", ")
		}
	}

	if len(o.Rows) > 0 && len(o.Columns) > 0 {
		rowValues := make([]string, 0, 1)
		row := o.Rows[0]
		for i := 0; i < len(o.Columns) && i < 1; i++ {
			if val, ok := row.ValueMap[o.Columns[i]]; ok {
				rowValues = append(rowValues, val.String())
			}
		}
		if len(o.Columns) > 1 {
			rowsDisplay = fmt.Sprintf("%s, ...", rowValues[0])
		} else {
			rowsDisplay = strings.Join(rowValues, ", ")
		}
	}

	// Build the base format string
	format := `ALE{
    Columns: %v [%v],
    Rows: %v [%v]`

	// Add fields when they are defined
	if o.FieldDelimiter.GetValue() != "" {
		format += `,
    FieldDelimiter: %v`
	}
	if o.VideoFormat.GetValue() != "" {
		format += `,
    VideoFormat: %v`
	}
	if o.AudioFormat.GetValue() != "" {
		format += `,
    AudioFormat: %v`
	}
	if o.FilmFormat.GetValue() != "" {
		format += `,
    FilmFormat: %v`
	}
	if o.Tape.GetValue() != "" {
		format += `,
    Tape: %v`
	}
	if o.FPS.GetValue() != "" {
		format += `,
    FPS: %v`
	}
	format += "\n}"

	// Build args slice starting with required fields
	args := []interface{}{
		len(o.Columns),
		columnsDisplay,
		len(o.Rows),
		rowsDisplay,
	}

	// Add optional fields to args if they are defined
	if o.FieldDelimiter.GetValue() != "" {
		args = append(args, o.FieldDelimiter.GetValue())
	}
	if o.VideoFormat.GetValue() != "" {
		args = append(args, o.VideoFormat.GetValue())
	}
	if o.AudioFormat.GetValue() != "" {
		args = append(args, o.AudioFormat.GetValue())
	}
	if o.FilmFormat.GetValue() != "" {
		args = append(args, o.FilmFormat.GetValue())
	}
	if o.Tape.GetValue() != "" {
		args = append(args, o.Tape.GetValue())
	}
	if o.FPS.GetValue() != "" {
		args = append(args, o.FPS.GetValue())
	}

	return fmt.Sprintf(format, args...)
}

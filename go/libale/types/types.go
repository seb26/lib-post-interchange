// Package types provides core type definitions for the ALE format.
package types

import "fmt"

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

// String returns a string representation of the Object.
func (o Object) String() string {
	return fmt.Sprintf("ALE{HeaderFields: %v, FieldDelimiter: %v, VideoFormat: %v, AudioFormat: %v, FPS: %v, FilmFormat: %v, Tape: %v, Columns: %v, Rows: %v}",
		o.HeaderFields, o.FieldDelimiter, o.VideoFormat, o.AudioFormat, o.FPS, o.FilmFormat, o.Tape, len(o.Columns), len(o.Rows))
}

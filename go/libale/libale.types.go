package libale

import "fmt"

// ALEField interface defines the common behavior for all ALE field types
type ALEField interface {
	GetKey() string
	GetValue() string
}

// ALEBaseField represents a base type for header fields
type ALEBaseField struct {
	Key   string
	Value string
}

// GetKey returns the key of the header field
func (f ALEBaseField) GetKey() string {
	return f.Key
}

// GetValue returns the value of the header field
func (f ALEBaseField) GetValue() string {
	return f.Value
}

// ALEFieldDelimiter represents the field delimiter value in the header
type ALEFieldDelimiter struct {
	Key   string
	Value string
}

// GetKey returns the key of the field delimiter
func (f ALEFieldDelimiter) GetKey() string {
	return f.Key
}

// GetValue returns the value of the field delimiter
func (f ALEFieldDelimiter) GetValue() string {
	return f.Value
}

// ALEFilmFormat represents the film format value in the header
type ALEFilmFormat struct {
	Key   string
	Value string
}

// GetKey returns the key of the film format
func (f ALEFilmFormat) GetKey() string {
	return f.Key
}

// GetValue returns the value of the film format
func (f ALEFilmFormat) GetValue() string {
	return f.Value
}

// ALETape represents the tape name value in the header
type ALETape struct {
	Key   string
	Value string
}

// GetKey returns the key of the tape name
func (f ALETape) GetKey() string {
	return f.Key
}

// GetValue returns the value of the tape name
func (f ALETape) GetValue() string {
	return f.Value
}

// ALEVideoFormat represents the video format value in the header
type ALEVideoFormat struct {
	Key   string
	Value string
}

// GetKey returns the key of the video format
func (f ALEVideoFormat) GetKey() string {
	return f.Key
}

// GetValue returns the value of the video format
func (f ALEVideoFormat) GetValue() string {
	return f.Value
}

// ALEAudioFormat represents the audio format value in the header
type ALEAudioFormat struct {
	Key   string
	Value string
}

// GetKey returns the key of the audio format
func (f ALEAudioFormat) GetKey() string {
	return f.Key
}

// GetValue returns the value of the audio format
func (f ALEAudioFormat) GetValue() string {
	return f.Value
}

// ALEFrameRate represents the framerate value in the header
type ALEFrameRate struct {
	Key   string
	Value string
}

// GetKey returns the key of the frame rate
func (f ALEFrameRate) GetKey() string {
	return f.Key
}

// GetValue returns the value of the frame rate
func (f ALEFrameRate) GetValue() string {
	return f.Value
}

// ToType returns a function that creates an ALEBaseField instance
func ToType(key string) (func(string) ALEBaseField, error) {
	switch key {
	case "FIELD_DELIM":
		return func(value string) ALEBaseField {
			return ALEBaseField{Key: "FIELD_DELIM", Value: value}
		}, nil
	case "VIDEO_FORMAT":
		return func(value string) ALEBaseField {
			return ALEBaseField{Key: "VIDEO_FORMAT", Value: value}
		}, nil
	case "AUDIO_FORMAT":
		return func(value string) ALEBaseField {
			return ALEBaseField{Key: "AUDIO_FORMAT", Value: value}
		}, nil
	case "FPS":
		return func(value string) ALEBaseField {
			return ALEBaseField{Key: "FPS", Value: value}
		}, nil
	case "FILM_FORMAT":
		return func(value string) ALEBaseField {
			return ALEBaseField{Key: "FILM_FORMAT", Value: value}
		}, nil
	case "TAPE":
		return func(value string) ALEBaseField {
			return ALEBaseField{Key: "TAPE", Value: value}
		}, nil
	default:
		return func(value string) ALEBaseField {
			return ALEBaseField{Key: key, Value: value}
		}, nil
	}
}

// ALEObject is a structured representation of an Avid Log Exchange file
type ALEObject struct {
	HeaderFields   []ALEField
	FieldDelimiter ALEFieldDelimiter
	VideoFormat    ALEVideoFormat
	AudioFormat    ALEAudioFormat
	FPS            ALEFrameRate
	FilmFormat     ALEFilmFormat
	Tape           ALETape
	Columns        []ALEColumn
	Rows           []ALERow
}

// String returns the string representation of the ALEObject
func (ale ALEObject) String() string {
	return fmt.Sprintf("ALEObject{HeaderFields: %v, FieldDelimiter: %v, VideoFormat: %v, AudioFormat: %v, FPS: %v, FilmFormat: %v, Tape: %v, Columns: %v, Rows: %v}",
		ale.HeaderFields, ale.FieldDelimiter, ale.VideoFormat, ale.AudioFormat, ale.FPS, ale.FilmFormat, ale.Tape, len(ale.Columns), len(ale.Rows))
}

// AssignHeaderFieldsToObject() assigns the header fields to the outside
// of the ALE object.
func AssignHeaderFieldsToObject(ale ALEObject) ALEObject {
	for _, field := range ale.HeaderFields {
		switch field.GetKey() {
		case "FIELD_DELIM":
			ale.FieldDelimiter = ALEFieldDelimiter{Key: field.GetKey(), Value: field.GetValue()}
		case "FPS":
			ale.FPS = ALEFrameRate{Key: field.GetKey(), Value: field.GetValue()}
		case "AUDIO_FORMAT":
			ale.AudioFormat = ALEAudioFormat{Key: field.GetKey(), Value: field.GetValue()}
		case "VIDEO_FORMAT":
			ale.VideoFormat = ALEVideoFormat{Key: field.GetKey(), Value: field.GetValue()}
		case "FILM_FORMAT":
			ale.FilmFormat = ALEFilmFormat{Key: field.GetKey(), Value: field.GetValue()}
		case "TAPE":
			ale.Tape = ALETape{Key: field.GetKey(), Value: field.GetValue()}
		}
	}
	return ale
}

// ALEColumn represents a column in the ALE data table
type ALEColumn struct {
	Name  string
	Order int
}

// ALERow represents a row in the ALE data table
type ALERow struct {
	Columns  []ALEColumn
	ValueMap map[ALEColumn]ALEValueString
	Order    int
}

// ALEValueString represents a string value
type ALEValueString struct {
	Column ALEColumn
	Value  string
}

// String returns the string value
func (v ALEValueString) String() string {
	return v.Value
}

// ALEValueInt represents an int value
type ALEValueInt struct {
	Column ALEColumn
	Value  int
}

// String returns the string value
func (v ALEValueInt) String() string {
	return fmt.Sprintf("%d", v.Value)
}

package libale

import "fmt"

// ALEHeaderField represents a base type for header fields
type ALEHeaderField struct {
	Key   string
	Value string
}

// ALEFieldDelimiter represents the field delimiter value in the header
type ALEFieldDelimiter struct {
	ALEHeaderField
}

// ToType returns a function that creates an ALEHeaderField instance
func ToType(key string) (func(string) ALEHeaderField, error) {
	switch key {
	case "FIELD_DELIM":
		return func(value string) ALEHeaderField {
			return NewALEFieldDelimiter(value).ALEHeaderField
		}, nil
	case "VIDEO_FORMAT":
		return func(value string) ALEHeaderField {
			return NewALEVideoFormat(value).ALEHeaderField
		}, nil
	case "AUDIO_FORMAT":
		return func(value string) ALEHeaderField {
			return NewALEAudioFormat(value).ALEHeaderField
		}, nil
	case "FPS":
		return func(value string) ALEHeaderField {
			return NewALEFrameRate(value).ALEHeaderField
		}, nil
	default:
		return func(value string) ALEHeaderField {
			return ALEHeaderField{Key: key, Value: value}
		}, nil
	}
}

// AssignHeaderFieldsToObject() assigns the header fields to the outside
// of the ALE object.
func AssignHeaderFieldsToObject(ale ALEObject) ALEObject {
	for _, field := range ale.HeaderFields {
		switch field.Key {
		case "FIELD_DELIM":
			ale.FieldDelimiter = ALEFieldDelimiter{ALEHeaderField: field}
		case "FPS":
			ale.FPS = ALEFrameRate{ALEHeaderField: field}
		case "AUDIO_FORMAT":
			ale.AudioFormat = ALEAudioFormat{ALEHeaderField: field}
		case "VIDEO_FORMAT":
			ale.VideoFormat = ALEVideoFormat{ALEHeaderField: field}
		}
	}
	return ale
}

// NewALEFieldDelimiter creates a new ALEFieldDelimiter with the Key set to "FIELD_DELIM"
func NewALEFieldDelimiter(value string) ALEFieldDelimiter {
	return ALEFieldDelimiter{
		ALEHeaderField: ALEHeaderField{
			Key:   "FIELD_DELIM",
			Value: value,
		},
	}
}

// ALEVideoFormat represents the video format value in the header
type ALEVideoFormat struct {
	ALEHeaderField
}

// NewALEVideoFormat creates a new ALEVideoFormat with the Key set to "VIDEO_FORMAT"
func NewALEVideoFormat(value string) ALEVideoFormat {
	return ALEVideoFormat{
		ALEHeaderField: ALEHeaderField{
			Key:   "VIDEO_FORMAT",
			Value: value,
		},
	}
}

// ALEAudioFormat represents the audio format value in the header
type ALEAudioFormat struct {
	ALEHeaderField
}

// NewALEAudioFormat creates a new ALEAudioFormat with the Key set to "AUDIO_FORMAT"
func NewALEAudioFormat(value string) ALEAudioFormat {
	return ALEAudioFormat{
		ALEHeaderField: ALEHeaderField{
			Key:   "AUDIO_FORMAT",
			Value: value,
		},
	}
}

// ALEFrameRate represents the framerate value in the header
type ALEFrameRate struct {
	ALEHeaderField
}

// NewALEFrameRate creates a new ALEFrameRate with the Key set to "FPS"
func NewALEFrameRate(value string) ALEFrameRate {
	return ALEFrameRate{
		ALEHeaderField: ALEHeaderField{
			Key:   "FPS",
			Value: value,
		},
	}
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

// ALEObject is a structured representation of an Avid Log Exchange file
type ALEObject struct {
	HeaderFields   []ALEHeaderField
	FieldDelimiter ALEFieldDelimiter
	VideoFormat    ALEVideoFormat
	AudioFormat    ALEAudioFormat
	FPS            ALEFrameRate
	Columns        []ALEColumn
	Rows           []ALERow
}

// String returns the string representation of the ALEObject
func (ale ALEObject) String() string {
	return fmt.Sprintf("ALEObject{HeaderFields: %v, FieldDelimiter: %v, VideoFormat: %v, AudioFormat: %v, FPS: %v, Columns: %v, Rows: %v}",
		ale.HeaderFields, ale.FieldDelimiter, ale.VideoFormat, ale.AudioFormat, ale.FPS, len(ale.Columns), len(ale.Rows))
}

package libale

// ALEHeadingWord represents the first word in the ALE file
const ALEHeadingWord = "Heading"

// ALEColumnWord represents the columns heading in the ALE file
const ALEColumnWord = "Column"

// ALEDataWord represents the data heading in the ALE file
const ALEDataWord = "Data"

// ALEHeadingWordPattern represents regexp for this section
const ALEHeadingWordPattern = `(?ms)^Heading(\r?\n|\r)(?P<fields>.*)(\r?\n|\r)(\r?\n|\r)Column(\r?\n|\r)(?P<columns>.*)(\r?\n|\r)(\r?\n|\r)Data$`

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

// String returns the string value
func (fd ALEFieldDelimiter) String() string {
	if fd.Value == "TABS" {
		return "TABS"
	} else {
		if LIBALE_OUTPUT_STRICT_VALUES {
			// For greater compatibility, always specify tabs.
			return "TABS"
		} else {
			return fd.Value
		}
	}
}

// ALE video formats
//const (
//	ALE_VIDEO_FORMAT_CUSTOM ALEVideoFormat = "CUSTOM"
//	ALE_VIDEO_FORMAT_1080   ALEVideoFormat = "1080"
//)

// String returns the string value
func (vf ALEVideoFormat) String() string {
	if vf.Value == "1080" {
		return vf.Value
	} else if vf.Value == "CUSTOM" {
		return "CUSTOM"
	} else {
		if LIBALE_OUTPUT_STRICT_VALUES {
			// For greater compatibility, return CUSTOM.
			return "CUSTOM"
		}
	}
	return vf.Value
}

// ALE audio formats
//const (
//	ALE_AUDIO_FORMAT_CUSTOM ALEAudioFormat = "CUSTOM"
//	ALE_AUDIO_FORMAT_48KHZ  ALEAudioFormat = "48kHz"
//	ALE_AUDIO_FORMAT_NONE   ALEAudioFormat = "NONE"
//)

// String returns the string value
func (af ALEAudioFormat) String() string {
	if af.Value == "48kHz" {
		return "48kHz"
	} else if af.Value == "CUSTOM" {
		return "CUSTOM"
	} else if af.Value == "NONE" {
		return "NONE"
	} else {
		if LIBALE_OUTPUT_STRICT_VALUES {
			// For greater compatibility, return CUSTOM.
			return "CUSTOM"
		}
	}
	return af.Value
}

// ALE frame rates
//const (
//	ALE_FRAME_RATE_CUSTOM ALEFrameRate = iota
//	ALE_FRAME_RATE_60
//	ALE_FRAME_RATE_5994_NDF
//	ALE_FRAME_RATE_5994_DF
//	ALE_FRAME_RATE_50
//	ALE_FRAME_RATE_48
//	ALE_FRAME_RATE_30_NDF
//	ALE_FRAME_RATE_30_DF
//	ALE_FRAME_RATE_2997_NDF
//	ALE_FRAME_RATE_2997_DF
//	ALE_FRAME_RATE_25
//	ALE_FRAME_RATE_24
//	ALE_FRAME_RATE_23976
//)

// String returns the string value
func (fr ALEFrameRate) String() string {
	switch fr.Value {
	case "60":
		return "60"
	case "59.94NDF":
		return "59.94NDF"
	case "59.94DF":
		return "59.94DF"
	case "50":
		return "50"
	case "48":
		return "48"
	case "30NDF":
		return "30NDF"
	case "30DF":
		return "30DF"
	case "29.97NDF":
		return "29.97NDF"
	case "29.97DF":
		return "29.97DF"
	case "25":
		return "25"
	case "24":
		return "24"
	case "23.976":
		return "23.976"
	default:
		return "CUSTOM"
	}
}

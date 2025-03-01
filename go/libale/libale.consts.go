package libale

// ALEHeadingWord represents the first word in the ALE file
const ALEHeadingWord = "Heading"

// ALEColumnWord represents the columns heading in the ALE file
const ALEColumnWord = "Column"

// ALEDataWord represents the data heading in the ALE file
const ALEDataWord = "Data"

// Predefined frame rates
// TODO - DROPFRAMES
var (
	FrameRate23976 = ALEFrameRate{Key: "FPS", Value: "23.976"}
	FrameRate24    = ALEFrameRate{Key: "FPS", Value: "24"}
	FrameRate25    = ALEFrameRate{Key: "FPS", Value: "25"}
	FrameRate2997  = ALEFrameRate{Key: "FPS", Value: "29.97"}
	FrameRate30    = ALEFrameRate{Key: "FPS", Value: "30"}
	FrameRate48    = ALEFrameRate{Key: "FPS", Value: "48"}
	FrameRate50    = ALEFrameRate{Key: "FPS", Value: "50"}
	FrameRate5994  = ALEFrameRate{Key: "FPS", Value: "59.94"}
	FrameRate60    = ALEFrameRate{Key: "FPS", Value: "60"}
)

// Predefined video formats
var (
	VideoFormatHD1080 = ALEVideoFormat{Key: "VIDEO_FORMAT", Value: "1080"}
	VideoFormatPAL    = ALEVideoFormat{Key: "VIDEO_FORMAT", Value: "PAL"}
	VideoFormatNTSC   = ALEVideoFormat{Key: "VIDEO_FORMAT", Value: "NTSC"}
	VideoFormatCUSTOM = ALEVideoFormat{Key: "VIDEO_FORMAT", Value: "CUSTOM"}
)

// Predefined audio formats
var (
	AudioFormatPCM48 = ALEAudioFormat{Key: "AUDIO_FORMAT", Value: "48kHz"}
)

// Predefined field delimiters
var (
	FieldDelimiterTab = ALEFieldDelimiter{Key: "FIELD_DELIM", Value: "TABS"}
)

// Predefined film formats
var (
	FilmFormat16mm = ALEFilmFormat{Key: "FILM_FORMAT", Value: "16 mm"}
	FilmFormat35mm = ALEFilmFormat{Key: "FILM_FORMAT", Value: "35 mm"}
	FilmFormat65mm = ALEFilmFormat{Key: "FILM_FORMAT", Value: "65 mm"}
)

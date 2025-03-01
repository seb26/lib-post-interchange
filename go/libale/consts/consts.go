// Package consts provides constant definitions for the libale package.
// These constants define standard values and formats used in ALE files.
package consts

import "lib-post-interchange/libale/types"

// ALEHeadingWord represents the first word in the ALE file
const ALEHeadingWord = "Heading"

// ALEColumnWord represents the columns heading in the ALE file
const ALEColumnWord = "Column"

// ALEDataWord represents the data heading in the ALE file
const ALEDataWord = "Data"

// Predefined frame rates
// TODO - DROPFRAMES
var (
	FrameRate23976 = types.ALEFrameRate{Key: "FPS", Value: "23.976"}
	FrameRate24    = types.ALEFrameRate{Key: "FPS", Value: "24"}
	FrameRate25    = types.ALEFrameRate{Key: "FPS", Value: "25"}
	FrameRate2997  = types.ALEFrameRate{Key: "FPS", Value: "29.97"}
	FrameRate30    = types.ALEFrameRate{Key: "FPS", Value: "30"}
	FrameRate48    = types.ALEFrameRate{Key: "FPS", Value: "48"}
	FrameRate50    = types.ALEFrameRate{Key: "FPS", Value: "50"}
	FrameRate5994  = types.ALEFrameRate{Key: "FPS", Value: "59.94"}
	FrameRate60    = types.ALEFrameRate{Key: "FPS", Value: "60"}
)

// Predefined video formats
var (
	VideoFormatHD1080 = types.ALEVideoFormat{Key: "VIDEO_FORMAT", Value: "1080"}
	VideoFormatPAL    = types.ALEVideoFormat{Key: "VIDEO_FORMAT", Value: "PAL"}
	VideoFormatNTSC   = types.ALEVideoFormat{Key: "VIDEO_FORMAT", Value: "NTSC"}
	VideoFormatCUSTOM = types.ALEVideoFormat{Key: "VIDEO_FORMAT", Value: "CUSTOM"}
)

// Predefined audio formats
var (
	AudioFormatPCM48 = types.ALEAudioFormat{Key: "AUDIO_FORMAT", Value: "48kHz"}
)

// Predefined field delimiters
var (
	FieldDelimiterTab = types.ALEFieldDelimiter{Key: "FIELD_DELIM", Value: "TABS"}
)

// Predefined film formats
var (
	FilmFormat16mm = types.ALEFilmFormat{Key: "FILM_FORMAT", Value: "16 mm"}
	FilmFormat35mm = types.ALEFilmFormat{Key: "FILM_FORMAT", Value: "35 mm"}
	FilmFormat65mm = types.ALEFilmFormat{Key: "FILM_FORMAT", Value: "65 mm"}
)

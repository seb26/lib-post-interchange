// Package format provides predefined format constants and helpers for ALE files.
package format

import "lib-post-interchange/libale/types"

// Section headers
const (
	Heading = "Heading"
	Column  = "Column"
	Data    = "Data"
)

// Frame rates
var (
	FPS23_976 = types.FrameRate{BaseField: types.BaseField{Key: "FPS", Value: "23.976"}}
	FPS24     = types.FrameRate{BaseField: types.BaseField{Key: "FPS", Value: "24"}}
	FPS25     = types.FrameRate{BaseField: types.BaseField{Key: "FPS", Value: "25"}}
	FPS29_97  = types.FrameRate{BaseField: types.BaseField{Key: "FPS", Value: "29.97"}}
	FPS30     = types.FrameRate{BaseField: types.BaseField{Key: "FPS", Value: "30"}}
	FPS48     = types.FrameRate{BaseField: types.BaseField{Key: "FPS", Value: "48"}}
	FPS50     = types.FrameRate{BaseField: types.BaseField{Key: "FPS", Value: "50"}}
	FPS59_94  = types.FrameRate{BaseField: types.BaseField{Key: "FPS", Value: "59.94"}}
	FPS60     = types.FrameRate{BaseField: types.BaseField{Key: "FPS", Value: "60"}}
)

// Video formats
var (
	VideoHD1080 = types.VideoFormat{BaseField: types.BaseField{Key: "VIDEO_FORMAT", Value: "1080"}}
	VideoPAL    = types.VideoFormat{BaseField: types.BaseField{Key: "VIDEO_FORMAT", Value: "PAL"}}
	VideoNTSC   = types.VideoFormat{BaseField: types.BaseField{Key: "VIDEO_FORMAT", Value: "NTSC"}}
	VideoCustom = types.VideoFormat{BaseField: types.BaseField{Key: "VIDEO_FORMAT", Value: "CUSTOM"}}
)

// Audio formats
var (
	AudioPCM48 = types.AudioFormat{BaseField: types.BaseField{Key: "AUDIO_FORMAT", Value: "48kHz"}}
)

// Field delimiters
var (
	DelimiterTab = types.FieldDelimiter{BaseField: types.BaseField{Key: "FIELD_DELIM", Value: "TABS"}}
)

// Film formats
var (
	Film16mm = types.FilmFormat{BaseField: types.BaseField{Key: "FILM_FORMAT", Value: "16 mm"}}
	Film35mm = types.FilmFormat{BaseField: types.BaseField{Key: "FILM_FORMAT", Value: "35 mm"}}
	Film65mm = types.FilmFormat{BaseField: types.BaseField{Key: "FILM_FORMAT", Value: "65 mm"}}
)

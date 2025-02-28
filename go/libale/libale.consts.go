package libale

// ALEHeadingWord represents the first word in the ALE file
const ALEHeadingWord = "Heading"

// ALEColumnWord represents the columns heading in the ALE file
const ALEColumnWord = "Column"

// ALEDataWord represents the data heading in the ALE file
const ALEDataWord = "Data"

// ALEHeadingWordPattern represents regexp from beginning of file
// to the end of the Data section including subsequent blank line.
const ALEHeadingWordPattern = `(?ms)^Heading(\r?\n|\r)(?P<fields>.*)(\r?\n|\r)(\r?\n|\r)Column(\r?\n|\r)(?P<columns>.*)(\r?\n|\r)(\r?\n|\r)(?P<data_header>Data)$`

// ALE video formats
//const (
//	ALE_VIDEO_FORMAT_CUSTOM ALEVideoFormat = "CUSTOM"
//	ALE_VIDEO_FORMAT_1080   ALEVideoFormat = "1080"
//)
// TODO: add PAL, NTSC, and gather other possibilities

// ALE audio formats
//const (
//	ALE_AUDIO_FORMAT_CUSTOM ALEAudioFormat = "CUSTOM"
//	ALE_AUDIO_FORMAT_48KHZ  ALEAudioFormat = "48kHz"
//	ALE_AUDIO_FORMAT_NONE   ALEAudioFormat = "NONE"
//)
// TODO: add 44khz and others

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

// TODO: HEADER FIELDS TO ADD
// 'TAPE' - string, any custom string tape name
// 'FILM_FORMAT' - string, "35 mm"

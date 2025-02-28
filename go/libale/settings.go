package libale

// LIBALE_OUTPUT_STRICT_VALUES is a runtime flag to control certain values.
// By default (true), the most compatible and generic value for the following
// fields is used, instead of the actual user defined value. This will prevent
// a failure when opening the file in legacy applications that do not expect
// the range of possible values a user could specify using this library.
const (
	LIBALE_OUTPUT_STRICT_VALUES = true
)

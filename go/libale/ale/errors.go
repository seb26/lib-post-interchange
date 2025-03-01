package ale

import "fmt"

// ErrorCategory represents the main category of an error
type ErrorCategory int32

// Error categories
const (
	CategorySection ErrorCategory = iota + 1
	CategoryParse
	CategoryFormat
	CategoryField
)

// Section error subcategories
const (
	SectionMissing = iota
	SectionMalformed
	SectionIncomplete
)

// Parse error subcategories
const (
	ParseFailed = iota
	ParseInvalid
	ParseIncomplete
)

// Format error subcategories
const (
	FormatInvalid = iota
	FormatUnsupported
	FormatMalformed
)

// Field error subcategories
const (
	FieldMissing = iota
	FieldInvalid
	FieldDuplicate
)

// Error represents an ALE parsing error.
type Error struct {
	Line    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ale: [%d.0] %s", e.Line, e.Message)
}

// Code returns the unique error code
func (e *Error) Code() int32 {
	return int32(e.Line)
}

// Error definitions for ALE parsing
var (
	// Section errors
	ErrSectionMissingHeading = &Error{
		Line:    1,
		Message: "missing 'Heading' section",
	}
	ErrSectionMissingColumn = &Error{
		Line:    2,
		Message: "missing 'Column' section",
	}
	ErrSectionMissingData = &Error{
		Line:    3,
		Message: "missing 'Data' section",
	}
	ErrSectionIncompleteColumn = &Error{
		Line:    2,
		Message: "incomplete 'Column' section",
	}

	// Format errors
	ErrFormatMalformedColumn = &Error{
		Line:    2,
		Message: "malformed column section",
	}

	// Parse errors
	ErrParseFailedHeader = &Error{
		Line:    2,
		Message: "failed to parse header fields",
	}
	ErrParseFailedColumns = &Error{
		Line:    2,
		Message: "failed to parse columns",
	}
	ErrParseFailedData = &Error{
		Line:    3,
		Message: "failed to parse data rows",
	}
	ErrParseFailedContent = &Error{
		Line:    1,
		Message: "failed to parse file content",
	}

	// Field errors
	ErrFieldInvalidHeader = &Error{
		Line:    2,
		Message: "invalid header field type",
	}
)

// IsCategory checks if an error belongs to a specific category
func IsCategory(err error, category ErrorCategory) bool {
	if aleErr, ok := err.(*Error); ok {
		return aleErr.Line == int(category)
	}
	return false
}

// IsError checks if an error matches a specific category and subcategory
func IsError(err error, category ErrorCategory, subCategory int32) bool {
	if aleErr, ok := err.(*Error); ok {
		return aleErr.Line == int(category) && aleErr.Line == int(subCategory)
	}
	return false
}

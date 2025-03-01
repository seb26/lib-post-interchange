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
	ParseEmpty
	ParseMismatch
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
	Category    ErrorCategory
	SubCategory int32
	Message     string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ale: [%d.%d] %s", e.Category, e.SubCategory, e.Message)
}

// Code returns the unique error code
func (e *Error) Code() int32 {
	return int32(e.Category)*1000 + e.SubCategory
}

// WithContext returns a new Error with additional context appended to the message
func (e *Error) WithContext(context string) *Error {
	return &Error{
		Category:    e.Category,
		SubCategory: e.SubCategory,
		Message:     e.Message + ": " + context,
	}
}

// Error definitions for ALE parsing
var (
	// Section errors
	ErrSectionMissingHeading = &Error{
		Category:    CategorySection,
		SubCategory: SectionMissing,
		Message:     "missing 'Heading' section",
	}
	ErrSectionMissingColumn = &Error{
		Category:    CategorySection,
		SubCategory: SectionMissing,
		Message:     "missing 'Column' section",
	}
	ErrSectionMissingData = &Error{
		Category:    CategorySection,
		SubCategory: SectionMissing,
		Message:     "missing 'Data' section",
	}
	ErrSectionIncompleteColumn = &Error{
		Category:    CategorySection,
		SubCategory: SectionIncomplete,
		Message:     "incomplete 'Column' section",
	}

	// Format errors
	ErrFormatMalformedColumn = &Error{
		Category:    CategoryFormat,
		SubCategory: FormatMalformed,
		Message:     "malformed column section",
	}

	// Parse errors
	ErrParseFailedHeader = &Error{
		Category:    CategoryParse,
		SubCategory: ParseFailed,
		Message:     "failed to parse header fields",
	}
	ErrParseFailedColumns = &Error{
		Category:    CategoryParse,
		SubCategory: ParseFailed,
		Message:     "failed to parse columns",
	}
	ErrParseFailedData = &Error{
		Category:    CategoryParse,
		SubCategory: ParseFailed,
		Message:     "failed to parse data rows",
	}
	ErrParseFailedContent = &Error{
		Category:    CategoryParse,
		SubCategory: ParseFailed,
		Message:     "failed to parse file content",
	}
	ErrParseFailedRows = &Error{
		Category:    CategoryParse,
		SubCategory: ParseFailed,
		Message:     "failed to create rows",
	}
	ErrParseEmptyInput = &Error{
		Category:    CategoryParse,
		SubCategory: ParseEmpty,
		Message:     "empty input",
	}
	ErrParseMismatchedColumns = &Error{
		Category:    CategoryParse,
		SubCategory: ParseMismatch,
		Message:     "row has mismatched column count",
	}

	// Field errors
	ErrFieldInvalidHeader = &Error{
		Category:    CategoryField,
		SubCategory: FieldInvalid,
		Message:     "invalid header field type",
	}
)

// IsCategory checks if an error belongs to a specific category
func IsCategory(err error, category ErrorCategory) bool {
	if aleErr, ok := err.(*Error); ok {
		return aleErr.Code()/1000 == int32(category)
	}
	return false
}

// IsError checks if an error matches a specific category and subcategory
func IsError(err error, category ErrorCategory, subCategory int32) bool {
	if aleErr, ok := err.(*Error); ok {
		code := aleErr.Code()
		return code/1000 == int32(category) && code%1000 == subCategory
	}
	return false
}

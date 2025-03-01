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
	CategoryValue
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

// Value error subcategories
const (
	ValueInvalid = iota
	ValueOutOfRange
	ValueUnsupported
)

// Error represents the base error type for all ALE-related errors
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

// Section errors
var (
	ErrSectionMissingHeading = &Error{
		Category:    CategorySection,
		SubCategory: SectionMissing,
		Message:     "missing Heading section",
	}
	ErrSectionMissingColumn = &Error{
		Category:    CategorySection,
		SubCategory: SectionMissing,
		Message:     "missing Column section",
	}
	ErrSectionMissingData = &Error{
		Category:    CategorySection,
		SubCategory: SectionMissing,
		Message:     "missing Data section",
	}
	ErrSectionIncompleteColumn = &Error{
		Category:    CategorySection,
		SubCategory: SectionIncomplete,
		Message:     "missing column names",
	}
)

// Parse errors
var (
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
		Message:     "error reading file content",
	}
)

// Format errors
var (
	ErrFormatMalformedColumn = &Error{
		Category:    CategoryFormat,
		SubCategory: FormatMalformed,
		Message:     "expected empty line after columns",
	}
)

// Field errors
var (
	ErrFieldInvalidHeader = &Error{
		Category:    CategoryField,
		SubCategory: FieldInvalid,
		Message:     "invalid header field type",
	}
)

// IsCategory checks if an error belongs to a specific category
func IsCategory(err error, category ErrorCategory) bool {
	if aleErr, ok := err.(*Error); ok {
		return aleErr.Category == category
	}
	return false
}

// IsError checks if an error matches a specific category and subcategory
func IsError(err error, category ErrorCategory, subCategory int32) bool {
	if aleErr, ok := err.(*Error); ok {
		return aleErr.Category == category && aleErr.SubCategory == subCategory
	}
	return false
}

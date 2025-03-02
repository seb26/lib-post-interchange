package errors

import "fmt"

// ErrorCategory represents the main category of an error
type ErrorCategory int32

// Error categories
const (
	CategoryInput ErrorCategory = iota + 1
	CategoryOutput
)

// Error represents an ALE error.
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
	// Input errors
	ErrInputMissingHeading = &Error{
		Category: CategoryInput,
		Message:  "missing 'Heading' section",
	}
	ErrInputMissingColumn = &Error{
		Category: CategoryInput,
		Message:  "missing 'Column' section",
	}
	ErrInputMissingData = &Error{
		Category: CategoryInput,
		Message:  "missing 'Data' section",
	}
	ErrInputIncompleteColumn = &Error{
		Category: CategoryInput,
		Message:  "incomplete 'Column' section",
	}
	ErrInputMalformedColumn = &Error{
		Category: CategoryInput,
		Message:  "malformed column section",
	}
	ErrInputFailedHeading = &Error{
		Category: CategoryInput,
		Message:  "failed to parse header fields",
	}
	ErrInputFailedColumns = &Error{
		Category: CategoryInput,
		Message:  "failed to parse columns",
	}
	ErrInputFailedData = &Error{
		Category: CategoryInput,
		Message:  "failed to parse data rows",
	}
	ErrInputFailedContent = &Error{
		Category: CategoryInput,
		Message:  "failed to parse file content",
	}
	ErrInputFailedRows = &Error{
		Category: CategoryInput,
		Message:  "failed to create rows",
	}
	ErrInputEmpty = &Error{
		Category: CategoryInput,
		Message:  "empty input",
	}
	ErrInputMismatchedColumns = &Error{
		Category: CategoryInput,
		Message:  "row has mismatched column count",
	}

	// Output errors
	ErrOutputNilObject = &Error{
		Category: CategoryOutput,
		Message:  "cannot write nil ALE object",
	}
	ErrOutputNilColumns = &Error{
		Category: CategoryOutput,
		Message:  "columns cannot be nil",
	}
	ErrOutputNilRows = &Error{
		Category: CategoryOutput,
		Message:  "rows cannot be nil",
	}
	ErrOutputEmptyColumnName = &Error{
		Category: CategoryOutput,
		Message:  "column name cannot be empty",
	}
	ErrOutputNilRowMap = &Error{
		Category: CategoryOutput,
		Message:  "row value map cannot be nil",
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

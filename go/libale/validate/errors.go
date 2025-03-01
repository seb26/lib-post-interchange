// Package validate provides validation-related errors and functions.
package validate

import "fmt"

// ErrorCategory represents the main category of an error
type ErrorCategory int32

// Error categories
const (
	CategoryValue ErrorCategory = iota + 1
)

// Value error subcategories
const (
	ValueNil = iota
	ValueInvalid
)

// Error represents the base error type for all validation-related errors
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

// Value errors
var (
	ErrValueNilColumns = &Error{
		Category:    CategoryValue,
		SubCategory: ValueNil,
		Message:     "columns cannot be nil",
	}
	ErrValueNilRows = &Error{
		Category:    CategoryValue,
		SubCategory: ValueNil,
		Message:     "rows cannot be nil",
	}
	ErrValueEmptyColumnName = &Error{
		Category:    CategoryValue,
		SubCategory: ValueInvalid,
		Message:     "column name cannot be empty",
	}
	ErrValueNilRowMap = &Error{
		Category:    CategoryValue,
		SubCategory: ValueNil,
		Message:     "row value map cannot be nil",
	}
)

// IsCategory checks if an error belongs to a specific category
func IsCategory(err error, category ErrorCategory) bool {
	if valErr, ok := err.(*Error); ok {
		return valErr.Category == category
	}
	return false
}

// IsError checks if an error matches a specific category and subcategory
func IsError(err error, category ErrorCategory, subCategory int32) bool {
	if valErr, ok := err.(*Error); ok {
		return valErr.Category == category && valErr.SubCategory == subCategory
	}
	return false
}

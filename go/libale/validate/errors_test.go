package validate

import "testing"

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "nil columns error",
			err: &Error{
				Category:    CategoryValue,
				SubCategory: ValueNil,
				Message:     "columns cannot be nil",
			},
			want: "ale: [1.0] columns cannot be nil",
		},
		{
			name: "invalid column name error",
			err: &Error{
				Category:    CategoryValue,
				SubCategory: ValueInvalid,
				Message:     "column name cannot be empty",
			},
			want: "ale: [1.1] column name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Code(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want int32
	}{
		{
			name: "value nil error code",
			err: &Error{
				Category:    CategoryValue,
				SubCategory: ValueNil,
			},
			want: 1000, // Category(1) * 1000 + SubCategory(0)
		},
		{
			name: "value invalid error code",
			err: &Error{
				Category:    CategoryValue,
				SubCategory: ValueInvalid,
			},
			want: 1001, // Category(1) * 1000 + SubCategory(1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Code(); got != tt.want {
				t.Errorf("Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCategory(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		category ErrorCategory
		want     bool
	}{
		{
			name:     "matching category",
			err:      ErrValueNilColumns,
			category: CategoryValue,
			want:     true,
		},
		{
			name:     "non-matching category",
			err:      ErrValueNilColumns,
			category: ErrorCategory(999),
			want:     false,
		},
		{
			name:     "nil error",
			err:      nil,
			category: CategoryValue,
			want:     false,
		},
		{
			name:     "non-validate error",
			err:      &Error{Category: ErrorCategory(999)},
			category: CategoryValue,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCategory(tt.err, tt.category); got != tt.want {
				t.Errorf("IsCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		category    ErrorCategory
		subCategory int32
		want        bool
	}{
		{
			name:        "exact match",
			err:         ErrValueNilColumns,
			category:    CategoryValue,
			subCategory: ValueNil,
			want:        true,
		},
		{
			name:        "wrong category",
			err:         ErrValueNilColumns,
			category:    ErrorCategory(999),
			subCategory: ValueNil,
			want:        false,
		},
		{
			name:        "wrong subcategory",
			err:         ErrValueNilColumns,
			category:    CategoryValue,
			subCategory: ValueInvalid,
			want:        false,
		},
		{
			name:        "nil error",
			err:         nil,
			category:    CategoryValue,
			subCategory: ValueNil,
			want:        false,
		},
		{
			name:        "non-validate error",
			err:         &Error{Category: ErrorCategory(999), SubCategory: 999},
			category:    CategoryValue,
			subCategory: ValueNil,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsError(tt.err, tt.category, tt.subCategory); got != tt.want {
				t.Errorf("IsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "ErrValueNilColumns",
			err:  ErrValueNilColumns,
			want: "ale: [1.0] columns cannot be nil",
		},
		{
			name: "ErrValueNilRows",
			err:  ErrValueNilRows,
			want: "ale: [1.0] rows cannot be nil",
		},
		{
			name: "ErrValueEmptyColumnName",
			err:  ErrValueEmptyColumnName,
			want: "ale: [1.1] column name cannot be empty",
		},
		{
			name: "ErrValueNilRowMap",
			err:  ErrValueNilRowMap,
			want: "ale: [1.0] row value map cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

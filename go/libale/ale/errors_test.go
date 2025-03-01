package ale

import "testing"

func TestErrorWithContext(t *testing.T) {
	tests := []struct {
		name    string
		err     *Error
		context string
		want    string
	}{
		{
			name: "add context to parse error",
			err: &Error{
				Category:    CategoryParse,
				SubCategory: ParseMismatch,
				Message:     "row has mismatched column count",
			},
			context: "row 1 has 2 columns, expected 3",
			want:    "ale: [2.4] row has mismatched column count: row 1 has 2 columns, expected 3",
		},
		{
			name: "add context to empty input error",
			err: &Error{
				Category:    CategoryParse,
				SubCategory: ParseEmpty,
				Message:     "empty input",
			},
			context: "no data provided",
			want:    "ale: [2.3] empty input: no data provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.WithContext(tt.context)
			if got.Error() != tt.want {
				t.Errorf("Error() = %v, want %v", got.Error(), tt.want)
			}
			// Check that category and subcategory are preserved
			if got.Category != tt.err.Category {
				t.Errorf("Category = %v, want %v", got.Category, tt.err.Category)
			}
			if got.SubCategory != tt.err.SubCategory {
				t.Errorf("SubCategory = %v, want %v", got.SubCategory, tt.err.SubCategory)
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
			name:     "parse error matches parse category",
			err:      ErrParseEmptyInput,
			category: CategoryParse,
			want:     true,
		},
		{
			name:     "section error matches section category",
			err:      ErrSectionMissingHeading,
			category: CategorySection,
			want:     true,
		},
		{
			name:     "error with context matches category",
			err:      ErrParseMismatchedColumns.WithContext("row 1 has 2 columns, expected 3"),
			category: CategoryParse,
			want:     true,
		},
		{
			name:     "wrong category does not match",
			err:      ErrParseEmptyInput,
			category: CategorySection,
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
			err:         ErrParseEmptyInput,
			category:    CategoryParse,
			subCategory: ParseEmpty,
			want:        true,
		},
		{
			name:        "with context still matches",
			err:         ErrParseEmptyInput.WithContext("no data"),
			category:    CategoryParse,
			subCategory: ParseEmpty,
			want:        true,
		},
		{
			name:        "wrong category",
			err:         ErrParseEmptyInput,
			category:    CategorySection,
			subCategory: ParseEmpty,
			want:        false,
		},
		{
			name:        "wrong subcategory",
			err:         ErrParseEmptyInput,
			category:    CategoryParse,
			subCategory: ParseFailed,
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

// Package ale implements reading and writing of Avid Log Exchange (ALE) files.
package ale

import "lib-post-interchange/libale/types"

// Section headers
const (
	Heading = "Heading"
	Column  = "Column"
	Data    = "Data"
)

// Handler provides the main interface for interacting with ALE files.
// It encapsulates all operations related to reading, writing, and manipulating ALE data.
type Handler struct{}

// New creates a new Handler instance that provides access to all ALE operations.
func New() *Handler {
	return &Handler{}
}

// ReadFile reads and parses an ALE file from the filesystem.
func (h *Handler) ReadFile(filepath string) (*types.Object, error) {
	return ReadFile(filepath)
}

// Read parses ALE data from a string.
func (h *Handler) Read(input string) (*types.Object, error) {
	return Read(input)
}

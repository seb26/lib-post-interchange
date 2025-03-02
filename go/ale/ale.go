package ale

import (
	"ale/types"
)

// Handler provides the main interface for interacting with ALE files.
// It encapsulates all operations related to reading, writing, and manipulating ALE data.
type Handler struct{}

// New creates a new Handler instance that provides access to all ALE operations.
func New() *Handler {
	return &Handler{}
}

// ReadFile provides the main entry point for loading ALE data from the filesystem.
// It returns a structured representation of the ALE file's contents.
func (h *Handler) ReadFile(filepath string) (*types.Object, error) {
	return ReadFile(filepath)
}

// Read serves as the primary interface for parsing ALE data from any string source.
func (h *Handler) Read(input string) (*types.Object, error) {
	return Read(input)
}

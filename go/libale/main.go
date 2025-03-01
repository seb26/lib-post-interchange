package libale

import "lib-post-interchange/libale/types"

// ALEHandler provides the main interface for interacting with Avid Log Exchange files.
// It encapsulates all operations related to reading, writing, and manipulating ALE data.
type ALEHandler struct{}

// New creates a new ALEHandler instance that provides access to all ALE operations.
func New() *ALEHandler {
	return &ALEHandler{}
}

// ReadFile provides the main entry point for loading ALE data from the filesystem.
// It returns a structured representation of the ALE file's contents.
func (h *ALEHandler) ReadFile(filepath string) (*types.ALEObject, error) {
	return ReadFile(filepath)
}

// Read serves as the primary interface for parsing ALE data from any string source.
func (h *ALEHandler) Read(input string) (*types.ALEObject, error) {
	return Read(input)
}

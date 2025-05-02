// Package read implements a file reading tool for Gendo.
// It provides functionality to read content from files, with support for
// base path configuration and proper error handling.
package read

import (
	"fmt"
	"os"
	"path/filepath"

	"gendo/pkg/log"
)

// ReadTool implements the tools.Tool interface for file reading
type ReadTool struct {
	basePath string
}

// NewReadTool creates a new file reading tool
func NewReadTool(basePath string) *ReadTool {
	log.Debug("Creating new read tool with base path: %q", basePath)
	return &ReadTool{
		basePath: basePath,
	}
}

// Process implements the tools.Tool interface for ReadTool
func (t *ReadTool) Process(input string) (string, error) {
	log.Debug("Processing read input: %q", input)

	if input == "" {
		log.Debug("Empty input provided")
		return "", fmt.Errorf("no file path provided")
	}

	filePath := input
	if t.basePath != "" {
		filePath = filepath.Join(t.basePath, input)
		log.Debug("Using full file path: %q", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Debug("Failed to read file %q: %v", filePath, err)
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	log.Debug("Successfully read %d bytes from %q", len(content), filePath)
	return string(content), nil
}

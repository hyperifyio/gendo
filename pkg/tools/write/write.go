package write

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gendo/pkg/log"
)

// WriteTool implements the tools.Tool interface for file writing
type WriteTool struct {
	basePath string
}

// NewWriteTool creates a new file writing tool
func NewWriteTool(basePath string) *WriteTool {
	log.Debug("Creating new write tool with base path: %q", basePath)
	return &WriteTool{
		basePath: basePath,
	}
}

// Process implements the tools.Tool interface for WriteTool
func (t *WriteTool) Process(input string) (string, error) {
	log.Debug("Processing write input: %q", input)

	// Split input into file path and content
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		log.Debug("Invalid input format")
		return "", fmt.Errorf("invalid input format: expected 'path content'")
	}

	filePath := parts[0]
	content := parts[1]

	if t.basePath != "" {
		filePath = filepath.Join(t.basePath, filePath)
		log.Debug("Using full file path: %q", filePath)
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Debug("Failed to write to file %q: %v", filePath, err)
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	log.Debug("Successfully wrote %d bytes to %q", len(content), filePath)
	return fmt.Sprintf("Successfully wrote to %s", filePath), nil
}

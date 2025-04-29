package rand

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"gendo/pkg/log"
)

// Tool implements the tools.Tool interface for random number generation
type Tool struct {
	rand *rand.Rand
}

// New creates a new random number generator tool
func New() *Tool {
	log.Debug("Creating new random number generator tool")
	source := rand.NewSource(time.Now().UnixNano())
	return &Tool{
		rand: rand.New(source),
	}
}

// Process implements the tools.Tool interface
func (t *Tool) Process(input string) (string, error) {
	log.Debug("Processing random input: %q", input)
	
	// Parse the max number
	max, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		log.Debug("Failed to parse max number %q: %v", input, err)
		return "", fmt.Errorf("invalid max number: %v", err)
	}
	
	if max <= 0 {
		log.Debug("Invalid max number: %d (must be positive)", max)
		return "", fmt.Errorf("max number must be positive")
	}
	
	// Generate random number
	result := t.rand.Int63n(max)
	output := fmt.Sprintf("%d", result)
	
	log.Debug("Generated random number: %s (max: %d)", output, max)
	return output, nil
} 
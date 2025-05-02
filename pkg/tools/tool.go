// Package tools provides the interface and registry for Gendo tools.
// It defines the Tool interface that all tools must implement and provides
// a Registry for managing and accessing tools by name.
package tools

// Tool represents a Gendo tool that can process input and produce output
type Tool interface {
	// Process takes input text and returns output text and an optional error
	Process(input string) (string, error)
}

// Registry is a map of tool names to their implementations
type Registry map[string]Tool

// NewRegistry creates a new empty tool registry
func NewRegistry() Registry {
	return make(Registry)
}

// Register adds a tool to the registry
func (r Registry) Register(name string, tool Tool) {
	r[name] = tool
}

// Get returns a tool by name, or nil if not found
func (r Registry) Get(name string) Tool {
	return r[name]
}

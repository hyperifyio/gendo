// Package llm provides the interface and registry for language models in Gendo.
// It defines the LLM interface that all language model implementations must follow
// and provides a Registry for managing and accessing different LLM implementations.
package llm

// LLM represents a language model that can process prompts and inputs
type LLM interface {
	// Process takes a system prompt and user input, returns the model's response
	Process(prompt, input string) (string, error)
}

// Registry is a map of LLM names to their implementations
type Registry map[string]LLM

// NewRegistry creates a new empty LLM registry
func NewRegistry() Registry {
	return make(Registry)
}

// Register adds an LLM to the registry
func (r Registry) Register(name string, llm LLM) {
	r[name] = llm
}

// Get returns an LLM by name, or nil if not found
func (r Registry) Get(name string) LLM {
	return r[name]
}

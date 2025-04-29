package gendo

import (
	"testing"

	"gendo/pkg/llm"
	"gendo/pkg/tools"
)

// mockLLM implements the llm.LLM interface for testing
type mockLLM struct {
	lastPrompt string
	lastInput  string
	response   string
}

func (m *mockLLM) Process(prompt, input string) (string, error) {
	m.lastPrompt = prompt
	m.lastInput = input
	return m.response, nil
}

func TestProcessNode(t *testing.T) {
	tests := []struct {
		name           string
		node           Node
		input          string
		expectedPrompt string
		expectedInput  string
	}{
		{
			name: "Full prompt with instruction and examples",
			node: Node{
				ID:     0,
				Prompt: "Extract the mathematical operation from the text. Return only the mathematical expression, removing any natural language. For example: \"What is 1 plus 1?\" -> \"1 + 1\", \"Calculate 5 times 3\" -> \"5 * 3\", \"Divide 10 by 2\" -> \"10 / 2\"",
			},
			input:          "What is 1+1?",
			expectedPrompt: "Extract the mathematical operation from the text. Return only the mathematical expression, removing any natural language. For example: \"What is 1 plus 1?\" -> \"1 + 1\", \"Calculate 5 times 3\" -> \"5 * 3\", \"Divide 10 by 2\" -> \"10 / 2\"",
			expectedInput:  "What is 1+1?",
		},
		{
			name: "Prompt with only instruction",
			node: Node{
				ID:     1,
				Prompt: "Format the calculation result in a natural language response.",
			},
			input:          "1 + 1 = 2",
			expectedPrompt: "Format the calculation result in a natural language response.",
			expectedInput:  "1 + 1 = 2",
		},
		{
			name: "Prompt with only examples",
			node: Node{
				ID:     2,
				Prompt: "For example: \"1 + 1 = 2\" -> \"The sum of 1 and 1 is 2.\"",
			},
			input:          "5 * 3 = 15",
			expectedPrompt: "For example: \"1 + 1 = 2\" -> \"The sum of 1 and 1 is 2.\"",
			expectedInput:  "5 * 3 = 15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock LLM
			mock := &mockLLM{response: "test response"}
			llmRegistry := llm.NewRegistry()
			llmRegistry.Register("openai", mock)

			// Create empty tool registry
			toolRegistry := tools.NewRegistry()

			// Process the node
			_, err := processNode(tt.node, tt.input, toolRegistry, llmRegistry)
			if err != nil {
				t.Errorf("processNode() error = %v", err)
				return
			}

			// Verify the prompt was passed correctly
			if mock.lastPrompt != tt.expectedPrompt {
				t.Errorf("processNode() prompt = %q, want %q", mock.lastPrompt, tt.expectedPrompt)
			}

			// Verify the input was passed correctly
			if mock.lastInput != tt.expectedInput {
				t.Errorf("processNode() input = %q, want %q", mock.lastInput, tt.expectedInput)
			}
		})
	}
} 
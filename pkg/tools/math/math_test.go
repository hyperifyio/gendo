package math

import (
	"testing"
	"strings"
)

func TestExtractFirstExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple expression",
			input:    "1+1",
			expected: "1+1",
		},
		{
			name:     "Expression with spaces",
			input:    "1 + 1",
			expected: "1+1",
		},
		{
			name:     "Expression with natural language prefix",
			input:    "What is 1 + 1?",
			expected: "1+1",
		},
		{
			name:     "Expression with natural language suffix",
			input:    "1 + 1 equals 2",
			expected: "1+1",
		},
		{
			name:     "Expression with word operators",
			input:    "1 plus 1",
			expected: "1+1",
		},
		{
			name:     "Expression with mixed operators",
			input:    "1 plus 2 * 3",
			expected: "1+2*3",
		},
		{
			name:     "Expression with negative numbers",
			input:    "-1 + -2",
			expected: "-1+-2",
		},
		{
			name:     "Expression with decimal numbers",
			input:    "1.5 * 2.3",
			expected: "1.5*2.3",
		},
		{
			name:     "Expression with quotes",
			input:    "\"1 + 1\" is 2",
			expected: "1+1",
		},
		{
			name:     "Expression with division",
			input:    "10 divided by 2",
			expected: "10/2",
		},
		{
			name:     "Expression with multiplication words",
			input:    "5 times 3",
			expected: "5*3",
		},
		{
			name:     "Expression with subtraction",
			input:    "5 minus 3",
			expected: "5-3",
		},
		{
			name:     "No valid expression",
			input:    "Hello world",
			expected: "",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFirstExpression(tt.input)
			if result != tt.expected {
				t.Errorf("extractFirstExpression(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantNum1       float64
		wantNum2       float64
		wantOperator   rune
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:         "Simple addition",
			input:        "1+2",
			wantNum1:     1,
			wantNum2:     2,
			wantOperator: '+',
		},
		{
			name:         "Addition with spaces",
			input:        "1 + 2",
			wantNum1:     1,
			wantNum2:     2,
			wantOperator: '+',
		},
		{
			name:         "Subtraction",
			input:        "5-3",
			wantNum1:     5,
			wantNum2:     3,
			wantOperator: '-',
		},
		{
			name:         "Multiplication",
			input:        "4*6",
			wantNum1:     4,
			wantNum2:     6,
			wantOperator: '*',
		},
		{
			name:         "Division",
			input:        "8/2",
			wantNum1:     8,
			wantNum2:     2,
			wantOperator: '/',
		},
		{
			name:         "Negative numbers",
			input:        "-1+-2",
			wantNum1:     -1,
			wantNum2:     -2,
			wantOperator: '+',
		},
		{
			name:           "No operator",
			input:          "123",
			wantErr:        true,
			wantErrMessage: "no valid operator found",
		},
		{
			name:           "Invalid first number",
			input:          "abc+2",
			wantErr:        true,
			wantErrMessage: "invalid first number",
		},
		{
			name:           "Invalid second number",
			input:          "1+def",
			wantErr:        true,
			wantErrMessage: "invalid second number",
		},
		{
			name:           "Empty input",
			input:          "",
			wantErr:        true,
			wantErrMessage: "no valid operator found",
		},
		{
			name:           "Multiple operators",
			input:          "1+2+3",
			wantNum1:       1,
			wantNum2:       2,
			wantOperator:   '+',
		},
		{
			name:         "Decimal numbers",
			input:        "1.5*2.3",
			wantNum1:     1.5,
			wantNum2:     2.3,
			wantOperator: '*',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			num1, num2, operator, err := parseExpression(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseExpression(%q) expected error containing %q, got nil", tt.input, tt.wantErrMessage)
					return
				}
				if tt.wantErrMessage != "" && !strings.Contains(err.Error(), tt.wantErrMessage) {
					t.Errorf("parseExpression(%q) error = %v, want error containing %q", tt.input, err, tt.wantErrMessage)
				}
				return
			}
			if err != nil {
				t.Errorf("parseExpression(%q) unexpected error: %v", tt.input, err)
				return
			}
			if num1 != tt.wantNum1 {
				t.Errorf("parseExpression(%q) num1 = %v, want %v", tt.input, num1, tt.wantNum1)
			}
			if num2 != tt.wantNum2 {
				t.Errorf("parseExpression(%q) num2 = %v, want %v", tt.input, num2, tt.wantNum2)
			}
			if operator != tt.wantOperator {
				t.Errorf("parseExpression(%q) operator = %q, want %q", tt.input, operator, tt.wantOperator)
			}
		})
	}
} 
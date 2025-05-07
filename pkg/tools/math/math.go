// Package math implements a mathematical expression evaluation tool for Gendo.
// It provides functionality to parse and evaluate basic arithmetic expressions,
// supporting addition, subtraction, multiplication, and division operations.
// The tool can extract mathematical expressions from natural language input.
package math

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"gendo/pkg/log"
)

// Tool implements the tools.Tool interface for math operations
type Tool struct{}

// NewTool creates a new math tool
func NewTool() *Tool {
	log.Debug("Creating new math tool")
	return &Tool{}
}

// extractFirstExpression extracts the first mathematical expression from the input
func extractFirstExpression(input string) string {
	input = strings.TrimSpace(input)

	// Remove quotes if present
	input = strings.Trim(input, "\"")

	// Try to extract a valid expression
	if expr, ok := tryExtractExpression(input); ok {
		return expr
	}

	return ""
}

// tryExtractExpression attempts to extract a valid mathematical expression from the input
func tryExtractExpression(input string) (string, bool) {
	var builder strings.Builder
	var lastChar rune
	var inNumber bool
	var hasOperator bool
	var foundDigit bool

	for i, char := range input {
		switch {
		case unicode.IsDigit(char) || char == '.':
			builder.WriteRune(char)
			inNumber = true
			foundDigit = true
		case char == '-':
			// Allow minus sign at start or after another operator
			if i == 0 || !unicode.IsDigit(rune(input[i-1])) {
				builder.WriteRune(char)
				inNumber = false
			} else if inNumber {
				builder.WriteRune(char)
				hasOperator = true
				inNumber = false
			}
		case char == '+' || char == '*' || char == '/':
			if foundDigit {
				builder.WriteRune(char)
				hasOperator = true
				inNumber = false
			}
		case unicode.IsSpace(char):
			continue
		default:
			if foundDigit && !inNumber && !hasOperator {
				continue
			}
			if !foundDigit {
				continue
			}
			if hasOperator && !inNumber {
				return "", false
			}
		}
		lastChar = char
	}

	result := builder.String()
	if result == "" {
		return "", false
	}

	// Remove trailing operator if present
	if lastChar == '+' || lastChar == '-' || lastChar == '*' || lastChar == '/' {
		result = result[:len(result)-1]
	}

	// Validate the expression
	if _, _, _, err := parseExpression(result); err != nil {
		return "", false
	}

	return result, true
}

// parseExpression parses a mathematical expression and returns the operands and operator
func parseExpression(expr string) (float64, float64, rune, error) {
	// Remove all spaces and quotes
	expr = strings.ReplaceAll(expr, " ", "")
	expr = strings.ReplaceAll(expr, "\"", "")

	// Find the first operator that's not a leading minus sign
	var operator rune
	var operatorIndex int = -1

	for i := 0; i < len(expr); i++ {
		c := rune(expr[i])
		if c == '+' || c == '*' || c == '/' || (c == '-' && i > 0 && expr[i-1] >= '0' && expr[i-1] <= '9') {
			operator = c
			operatorIndex = i
			break
		}
	}

	if operatorIndex == -1 {
		return 0, 0, 0, fmt.Errorf("no valid operator found")
	}

	// Split into operands
	first := expr[:operatorIndex]
	second := expr[operatorIndex+1:]

	// Parse operands
	num1, err := strconv.ParseFloat(first, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid first number: %v", err)
	}

	// Find the next operator in the second part
	nextOpIndex := -1
	for i := 0; i < len(second); i++ {
		c := rune(second[i])
		if c == '+' || c == '*' || c == '/' || (c == '-' && i > 0 && second[i-1] >= '0' && second[i-1] <= '9') {
			nextOpIndex = i
			break
		}
	}

	// If there's another operator, only take up to that point
	if nextOpIndex != -1 {
		second = second[:nextOpIndex]
	}

	num2, err := strconv.ParseFloat(second, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid second number: %v", err)
	}

	return num1, num2, operator, nil
}

// Process implements the tools.Tool interface
func (t *Tool) Process(input string) (string, error) {
	log.Debug("Processing math input: %q", input)

	// Extract expression
	expr := extractFirstExpression(input)
	if expr == "" {
		log.Debug("No valid mathematical expression found")
		return "", fmt.Errorf("no valid mathematical expression found")
	}

	log.Debug("Extracted expression: %q", expr)

	// Parse expression
	num1, num2, operator, err := parseExpression(expr)
	if err != nil {
		log.Debug("Failed to parse expression: %v", err)
		return "", err
	}

	// Perform calculation
	var result float64
	switch operator {
	case '+':
		result = num1 + num2
	case '-':
		result = num1 - num2
	case '*':
		result = num1 * num2
	case '/':
		if num2 == 0 {
			log.Debug("Division by zero attempted")
			return "", fmt.Errorf("division by zero")
		}
		result = num1 / num2
	default:
		log.Debug("Unsupported operator: %c", operator)
		return "", fmt.Errorf("unsupported operator: %c", operator)
	}

	output := fmt.Sprintf("%g", result)
	log.Debug("Math result: %s", output)
	return output, nil
}

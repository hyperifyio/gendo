// Package rand_test contains test cases for the random number generation tool.
// It includes unit tests for number generation, input validation, and
// distribution analysis to ensure proper random number generation.
package rand

import (
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	tool := New()
	if tool == nil {
		t.Error("New() returned nil")
	}
	if tool.rand == nil {
		t.Error("New() returned tool with nil rand")
	}
}

func TestProcess(t *testing.T) {
	tool := New()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid positive number",
			input:   "100",
			wantErr: false,
		},
		{
			name:    "zero input",
			input:   "0",
			wantErr: true,
		},
		{
			name:    "negative input",
			input:   "-10",
			wantErr: true,
		},
		{
			name:    "invalid input",
			input:   "not a number",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Process(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("Process() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Process() unexpected error: %v", err)
				return
			}

			// Verify the result is a valid number within range
			num, err := strconv.ParseInt(result, 10, 64)
			if err != nil {
				t.Errorf("Process() returned invalid number: %v", err)
				return
			}

			max, _ := strconv.ParseInt(tt.input, 10, 64)
			if num < 0 || num >= max {
				t.Errorf("Process() returned number %d outside valid range [0, %d)", num, max)
			}
		})
	}
}

func TestProcessDistribution(t *testing.T) {
	tool := New()
	max := int64(10)
	iterations := 1000
	counts := make(map[int64]int)

	// Generate many random numbers and count their distribution
	for i := 0; i < iterations; i++ {
		result, err := tool.Process(strconv.FormatInt(max, 10))
		if err != nil {
			t.Errorf("Process() unexpected error: %v", err)
			continue
		}

		num, _ := strconv.ParseInt(result, 10, 64)
		counts[num]++
	}

	// Check if the distribution is roughly uniform
	expectedCount := float64(iterations) / float64(max)
	tolerance := 0.3 // Allow 30% deviation from expected count

	for num, count := range counts {
		if num < 0 || num >= max {
			t.Errorf("Got number %d outside valid range [0, %d)", num, max)
			continue
		}

		deviation := float64(count) - expectedCount
		if deviation < 0 {
			deviation = -deviation
		}

		if deviation > expectedCount*tolerance {
			t.Errorf("Number %d appeared %d times, expected around %.1f (deviation too large)",
				num, count, expectedCount)
		}
	}
}

// Package write_test contains test cases for the file writing tool.
// It includes unit tests for various write operations and benchmarks
// for different file sizes. Tests use temporary directories to avoid
// affecting the actual filesystem.
package write

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteTool(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "gendo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		input    string
		basePath string
		want     string
		wantErr  bool
		check    func(t *testing.T, path string)
	}{
		{
			name:     "write new file",
			input:    "test.txt test content",
			basePath: tmpDir,
			want:     "Successfully wrote to",
			wantErr:  false,
			check: func(t *testing.T, path string) {
				content, err := os.ReadFile(filepath.Join(path, "test.txt"))
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
				}
				if string(content) != "test content" {
					t.Errorf("Written content = %v, want %v", string(content), "test content")
				}
			},
		},
		{
			name:     "invalid input - no content",
			input:    "test.txt",
			basePath: tmpDir,
			want:     "",
			wantErr:  true,
			check:    func(t *testing.T, path string) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := NewWriteTool(tt.basePath)
			got, err := tool.Process(tt.input)

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

			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("Process() = %v, want prefix %v", got, tt.want)
			}
			tt.check(t, tt.basePath)
		})
	}
}

func BenchmarkWriteTool(b *testing.B) {
	// Create a temporary directory for benchmark files
	tmpDir, err := os.MkdirTemp("", "gendo-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test data sizes
	sizes := map[string]int{
		"small":  100,     // 100 bytes
		"medium": 10000,   // 10KB
		"large":  1000000, // 1MB
	}

	tool := NewWriteTool(tmpDir)

	for name, size := range sizes {
		content := strings.Repeat("x", size)
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				filename := filepath.Join(tmpDir, fmt.Sprintf("bench_%d.txt", i))
				_, _ = tool.Process(fmt.Sprintf("%s %s", filename, content))
			}
		})
	}
}

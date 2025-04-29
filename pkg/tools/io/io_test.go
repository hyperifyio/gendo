package io

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadTool(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "gendo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testContent := "test content"
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		basePath string
		want     string
	}{
		{
			name:     "read existing file",
			input:    "test.txt",
			basePath: tmpDir,
			want:     testContent,
		},
		{
			name:     "read non-existent file",
			input:    "nonexistent.txt",
			basePath: tmpDir,
			want:     "ERROR: Failed to read file:",
		},
		{
			name:     "invalid input",
			input:    "file1.txt file2.txt",
			basePath: tmpDir,
			want:     "ERROR: Read tool requires a filename",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := NewReadTool(tt.basePath)
			got := tool.Process(tt.input)
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("ReadTool.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		check    func(t *testing.T, path string)
	}{
		{
			name:     "write new file",
			input:    "test.txt test content",
			basePath: tmpDir,
			want:     "Written to",
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
			want:     "ERROR: Write tool requires a filename and content",
			check:    func(t *testing.T, path string) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := NewWriteTool(tt.basePath)
			got := tool.Process(tt.input)
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("WriteTool.Process() = %v, want %v", got, tt.want)
			}
			tt.check(t, tt.basePath)
		})
	}
}

func BenchmarkReadTool(b *testing.B) {
	// Create a temporary directory for benchmark files
	tmpDir, err := os.MkdirTemp("", "gendo-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files of different sizes
	files := map[string]int{
		"small.txt":  100,     // 100 bytes
		"medium.txt": 10000,   // 10KB
		"large.txt":  1000000, // 1MB
	}

	for filename, size := range files {
		content := strings.Repeat("x", size)
		filepath := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			b.Fatalf("Failed to create benchmark file: %v", err)
		}
	}

	tool := NewReadTool(tmpDir)

	for filename := range files {
		b.Run(filename, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tool.Process(filename)
			}
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
				tool.Process(fmt.Sprintf("%s %s", filename, content))
			}
		})
	}
} 
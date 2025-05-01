package read

import (
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
		wantErr  bool
	}{
		{
			name:     "read existing file",
			input:    "test.txt",
			basePath: tmpDir,
			want:     testContent,
			wantErr:  false,
		},
		{
			name:     "read non-existent file",
			input:    "nonexistent.txt",
			basePath: tmpDir,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "invalid input",
			input:    "file1.txt file2.txt",
			basePath: tmpDir,
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := NewReadTool(tt.basePath)
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

			if got != tt.want {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
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
				_, _ = tool.Process(filename)
			}
		})
	}
}

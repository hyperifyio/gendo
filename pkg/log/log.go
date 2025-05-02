// Package log provides logging functionality for the Gendo tool.
// It supports different log levels (Debug, Info, Error) and includes
// caller context information in log messages. The package allows
// configuration of verbosity and output destination.
package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

var (
	verbose bool
	output  io.Writer = os.Stderr
)

// SetVerbose enables or disables verbose logging
func SetVerbose(v bool) {
	verbose = v
}

// SetOutput sets the output writer for logging
func SetOutput(w io.Writer) {
	output = w
}

// getCallerContext returns the file name and line number of the caller
func getCallerContext() string {
	_, file, line, ok := runtime.Caller(2) // Skip 2 frames to get the actual caller
	if !ok {
		return "unknown:0"
	}
	// Get just the file name without the full path
	file = filepath.Base(file)
	return fmt.Sprintf("%s:%d", file, line)
}

// Debug logs a debug message if verbose mode is enabled
func Debug(format string, args ...interface{}) {
	if verbose {
		context := getCallerContext()
		fmt.Fprintf(output, "DEBUG [%s]: "+format+"\n", append([]interface{}{context}, args...)...)
	}
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	context := getCallerContext()
	fmt.Fprintf(output, "INFO [%s]: "+format+"\n", append([]interface{}{context}, args...)...)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	context := getCallerContext()
	fmt.Fprintf(output, "ERROR [%s]: "+format+"\n", append([]interface{}{context}, args...)...)
}

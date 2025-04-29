package main

import (
	"flag"
	"fmt"
	"os"

	"gendo/internal/gendo"
	"gendo/pkg/log"
)

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	model := flag.String("model", "", "Model to use for LLM (overrides GENDO_MODEL environment variable)")
	flag.StringVar(model, "m", "", "Model to use for LLM (shorthand)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-verbose] [-m model] <script>\n", os.Args[0])
		os.Exit(1)
	}

	log.SetVerbose(*verbose)
	log.Debug("Verbose logging enabled")

	if err := gendo.Run(args[0], *model); err != nil {
		log.Error("Failed to run script: %v", err)
		os.Exit(1)
	}
} 
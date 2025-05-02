package gendo

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"gendo/pkg/llm"
	"gendo/pkg/llm/openai"
	"gendo/pkg/log"
	"gendo/pkg/parser"
	"gendo/pkg/tools"
	"gendo/pkg/tools/math"
	"gendo/pkg/tools/rand"
	readtool "gendo/pkg/tools/read"
	writetool "gendo/pkg/tools/write"
)

// NodeType represents the type of a node
type NodeType string

const (
	NodeTypeTool    NodeType = "tool"
	NodeTypeIn      NodeType = "in"
	NodeTypeOut     NodeType = "out"
	NodeTypeErr     NodeType = "err"
	NodeTypeDefault NodeType = ""
)

// Node represents a Gendo node with its ID, references, and optional prompt
type Node struct {
	ID     int
	Refs   []int
	Prompt string
	Tool   string
	Type   NodeType
}

// processNode processes input through a node, either using OpenAI API, tool, or passthrough
func processNode(node Node, input string, toolRegistry tools.Registry, llmRegistry llm.Registry) (string, error) {
	log.Debug("Processing node %d with input: %q", node.ID, input)

	switch node.Type {
	case NodeTypeTool:
		if tool := toolRegistry.Get(node.Tool); tool != nil {
			log.Debug("Using tool %q for node %d", node.Tool, node.ID)
			result, err := tool.Process(input)
			if err != nil {
				log.Debug("Tool %q failed: %v", node.Tool, err)
				return "", fmt.Errorf("tool %q failed: %v", node.Tool, err)
			}
			log.Debug("Tool %q returned: %q", node.Tool, result)
			return result, nil
		}
		return "", fmt.Errorf("unknown tool: %s", node.Tool)
	case NodeTypeIn:
		// Input nodes are handled separately in processInput
		return input, nil
	case NodeTypeOut:
		// Output nodes are handled separately in processInput
		return input, nil
	case NodeTypeErr:
		// Error nodes are handled separately in processInput
		return input, nil
	default:
		if node.Prompt != "" {
			// Use the OpenAI LLM for processing
			if llm := llmRegistry.Get("openai"); llm != nil {
				log.Debug("Using OpenAI LLM for node %d with prompt: %q", node.ID, node.Prompt)
				result, err := llm.Process(node.Prompt, input)
				if err != nil {
					log.Debug("OpenAI LLM failed: %v", err)
					return "", err
				}
				log.Debug("OpenAI LLM returned: %q", result)
				return result, nil
			}
			return "", fmt.Errorf("no LLM available")
		}
		log.Debug("Node %d is a passthrough node", node.ID)
		return input, nil // Passthrough for non-AI nodes
	}
}

// processInput processes a single input line according to Gendo rules
func processInput(line string, nodes map[int]Node, toolRegistry tools.Registry, llmRegistry llm.Registry, stdoutDefault, stderrDefault int, stdout, stderr io.Writer) error {
	log.Debug("Processing input line: %q", line)

	// Set up default I/O if not provided
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	// Find I/O nodes
	var inNode, outNode, errNode *Node
	for _, node := range nodes {
		switch node.Type {
		case NodeTypeIn:
			inNode = &node
		case NodeTypeOut:
			outNode = &node
		case NodeTypeErr:
			errNode = &node
		}
	}

	// Process through input node if defined
	if inNode != nil {
		log.Debug("Processing through input node %d", inNode.ID)
		output, err := processNode(*inNode, line, toolRegistry, llmRegistry)
		if err != nil {
			log.Error("Input node failed: %v", err)
			fmt.Fprintf(stderr, "Error: %v\n", err)
			return err
		}
		line = output
	}

	// Process through output node if defined
	if outNode != nil {
		log.Debug("Processing through output node %d", outNode.ID)
		output, err := processNode(*outNode, line, toolRegistry, llmRegistry)
		if err != nil {
			log.Error("Output node failed: %v", err)
			fmt.Fprintf(stderr, "Error: %v\n", err)
			return err
		}
		line = output
	}

	// Process through the chain of nodes defined in the script
	for nodeID := 3; nodeID >= 1; nodeID-- {
		if node, ok := nodes[nodeID]; ok {
			log.Debug("Processing through node %d", nodeID)
			output, err := processNode(node, line, toolRegistry, llmRegistry)
			if err != nil {
				if errNode != nil {
					log.Debug("Processing error through error node %d", errNode.ID)
					errOutput, _ := processNode(*errNode, err.Error(), toolRegistry, llmRegistry)
					fmt.Fprintln(stderr, errOutput)
				} else if stderrDefault > 0 {
					if errNode, ok := nodes[stderrDefault]; ok {
						log.Debug("Processing error through default error node %d", stderrDefault)
						errOutput, _ := processNode(errNode, err.Error(), toolRegistry, llmRegistry)
						fmt.Fprintln(stderr, errOutput)
					}
				}
				return err
			}
			line = output
		}
	}

	log.Debug("Final output: %q", line)
	fmt.Fprintln(stdout, line)
	return nil
}

// Run executes a Gendo script from a file
func Run(filename string, model string) error {
	log.Debug("Running script: %s", filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Error("Failed to open script: %v", err)
		return fmt.Errorf("failed to open script: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		log.Debug("Read line: %q", line)
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Error("Failed to read script: %v", err)
		return fmt.Errorf("failed to read script: %v", err)
	}

	log.Debug("Read %d lines from script", len(lines))

	// Parse script using the new parser
	nodes := make(map[int]Node)
	var inputLines []string
	defaultErrorNode := 0

	for _, line := range lines {
		result, ok := parser.ParseLine(line)
		if !ok {
			log.Debug("Failed to parse line: %q", line)
			continue
		}

		switch r := result.(type) {
		case *parser.NodeDefinition:
			log.Debug("Parsed node definition: ID=%d, Tool=%q", r.ID, r.Tool)
			nodes[r.ID] = Node{
				ID:     r.ID,
				Refs:   r.RefIDs,
				Prompt: r.Prompt,
				Tool:   r.Tool,
				Type:   NodeType(r.Type),
			}
		case *parser.RouteDefinition:
			if r.Source == 0 && r.Dest == 0 && r.ErrorDest == 0 {
				// This is an input line
				log.Debug("Parsed input line: %q", r.Input)
				inputLines = append(inputLines, r.Input)
			} else if r.ErrorDest > 0 {
				log.Debug("Setting default error node to %d", r.ErrorDest)
				defaultErrorNode = r.ErrorDest
			}
		}
	}

	// Initialize tool registry
	log.Debug("Initializing tool registry")
	toolRegistry := tools.NewRegistry()
	toolRegistry.Register("read", readtool.NewReadTool(""))
	toolRegistry.Register("write", writetool.NewWriteTool(""))
	toolRegistry.Register("math", math.NewTool())
	toolRegistry.Register("rand", rand.New())

	// Initialize LLM registry
	log.Debug("Initializing LLM registry")
	llmRegistry := llm.NewRegistry()
	llmRegistry.Register("openai", openai.New("", "bitnet"))

	// Check if we have input from stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		log.Debug("Reading input from pipe")
		// Input from pipe
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if err := processInput(scanner.Text(), nodes, toolRegistry, llmRegistry, 1, defaultErrorNode, os.Stdout, os.Stderr); err != nil {
				return err
			}
		}
		if err := scanner.Err(); err != nil {
			log.Error("Failed to read stdin: %v", err)
			return fmt.Errorf("failed to read stdin: %v", err)
		}
	} else {
		log.Debug("Processing %d script input lines", len(inputLines))
		// Process script input lines
		for _, line := range inputLines {
			if err := processInput(line, nodes, toolRegistry, llmRegistry, 1, defaultErrorNode, os.Stdout, os.Stderr); err != nil {
				return err
			}
		}
	}

	return nil
}

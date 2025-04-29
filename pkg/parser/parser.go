package parser

import (
	"strconv"
	"strings"
)

// NodeDefinition represents a parsed node from a script line
type NodeDefinition struct {
	ID      int
	RefIDs  []int    // Reference IDs this node can call
	Prompt  string   // Optional prompt text
	IsTool  bool     // Whether this is a tool node
	Tool    string   // Tool name if IsTool is true
}

// RouteDefinition represents a routing between nodes with optional error handling
type RouteDefinition struct {
	Source    int    // Source node ID
	Dest      int    // Destination node ID (for stdout)
	ErrorDest int    // Error destination node ID (optional)
	Input     string // Input text
}

// ParseLine parses a single line from a Gendo script
func ParseLine(line string) (interface{}, bool) {
	// Trim spaces and skip empty lines or comments
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil, false
	}

	// Check if this is a node definition (contains : but not <)
	if strings.Contains(line, ":") && !strings.Contains(line, "<") {
		return parseNodeDefinition(line)
	}

	// Otherwise it's a routing line
	return parseRouting(line)
}

// parseNodeDefinition parses a node definition line: nodeID : refID refID … [: prompt text]
func parseNodeDefinition(line string) (*NodeDefinition, bool) {
	// Split into parts by first colon
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return nil, false
	}

	// Parse node ID
	id, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, false
	}

	// Get the rest of the line after the first colon
	rest := strings.TrimSpace(parts[1])
	node := &NodeDefinition{ID: id}

	// Check if this is a tool definition
	if strings.HasPrefix(rest, "tool ") {
		node.IsTool = true
		node.Tool = strings.TrimSpace(strings.TrimPrefix(rest, "tool "))
		return node, true
	}

	// If there's a second colon, everything after it is the prompt
	if idx := strings.Index(rest, ":"); idx != -1 {
		refs := strings.Fields(rest[:idx])
		for _, ref := range refs {
			if refID, err := strconv.Atoi(ref); err == nil {
				node.RefIDs = append(node.RefIDs, refID)
			}
		}
		// Take everything after the first colon as the prompt
		node.Prompt = strings.TrimSpace(rest[idx+1:])
	} else {
		// No prompt, just refs
		refs := strings.Fields(rest)
		for _, ref := range refs {
			if refID, err := strconv.Atoi(ref); err == nil {
				node.RefIDs = append(node.RefIDs, refID)
			}
		}
	}

	return node, true
}

// parseRouting parses a routing line: [errorDest !] [dest <] src input text
func parseRouting(line string) (*RouteDefinition, bool) {
	route := &RouteDefinition{}

	// Check for error destination
	if idx := strings.Index(line, "!"); idx != -1 {
		errPart := strings.TrimSpace(line[:idx])
		if errID, err := strconv.Atoi(errPart); err == nil {
			route.ErrorDest = errID
		}
		line = strings.TrimSpace(line[idx+1:])
	}

	// Check for output destination
	if idx := strings.Index(line, "<"); idx != -1 {
		destPart := strings.TrimSpace(line[:idx])
		if destID, err := strconv.Atoi(destPart); err == nil {
			route.Dest = destID
		}
		line = strings.TrimSpace(line[idx+1:])
	}

	// If there's no source or input, this is a default handler
	if line == "" {
		return route, true
	}

	// The remaining parts are source and input
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 1 {
		return nil, false
	}

	// Parse source ID
	srcID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, false
	}
	route.Source = srcID

	// Get input text if any
	if len(parts) > 1 {
		route.Input = strings.TrimSpace(parts[1])
	}

	return route, true
} 
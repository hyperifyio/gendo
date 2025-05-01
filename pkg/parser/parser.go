package parser

import (
	"strconv"
	"strings"
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

// NodeDefinition represents a parsed node from a script line
type NodeDefinition struct {
	ID     int
	RefIDs []int  // Reference IDs this node can call
	Prompt string // Optional prompt text
	IsTool bool   // Whether this is a tool node
	Tool   string // Tool name if IsTool is true
	Type   NodeType
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

// parseNodeDefinition parses a node definition line
func parseNodeDefinition(line string) (*NodeDefinition, bool) {
	// Remove any leading/trailing whitespace
	line = strings.TrimSpace(line)

	// Split by colon
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return nil, false
	}

	// Parse node ID
	idStr := strings.TrimSpace(parts[0])
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, false
	}

	// Parse the rest of the definition
	rest := strings.TrimSpace(parts[1])
	var refIDs []int
	var prompt string
	var tool string
	var nodeType NodeType

	// Check for special node types first
	if rest == "in" {
		nodeType = NodeTypeIn
	} else if rest == "out" {
		nodeType = NodeTypeOut
	} else if rest == "err" {
		nodeType = NodeTypeErr
	} else {
		// Parse tool or prompt
		if strings.HasPrefix(rest, "tool") {
			nodeType = NodeTypeTool
			tool = strings.TrimSpace(strings.TrimPrefix(rest, "tool"))
		} else {
			// Check for references
			refParts := strings.Split(rest, " ")
			for _, ref := range refParts {
				if refID, err := strconv.Atoi(ref); err == nil {
					refIDs = append(refIDs, refID)
				} else {
					prompt = rest
					break
				}
			}
		}
	}

	return &NodeDefinition{
		ID:     id,
		RefIDs: refIDs,
		Prompt: prompt,
		Tool:   tool,
		Type:   nodeType,
	}, true
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

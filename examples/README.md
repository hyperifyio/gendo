# Gendo Script Examples

This directory contains example Gendo scripts demonstrating various features and use cases.

## Configuration

Gendo can be configured using environment variables:

- `OPENAI_API_BASE`: Base URL for the OpenAI API (default: https://api.openai.com/v1)
  - For local testing: `OPENAI_API_BASE=http://localhost:18080/v1`
- `OPENAI_API_KEY`: API key for authentication (optional)

## Script Format

Gendo scripts consist of:
1. Node definitions (numbered)
2. Input lines with routing
3. Comments (lines starting with #)

Node types:
- Empty node (`1:`) - Passthrough
- AI node (`1: prompt text`) - Uses AI to process input
- Tool node (`1: tool toolname`) - Uses a specific tool

Routing syntax:
- `< N` - Route output to node N
- `! N` - Route errors to node N

## Examples

### hello.gendo
A simple example showing basic node definitions and AI integration:
- Demonstrates passthrough nodes
- Shows AI-powered text generation
- Basic routing

### calculator.gendo
Demonstrates the math tool with error handling:
- Math operations
- Error routing
- Result formatting
- Multiple node pipeline

### file_processor.gendo
Shows file I/O operations:
- Reading files
- Writing files
- Content processing
- Error handling
- Multi-step processing

### random_story.gendo
Complex example combining multiple features:
- Random number generation
- Multi-stage processing
- AI-powered content generation
- State management between nodes

## Running Examples

To run any example:
```bash
# Using default OpenAI API
./gendo examples/hello.gendo

# Using local server
OPENAI_API_BASE=http://localhost:18080/v1 ./gendo examples/hello.gendo
```

## Error Handling

Gendo will:
- Print error messages to stderr
- Set non-zero exit code on errors
- Route errors to specified error handlers using `! N` syntax
- Continue processing if possible, stop on fatal errors 
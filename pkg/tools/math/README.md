# Math Tool

The Math Tool allows you to perform basic arithmetic operations in the gendo language.

## Usage

First, define the tool nodes:
```gendo
# Node 1 echoes everything back
1 : tool : echo

# Node 50 runs the host math evaluator
50 : tool : math
```

Then invoke them:
```gendo
1 < 50 3 * (2 + 5)
```

## Examples

Basic arithmetic operations:
```gendo
# Define the nodes
1 : tool : echo
50 : tool : math

# Evaluate expressions
1 < 50 5 + 3
1 < 50 10 - 4
1 < 50 6 * 7
1 < 50 20 / 4
```

## Supported Operations

- Addition: `+`
- Subtraction: `-`
- Multiplication: `*`
- Division: `/`

## Notes

- The tool supports decimal numbers
- Division by zero will return an error
- The tool can handle expressions with or without spaces
- Results are returned as floating-point numbers
- The tool automatically handles negative numbers
- Only one operation can be performed at a time
- Tool nodes are sandboxed and only execute their designated operation 
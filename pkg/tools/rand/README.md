# Random Number Generator Tool

The Random Number Generator Tool allows you to generate random numbers in the gendo language.

## Usage

First, define the tool nodes:
```gendo
# Node 1 echoes everything back
1 : tool : echo

# Node 51 runs the host RNG
51 : tool : rand
```

Then invoke them:
```gendo
1 < 51 1 100
```

## Examples

1. Generate a random number between 0 and 10:
```gendo
# Define the nodes
1 : tool : echo
50 : tool : rand

# Generate number
1 < 50 0 10
```

2. Generate a random number between 0 and 100:
```gendo
# Define the nodes
1 : tool : echo
51 : tool : rand

# Generate number
1 < 51 0 100
```

3. Generate a random number between 0 and 1000:
```gendo
# Define the nodes
1 : tool : echo
52 : tool : rand

# Generate number
1 < 52 0 1000
```

## Notes

- The tool generates a random 64-bit integer
- The maximum number must be positive
- The generated number will be between the specified bounds (inclusive)
- The tool uses a cryptographically secure random number generator
- The seed is automatically initialized using the current timestamp
- Tool nodes are sandboxed and only execute their designated operation 
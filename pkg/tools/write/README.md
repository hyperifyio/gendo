# Write Tool

The Write Tool allows you to write content to files in the gendo language.

## Usage

First, define the tool nodes:
```gendo
# Node 10 echoes everything back (without LLM)
10 : tool : echo

# Node 61 writes to "results.txt"
61 : tool : write : results.txt
```

Then invoke them:
```gendo
61 < 10 Some content to write
```

## Examples

1. Write a simple text file:
```gendo
# Define the nodes
10 : tool : echo
60 : tool : write : greeting.txt

# Write content
60 < 10 Welcome to gendo!
```

2. Write JSON content:
```gendo
# Define the nodes
10 : tool : echo
61 : tool : write : config.json

# Write content
61 < 10 {"name": "gendo", "version": "1.0.0"}
```

3. Write to a nested directory:
```gendo
# Define the nodes
10 : tool : echo
62 : tool : write : data/logs/app.log

# Write content
62 < 10 Application started at 2024-03-20
```

## Notes

- The tool will automatically create any necessary directories in the path
- File permissions follow the user's umask settings
- The tool returns a success message upon completion
- If the write operation fails, an error message will be returned
- Tool nodes are sandboxed and only execute their designated operation
- Filenames are sandboxed and isolated per program; no arbitrary paths allowed 

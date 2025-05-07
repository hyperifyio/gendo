# Read Tool

The Read Tool allows you to read the contents of files in the gendo language.

## Usage

First, define the tool node:
```gendo
# Node 60 reads "config.json"
60 : tool : read : config.json
```

Then invoke it:
```gendo
1 < 60
```

## Examples

1. Read a text file:
```gendo
# Define the nodes
1 : out
61 : tool : read : greeting.txt

# Read content
1 < 61
```

2. Read a configuration file:
```gendo
# Define the nodes
1 : out
62 : tool : read : config.json

# Read content
1 < 62
```

3. Read from a nested directory:
```gendo
# Define the nodes
1 : out
63 : tool : read : data/logs/app.log

# Read content
1 < 63
```

## Notes

- The tool returns the file contents as a string
- If the file doesn't exist or can't be read, an error message will be returned
- The tool supports reading any text-based file format
- Tool nodes are sandboxed and only execute their designated operation
- Filenames are sandboxed and isolated per program; no arbitrary paths allowed 
# Gendo Language Specification v0.2

## 1. Introduction

Gendo is a minimalist, prompt-based programming language designed for live, incremental code generation and execution via small, local AI models. Programs consist of self-contained **nodes**—each defining behavior or invoking AI prompts—that pass plain-text streams to each other, enabling rapid composition of functionality without mutable global state or hidden dependencies.

## 2. Core Concepts

### 2.1 Nodes

- **Definition**: `nodeID : refID refID … [: prompt text]`
  - `nodeID`: unique integer identifier
  - `refID`s: list of nodeIDs this node can call
  - `prompt text` *(optional)*: instructions for the AI when this node is invoked
- **Invocation**: `[errorDest !] [dest <] src input text`
  - Routes stdout (`dest`) and stderr (`errorDest`) to designated nodes
  - Defaults: stdout→1, stderr→2

### 2.2 Streams

- Nodes exchange plain-text. AI-enabled nodes transform input via their prompt; passthrough nodes output input verbatim.
- Errors are first-class data, buffered and routed like stdout.

### 2.3 Default Handlers

You can set **default** destinations for stdout and stderr across subsequent invocations by writing a line with only the handler syntax. Whitespace may be used to indent purely for readability; it has no semantic effect.

```gendo
# Only redefine stdout default to node 3 (errors still go to node 2)
3 <

# Errors still go to previously set default (node 2)
5 Another input

# You can override the default by specifying both handlers on a command:
# Here, errors→5, stdout→6 for this line only
5 ! 6 < Overridden command text
```

You can individually redefine defaults:

```gendo
# Only redefine stdout default to node 3 (errors still go to node 2)
3 <

  # Errors still go to previously set default (node 2)
  5 Another input

  # You can override the default by specifying it
  5 < 6 Second command text
```

The default handlers remain in effect until redefined or the script ends.

## 3. Structured Control Flow

*(Looping and conditionals TBD—let's agree on design here before fleshing out.)*

## 4. Modular Units & Files

*(Modular units, namespaces, and imports TBD—let's agree on design before fleshing this out.)*

## 5. Built-in Utilities

> **Note:** Each tool-backed node requires enabling the corresponding tool in the Gendo runtime configuration. If a tool (e.g., `math`, `rand`, `read`, `write`) is not enabled, attempting to invoke its node will result in an error.


### 5.1 Math

Gendo uses explicit **tool nodes** for arithmetic. If a node’s ref list includes the special `tool` directive, the runtime connects it to the math evaluator.

- **Definition Syntax**: `nodeID : tool : math [config...]`
  - `tool` marks a tool-backed node.
  - Optional `config` may specify precision or mode (e.g., `float`).

**Example Definition**
```gendo
# Node 50 runs the host math evaluator
50 : tool : math
```

**Example Invocation**
```gendo
# Evaluate an expression
< 50 3 * (2 + 5)
# → 21
```

Tool nodes are sandboxed and only execute their designated operation.

### 5.2 Random

Gendo defines **tool nodes** for randomness. Including `tool` with `rand` uses the host RNG.

- **Definition Syntax**: `nodeID : tool : rand [config...]`
  - `config` may specify distribution (`uniform`, `normal`) or bounds.

**Example Definition**
```gendo
# Node 51 runs the host RNG
51 : tool : rand
```

**Example Invocation**
```gendo
# Generate a random integer in [1,100]
< 51 1 100
# → 73 (example)
```

Tool nodes are sandboxed and only execute their designated operation.

### 5.3 I/O & Persistence I/O & Persistence

Gendo also uses **tool nodes** for safe, sandboxed file operations. Include `tool` in the ref list and specify `read` or `write` as the tool name.

- **Definition Syntax**: `nodeID : tool : read|write [filename]`
  - `read` nodes take no input arguments and output the contents of the named file.
  - `write` nodes accept stdin and save it to the named file, returning a confirmation message.

**Example Definitions**
```gendo

# Note 10 prints in
10 : tool : echo

# Node 60 reads "config.json"
60 : tool : read : config.json

# Node 61 writes to "results.txt"
61 : tool : write : results.txt
```

**Example Invocations**
```gendo
# Load configuration and send to node 61
# → {"threshold":10}
61 < 60

# Write to node 61
61 < 10 Some computed output text
# → "Written to results.txt"
```

- Filenames are sandboxed and isolated per program; no arbitrary paths allowed.

## 6. Safety & Concurrency

Gendo emphasizes reliability and performance:

- **Stateless Nodes**: By default, nodes have no hidden state; all side effects occur through explicit tool nodes (e.g., I/O), ensuring predictable behavior.
- **Error Handling**: Errors are treated as first-class data. You choose where to route error messages via the `errorDest !` syntax; unhandled errors by default go to node 2. This allows logging, retries, or feeding errors into AI prompts for recovery.
- **Concurrency and Parallelism**: The runtime can execute independent node invocations in parallel when there are no data dependencies. This lets you leverage multi-core CPUs without adding complex syntax.
- **Sandboxing**: Tool nodes (math, rand, read, write) are isolated from arbitrary host resources. Filesystem and network access occur only through sandboxed APIs, preventing unauthorized operations.

## 7. Data Model

Gendo operates purely on plain text streams. Each node receives a string and returns a string. For structured data (e.g., JSON), simply define your prompts or AI nodes to parse and emit valid JSON. Gendo does not enforce data schemas, offering maximum flexibility.

## 8. Community and Next Steps

Gendo invites developers to build small, focused units that grow at runtime via AI. Its minimal core encourages experimentation:

- **Extensibility**: Community-contributed tools and node libraries can add capabilities (e.g., HTTP, database connectors) without altering the core.
- **Safety**: All extensions must register as explicit tools and respect sandbox rules.
- **Example Library**: Curated sets of nodes for common tasks (e.g., data processing pipelines, chat bots).

*Gendo makes it so.*

# Tool Specification - list-files-tool.go

This document provides a map of the code structure in `list-files-tool.go` to facilitate future edits and understanding.

## Code Blocks Map

| Line Range | Function / Block | Description |
|------------|------------------|-------------|
| 1-8 | Package & Imports | Defines the package and imports necessary libraries (bytes, encoding/json, fmt, os). |
| 10-16 | `listFilesTool` | Defines the `Tool` object for listing files, including its definition, guidelines, and execution logic. |
| 17-30 | `listFilesToolDef` | Defines the `APITool` structure, including the function name, description, and required parameters (path). |
| 32-61 | `runListFiles` | The core execution logic. It unmarshals the arguments, reads the directory contents, and formats the output as a list of [DIR] or [FILE] entries with their sizes. |

## Key Logic Details

### Argument Parsing (Lines 32-36)
The function accepts a JSON string as `arguments`, unmarshals it into a `ListFilesArgs` struct, and returns an error message if the JSON is invalid.

### Directory Reading (Lines 38-41)
It uses `os.ReadDir` to retrieve the entries of the specified path. If an error occurs (e.g., directory doesn't exist), it returns a descriptive error message.

### Result Formatting (Lines 47-60)
Iterates through the directory entries:
- **Directories**: Prefixed with `[DIR]`.
- **Files**: Prefixed with `[FILE]`, followed by the file name and its size in bytes.
- **Empty Directories**: Returns a specific "Directory is empty" message if no entries are found.

# Tool Specification - read-file-tool.go

This document provides a map of the code structure in `read-file-tool.go` to facilitate future edits and understanding.

## Code Blocks Map

| Line Range | Function / Block | Description |
|------------|------------------|-------------|
| 1-9 | Package & Imports | Defines the package and imports necessary libraries (bufio, bytes, encoding/json, fmt, os). |
| 11-21 | `readToolGuideline` | A multi-line string defining the usage guidelines and instructions for the `read_file` tool. |
| 23-28 | `readFilesTool` | Defines the `Tool` object for reading files, including its definition, guidelines, and execution logic. |
| 30-54 | `readFileToolDef` | Defines the `APITool` structure, including the function name, description, and parameters (path, start_line, end_line). |
| 56-102 | `runReadFile` | The core execution logic. It unmarshals arguments, opens the file, handles line range selection, and scans the file line by line. |

## Key Logic Details

### Argument Parsing (Lines 56-60)
The function accepts a JSON string as `arguments`, unmarshals it into a `ReadFileArgs` struct, and returns an error message if the JSON is invalid.

### File Handling (Lines 62-77)
- **Opening**: Opens the file using `os.Open` and ensures it's closed with `defer`.
- **Range Selection**: Determines the `start` and `end` line numbers based on optional `start_line` and `end_line` parameters.

### Scanning Logic (Lines 78-91)
Uses a `bufio.Scanner` to iterate through the file:
- **Line Numbering**: Keeps track of the current line number.
- **Filtering**: Only writes to the buffer if the current line is within the requested `start` and `end` range.
- **Termination**: Breaks the loop early if the `end` line is reached.

### Output Generation (Lines 93-102)
- **Error Handling**: Checks for scanner errors during the process.
- **Empty File Check**: If no content is collected, returns a specific message indicating the file was empty.
- **Formatting**: Prepends the line number to each line of the content (e.g., `1 | content`).

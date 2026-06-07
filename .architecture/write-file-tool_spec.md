# Tool Specification - write-file-tool.go

This document provides a map of the code structure in `write-file-tool.go` to facilitate future edits and understanding.

## Code Blocks Map

| Line Range | Function / Block | Description |
|------------|------------------|-------------|
| 1-8 | Package & Imports | Defines the package and imports necessary libraries (encoding/json, fmt, os, path/filepath). |
| 10-16 | `writeFileTool` | Defines the `Tool` object for writing files, including its definition, guidelines, and execution logic. |
| 18-25 | `writeToolGuideline` | A multi-line string defining the usage guidelines and instructions for the `write_file` tool. |
| 26-46 | `writeFileToolDef` | Defines the `APITool` structure, including the function name, description, and parameters (path, content). |
| 48-81 | `runWriteFile` | The core execution logic. It handles backup of existing content, directory creation, and writing the new content to the file. |

## Key Logic Details

### Argument Parsing (Lines 48-52)
The function accepts a JSON string as `arguments`, unmarshals it into a `WriteFileArgs` struct, and returns an error message if the JSON is invalid.

### Backup Mechanism (Lines 54-66)
- **Existing File**: If the file already exists, the tool reads its current content and stores it in `lastBackupContent`. `wasNewlyCreated` is set to `false`.
- **New File**: If the file does not exist, `lastBackupContent` is set to `nil` and `wasNewlyCreated` is set to `true`.

### Directory Creation (Lines 71-74)
Uses `filepath.Dir` to identify the parent directory and `os.MkdirAll` to create any missing parent directories with `0755` permissions.

### File Writing (Lines 76-81)
Writes the provided content to the file using `os.WriteFile` with `0644` permissions. Returns a success message including the number of bytes written.

# Tool Specification - revert-file-tool.go

This document provides a map of the code structure in `revert-file-tool.go` to facilitate future edits and understanding.

## Code Blocks Map

| Line Range | Function / Block | Description |
|------------|------------------|-------------|
| 1-6 | Package & Imports | Defines the package and imports necessary libraries (fmt, os). |
| 8-13 | Global State | Manages variables for tracking the last file operation: `hasBackup`, `lastBackupPath`, `lastBackupContent`, and `wasNewlyCreated`. |
| 15-21 | `revertFileTool` | Defines the `Tool` object for reverting a file, including its definition, guidelines, and execution logic. |
| 23-33 | `revertFileToolDef` | Defines the `APITool` structure, including the function name, description, and parameters (none). |
| 35-53 | `runRevertFile` | The core execution logic. It checks for a previous backup and either deletes a newly created file or restores the previous content. |

## Key Logic Details

### State Management (Lines 8-13)
The tool uses global variables to maintain the state of the last file operation. This state is shared across all tool executions.

### Revert Logic (Lines 35-53)
The function `runRevertFile` handles two scenarios:
- **Newly Created File**: If the file was newly created during the last operation, it deletes the file from the filesystem.
- **Modified File**: If the file already existed, it restores the previous content from the `lastBackupContent` buffer.

### State Reset (Lines 44-52)
After a successful revert, the `hasBackup` flag is set to `false` to ensure that subsequent reverts do not attempt to use stale data.

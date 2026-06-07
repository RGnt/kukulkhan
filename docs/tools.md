# Tools Documentation

This document provides specifications and details for the tools available to the LLM.

## Tools Registry

The following tools are registered and available for execution:

### 1. `list_files`
- **Description**: List all files and directories in the specified path.
- **Parameters**:
    - `path` (string): The absolute or relative path to list.
- **Returns**: A formatted list of files and directories with their sizes.

### 2. `read_file`
- **Description**: Read the contents of a file. You can optionally specify a start_line and end_line to read a specific chunk. Lines are 1-indexed.
- **Parameters**:
    - `path` (string): The absolute or relative path to the file.
    - `start_line` (integer, optional): The line number to start reading from (1-indexed).
    - `end_line` (integer, optional): The line number to stop reading at (inclusive).
- **Returns**: The content of the file with line numbers.

### 3. `write_file`
- **Description**: Write entire content to a file. Overwrites existing files and creates missing directories. Automatically backs up the previous state.
- **Parameters**:
    - `path` (string): The absolute or relative path to the file.
    - `content` (string): The complete new content of the file.
- **Returns**: A success message with the number of bytes written.

### 4. `revert_file`
- **Description**: Reverts the last `write_file` operation. Use this immediately if you realize your last file write was incorrect.
- **Parameters**: None.
- **Returns**: A success message indicating the revert action.

## Implementation Details

- **Routing**: The `executeTool` function acts as the central router, mapping tool names to their respective implementation functions.
- **Backup Mechanism**: The `write_file` tool maintains a global state (`hasBackup`, `lastBackupPath`, `lastBackupContent`, `wasNewlyCreated`) to allow for the `revert_file` operation.
- **File Handling**:
    - `list_files` uses `os.ReadDir` and provides directory/file distinction.
    - `read_file` uses `bufio.Scanner` for efficient line-by-line reading.
    - `write_file` ensures parent directories exist using `os.MkdirAll` before writing.

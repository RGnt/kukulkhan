# Tools Specification - tools.go

This document provides a map of the code blocks in `tools.go` to assist with future edits and tool integration.

## Overview
The `tools.go` file contains the definitions and implementations for the core tools provided to the LLM, including file system operations and a backup/revert mechanism.

## Code Map

### Global Variables & Tool Definitions
- **Lines 12-17**: Global state variables for the file backup system (`hasBackup`, `lastBackupPath`, `lastBackupContent`, `wasNewlyCreated`).
- **Lines 19-32**: Definition of the `list_files` tool (APITool).
- **Lines 34-58**: Definition of the `read_file` tool (APITool).
- **Lines 60-80**: Definition of the `write_file` tool (APITool).
- **Lines 82-93**: Definition of the `revert_file` tool (APITool).
- **Lines 94-184**: Commented out registry of tools and a commented out `executeTool` router.

### Tool Implementations

#### `runListFiles`
- **Lines 185-214**: Implementation of the `list_files` tool.
  - **Description**: Reads the directory at the specified path and returns a formatted string listing files (`[FILE]`) and directories (`[DIR]`) with their sizes.

#### `runReadFile`
- **Lines 216-262**: Implementation of the `read_file` tool.
  - **Description**: Reads the content of a file. Supports optional `start_line` and `end_line` parameters to return a specific range of lines with line numbers.

#### `runWriteFile`
- **Lines 264-297**: Implementation of the `write_file` tool.
  - **Description**: Writes content to a file, creating directories if necessary. It automatically backs up the existing content (if any) to a global variable to allow for reversion.

#### `runRevertFile`
- **Lines 299-317**: Implementation of the `revert_file` tool.
  - **Description**: Reverts the last `write_file` operation. If the file was newly created, it deletes it. If it was an overwrite, it restores the previously backed-up content.

# Kukulkhan agent harness

A Go-based application that enables an LLM to interact with local tools, such as file system operations and mathematical calculations, through a structured tool-calling interface.

## Project Structure

- `cmd/`: Contains the application logic and tool implementations.
  - `main.go`: The entry point of the application.
  - `agent.go`: Contains the core logic for the agent's behavior and interaction.
  - `types.go`: Contains shared data structures and types used across the project.
  - `list-files-tool.go`: Implementation of the `list_files` tool.
  - `read-file-tool.go`: Implementation of the `read_file` tool.
  - `write-file-tool.go`: Implementation of the `write_file` tool.
  - `revert-file-tool.go`: Implementation of the `revert_file` tool.
  - `todo-tool.go`: Implementation of the `todo` tool.
- `docs/`: Contains documentation for the tools.

## Getting Started

### Prerequisites

- Go (version 1.18 or higher recommended)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-name>
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Application

To run the application directly:
```bash
go run cmd/main.go
```

To build the binary:
```bash
go build -o agent ./cmd
./agent
```

## Documentation

For detailed specifications of the available tools, including parameters and return types, please refer to the documentation:

[Tool Documentation](./docs/tools.md)

## License

This project is licensed under the [MIT License](LICENSE).

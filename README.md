1 | # Kukulkhan agent harness
2 | 
3 | A Go-based application that enables an LLM to interact with local tools, such as file system operations and mathematical calculations, through a structured tool-calling interface.
4 | 
5 | ## Project Structure
6 | 
7 | - `cmd/`: Contains the application logic and tool implementations.
8 |   - `main.go`: The entry point of the application.
9 |   - `agent.go`: Contains the core logic for the agent's behavior and interaction.
10 |   - `types.go`: Contains shared data structures and types used across the project.
11 |   - `list-files-tool.go`: Implementation of the `list_files` tool.
12 |   - `read-file-tool.go`: Implementation of the `read_file` tool.
13 |   - `write-file-tool.go`: Implementation of the `write_file` tool.
14 |   - `revert-file-tool.go`: Implementation of the `revert_file` tool.
15 |   - `todo-tool.go`: Implementation of the `todo` tool.
16 | - `docs/`: Contains documentation for the tools.
17 | 
18 | ## Getting Started
19 | 
20 | ### Prerequisites
21 | 
22 | - Go (version 1.18 or higher recommended)
23 | 
24 | ### Installation
25 | 
26 | 1. Clone the repository:
27 |    ```bash
28 |    git clone <repository-url>
29 |    cd <repository-name>
30 |    ```
31 | 
32 | 2. Install dependencies:
33 |    ```bash
34 |    go mod tidy
35 |    ```
36 | 
37 | ### Running the Application
38 | 
39 | To run the application directly:
40 | ```bash
41 | go run cmd/main.go
42 | ```
43 | 
44 | To build the binary:
45 | ```bash
46 | go build -o agent ./cmd
47 | ./agent
48 | ```
49 | 
50 | ## Documentation
51 | 
52 | For detailed specifications of the available tools, including parameters and return types, please refer to the documentation:
53 | 
54 | [Tool Documentation](./docs/tools.md)
55 | 
56 | ## License
57 | 
58 | This project is licensed under the [MIT License](LICENSE).

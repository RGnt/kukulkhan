# Agent Specification - agent.go

This document provides a map of the code structure in `agent.go` to facilitate future edits and understanding.

## Code Blocks Map

| Line Range | Function / Block | Description |
|------------|------------------|-------------|
| 1-12 | Package & Imports | Defines the package and imports necessary libraries (bufio, bytes, json, fmt, log, net/http, strings, time). |
| 14-36 | `NewAgent` | Constructor for the `Agent` struct. It initializes the agent's role, model, temperature, and a map of available tools. It also builds the system prompt by including tool guidelines. |
| 38-47 | `ExecuteTask` | A high-level method to execute a single task. It wraps the prompt into a history with a system message and returns the content of the final message in the execution chain. |
| 49-123 | `Run` | The core ReAct (Reasoning and Acting) loop. It manages the interaction with the LLM, handles tool calls by executing them and feeding results back into the history, and repeats until a final answer is reached or max steps are hit. |
| 125-212 | `streamAndBuildMessage` | A helper function that processes the streaming response from the LLM. It parses the `text/event-stream` format, reconstructs tool call fragments into complete tool calls, and builds the final `Message` object. |

## Key Logic Details

### Tool Execution (Lines 87-111)
Within the `Run` loop, if the LLM returns `ToolCalls`, the agent:
1. Identifies the tool by name.
2. Executes the tool's `Execute` method with the provided arguments.
3. Appends the tool's output back into the message history as a `tool` role message.

### Streaming Parser (Lines 125-212)
The `streamAndBuildMessage` function handles the complexity of streaming:
- **Content**: Concatenates text fragments as they arrive.
- **Tool Calls**: Uses a map of `toolCallBuilder` structs to aggregate fragments of tool calls (ID, Name, Arguments) because they may arrive in separate chunks.
- **Assembly**: Once the stream ends, it converts the builders into the final `ToolCall` slice.

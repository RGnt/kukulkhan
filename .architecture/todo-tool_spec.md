# Tool Specification - todo-tool.go

This document provides a map of the code structure in `todo-tool.go` to facilitate future edits and understanding.

## Code Blocks Map

| Line Range | Function / Block | Description |
|------------|------------------|-------------|
| 3-181 | `todoToolDescription` | A large multi-line string containing comprehensive guidelines, examples, and rules for using the todo list tool. |

## Key Logic Details

### Guidelines and Examples (Lines 3-181)
This file primarily serves as a prompt/guideline for the LLM agent:
- **When to Use**: Outlines scenarios for complex, multi-step, or non-trivial tasks where proactive tracking is beneficial.
- **When NOT to Use**: Provides clear boundaries for trivial, informational, or single-step tasks.
- **Examples**: Includes several structured examples showing the correct way to use the todo list (and when to skip it) along with the internal reasoning for those decisions.
- **Task States**: Defines the lifecycle of a task (`pending`, `in_progress`, `completed`) and strict rules for managing these states (e.g., only one task in progress at a time).
- **Task Breakdown**: Instructions on creating actionable, descriptive items with both imperative and present continuous forms.

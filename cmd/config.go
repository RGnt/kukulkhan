package main

import (
	"os"
	"strconv"
)

// AppConfig holds the global configuration for the application
type AppConfig struct {
	WorkspaceDir   string
	StreamResponse bool
}

// Config is the globally accessible configuration instance
var Config AppConfig

// LoadConfig initializes the configuration with defaults and environment variables
func LoadConfig() {
	// 1. Set Defaults
	Config = AppConfig{
		WorkspaceDir:   ".",   // Default to current directory
		StreamResponse: false, // Default to the "ChatGPT-like" experience
	}

	// 2. Override Workspace if AGENT_WORKSPACE is set
	if envWorkspace := os.Getenv("AGENT_WORKSPACE"); envWorkspace != "" {
		Config.WorkspaceDir = envWorkspace
	} else {
		// Try to resolve the absolute path of the current directory
		if cwd, err := os.Getwd(); err == nil {
			Config.WorkspaceDir = cwd
		}
	}

	// 3. Override Streaming if AGENT_STREAM is set
	if envStream := os.Getenv("AGENT_STREAM"); envStream != "" {
		if parsedBool, err := strconv.ParseBool(envStream); err == nil {
			Config.StreamResponse = parsedBool
		}
	}
}

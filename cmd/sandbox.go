package main

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Sandbox struct {
	cli         *client.Client
	containerID string
	workspace   string
}

// NewSandbox starts a persistent, long-running Docker container
func NewSandbox(ctx context.Context, hostWorkspace string) (*Sandbox, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker: %v", err)
	}

	// Ensure the workspace is an absolute path for Docker bind mounts
	absWorkspace, err := filepath.Abs(hostWorkspace)
	if err != nil {
		return nil, err
	}

	// 1. Define the container. We use 'sleep infinity' to keep it running forever
	// so the Agent can exec into it repeatedly without startup latency.
	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:      "agent-sandbox", // Or "ubuntu", "golang", "python", etc.
			Cmd:        []string{"sleep", "infinity"},
			WorkingDir: "/workspace",
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: absWorkspace,
					Target: "/workspace",
				},
			},
		}, nil, nil, "")

	if err != nil {
		return nil, fmt.Errorf("failed to create sandbox container: %v", err)
	}

	// 2. Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start sandbox: %v", err)
	}

	return &Sandbox{
		cli:         cli,
		containerID: resp.ID,
		workspace:   absWorkspace,
	}, nil
}

// Execute runs a shell command inside the running sandbox and returns stdout/stderr
func (s *Sandbox) Execute(ctx context.Context, command string) (string, error) {
	// 1. Create the exec instance
	execResp, err := s.cli.ContainerExecCreate(ctx, s.containerID, container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"sh", "-c", command}, // Use "bash" if using Ubuntu
	})
	if err != nil {
		return "", fmt.Errorf("failed to create exec: %v", err)
	}

	// 2. Attach to the exec instance to read the output streams
	attachResp, err := s.cli.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec: %v", err)
	}
	defer attachResp.Close()

	var stdout, stderr bytes.Buffer

	// 3. Docker multiplexes stdout/stderr over a single TCP connection.
	// StdCopy cleanly splits them back apart into our buffers.
	_, err = stdcopy.StdCopy(&stdout, &stderr, attachResp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %v", err)
	}

	// 4. Combine output to feed back to the LLM
	result := stdout.String()
	if stderr.Len() > 0 {
		result += fmt.Sprintf("\n[Stderr]:\n%s", stderr.String())
	}

	return result, nil
}

// Close gracefully kills and removes the container when the agent shuts down
func (s *Sandbox) Close(ctx context.Context) error {
	fmt.Println("\n[Sandbox] Tearing down secure environment...")
	// Force remove bypasses needing to stop it first
	return s.cli.ContainerRemove(ctx, s.containerID, container.RemoveOptions{Force: true})
}

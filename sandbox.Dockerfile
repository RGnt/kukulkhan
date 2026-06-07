# Start with the official Go image
FROM golang:1.22-alpine

# Install the Agent Toolbelt
RUN apk update && apk add --no-cache \
    # GNU utility overrides (prevents Alpine/Busybox flag errors)
    coreutils \
    findutils \
    grep \
    bash \
    # Vision and Search
    tree \
    ripgrep \
    # Data Parsing
    jq \
    # Network and Version Control
    curl \
    wget \
    git \
    # Compilation basics (gcc, make, libc-dev)
    build-base 

# (Optional) Install golangci-lint so the agent can self-check its code
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.1

# Ensure bash is the default shell for exec commands
ENV SHELL=/bin/bash

WORKDIR /workspace

# Keep container alive
CMD ["sleep", "infinity"]
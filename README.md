# Tududi MCP (HTTP)

[![CI](https://github.com/nouchka/mcp-tududi/actions/workflows/ci.yml/badge.svg)](https://github.com/nouchka/mcp-tududi/actions/workflows/ci.yml)
[![Docker](https://github.com/nouchka/mcp-tududi/actions/workflows/docker.yml/badge.svg)](https://github.com/nouchka/mcp-tududi/actions/workflows/docker.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue)](https://golang.org/)

A Model Context Protocol (MCP) HTTP server written in Go that integrates [Tududi](https://tududi.com/) task management with AI-powered development tools.

## Overview

Tududi MCP HTTP enables AI agents and developers to interact with Tududi tasks, projects, and areas directly through HTTP endpoints. Manage your tasks seamlessly while coding, powered by the Model Context Protocol standard.

### Features

- **Task Management**: Create, read, update, and delete tasks
- **Subtask Management**: Full support for hierarchical subtasks
- **Project Organization**: Manage projects and organize work
- **Area Management**: Organize tasks by areas
- **HTTP Interface**: RESTful API endpoints for all operations
- **AI-Ready**: Works with AI agents like GitHub Copilot
- **Docker Support**: Pre-built Docker images available
- **Go Performance**: Fast and efficient implementation in Go

## Quick Start

### Option 1: Docker (Recommended)

```bash
# Pull the latest image
docker pull ghcr.io/nouchka/mcp-tududi:latest

# Run with environment variables
docker run -d \
  -e TUDUDI_API_URL=http://your-tududi-instance:3000 \
  -e TUDUDI_API_KEY=your-api-key \
  -e HTTP_PORT=8080 \
  -p 8080:8080 \
  ghcr.io/nouchka/mcp-tududi:latest
```

### Option 2: Docker Compose

```bash
git clone https://github.com/nouchka/mcp-tududi.git
cd mcp-tududi

# Copy and edit the environment file
cp .env.example .env
# Edit .env with your Tududi API credentials

# Run with Docker Compose
docker-compose up -d
```

#### Testing with CLI

A CLI test service is included to verify the MCP server is working correctly:

```bash
# Run the CLI test service (requires the mcp-tududi service to be running)
docker-compose --profile test run mcp-cli

# This will test:
# - Health endpoint
# - Tools listing endpoint
# - Tool calling endpoint (example: list_tasks)
```

The CLI test service can also be run with a custom MCP server URL:

```bash
docker-compose --profile test run -e MCP_SERVER_URL=http://custom-host:8080 mcp-cli
```

### Option 3: Local Installation

```bash
# Prerequisites: Go 1.24 or later

git clone https://github.com/nouchka/mcp-tududi.git
cd mcp-tududi

# Build the project
go build -o mcp-tududi ./cmd/server

# Run the server
TUDUDI_API_URL=http://localhost:3000 \
TUDUDI_API_KEY=your-api-key \
./mcp-tududi
```

## Configuration

The MCP server is configured via environment variables:

### Required Variables

- `TUDUDI_API_URL`: The base URL of your Tududi instance (e.g., `http://localhost:3000`)

### Authentication (Choose One)

- `TUDUDI_API_KEY`: API key for authentication (recommended)
- `TUDUDI_EMAIL` and `TUDUDI_PASSWORD`: Email and password for authentication (legacy)

### Optional Variables

- `HTTP_PORT`: HTTP port to listen on (default: `8080`)
- `LOG_LEVEL`: Logging level - `debug`, `info`, `warn`, `error` (default: `info`)

### Example .env File

```env
TUDUDI_API_URL=http://localhost:3000
TUDUDI_API_KEY=your-api-key-here
HTTP_PORT=8080
LOG_LEVEL=info
```

## API Endpoints

### Health Check

```
GET /health
```

Returns server health status.

### List Available Tools

```
GET /tools
```

Returns list of available MCP tools.

### Call Tool

```
POST /call_tool
```

Executes an MCP tool with provided parameters.

Example request:
```json
{
  "name": "tududi_list_tasks",
  "input": {}
}
```

## Available Tools

### Task Management

- `tududi_list_tasks` - List all tasks
- `tududi_create_task` - Create a new task
- `tududi_update_task` - Update an existing task
- `tududi_delete_task` - Delete a task
- `tududi_complete_task` - Mark a task as complete

### Subtask Management

- `tududi_list_subtasks` - List subtasks for a parent task
- `tududi_create_subtask` - Create a subtask under a parent task
- `tududi_update_subtask` - Update a subtask
- `tududi_delete_subtask` - Delete a subtask

### Project Management

- `tududi_list_projects` - List all projects
- `tududi_create_project` - Create a new project
- `tududi_update_project` - Update a project
- `tududi_delete_project` - Delete a project

### Area Management

- `tududi_list_areas` - List all areas
- `tududi_create_area` - Create a new area
- `tududi_update_area` - Update an area
- `tududi_delete_area` - Delete an area

## Development

### Prerequisites

- Go 1.24 or later
- `git`

### Setup

```bash
# Clone the repository
git clone https://github.com/nouchka/mcp-tududi.git
cd mcp-tududi

# Download dependencies
go mod download

# Build the project
go build -v ./cmd/server

# Build the CLI test tool
go build -v ./cmd/cli

# Run tests
go test -v ./...

# Lint the code
golangci-lint run
```

### CLI Test Tool

The CLI test tool can be used to verify the MCP server is working correctly:

```bash
# Build the CLI
go build -o mcp-cli ./cmd/cli

# Run the CLI test tool
# Default MCP server URL is http://localhost:8080
./mcp-cli

# Or with a custom MCP server URL
MCP_SERVER_URL=http://example.com:8080 ./mcp-cli
```

The CLI performs the following tests:
- **Health Check**: Verifies the server is responding (with retries)
- **Tools Endpoint**: Lists all available MCP tools
- **Tool Execution**: Calls a sample tool (list_tasks) to verify functionality

### Project Structure

```
.
├── cmd/
│   ├── server/          # Application entry point
│   └── cli/             # CLI test tool
├── internal/
│   ├── client/          # Tududi API client
│   ├── config/          # Configuration management
│   └── mcp/             # MCP server implementation
├── Dockerfile           # Docker container definition (multi-stage)
├── docker-compose.yml   # Docker Compose configuration
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
└── README.md            # This file
```

## Agent Integration

The Tududi MCP server can be integrated with various AI agents to provide task management capabilities directly within your development workflow.

### GitHub Copilot Integration

To add the Tududi MCP server to GitHub Copilot, update your Copilot configuration file (typically `.github/copilot/config.json` or your Copilot settings):

```json
{
  "mcpServers": [
    {
      "name": "tududi",
      "url": "http://localhost:8080"
    }
  ]
}
```

Or via environment configuration:
```bash
# Ensure the MCP server is running
docker-compose up -d

# The server will be available at http://localhost:8080
# Configure your agent to connect to this endpoint
```

### Generic MCP Integration

Any MCP-compatible agent can connect to the server using the HTTP interface:

1. **Start the MCP server** using Docker Compose or local installation
2. **Configure your agent** to connect to the server URL (default: `http://localhost:8080`)
3. **Agent can now access** all available Tududi tools through the MCP protocol

### Example: Calling Tools from an Agent

Once connected, agents can call tools using the HTTP endpoints:

```bash
# Get available tools
curl http://localhost:8080/tools

# Create a task
curl -X POST http://localhost:8080/call_tool \
  -H "Content-Type: application/json" \
  -d '{
    "name": "tududi_create_task",
    "input": {
      "title": "Implement feature X",
      "description": "Add support for feature X"
    }
  }'

# List all tasks
curl -X POST http://localhost:8080/call_tool \
  -H "Content-Type: application/json" \
  -d '{
    "name": "tududi_list_tasks",
    "input": {}
  }'
```

### Environment Variables for Agent Integration

When running the MCP server for agent integration, ensure these variables are properly configured:

- `TUDUDI_API_URL`: URL of your Tududi instance
- `TUDUDI_API_KEY`: Authentication key (recommended)
- `HTTP_PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging verbosity (default: info)

## Docker Image

Docker images are automatically built and published to GitHub Container Registry (ghcr.io) on every push to main and for all tags.

### Available Tags

- `latest` - Latest version from main branch
- `vX.Y.Z` - Specific version (semantic versioning)
- `main-{sha}` - Specific commit on main branch

### Example Usage

```bash
# Using latest
docker run ghcr.io/nouchka/mcp-tududi:latest

# Using specific version
docker run ghcr.io/nouchka/mcp-tududi:v1.0.0

# With environment configuration
docker run -e TUDUDI_API_URL=http://host.docker.internal:3000 \
           -e TUDUDI_API_KEY=your-key \
           -p 8080:8080 \
           ghcr.io/nouchka/mcp-tududi:latest
```

## Troubleshooting

### Connection Issues

If you get connection errors when connecting to Tududi:

1. Verify the `TUDUDI_API_URL` is correct and accessible
2. Check your `TUDUDI_API_KEY` or credentials are valid
3. Ensure no firewall is blocking the connection
4. Check Tududi server is running and healthy

### Authorization Errors

If you get 401/403 errors:

1. Verify your API key is valid
2. If using email/password, ensure both are correct
3. Check that your user has appropriate permissions in Tududi

### Docker Container Issues

If the container won't start:

1. Check logs: `docker logs <container-id>`
2. Verify environment variables are set: `docker inspect <container-id>`
3. Check port is not already in use: `netstat -tln | grep 8080`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## References

- [Tududi](https://tududi.com/)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [Go](https://golang.org/)

## Based On

This HTTP implementation was inspired by [tududi-mcp](https://github.com/jerrytunin/tududi-mcp) by Jerry Tunin.

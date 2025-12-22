# Quick Start Guide

This guide will help you get the MCP server up and running in minutes.

## Prerequisites

- Go 1.25.1 or later installed
- Basic knowledge of HTTP and JSON-RPC

## Step 1: Build the Server

```bash
make build
# or
go build -o bin/mcp-server main.go
```

## Step 2: Run the Server

```bash
make run
# or
./bin/mcp-server
```

You should see:
```
Starting MCP server on port 8080
MCP endpoint: http://localhost:8080/mcp
Health check: http://localhost:8080/health
```

## Step 3: Test the Server

Run the test script:

```bash
./test_server.sh
```

Or manually test with curl:

```bash
# Test health endpoint
curl http://localhost:8080/health

# Initialize connection
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }'

# List available tools
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'

# Call the echo tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "echo",
      "arguments": {
        "message": "Hello, World!"
      }
    }
  }'
```

## Understanding the Response

All responses follow the JSON-RPC 2.0 format:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    // Response data here
  }
}
```

Or in case of error:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32600,
    "message": "Error description"
  }
}
```

## Available Tools

### 1. Echo Tool

Echoes back your message:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "echo",
      "arguments": {
        "message": "Your message here"
      }
    }
  }'
```

### 2. Add Tool

Adds two numbers:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "add",
      "arguments": {
        "a": 10,
        "b": 20
      }
    }
  }'
```

### 3. Get Time Tool

Returns current server time:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "get_time"
    }
  }'
```

## Working with Resources

List available resources:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "resources/list"
  }'
```

Read a resource:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "resources/read",
    "params": {
      "uri": "server://info"
    }
  }'
```

## Working with Prompts

List available prompts:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "prompts/list"
  }'
```

Get a prompt with arguments:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "prompts/get",
    "params": {
      "name": "greeting",
      "arguments": {
        "name": "Alice"
      }
    }
  }'
```

## Stopping the Server

Press `Ctrl+C` in the terminal where the server is running. The server will gracefully shut down.

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Explore the source code in `main.go` to understand how tools, resources, and prompts are implemented
- Add your own custom tools, resources, and prompts
- Integrate the server with your applications

## Troubleshooting

### Server won't start

- Check if port 8080 is already in use: `lsof -i :8080`
- Try a different port by modifying `serverPort` in `main.go`

### Connection refused

- Make sure the server is running
- Check firewall settings
- Verify the URL is correct: `http://localhost:8080/mcp`

### Invalid JSON errors

- Ensure your JSON is properly formatted
- Use a JSON validator before sending requests
- Check that all required fields are present

## Support

For issues or questions, please refer to:
- [MCP Go Library Documentation](https://pkg.go.dev/github.com/mark3labs/mcp-go)
- [Model Context Protocol Specification](https://modelcontextprotocol.io/)



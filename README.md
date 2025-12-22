# MCP Server with Streamable HTTP (SSE)

This is a Model Context Protocol (MCP) server implementation using the mark3labs/mcp-go library with Server-Sent Events (SSE) for streamable HTTP transport.

> 🚀 **New to MCP?** Check out the [Quick Start Guide](QUICKSTART.md) to get up and running in minutes!

## Features

- **Streamable HTTP Transport**: Uses SSE (Server-Sent Events) for real-time communication
- **Tools**: Provides example tools for echo, addition, and time retrieval
- **Resources**: Exposes server information as a resource
- **Prompts**: Includes example prompts for greetings and code reviews

## Prerequisites

- Go 1.25.1 or later
- The `mcp-go` library (already included in go.mod)

## Installation

The dependencies are already configured in `go.mod`. To ensure all dependencies are downloaded:

```bash
go mod download
```

## Running the Server

Start the server with:

```bash
go run main.go
```

The server will start on port 8080 by default. You should see:

```
Starting MCP server on port 8080
MCP endpoint: http://localhost:8080/mcp
Health check: http://localhost:8080/health
```

## Endpoints

- **MCP Endpoint**: `http://localhost:8080/mcp` - Main MCP communication endpoint (uses Server-Sent Events)
- **Health Check**: `http://localhost:8080/health` - Simple health check endpoint

## Available Tools

### 1. Echo Tool
Echoes back the input text.

**Parameters:**
- `message` (string, required): The message to echo back

### 2. Add Tool
Adds two numbers together.

**Parameters:**
- `a` (number, required): First number
- `b` (number, required): Second number

### 3. Get Time Tool
Returns the current server time.

**Parameters:** None

## Available Resources

### Server Information
- **URI**: `server://info`
- **Type**: text/plain
- **Description**: Provides information about the MCP server

## Available Prompts

### 1. Greeting Prompt
Generates a personalized greeting.

**Arguments:**
- `name` (string, required): Name of the person to greet

### 2. Code Review Prompt
Generates a code review prompt template.

**Arguments:**
- `language` (string, required): Programming language for the code review

## Testing the Server

You can test the server using curl or any HTTP client that supports SSE:

```bash
# Health check
curl http://localhost:8080/health

# Connect to MCP endpoint (requires MCP client)
curl -N -H "Accept: text/event-stream" http://localhost:8080/mcp
```

## Configuration

You can modify the following constants in `main.go`:

- `serverName`: Name of the MCP server
- `serverVersion`: Version of the server
- `serverPort`: Port to run the server on (default: 8080)

### Stateless vs Stateful Mode

The server is configured to run in **stateless mode** by default:

```go
streamableServer := server.NewStreamableHTTPServer(s, server.WithStateLess(true))
```

**Stateless Mode:**
- Each request is independent
- No session management required
- Easier to test and scale horizontally
- Recommended for most use cases

**Stateful Mode:**
- Sessions are maintained across requests
- Requires cookie-based session management
- Useful for complex workflows that need state persistence

To switch to stateful mode, change the configuration to:

```go
streamableServer := server.NewStreamableHTTPServer(s)
// or explicitly
streamableServer := server.NewStreamableHTTPServer(s, server.WithStateLess(false))
```

## Graceful Shutdown

The server supports graceful shutdown. Press `Ctrl+C` to stop the server. It will:
1. Stop accepting new connections
2. Wait for existing connections to complete (up to 5 seconds)
3. Shut down cleanly

## Architecture

The server uses:
- **Streamable HTTP Server**: Provided by `server.NewStreamableHTTPServer()` for streamable HTTP transport using Server-Sent Events (SSE)
- **Tool Handlers**: Custom functions for each tool
- **Resource Handlers**: Individual handlers for each resource
- **Prompt Handlers**: Custom functions for each prompt

## Extending the Server

### Adding a New Tool

```go
newTool := mcp.NewTool("tool_name",
    mcp.WithDescription("Tool description"),
    mcp.WithString("param_name",
        mcp.Required(),
        mcp.Description("Parameter description"),
    ),
)

s.AddTool(newTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Your tool logic here
    return mcp.NewToolResultText("Result"), nil
})
```

### Adding a New Resource

```go
s.AddResource("resource://uri",
    mcp.WithResourceName("Resource Name"),
    mcp.WithResourceDescription("Resource description"),
    mcp.WithResourceMIMEType("text/plain"),
)

// Handle in the resource handler function
```

### Adding a New Prompt

```go
newPrompt := mcp.NewPrompt("prompt_name",
    mcp.WithPromptDescription("Prompt description"),
    mcp.WithPromptArgument("arg_name",
        mcp.WithPromptArgumentDescription("Argument description"),
        mcp.WithPromptArgumentRequired(true),
    ),
)

s.AddPrompt(newPrompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    // Your prompt logic here
    return &mcp.GetPromptResult{
        Description: "Prompt result",
        Messages: []mcp.PromptMessage{...},
    }, nil
})
```

## License

This is example code for demonstration purposes.

# chatgptApp

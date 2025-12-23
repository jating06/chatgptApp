# API Reference

Complete API reference for the MCP Server.

## Base URL

```
http://localhost:8080
```

## Endpoints

### Health Check

```
GET /health
```

Returns server health status.

**Response:**
```
OK
```

**Status Code:** 200

---

### MCP Endpoint

```
POST /mcp
```

Main endpoint for all MCP operations using JSON-RPC 2.0.

**Headers:**
- `Content-Type: application/json`
- `Accept: application/json`

---

## JSON-RPC Methods

All requests follow JSON-RPC 2.0 format:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "method_name",
  "params": {}
}
```

All responses follow:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {}
}
```

Or in case of error:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32600,
    "message": "Error message"
  }
}
```

---

## Initialize

Initialize the MCP connection.

**Method:** `initialize`

**Parameters:**
```json
{
  "protocolVersion": "2024-11-05",
  "capabilities": {
    "tools": {},
    "resources": {},
    "prompts": {}
  },
  "clientInfo": {
    "name": "client-name",
    "version": "1.0.0"
  }
}
```

**Response:**
```json
{
  "protocolVersion": "2024-11-05",
  "capabilities": {
    "prompts": {
      "listChanged": true
    },
    "resources": {
      "subscribe": true
    },
    "tools": {
      "listChanged": true
    }
  },
  "serverInfo": {
    "name": "example-mcp-server",
    "version": "1.0.0"
  }
}
```

---

## Tools

### List Tools

List all available tools.

**Method:** `tools/list`

**Parameters:** None

**Response:**
```json
{
  "tools": [
    {
      "name": "echo",
      "description": "Echoes back the input text",
      "inputSchema": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string",
            "description": "The message to echo back"
          }
        },
        "required": ["message"]
      },
      "annotations": {
        "readOnlyHint": false,
        "destructiveHint": true,
        "idempotentHint": false,
        "openWorldHint": true
      }
    }
  ]
}
```

---

### Call Tool

Execute a tool with given arguments.

**Method:** `tools/call`

**Parameters:**
```json
{
  "name": "tool_name",
  "arguments": {
    "param1": "value1",
    "param2": "value2"
  }
}
```

**Response:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "Tool result"
    }
  ]
}
```

---

### Echo Tool

Echoes back the input message.

**Tool Name:** `echo`

**Arguments:**
- `message` (string, required) - The message to echo back

**Example Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "message": "Hello, World!"
    }
  }
}
```

**Example Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Echo: Hello, World!"
      }
    ]
  }
}
```

---

### Add Tool

Adds two numbers together.

**Tool Name:** `add`

**Arguments:**
- `a` (number, required) - First number
- `b` (number, required) - Second number

**Example Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "add",
    "arguments": {
      "a": 42,
      "b": 58
    }
  }
}
```

**Example Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Result: 100.00"
      }
    ]
  }
}
```

---

### Get Time Tool

Returns the current server time.

**Tool Name:** `get_time`

**Arguments:** None

**Example Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get_time"
  }
}
```

**Example Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Current time: 2025-12-22T16:45:50+05:30"
      }
    ]
  }
}
```

---

## Resources

### List Resources

List all available resources.

**Method:** `resources/list`

**Parameters:** None

**Response:**
```json
{
  "resources": [
    {
      "uri": "server://info",
      "name": "Server Information",
      "description": "Information about this MCP server",
      "mimeType": "text/plain"
    }
  ]
}
```

---

### Read Resource

Read a resource by URI.

**Method:** `resources/read`

**Parameters:**
```json
{
  "uri": "resource://uri"
}
```

**Response:**
```json
{
  "contents": [
    {
      "uri": "resource://uri",
      "mimeType": "text/plain",
      "text": "Resource content"
    }
  ]
}
```

---

### Server Info Resource

Provides server information.

**URI:** `server://info`

**MIME Type:** `text/plain`

**Example Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "resources/read",
  "params": {
    "uri": "server://info"
  }
}
```

**Example Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "contents": [
      {
        "uri": "server://info",
        "mimeType": "text/plain",
        "text": "Server: example-mcp-server\nVersion: 1.0.0\nTime: 2025-12-22T16:45:50+05:30"
      }
    ]
  }
}
```

---

## Prompts

### List Prompts

List all available prompts.

**Method:** `prompts/list`

**Parameters:** None

**Response:**
```json
{
  "prompts": [
    {
      "name": "greeting",
      "description": "Generate a personalized greeting",
      "arguments": [
        {
          "name": "name",
          "description": "Name of the person to greet",
          "required": true
        }
      ]
    }
  ]
}
```

---

### Get Prompt

Get a prompt with arguments.

**Method:** `prompts/get`

**Parameters:**
```json
{
  "name": "prompt_name",
  "arguments": {
    "arg1": "value1",
    "arg2": "value2"
  }
}
```

**Response:**
```json
{
  "description": "Prompt description",
  "messages": [
    {
      "role": "user",
      "content": {
        "type": "text",
        "text": "Prompt text"
      }
    }
  ]
}
```

---

### Greeting Prompt

Generates a personalized greeting.

**Prompt Name:** `greeting`

**Arguments:**
- `name` (string, required) - Name of the person to greet

**Example Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "prompts/get",
  "params": {
    "name": "greeting",
    "arguments": {
      "name": "Alice"
    }
  }
}
```

**Example Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "description": "A personalized greeting",
    "messages": [
      {
        "role": "user",
        "content": {
          "type": "text",
          "text": "Hello, Alice! Welcome to our MCP server."
        }
      }
    ]
  }
}
```

---

### Code Review Prompt

Generates a code review prompt template.

**Prompt Name:** `code_review`

**Arguments:**
- `language` (string, required) - Programming language for the code review

**Example Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "prompts/get",
  "params": {
    "name": "code_review",
    "arguments": {
      "language": "Go"
    }
  }
}
```

**Example Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "description": "Code review guidelines",
    "messages": [
      {
        "role": "user",
        "content": {
          "type": "text",
          "text": "Please review the following Go code for:\n1. Best practices\n2. Security issues\n3. Performance concerns\n4. Code style"
        }
      }
    ]
  }
}
```

---

## Error Codes

Standard JSON-RPC 2.0 error codes:

| Code | Message | Description |
|------|---------|-------------|
| -32700 | Parse error | Invalid JSON was received |
| -32600 | Invalid Request | The JSON sent is not a valid Request object |
| -32601 | Method not found | The method does not exist |
| -32602 | Invalid params | Invalid method parameter(s) |
| -32603 | Internal error | Internal JSON-RPC error |

**Example Error Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32602,
    "message": "Invalid params: message must be a string"
  }
}
```

---

## Rate Limiting

Currently, there is no rate limiting implemented. For production use, consider adding:
- Request rate limiting per client
- Concurrent request limits
- Timeout enforcement

---

## Authentication

Currently, there is no authentication implemented. For production use, consider adding:
- API key authentication
- JWT tokens
- OAuth 2.0
- mTLS

---

## Best Practices

1. **Always check for errors** in responses
2. **Use sequential request IDs** for easier debugging
3. **Validate arguments** before sending requests
4. **Handle timeouts** appropriately
5. **Log requests and responses** for debugging
6. **Use proper types** (numbers as numbers, not strings)
7. **Close connections** properly when done

---

## Testing

Use the provided test script:

```bash
./test_server.sh
```

Or test individual endpoints with curl:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

---

## Support

For issues or questions:
- Check the [README.md](README.md)
- Review [EXAMPLES.md](EXAMPLES.md)
- Consult the [MCP Specification](https://modelcontextprotocol.io/)
- Visit [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)




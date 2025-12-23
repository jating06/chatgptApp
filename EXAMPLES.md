# Examples

This document provides practical examples of using the MCP server.

## Table of Contents

- [Running the Examples](#running-the-examples)
- [Simple Go Client](#simple-go-client)
- [cURL Examples](#curl-examples)
- [Python Client Example](#python-client-example)
- [JavaScript Client Example](#javascript-client-example)

## Running the Examples

### Go Client

The repository includes a complete Go client example:

```bash
# Make sure the server is running first
./bin/mcp-server

# In another terminal, run the client
cd examples
go run simple_client.go
```

### Test Script

Run the automated test script:

```bash
./test_server.sh
```

## Simple Go Client

See `examples/simple_client.go` for a complete, working example of a Go client that:

- Initializes a connection
- Lists available tools, resources, and prompts
- Calls tools with arguments
- Reads resources
- Gets prompts with arguments

The client demonstrates proper error handling and JSON-RPC communication.

## cURL Examples

### Initialize Connection

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {
        "tools": {},
        "resources": {},
        "prompts": {}
      },
      "clientInfo": {
        "name": "curl-client",
        "version": "1.0.0"
      }
    }
  }'
```

### List and Call Tools

```bash
# List all tools
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'

# Call echo tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "echo",
      "arguments": {
        "message": "Hello, MCP!"
      }
    }
  }'

# Call add tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "tools/call",
    "params": {
      "name": "add",
      "arguments": {
        "a": 42,
        "b": 58
      }
    }
  }'

# Call get_time tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 5,
    "method": "tools/call",
    "params": {
      "name": "get_time"
    }
  }'
```

### Work with Resources

```bash
# List resources
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 6,
    "method": "resources/list"
  }'

# Read server info resource
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 7,
    "method": "resources/read",
    "params": {
      "uri": "server://info"
    }
  }'
```

### Work with Prompts

```bash
# List prompts
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 8,
    "method": "prompts/list"
  }'

# Get greeting prompt
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 9,
    "method": "prompts/get",
    "params": {
      "name": "greeting",
      "arguments": {
        "name": "Alice"
      }
    }
  }'

# Get code review prompt
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 10,
    "method": "prompts/get",
    "params": {
      "name": "code_review",
      "arguments": {
        "language": "Go"
      }
    }
  }'
```

## Python Client Example

Here's a simple Python client:

```python
import requests
import json

class MCPClient:
    def __init__(self, base_url):
        self.base_url = base_url
        self.request_id = 0
        
    def send_request(self, method, params=None):
        self.request_id += 1
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method
        }
        if params:
            request["params"] = params
            
        response = requests.post(
            self.base_url,
            json=request,
            headers={"Content-Type": "application/json"}
        )
        
        result = response.json()
        if "error" in result:
            raise Exception(f"MCP Error: {result['error']}")
        return result.get("result")
    
    def initialize(self):
        return self.send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "capabilities": {
                "tools": {},
                "resources": {},
                "prompts": {}
            },
            "clientInfo": {
                "name": "python-client",
                "version": "1.0.0"
            }
        })
    
    def list_tools(self):
        return self.send_request("tools/list")
    
    def call_tool(self, name, arguments):
        return self.send_request("tools/call", {
            "name": name,
            "arguments": arguments
        })

# Usage
client = MCPClient("http://localhost:8080/mcp")

# Initialize
print("Initializing...")
init_result = client.initialize()
print(json.dumps(init_result, indent=2))

# List tools
print("\nListing tools...")
tools = client.list_tools()
print(json.dumps(tools, indent=2))

# Call echo tool
print("\nCalling echo tool...")
echo_result = client.call_tool("echo", {"message": "Hello from Python!"})
print(json.dumps(echo_result, indent=2))

# Call add tool
print("\nCalling add tool...")
add_result = client.call_tool("add", {"a": 10, "b": 20})
print(json.dumps(add_result, indent=2))
```

## JavaScript Client Example

Here's a simple JavaScript/Node.js client:

```javascript
const fetch = require('node-fetch');

class MCPClient {
    constructor(baseUrl) {
        this.baseUrl = baseUrl;
        this.requestId = 0;
    }
    
    async sendRequest(method, params = null) {
        this.requestId++;
        const request = {
            jsonrpc: "2.0",
            id: this.requestId,
            method: method
        };
        
        if (params) {
            request.params = params;
        }
        
        const response = await fetch(this.baseUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(request)
        });
        
        const result = await response.json();
        if (result.error) {
            throw new Error(`MCP Error: ${JSON.stringify(result.error)}`);
        }
        return result.result;
    }
    
    async initialize() {
        return this.sendRequest("initialize", {
            protocolVersion: "2024-11-05",
            capabilities: {
                tools: {},
                resources: {},
                prompts: {}
            },
            clientInfo: {
                name: "javascript-client",
                version: "1.0.0"
            }
        });
    }
    
    async listTools() {
        return this.sendRequest("tools/list");
    }
    
    async callTool(name, arguments) {
        return this.sendRequest("tools/call", {
            name: name,
            arguments: arguments
        });
    }
    
    async listResources() {
        return this.sendRequest("resources/list");
    }
    
    async listPrompts() {
        return this.sendRequest("prompts/list");
    }
}

// Usage
(async () => {
    const client = new MCPClient("http://localhost:8080/mcp");
    
    // Initialize
    console.log("Initializing...");
    const initResult = await client.initialize();
    console.log(JSON.stringify(initResult, null, 2));
    
    // List tools
    console.log("\nListing tools...");
    const tools = await client.listTools();
    console.log(JSON.stringify(tools, null, 2));
    
    // Call echo tool
    console.log("\nCalling echo tool...");
    const echoResult = await client.callTool("echo", {
        message: "Hello from JavaScript!"
    });
    console.log(JSON.stringify(echoResult, null, 2));
    
    // Call add tool
    console.log("\nCalling add tool...");
    const addResult = await client.callTool("add", { a: 15, b: 25 });
    console.log(JSON.stringify(addResult, null, 2));
})();
```

## Error Handling

All clients should handle errors properly. The server returns JSON-RPC error responses:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32600,
    "message": "Invalid Request"
  }
}
```

Common error codes:
- `-32700`: Parse error
- `-32600`: Invalid Request
- `-32601`: Method not found
- `-32602`: Invalid params
- `-32603`: Internal error

## Best Practices

1. **Always initialize**: Call the `initialize` method before any other operations
2. **Handle errors**: Check for error responses in all client implementations
3. **Validate arguments**: Ensure tool arguments match the expected schema
4. **Use proper types**: Numbers should be sent as numbers, not strings
5. **Close connections**: Clean up resources when done (especially in stateful mode)
6. **Timeout handling**: Set appropriate timeouts for long-running operations

## Next Steps

- Explore the source code in `main.go` to understand server implementation
- Create custom tools for your specific use case
- Add authentication and authorization if needed
- Deploy to production with proper monitoring and logging




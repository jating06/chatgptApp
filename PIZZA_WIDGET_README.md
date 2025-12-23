# Pizza List Widget - MCP Tool Implementation

This implementation adds a `pizza_list` tool to the Go MCP server, inspired by OpenAI's pizzaz example from the `openai-apps-sdk-examples` repository.

## Overview

The pizza list widget displays an interactive list of the best pizzerias with ratings, locations, and images. This demonstrates how to create MCP tools that return HTML widgets compatible with ChatGPT's widget system.

## Implementation Details

### 1. Tool Definition (`pizza_list`)

**Location**: `main.go` - `registerTools()` function

The tool accepts an optional `pizzaTopping` parameter and returns:
- A text response with the list of pizzerias
- Structured JSON data with pizza information
- A reference to the widget resource

**Example Tool Call**:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "pizza_list",
    "arguments": {
      "pizzaTopping": "pepperoni"
    }
  }
}
```

### 2. Widget Resource

**Location**: `main.go` - `registerResources()` function

- **URI**: `widget://pizza-list`
- **MIME Type**: `text/html+skybridge`
- **HTML File**: `ui/pizza-list.html`

**Metadata**:
- `openai/widgetCSP`: Content Security Policy for loading images from `persistent.oaistatic.com`
- `openai/outputTemplate`: Points to the widget URI
- `openai/toolInvocation/invoking`: Status text while loading ("Hand-tossing a list")
- `openai/toolInvocation/invoked`: Status text after loaded ("Served a fresh list")
- `openai/widgetAccessible`: Marks the widget as accessible

### 3. HTML Widget

**Location**: `ui/pizza-list.html`

A self-contained HTML file with:
- Embedded CSS for styling (responsive design)
- JavaScript for rendering the pizza list
- Interactive elements (buttons, hover effects)
- Sample pizza data hardcoded in the JavaScript

**Features**:
- Responsive design (mobile and desktop)
- Pizza images from persistent.oaistatic.com
- Star ratings display
- Location information
- Add to list buttons
- Save list functionality

## Pizza Data Structure

Each pizza entry contains:
```json
{
  "id": "nova-slice-lab",
  "name": "Nova Slice Lab",
  "city": "North Beach",
  "rating": 4.8,
  "thumbnail": "https://persistent.oaistatic.com/pizzaz/pizzaz-1.png"
}
```

## Testing

### Prerequisites
- Server must be running on port 8080
- `jq` installed for JSON formatting

### Run Tests

```bash
# Start the server
./build/mcp-server

# In another terminal, run the test script
./test_pizza_widget.sh
```

### Manual Testing

1. **List Tools**:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | jq '.'
```

2. **Call Pizza List Tool**:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/call",
    "params": {
      "name": "pizza_list",
      "arguments": {"pizzaTopping": "mushrooms"}
    }
  }' | jq '.'
```

3. **Read Widget Resource**:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "resources/read",
    "params": {"uri": "widget://pizza-list"}
  }' | jq '.result.contents[0]._meta'
```

## Key Differences from OpenAI's Node.js Example

1. **Language**: Go instead of Node.js/TypeScript
2. **MCP Library**: Uses `github.com/mark3labs/mcp-go` instead of `@modelcontextprotocol/sdk`
3. **Transport**: HTTP/JSON-RPC instead of SSE (Server-Sent Events)
4. **Widget HTML**: Self-contained with embedded data instead of React components
5. **Stateless**: Simpler stateless HTTP server instead of session-based SSE

## Files Modified/Created

- ✅ `main.go` - Added `pizza_list` tool and resource
- ✅ `ui/pizza-list.html` - New HTML widget file
- ✅ `test_pizza_widget.sh` - Test script

## Integration with ChatGPT

When integrated with ChatGPT:
1. User asks: "Show me the best pizza places"
2. ChatGPT calls the `pizza_list` tool
3. Server returns structured data and widget reference
4. ChatGPT displays the interactive widget
5. User can interact with the list (add pizzerias, save list)

## References

- OpenAI Apps SDK Examples: https://github.com/openai/openai-apps-sdk-examples
- MCP Go Library: https://github.com/mark3labs/mcp-go
- MCP Specification: https://spec.modelcontextprotocol.io/

## Next Steps

To extend this implementation:
1. Add backend API to fetch real pizza data
2. Implement save/add functionality with database
3. Add more widget variations (map view, carousel)
4. Add filtering and search capabilities
5. Integrate with real restaurant APIs (Yelp, Google Places)


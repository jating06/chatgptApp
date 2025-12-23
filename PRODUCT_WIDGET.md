# Product Widget Tool

This document describes the `list_products` tool and its associated HTML widget resource.

## Overview

The product widget is a simple MCP tool that demonstrates how to create an interactive HTML widget for displaying and selecting products. It consists of two main components:

1. **Tool**: `list_products` - Returns product data and references the widget
2. **Resource**: `widget://list-products` - Provides the HTML+JavaScript widget

## Tool: list_products

### Description
Display an interactive product selection widget

### Parameters
None

### Response
Returns a text response containing:
- List of available products with names and prices
- JSON data structure with product information
- Reference to the widget resource URI

### Example Usage

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "list_products",
    "arguments": {}
  }
}
```

### Example Response

```json
{
  "content": [
    {
      "type": "text",
      "text": "Product selection widget loaded.\n\nAvailable products:\n- Premium Widget: $99.99\n- Standard Package: $49.99\n- Basic Starter: $29.99\n- Enterprise Solution: $199.99\n\nWidget data: {...}\nWidget resource: widget://list-products"
    }
  ]
}
```

## Resource: widget://list-products

### Description
Interactive HTML widget for selecting products

### MIME Type
`text/html+skybridge`

### Features
- Displays a list of products with checkboxes
- Allows multiple product selection
- Shows success/error messages on form submission
- Styled with inline CSS for better presentation
- Responds to OpenAI's `set_globals` event for dynamic data

### HTML Structure

The widget includes:
- Product list with checkboxes
- Submit button
- Result display area
- JavaScript event handlers

### Event Handling

The widget listens for the `openai:set_globals` event to receive product data:

```javascript
window.addEventListener("openai:set_globals", handleSetGlobal, {
  passive: true,
});
```

Expected data format:
```json
{
  "products": [
    {
      "name": "Product Name",
      "price": "99.99",
      "priceId": "price_id_here"
    }
  ]
}
```

## Testing

Run the test script to verify the tool and resource:

```bash
./test_product_widget_simple.sh
```

This will test:
1. Connection initialization
2. Tool listing (verifies `list_products` is available)
3. Tool execution (calls `list_products`)
4. Resource listing (verifies widget resource exists)
5. Resource reading (retrieves the HTML widget)

## Implementation Details

### File Structure

```
chatgptApp/
├── main.go                          # Server with tool and resource registration
├── ui/
│   └── list-products.html           # HTML widget file
├── test_product_widget_simple.sh    # Test script
└── PRODUCT_WIDGET.md               # This documentation
```

### Code Location

- **Tool Registration**: `main.go` - `registerTools()` function
- **Resource Registration**: `main.go` - `registerResources()` function
- **Widget HTML**: `ui/list-products.html`

### Sample Products

The tool returns these sample products:
- Premium Widget - $99.99
- Standard Package - $49.99
- Basic Starter - $29.99
- Enterprise Solution - $199.99

## Customization

### Adding More Products

Edit the `products` slice in the `list_products` tool handler in `main.go`:

```go
products := []map[string]interface{}{
    {
        "name":    "Your Product Name",
        "price":   "XX.XX",
        "priceId": "your_price_id",
    },
    // Add more products...
}
```

### Modifying the Widget UI

Edit `ui/list-products.html` to customize:
- Styling (inline CSS in the `renderProduct` and `renderApp` functions)
- Layout and structure
- Form submission behavior
- Event handlers

### Adding CSP Domains

If your widget needs to connect to external domains, update the metadata in the resource handler:

```go
metadata := map[string]interface{}{
    "openai/widgetCSP": map[string]interface{}{
        "connect_domains":  []string{"https://example.com"},
        "resource_domains": []string{"https://cdn.example.com"},
    },
}
```

## Integration

This widget can be integrated with:
- ChatGPT with MCP support
- Claude Desktop with MCP support
- Any MCP-compatible client

The client should:
1. Call the `list_products` tool
2. Read the `widget://list-products` resource
3. Render the HTML widget
4. Send product data via the `openai:set_globals` event

## Notes

- The widget uses inline styles for simplicity and portability
- Default products are shown if no globals are set (for testing)
- The MIME type `text/html+skybridge` indicates this is an interactive widget
- The widget is self-contained with no external dependencies




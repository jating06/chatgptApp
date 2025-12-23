#!/bin/bash
# Test script for pizza-list widget

echo "=== Testing Pizza List Widget ==="
echo ""

echo "1. Testing server health check..."
curl -s http://localhost:8080/health
echo -e "\n"

echo "2. Getting server info..."
curl -s http://localhost:8080/mcp | jq '.tools[] | select(.name == "pizza_list")'
echo -e "\n"

echo "3. Listing all tools..."
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/list"
  }' | jq '.result.tools[] | select(.name == "pizza_list")'
echo -e "\n"

echo "4. Calling pizza_list tool with mushrooms topping..."
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/call",
    "params": {
      "name": "pizza_list",
      "arguments": {
        "pizzaTopping": "mushrooms"
      }
    }
  }' | jq '.result.content[0].text' | head -15
echo -e "\n"

echo "5. Listing resources..."
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "resources/list"
  }' | jq '.result.resources[] | select(.uri == "widget://pizza-list")'
echo -e "\n"

echo "6. Reading pizza-list widget resource metadata..."
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 4,
    "method": "resources/read",
    "params": {
      "uri": "widget://pizza-list"
    }
  }' | jq '.result.contents[0] | {uri, mimeType, _meta}'
echo -e "\n"

echo "=== All tests completed ==="


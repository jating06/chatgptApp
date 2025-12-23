#!/bin/bash
# Quick Start Guide for Pizza List Widget

echo "🍕 Pizza List Widget - Quick Start Guide"
echo "=========================================="
echo ""

# Check if server is running
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Server is already running on port 8080"
else
    echo "❌ Server is not running. Starting server..."
    echo ""
    echo "To start the server manually, run:"
    echo "  ./build/mcp-server"
    echo ""
    echo "Or build and run with:"
    echo "  go build -o build/mcp-server main.go && ./build/mcp-server"
    exit 1
fi

echo ""
echo "📋 Available Commands:"
echo "====================="
echo ""

echo "1️⃣  List all tools:"
echo "   curl -X POST http://localhost:8080/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}' | jq '.result.tools[] | select(.name==\"pizza_list\")'"
echo ""

echo "2️⃣  Call pizza_list with pepperoni topping:"
echo "   curl -X POST http://localhost:8080/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"tools/call\",\"params\":{\"name\":\"pizza_list\",\"arguments\":{\"pizzaTopping\":\"pepperoni\"}}}' | jq '.'"
echo ""

echo "3️⃣  Get widget resource:"
echo "   curl -X POST http://localhost:8080/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":3,\"method\":\"resources/read\",\"params\":{\"uri\":\"widget://pizza-list\"}}' | jq '.result.contents[0]._meta'"
echo ""

echo "4️⃣  Run comprehensive tests:"
echo "   ./test_pizza_widget.sh"
echo ""

echo "📚 Documentation:"
echo "================"
echo "   - PIZZA_WIDGET_README.md - Full implementation details"
echo "   - IMPLEMENTATION_SUMMARY.md - Quick overview"
echo "   - API_REFERENCE.md - General MCP API reference"
echo ""

echo "🎨 Try it out:"
echo "=============="
read -p "Would you like to call the pizza_list tool now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo ""
    echo "Calling pizza_list with 'mushrooms' topping..."
    echo ""
    curl -s -X POST http://localhost:8080/mcp \
      -H 'Content-Type: application/json' \
      -d '{
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/call",
        "params": {
          "name": "pizza_list",
          "arguments": {
            "pizzaTopping": "mushrooms"
          }
        }
      }' | jq -r '.result.content[0].text'
    echo ""
fi

echo ""
echo "✨ Widget features:"
echo "   - 7 top-rated pizzerias"
echo "   - Star ratings and locations"
echo "   - Interactive add buttons"
echo "   - Responsive design"
echo "   - Images from persistent.oaistatic.com"
echo ""
echo "🚀 Ready to integrate with ChatGPT!"


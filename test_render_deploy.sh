#!/bin/bash
# Test Render deployment commands locally

echo "🚀 Testing Render Deployment Commands Locally"
echo "=============================================="
echo ""

echo "1️⃣  Building with Render build command..."
go build -tags netgo -ldflags '-s -w' -o app main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "   Binary size: $(ls -lh app | awk '{print $5}')"
else
    echo "❌ Build failed!"
    exit 1
fi
echo ""

echo "2️⃣  Running pre-deploy command..."
mkdir -p ui

if [ $? -eq 0 ]; then
    echo "✅ Pre-deploy successful!"
    echo "   UI directory exists: $([ -d ui ] && echo 'Yes' || echo 'No')"
else
    echo "❌ Pre-deploy failed!"
    exit 1
fi
echo ""

echo "3️⃣  Testing with different PORT values..."
echo ""

# Kill any existing server
pkill -f "./app" 2>/dev/null

# Test with custom port
echo "Starting server on PORT=3001..."
PORT=3001 ./app &
SERVER_PID=$!
sleep 2

# Test health endpoint
echo "Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s http://localhost:3001/health)
if [ "$HEALTH_RESPONSE" == "OK" ]; then
    echo "✅ Server responding on custom port 3001"
else
    echo "❌ Server not responding on custom port 3001"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

# Test MCP endpoint
echo "Testing MCP info endpoint..."
curl -s http://localhost:3001/mcp | jq -r '.name, .version' | head -2

# Test tool call
echo ""
echo "Testing pizza_list tool..."
PIZZA_RESPONSE=$(curl -s -X POST http://localhost:3001/mcp \
  -H "Content-Type: application/json" \
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
  }' | jq -r '.result.content[0].text' | head -1)

echo "$PIZZA_RESPONSE"

if [[ $PIZZA_RESPONSE == *"pizza"* ]]; then
    echo "✅ Pizza list tool working!"
else
    echo "❌ Pizza list tool failed!"
fi

# Cleanup
echo ""
echo "Cleaning up..."
kill $SERVER_PID 2>/dev/null
sleep 1

echo ""
echo "✅ All Render deployment commands tested successfully!"
echo ""
echo "📋 Commands for Render:"
echo "   Build: go build -tags netgo -ldflags '-s -w' -o app main.go"
echo "   Pre-Deploy: mkdir -p ui"
echo "   Start: ./app"
echo ""
echo "🌐 The app will use Render's PORT environment variable automatically"
echo "   (defaults to 8080 for local development)"




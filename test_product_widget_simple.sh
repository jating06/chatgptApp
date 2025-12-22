#!/bin/bash

# Simple test for the product widget tool and resource

set -e

SERVER_URL="http://localhost:8080"
MCP_ENDPOINT="${SERVER_URL}/mcp"

echo "================================"
echo "Product Widget Test"
echo "================================"
echo ""

# Store cookies to maintain session
COOKIE_JAR=$(mktemp)

# Test 1: Initialize
echo "Test 1: Initialize Connection"
echo "-----------------------------"
init_request='{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "tools": {},
      "resources": {}
    },
    "clientInfo": {
      "name": "product-widget-test",
      "version": "1.0.0"
    }
  }
}'

response=$(curl -s -c "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$init_request")

if echo "$response" | grep -q '"result"'; then
    echo "✓ Initialize successful"
else
    echo "✗ Initialize failed"
    echo "Response: $response"
    rm -f "$COOKIE_JAR"
    exit 1
fi
echo ""

# Test 2: List tools
echo "Test 2: List Tools (looking for list_products)"
echo "----------------------------------------------"
list_tools_request='{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}'

response=$(curl -s -b "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$list_tools_request")

if echo "$response" | grep -q 'list_products'; then
    echo "✓ list_products tool found!"
    echo "$response" | jq '.result.tools[] | select(.name == "list_products")'
else
    echo "✗ list_products tool not found"
    echo "Response: $response"
fi
echo ""

# Test 3: Call list_products tool
echo "Test 3: Call list_products Tool"
echo "-------------------------------"
list_products_request='{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "list_products",
    "arguments": {}
  }
}'

response=$(curl -s -b "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$list_products_request")

if echo "$response" | grep -q '"content"'; then
    echo "✓ list_products tool call successful"
    echo "Response:"
    echo "$response" | jq '.result'
else
    echo "✗ list_products tool call failed"
    echo "Response: $response"
fi
echo ""

# Test 4: List resources
echo "Test 4: List Resources (looking for widget://list-products)"
echo "-----------------------------------------------------------"
list_resources_request='{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "resources/list"
}'

response=$(curl -s -b "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$list_resources_request")

if echo "$response" | grep -q 'widget://list-products'; then
    echo "✓ Product widget resource found!"
    echo "$response" | jq '.result.resources[] | select(.uri == "widget://list-products")'
else
    echo "✗ Product widget resource not found"
    echo "Response: $response"
fi
echo ""

# Test 5: Read widget resource
echo "Test 5: Read Product Widget Resource"
echo "------------------------------------"
read_resource_request='{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "resources/read",
  "params": {
    "uri": "widget://list-products"
  }
}'

response=$(curl -s -b "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$read_resource_request")

if echo "$response" | grep -q '"contents"'; then
    echo "✓ Widget resource read successful"
    echo "Response (first 500 chars):"
    echo "$response" | jq '.result' | head -c 500
    echo "..."
else
    echo "✗ Widget resource read failed"
    echo "Response: $response"
fi
echo ""

# Cleanup
rm -f "$COOKIE_JAR"

echo "================================"
echo "Product Widget Tests Complete!"
echo "================================"



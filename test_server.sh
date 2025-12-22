#!/bin/bash

# Test script for MCP Server
# This script tests the basic functionality of the MCP server

set -e

SERVER_URL="http://localhost:8080"
MCP_ENDPOINT="${SERVER_URL}/mcp"
HEALTH_ENDPOINT="${SERVER_URL}/health"

echo "================================"
echo "MCP Server Test Suite"
echo "================================"
echo ""

# Test 1: Health Check
echo "Test 1: Health Check"
echo "--------------------"
response=$(curl -s -w "\n%{http_code}" "${HEALTH_ENDPOINT}")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n1)

if [ "$http_code" = "200" ] && [ "$body" = "OK" ]; then
    echo "✓ Health check passed"
else
    echo "✗ Health check failed (HTTP $http_code)"
    exit 1
fi
echo ""

# Test 2: Initialize Connection
echo "Test 2: Initialize Connection"
echo "-----------------------------"
init_request='{
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
      "name": "test-client",
      "version": "1.0.0"
    }
  }
}'

# Store cookies to maintain session
COOKIE_JAR=$(mktemp)
response=$(curl -s -c "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$init_request")

if echo "$response" | grep -q '"result"'; then
    echo "✓ Initialize request successful"
    echo "Response: $response" | head -c 200
    echo "..."
else
    echo "✗ Initialize request failed"
    echo "Response: $response"
    rm -f "$COOKIE_JAR"
    exit 1
fi
echo ""
echo ""

# Test 3: List Tools
echo "Test 3: List Tools"
echo "------------------"
list_tools_request='{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}'

response=$(curl -s -b "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$list_tools_request")

if echo "$response" | grep -q '"tools"'; then
    echo "✓ List tools successful"
    echo "Response: $response" | head -c 200
    echo "..."
else
    echo "✗ List tools failed"
    echo "Response: $response"
    rm -f "$COOKIE_JAR"
    exit 1
fi
echo ""
echo ""

# Test 4: Call Echo Tool
echo "Test 4: Call Echo Tool"
echo "----------------------"
echo_request='{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": {
      "message": "Hello MCP!"
    }
  }
}'

response=$(curl -s -b "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$echo_request")

if echo "$response" | grep -q '"content"'; then
    echo "✓ Echo tool call successful"
    echo "Response: $response" | head -c 200
    echo "..."
else
    echo "✗ Echo tool call failed"
    echo "Response: $response"
    rm -f "$COOKIE_JAR"
    exit 1
fi
echo ""
echo ""

# Test 5: Call Add Tool
echo "Test 5: Call Add Tool"
echo "---------------------"
add_request='{
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

response=$(curl -s -b "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$add_request")

if echo "$response" | grep -q '"content"'; then
    echo "✓ Add tool call successful"
    echo "Response: $response" | head -c 200
    echo "..."
else
    echo "✗ Add tool call failed"
    echo "Response: $response"
    rm -f "$COOKIE_JAR"
    exit 1
fi
echo ""
echo ""

# Test 6: List Resources
echo "Test 6: List Resources"
echo "----------------------"
list_resources_request='{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "resources/list"
}'

response=$(curl -s -b "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$list_resources_request")

if echo "$response" | grep -q '"resources"'; then
    echo "✓ List resources successful"
    echo "Response: $response" | head -c 200
    echo "..."
else
    echo "✗ List resources failed"
    echo "Response: $response"
    rm -f "$COOKIE_JAR"
    exit 1
fi
echo ""
echo ""

# Test 7: List Prompts
echo "Test 7: List Prompts"
echo "--------------------"
list_prompts_request='{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "prompts/list"
}'

response=$(curl -s -b "$COOKIE_JAR" -X POST "${MCP_ENDPOINT}" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -d "$list_prompts_request")

if echo "$response" | grep -q '"prompts"'; then
    echo "✓ List prompts successful"
    echo "Response: $response" | head -c 200
    echo "..."
else
    echo "✗ List prompts failed"
    echo "Response: $response"
    rm -f "$COOKIE_JAR"
    exit 1
fi
echo ""
echo ""

# Cleanup
rm -f "$COOKIE_JAR"

echo "================================"
echo "All tests passed! ✓"
echo "================================"


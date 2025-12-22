package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// This is a simple example client to demonstrate how to connect to the MCP server
// This file is for reference only and should not be run alongside main.go

type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *MCPError       `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func exampleClient() {
	serverURL := "http://localhost:8080/mcp"

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Example 1: Initialize connection
	fmt.Println("=== Example 1: Initialize Connection ===")
	initRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools":     map[string]interface{}{},
				"resources": map[string]interface{}{},
				"prompts":   map[string]interface{}{},
			},
			"clientInfo": map[string]string{
				"name":    "example-client",
				"version": "1.0.0",
			},
		},
	}

	sendRequest(client, serverURL, initRequest)

	// Example 2: List available tools
	fmt.Println("\n=== Example 2: List Tools ===")
	listToolsRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	sendRequest(client, serverURL, listToolsRequest)

	// Example 3: Call echo tool
	fmt.Println("\n=== Example 3: Call Echo Tool ===")
	echoRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "echo",
			"arguments": map[string]interface{}{
				"message": "Hello, MCP Server!",
			},
		},
	}

	sendRequest(client, serverURL, echoRequest)

	// Example 4: Call add tool
	fmt.Println("\n=== Example 4: Call Add Tool ===")
	addRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      4,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "add",
			"arguments": map[string]interface{}{
				"a": 42.5,
				"b": 57.5,
			},
		},
	}

	sendRequest(client, serverURL, addRequest)

	// Example 5: List resources
	fmt.Println("\n=== Example 5: List Resources ===")
	listResourcesRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      5,
		Method:  "resources/list",
	}

	sendRequest(client, serverURL, listResourcesRequest)

	// Example 6: Read resource
	fmt.Println("\n=== Example 6: Read Server Info Resource ===")
	readResourceRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      6,
		Method:  "resources/read",
		Params: map[string]interface{}{
			"uri": "server://info",
		},
	}

	sendRequest(client, serverURL, readResourceRequest)

	// Example 7: List prompts
	fmt.Println("\n=== Example 7: List Prompts ===")
	listPromptsRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      7,
		Method:  "prompts/list",
	}

	sendRequest(client, serverURL, listPromptsRequest)

	// Example 8: Get greeting prompt
	fmt.Println("\n=== Example 8: Get Greeting Prompt ===")
	getPromptRequest := MCPRequest{
		JSONRPC: "2.0",
		ID:      8,
		Method:  "prompts/get",
		Params: map[string]interface{}{
			"name": "greeting",
			"arguments": map[string]interface{}{
				"name": "Alice",
			},
		},
	}

	sendRequest(client, serverURL, getPromptRequest)
}

func sendRequest(client *http.Client, serverURL string, request MCPRequest) {
	// Marshal request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return
	}

	fmt.Printf("Request: %s\n", string(requestBody))

	// Create HTTP request
	req, err := http.NewRequestWithContext(context.Background(), "POST", serverURL, strings.NewReader(string(requestBody)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read SSE response
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			fmt.Printf("Response: %s\n", data)
			
			// Parse the response
			var mcpResp MCPResponse
			if err := json.Unmarshal([]byte(data), &mcpResp); err != nil {
				log.Printf("Error parsing response: %v", err)
				continue
			}

			if mcpResp.Error != nil {
				fmt.Printf("Error: %s (code: %d)\n", mcpResp.Error.Message, mcpResp.Error.Code)
			} else {
				// Pretty print the result
				var prettyJSON interface{}
				if err := json.Unmarshal(mcpResp.Result, &prettyJSON); err == nil {
					prettyBytes, _ := json.MarshalIndent(prettyJSON, "", "  ")
					fmt.Printf("Result:\n%s\n", string(prettyBytes))
				}
			}
			break // For simplicity, just read the first data event
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading response: %v", err)
	}
}


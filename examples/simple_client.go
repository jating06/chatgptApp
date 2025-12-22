package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MCPClient is a simple client for interacting with the MCP server
type MCPClient struct {
	baseURL    string
	httpClient *http.Client
	requestID  int
}

// MCPRequest represents a JSON-RPC request
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse represents a JSON-RPC response
type MCPResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *MCPError       `json:"error,omitempty"`
}

// MCPError represents a JSON-RPC error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewMCPClient creates a new MCP client
func NewMCPClient(baseURL string) *MCPClient {
	return &MCPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		requestID: 0,
	}
}

// sendRequest sends a request to the MCP server
func (c *MCPClient) sendRequest(method string, params interface{}) (*MCPResponse, error) {
	c.requestID++

	request := MCPRequest{
		JSONRPC: "2.0",
		ID:      c.requestID,
		Method:  method,
		Params:  params,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var mcpResp MCPResponse
	if err := json.Unmarshal(body, &mcpResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if mcpResp.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", mcpResp.Error.Code, mcpResp.Error.Message)
	}

	return &mcpResp, nil
}

// Initialize initializes the connection with the MCP server
func (c *MCPClient) Initialize() error {
	params := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools":     map[string]interface{}{},
			"resources": map[string]interface{}{},
			"prompts":   map[string]interface{}{},
		},
		"clientInfo": map[string]string{
			"name":    "simple-go-client",
			"version": "1.0.0",
		},
	}

	resp, err := c.sendRequest("initialize", params)
	if err != nil {
		return fmt.Errorf("initialize failed: %w", err)
	}

	fmt.Printf("Initialized successfully: %s\n", string(resp.Result))
	return nil
}

// ListTools lists all available tools
func (c *MCPClient) ListTools() (json.RawMessage, error) {
	resp, err := c.sendRequest("tools/list", nil)
	if err != nil {
		return nil, fmt.Errorf("list tools failed: %w", err)
	}

	return resp.Result, nil
}

// CallTool calls a tool with the given name and arguments
func (c *MCPClient) CallTool(name string, arguments map[string]interface{}) (json.RawMessage, error) {
	params := map[string]interface{}{
		"name":      name,
		"arguments": arguments,
	}

	resp, err := c.sendRequest("tools/call", params)
	if err != nil {
		return nil, fmt.Errorf("call tool failed: %w", err)
	}

	return resp.Result, nil
}

// ListResources lists all available resources
func (c *MCPClient) ListResources() (json.RawMessage, error) {
	resp, err := c.sendRequest("resources/list", nil)
	if err != nil {
		return nil, fmt.Errorf("list resources failed: %w", err)
	}

	return resp.Result, nil
}

// ReadResource reads a resource by URI
func (c *MCPClient) ReadResource(uri string) (json.RawMessage, error) {
	params := map[string]interface{}{
		"uri": uri,
	}

	resp, err := c.sendRequest("resources/read", params)
	if err != nil {
		return nil, fmt.Errorf("read resource failed: %w", err)
	}

	return resp.Result, nil
}

// ListPrompts lists all available prompts
func (c *MCPClient) ListPrompts() (json.RawMessage, error) {
	resp, err := c.sendRequest("prompts/list", nil)
	if err != nil {
		return nil, fmt.Errorf("list prompts failed: %w", err)
	}

	return resp.Result, nil
}

// GetPrompt gets a prompt with the given name and arguments
func (c *MCPClient) GetPrompt(name string, arguments map[string]string) (json.RawMessage, error) {
	params := map[string]interface{}{
		"name":      name,
		"arguments": arguments,
	}

	resp, err := c.sendRequest("prompts/get", params)
	if err != nil {
		return nil, fmt.Errorf("get prompt failed: %w", err)
	}

	return resp.Result, nil
}

func main() {
	// Create a new MCP client
	client := NewMCPClient("http://localhost:8080/mcp")

	// Initialize the connection
	fmt.Println("=== Initializing Connection ===")
	if err := client.Initialize(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println()

	// List available tools
	fmt.Println("=== Listing Tools ===")
	tools, err := client.ListTools()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Tools: %s\n\n", string(tools))

	// Call the echo tool
	fmt.Println("=== Calling Echo Tool ===")
	echoResult, err := client.CallTool("echo", map[string]interface{}{
		"message": "Hello from Go client!",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Echo result: %s\n\n", string(echoResult))

	// Call the add tool
	fmt.Println("=== Calling Add Tool ===")
	addResult, err := client.CallTool("add", map[string]interface{}{
		"a": 123.45,
		"b": 876.55,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Add result: %s\n\n", string(addResult))

	// Call the get_time tool
	fmt.Println("=== Calling Get Time Tool ===")
	timeResult, err := client.CallTool("get_time", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Time result: %s\n\n", string(timeResult))

	// List resources
	fmt.Println("=== Listing Resources ===")
	resources, err := client.ListResources()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Resources: %s\n\n", string(resources))

	// Read a resource
	fmt.Println("=== Reading Server Info Resource ===")
	serverInfo, err := client.ReadResource("server://info")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Server info: %s\n\n", string(serverInfo))

	// List prompts
	fmt.Println("=== Listing Prompts ===")
	prompts, err := client.ListPrompts()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Prompts: %s\n\n", string(prompts))

	// Get a prompt
	fmt.Println("=== Getting Greeting Prompt ===")
	greeting, err := client.GetPrompt("greeting", map[string]string{
		"name": "Alice",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Greeting prompt: %s\n\n", string(greeting))

	fmt.Println("=== All operations completed successfully! ===")
}


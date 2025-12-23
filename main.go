package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "example-mcp-server"
	serverVersion = "1.0.0"
)

func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
		server.WithPromptCapabilities(true),
	)

	// Register tools
	registerTools(s)

	// Register resources
	registerResources(s)

	// Register prompts
	registerPrompts(s)

	// Create streamable HTTP server (stateless mode for easier testing)
	streamableServer := server.NewStreamableHTTPServer(s, server.WithStateLess(true))

	// Setup HTTP server with custom mux
	mux := http.NewServeMux()
	
	// Wrap MCP handler to support GET requests with info page
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Return info page for GET requests
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			info := map[string]interface{}{
				"name":        serverName,
				"version":     serverVersion,
				"description": "MCP Server with Product Widget Tool",
				"endpoints": map[string]string{
					"mcp":    "/mcp (POST for MCP protocol)",
					"health": "/health (GET for health check)",
				},
				"tools": []string{
					"echo - Echoes back the input text",
					"add - Adds two numbers together",
					"get_time - Returns the current server time",
					"list_products - Display an interactive product selection widget",
					"pizza_list - Show Pizza List - Display an interactive list of the best pizzerias",
				},
				"resources": []string{
					"server://info - Server information",
					"widget://list-products - Product selection widget",
					"widget://pizza-list - Pizza list widget",
				},
				"usage": "Send POST requests with JSON-RPC 2.0 format to /mcp endpoint",
			}
			json.NewEncoder(w).Encode(info)
			return
		}
		// For POST requests, use the MCP handler
		streamableServer.ServeHTTP(w, r)
	})
	
	// Add a health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	serverPort := getServerPort()
	httpServer := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting MCP server on port %s", serverPort)
		log.Printf("MCP endpoint: http://localhost:%s/mcp", serverPort)
		log.Printf("Health check: http://localhost:%s/health", serverPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func registerTools(s *server.MCPServer) {
	// Example tool: Echo tool that returns the input
	echoTool := mcp.NewTool("echo",
		mcp.WithDescription("Echoes back the input text"),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("The message to echo back"),
		),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	)

	s.AddTool(echoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}
		
		message, ok := args["message"].(string)
		if !ok {
			return mcp.NewToolResultError("message must be a string"), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Echo: %s", message)), nil
	})

	// Example tool: Add numbers
	addTool := mcp.NewTool("add",
		mcp.WithDescription("Adds two numbers together"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number"),
		),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	)

	s.AddTool(addTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments"), nil
		}
		
		a, okA := args["a"].(float64)
		b, okB := args["b"].(float64)
		
		if !okA || !okB {
			return mcp.NewToolResultError("both a and b must be numbers"), nil
		}

		result := a + b
		return mcp.NewToolResultText(fmt.Sprintf("Result: %.2f", result)), nil
	})

	// Example tool: Get current time
	timeTool := mcp.NewTool("get_time",
		mcp.WithDescription("Returns the current server time"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	)

	s.AddTool(timeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		currentTime := time.Now().Format(time.RFC3339)
		return mcp.NewToolResultText(fmt.Sprintf("Current time: %s", currentTime)), nil
	})

	// Product listing tool with HTML widget
	listProductsTool := mcp.NewTool("list_products",
		mcp.WithDescription("Display an interactive product selection widget"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	)

	s.AddTool(listProductsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Sample product data with real stock images and descriptions
		products := []map[string]interface{}{
			{
				"name":        "Premium Widget",
				"price":       "99.99",
				"priceId":     "price_premium_widget",
				"description": "Our flagship product with advanced features and premium support",
				"image":       "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?w=150&h=150&fit=crop",
			},
			{
				"name":        "Standard Package",
				"price":       "49.99",
				"priceId":     "price_standard_package",
				"description": "Perfect for small teams with essential features included",
				"image":       "https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=150&h=150&fit=crop",
			},
			{
				"name":        "Basic Starter",
				"price":       "29.99",
				"priceId":     "price_basic_starter",
				"description": "Get started with our basic plan, ideal for individuals",
				"image":       "https://images.unsplash.com/photo-1484480974693-6ca0a78fb36b?w=150&h=150&fit=crop",
			},
			{
				"name":        "Enterprise Solution",
				"price":       "199.99",
				"priceId":     "price_enterprise_solution",
				"description": "Complete enterprise solution with dedicated support and custom features",
				"image":       "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=150&h=150&fit=crop",
			},
		}

		// Create structured content with products data
		productsJSON, err := json.Marshal(map[string]interface{}{
			"products": products,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal products: %v", err)), nil
		}

		// Return result with reference to the widget resource
		message := fmt.Sprintf("Product selection widget loaded.\n\nAvailable products:\n")
		for _, p := range products {
			message += fmt.Sprintf("- %s: $%s\n", p["name"], p["price"])
		}
		message += fmt.Sprintf("\nWidget data: %s\n", string(productsJSON))
		message += fmt.Sprintf("Widget resource: widget://list-products")

		return mcp.NewToolResultText(message), nil
	})

	// Pizza list tool with HTML widget (inspired by OpenAI's pizzaz example)
	pizzaListTool := mcp.NewTool("pizza_list",
		mcp.WithDescription("Show Pizza List - Display an interactive list of the best pizzerias"),
		mcp.WithString("pizzaTopping",
			mcp.Description("Topping to mention when rendering the widget (e.g., pepperoni, mushrooms, cheese)"),
		),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithOpenWorldHintAnnotation(false),
	)

	s.AddTool(pizzaListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]interface{})
		topping := "cheese" // default topping
		if ok {
			if t, exists := args["pizzaTopping"].(string); exists && t != "" {
				topping = t
			}
		}

		// Sample pizza data
		pizzas := []map[string]interface{}{
			{
				"id":        "nova-slice-lab",
				"name":      "Nova Slice Lab",
				"city":      "North Beach",
				"rating":    4.8,
				"thumbnail": "https://persistent.oaistatic.com/pizzaz/pizzaz-1.png",
			},
			{
				"id":        "midnight-marinara",
				"name":      "Midnight Marinara",
				"city":      "North Beach",
				"rating":    4.6,
				"thumbnail": "https://persistent.oaistatic.com/pizzaz/pizzaz-2.png",
			},
			{
				"id":        "cinder-oven-co",
				"name":      "Cinder Oven Co.",
				"city":      "Mission",
				"rating":    4.5,
				"thumbnail": "https://persistent.oaistatic.com/pizzaz/pizzaz-3.png",
			},
			{
				"id":        "neon-crust-works",
				"name":      "Neon Crust Works",
				"city":      "Alamo Square",
				"rating":    4.5,
				"thumbnail": "https://persistent.oaistatic.com/pizzaz/pizzaz-6.png",
			},
			{
				"id":        "luna-pie-collective",
				"name":      "Luna Pie Collective",
				"city":      "North Beach",
				"rating":    4.6,
				"thumbnail": "https://persistent.oaistatic.com/pizzaz/pizzaz-4.png",
			},
			{
				"id":        "bricklight-deep-dish",
				"name":      "Bricklight Deep Dish",
				"city":      "North Beach",
				"rating":    4.4,
				"thumbnail": "https://persistent.oaistatic.com/pizzaz/pizzaz-5.png",
			},
			{
				"id":        "garden-ember-pies",
				"name":      "Garden Ember Pies",
				"city":      "Lower Haight",
				"rating":    4.4,
				"thumbnail": "https://persistent.oaistatic.com/pizzaz/pizzaz-1.png",
			},
		}

		// Create structured content with pizzas data
		pizzasJSON, err := json.Marshal(map[string]interface{}{
			"pizzas":       pizzas,
			"pizzaTopping": topping,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal pizzas: %v", err)), nil
		}

		// Return result with reference to the widget resource
		message := fmt.Sprintf("🍕 Rendered a pizza list with %s topping!\n\n", topping)
		message += "National Best Pizza List:\n"
		for i, p := range pizzas {
			message += fmt.Sprintf("%d. %s - %s (Rating: %.1f)\n", i+1, p["name"], p["city"], p["rating"])
		}
		message += fmt.Sprintf("\nWidget data: %s\n", string(pizzasJSON))
		message += "Widget resource: widget://pizza-list"

		return mcp.NewToolResultText(message), nil
	})
}

func registerResources(s *server.MCPServer) {
	// Example resource: Server info
	serverInfoResource := mcp.Resource{
		URI:         "server://info",
		Name:        "Server Information",
		Description: "Information about this MCP server",
		MIMEType:    "text/plain",
	}

	s.AddResource(serverInfoResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		info := fmt.Sprintf("Server: %s\nVersion: %s\nTime: %s",
			serverName,
			serverVersion,
			time.Now().Format(time.RFC3339),
		)
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     info,
			},
		}, nil
	})

	// Product listing widget resource
	productWidgetResource := mcp.Resource{
		URI:         "widget://list-products",
		Name:        "Product Selection Widget",
		Description: "Interactive HTML widget for selecting products",
		MIMEType:    "text/html+skybridge",
	}

	s.AddResource(productWidgetResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Read the HTML widget file
		htmlContent, err := os.ReadFile("ui/list-products.html")
		if err != nil {
			return nil, fmt.Errorf("failed to read widget HTML: %v", err)
		}

		// Create metadata for CSP (Content Security Policy) to allow image loading
		metadata := map[string]interface{}{
			"openai/widgetCSP": map[string]interface{}{
				"connect_domains":  []string{"https://images.unsplash.com"},
				"resource_domains": []string{"https://images.unsplash.com"},
			},
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/html+skybridge",
				Text:     string(htmlContent),
				Meta:     metadata,
			},
		}, nil
	})

	// Pizza list widget resource (inspired by OpenAI's pizzaz example)
	pizzaListResource := mcp.Resource{
		URI:         "widget://pizza-list",
		Name:        "Pizza List Widget",
		Description: "Interactive HTML widget displaying the best pizzerias",
		MIMEType:    "text/html+skybridge",
	}

	s.AddResource(pizzaListResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Read the HTML widget file
		htmlContent, err := os.ReadFile("ui/pizza-list.html")
		if err != nil {
			return nil, fmt.Errorf("failed to read pizza widget HTML: %v", err)
		}

		// Create metadata for CSP (Content Security Policy) to allow image loading from persistent.oaistatic.com
		metadata := map[string]interface{}{
			"openai/widgetCSP": map[string]interface{}{
				"connect_domains":  []string{"https://persistent.oaistatic.com"},
				"resource_domains": []string{"https://persistent.oaistatic.com"},
			},
			"openai/outputTemplate":              "widget://pizza-list",
			"openai/toolInvocation/invoking":     "Hand-tossing a list",
			"openai/toolInvocation/invoked":      "Served a fresh list",
			"openai/widgetAccessible":            true,
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/html+skybridge",
				Text:     string(htmlContent),
				Meta:     metadata,
			},
		}, nil
	})
}

func registerPrompts(s *server.MCPServer) {
	// Example prompt: Greeting prompt
	greetingPrompt := mcp.Prompt{
		Name:        "greeting",
		Description: "Generate a personalized greeting",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "name",
				Description: "Name of the person to greet",
				Required:    true,
			},
		},
	}

	s.AddPrompt(greetingPrompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		name, ok := request.Params.Arguments["name"]
		if !ok || name == "" {
			return nil, fmt.Errorf("name argument is required")
		}

		return &mcp.GetPromptResult{
			Description: "A personalized greeting",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Hello, %s! Welcome to our MCP server.", name),
					},
				},
			},
		}, nil
	})

	// Example prompt: Code review prompt
	codeReviewPrompt := mcp.Prompt{
		Name:        "code_review",
		Description: "Generate a code review prompt",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "language",
				Description: "Programming language",
				Required:    true,
			},
		},
	}

	s.AddPrompt(codeReviewPrompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		language, ok := request.Params.Arguments["language"]
		if !ok || language == "" {
			return nil, fmt.Errorf("language argument is required")
		}

		return &mcp.GetPromptResult{
			Description: "Code review guidelines",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Please review the following %s code for:\n1. Best practices\n2. Security issues\n3. Performance concerns\n4. Code style", language),
					},
				},
			},
		}, nil
	})
}


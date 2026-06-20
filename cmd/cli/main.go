package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// Get the MCP server URL from environment or use default
	mcpURL := os.Getenv("MCP_SERVER_URL")
	if mcpURL == "" {
		mcpURL = "http://localhost:8080"
	}

	// Get max retry attempts and delay
	maxRetries := 30
	retryDelay := time.Second * 2

	fmt.Printf("Testing MCP server at %s\n", mcpURL)
	fmt.Println(strings.Repeat("=", 50))

	// Test health endpoint with retries
	if !testHealthWithRetries(mcpURL, maxRetries, retryDelay) {
		fmt.Println("FAILED: Health check failed after retries")
		os.Exit(1)
	}

	// Test tools endpoint
	if !testTools(mcpURL) {
		fmt.Println("FAILED: Tools endpoint test failed")
		os.Exit(1)
	}

	// Test call_tool endpoint
	if !testCallTool(mcpURL) {
		fmt.Println("FAILED: Call tool test failed")
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("SUCCESS: All tests passed!")
	os.Exit(0)
}

func testHealthWithRetries(mcpURL string, maxRetries int, delay time.Duration) bool {
	fmt.Println("\nTest 1: Health Check (with retries)")
	fmt.Printf("Checking health endpoint: %s/health\n", mcpURL)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := http.Get(mcpURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				fmt.Printf("✓ Health check passed (attempt %d/%d): %s\n", attempt, maxRetries, string(body))
				return true
			}
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempt < maxRetries {
			fmt.Printf("  Attempt %d failed, retrying in %v...\n", attempt, delay)
			time.Sleep(delay)
		}
	}

	fmt.Printf("✗ Health check failed after %d attempts\n", maxRetries)
	return false
}

func testTools(mcpURL string) bool {
	fmt.Println("\nTest 2: Get Tools")
	fmt.Printf("Checking tools endpoint: %s/tools\n", mcpURL)

	resp, err := http.Get(mcpURL + "/tools")
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("✗ Status code: %d (expected 200)\n", resp.StatusCode)
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("✗ Error reading response: %v\n", err)
		return false
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("✗ Error parsing JSON: %v\n", err)
		return false
	}

	tools, ok := result["tools"].([]interface{})
	if !ok {
		fmt.Printf("✗ Invalid response format\n")
		return false
	}

	fmt.Printf("✓ Tools endpoint returned %d tools\n", len(tools))
	for i, tool := range tools {
		if toolMap, ok := tool.(map[string]interface{}); ok {
			if name, ok := toolMap["name"].(string); ok {
				fmt.Printf("  - Tool %d: %s\n", i+1, name)
			}
		}
	}

	return true
}

func testCallTool(mcpURL string) bool {
	fmt.Println("\nTest 3: Call Tool (list_tasks)")
	fmt.Printf("Calling tool endpoint: %s/call_tool\n", mcpURL)

	// Create a request to list tasks
	toolRequest := map[string]interface{}{
		"name":  "tududi_list_tasks",
		"input": map[string]interface{}{},
	}

	requestBody, err := json.Marshal(toolRequest)
	if err != nil {
		fmt.Printf("✗ Error marshaling request: %v\n", err)
		return false
	}

	resp, err := http.Post(
		mcpURL+"/call_tool",
		"application/json",
		io.NopCloser(bytes.NewReader(requestBody)),
	)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("✗ Status code: %d\n", resp.StatusCode)
		fmt.Printf("  Response: %s\n", string(body))
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("✗ Error reading response: %v\n", err)
		return false
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("✗ Error parsing JSON: %v\n", err)
		return false
	}

	fmt.Printf("✓ Tool call successful\n")
	fmt.Printf("  Response: %s\n", string(body))

	return true
}

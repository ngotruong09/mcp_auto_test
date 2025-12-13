package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// MCP Protocol Structures
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// Browser context holder
type BrowserContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var browserCtx *BrowserContext

func main() {
	log.SetOutput(os.Stderr)
	log.Println("MCP ChromeDP Server starting...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req MCPRequest
		if err := json.Unmarshal(line, &req); err != nil {
			log.Printf("Error parsing request: %v", err)
			continue
		}

		resp := handleRequest(req)
		respJSON, _ := json.Marshal(resp)
		fmt.Println(string(respJSON))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}

	if browserCtx != nil {
		browserCtx.cancel()
	}
}

func handleRequest(req MCPRequest) MCPResponse {
	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		resp.Result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]bool{},
			},
			"serverInfo": map[string]string{
				"name":    "chromedp-mcp-server",
				"version": "1.0.0",
			},
		}

	case "tools/list":
		resp.Result = map[string]interface{}{
			"tools": getTools(),
		}

	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			resp.Error = &MCPError{Code: -32602, Message: "Invalid params"}
			return resp
		}

		result, err := executeTool(params.Name, params.Arguments)
		if err != nil {
			resp.Error = &MCPError{Code: -32000, Message: err.Error()}
		} else {
			resp.Result = map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": result,
					},
				},
			}
		}

	default:
		resp.Error = &MCPError{Code: -32601, Message: "Method not found"}
	}

	return resp
}

func getTools() []MCPTool {
	return []MCPTool{
		// Navigation & Browser Control
		{
			Name:        "playwright_navigate",
			Description: "Navigate to a URL",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]string{"type": "string", "description": "URL to navigate to"},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "playwright_goto",
			Description: "Navigate to a URL (alias for navigate)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]string{"type": "string", "description": "URL to navigate to"},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "playwright_go_back",
			Description: "Go back in browser history",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "playwright_go_forward",
			Description: "Go forward in browser history",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "playwright_reload",
			Description: "Reload the current page",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},

		// Screenshots & PDF
		{
			Name:        "playwright_screenshot",
			Description: "Take a screenshot of the page",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector for element (optional, full page if omitted)"},
					"fullPage": map[string]string{"type": "boolean", "description": "Capture full scrollable page"},
				},
			},
		},
		{
			Name:        "playwright_pdf",
			Description: "Generate PDF of the page",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},

		// Interaction
		{
			Name:        "playwright_click",
			Description: "Click an element",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector"},
				},
				"required": []string{"selector"},
			},
		},
		{
			Name:        "playwright_fill",
			Description: "Fill an input field",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector"},
					"value":    map[string]string{"type": "string", "description": "Value to fill"},
				},
				"required": []string{"selector", "value"},
			},
		},
		{
			Name:        "playwright_type",
			Description: "Type text into an element",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector"},
					"text":     map[string]string{"type": "string", "description": "Text to type"},
				},
				"required": []string{"selector", "text"},
			},
		},
		{
			Name:        "playwright_press",
			Description: "Press a keyboard key",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector (optional)"},
					"key":      map[string]string{"type": "string", "description": "Key to press (e.g., Enter, Escape)"},
				},
				"required": []string{"key"},
			},
		},
		{
			Name:        "playwright_select_option",
			Description: "Select an option from a dropdown",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector"},
					"value":    map[string]string{"type": "string", "description": "Option value to select"},
				},
				"required": []string{"selector", "value"},
			},
		},

		// Element Operations
		{
			Name:        "playwright_get_text",
			Description: "Get text content of an element",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector"},
				},
				"required": []string{"selector"},
			},
		},
		{
			Name:        "playwright_get_attribute",
			Description: "Get attribute value of an element",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector":  map[string]string{"type": "string", "description": "CSS selector"},
					"attribute": map[string]string{"type": "string", "description": "Attribute name"},
				},
				"required": []string{"selector", "attribute"},
			},
		},
		{
			Name:        "playwright_wait_for_selector",
			Description: "Wait for an element to appear",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector"},
					"timeout":  map[string]string{"type": "number", "description": "Timeout in milliseconds (default: 30000)"},
				},
				"required": []string{"selector"},
			},
		},
		{
			Name:        "playwright_query_selector",
			Description: "Check if an element exists",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector"},
				},
				"required": []string{"selector"},
			},
		},

		// Evaluation
		{
			Name:        "playwright_evaluate",
			Description: "Execute JavaScript in the page",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"script": map[string]string{"type": "string", "description": "JavaScript code to execute"},
				},
				"required": []string{"script"},
			},
		},
		{
			Name:        "playwright_console",
			Description: "Get console messages from the page",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

func ensureBrowser() error {
	if browserCtx == nil {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
		)
		allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		browserCtx = &BrowserContext{ctx: ctx, cancel: cancel}
	}
	return nil
}

func executeTool(name string, args map[string]interface{}) (string, error) {
	if err := ensureBrowser(); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(browserCtx.ctx, 30*time.Second)
	defer cancel()

	switch name {
	// Navigation & Browser Control
	case "playwright_navigate", "playwright_goto":
		url := args["url"].(string)
		err := chromedp.Run(ctx, chromedp.Navigate(url))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Navigated to %s", url), nil

	case "playwright_go_back":
		err := chromedp.Run(ctx, chromedp.NavigateBack())
		if err != nil {
			return "", err
		}
		return "Navigated back", nil

	case "playwright_go_forward":
		err := chromedp.Run(ctx, chromedp.NavigateForward())
		if err != nil {
			return "", err
		}
		return "Navigated forward", nil

	case "playwright_reload":
		err := chromedp.Run(ctx, chromedp.Reload())
		if err != nil {
			return "", err
		}
		return "Page reloaded", nil

	// Screenshots & PDF
	case "playwright_screenshot":
		var buf []byte
		selector, hasSelector := args["selector"].(string)

		if hasSelector && selector != "" {
			err := chromedp.Run(ctx, chromedp.Screenshot(selector, &buf, chromedp.NodeVisible))
			if err != nil {
				return "", err
			}
		} else {
			err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 90))
			if err != nil {
				return "", err
			}
		}

		encoded := base64.StdEncoding.EncodeToString(buf)
		return fmt.Sprintf("Screenshot captured (base64): %s", encoded[:100]+"..."), nil

	case "playwright_pdf":
		var buf []byte
		err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			return err
		}))
		if err != nil {
			return "", err
		}
		encoded := base64.StdEncoding.EncodeToString(buf)
		return fmt.Sprintf("PDF generated (base64): %s", encoded[:100]+"..."), nil

	// Interaction
	case "playwright_click":
		selector := args["selector"].(string)
		err := chromedp.Run(ctx, chromedp.Click(selector, chromedp.NodeVisible))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Clicked on %s", selector), nil

	case "playwright_fill":
		selector := args["selector"].(string)
		value := args["value"].(string)
		err := chromedp.Run(ctx,
			chromedp.Clear(selector),
			chromedp.SendKeys(selector, value, chromedp.NodeVisible),
		)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Filled %s with value", selector), nil

	case "playwright_type":
		selector := args["selector"].(string)
		text := args["text"].(string)
		err := chromedp.Run(ctx, chromedp.SendKeys(selector, text, chromedp.NodeVisible))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Typed text into %s", selector), nil

	case "playwright_press":
		key := args["key"].(string)
		selector, hasSelector := args["selector"].(string)

		if hasSelector && selector != "" {
			err := chromedp.Run(ctx, chromedp.SendKeys(selector, key, chromedp.NodeVisible))
			if err != nil {
				return "", err
			}
		} else {
			err := chromedp.Run(ctx, chromedp.KeyEvent(key))
			if err != nil {
				return "", err
			}
		}
		return fmt.Sprintf("Pressed key: %s", key), nil

	case "playwright_select_option":
		selector := args["selector"].(string)
		value := args["value"].(string)
		err := chromedp.Run(ctx, chromedp.SetValue(selector, value, chromedp.NodeVisible))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Selected option %s in %s", value, selector), nil

	// Element Operations
	case "playwright_get_text":
		selector := args["selector"].(string)
		var text string
		err := chromedp.Run(ctx, chromedp.Text(selector, &text, chromedp.NodeVisible))
		if err != nil {
			return "", err
		}
		return text, nil

	case "playwright_get_attribute":
		selector := args["selector"].(string)
		attribute := args["attribute"].(string)
		var value string
		err := chromedp.Run(ctx, chromedp.AttributeValue(selector, attribute, &value, nil, chromedp.NodeVisible))
		if err != nil {
			return "", err
		}
		return value, nil

	case "playwright_wait_for_selector":
		selector := args["selector"].(string)
		timeout := 30000
		if t, ok := args["timeout"].(float64); ok {
			timeout = int(t)
		}

		waitCtx, waitCancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		defer waitCancel()

		err := chromedp.Run(waitCtx, chromedp.WaitVisible(selector))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Element %s is visible", selector), nil

	case "playwright_query_selector":
		selector := args["selector"].(string)
		var nodes []*cdp.Node
		err := chromedp.Run(ctx, chromedp.Nodes(selector, &nodes))
		if err != nil {
			return "", err
		}
		if len(nodes) > 0 {
			return fmt.Sprintf("Element found: %s", selector), nil
		}
		return fmt.Sprintf("Element not found: %s", selector), nil

	// Evaluation
	case "playwright_evaluate":
		script := args["script"].(string)
		var result interface{}
		err := chromedp.Run(ctx, chromedp.Evaluate(script, &result))
		if err != nil {
			return "", err
		}
		resultJSON, _ := json.Marshal(result)
		return string(resultJSON), nil

	case "playwright_console":
		// Note: Console message collection requires setting up listeners
		// This is a simplified version
		return "Console messages: (feature requires console listener setup)", nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

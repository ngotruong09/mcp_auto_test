package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
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
	ctx             context.Context
	cancel          context.CancelFunc
	consoleMessages []string
	networkRequests []NetworkRequest
	dialogHandler   func(string) string
}

type NetworkRequest struct {
	URL    string `json:"url"`
	Method string `json:"method"`
	Status int    `json:"status"`
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

		// New Browser Tools
		{
			Name:        "browser_hover",
			Description: "Hover over an element",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector"},
				},
				"required": []string{"selector"},
			},
		},
		{
			Name:        "browser_resize",
			Description: "Resize browser window",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"width":  map[string]string{"type": "number", "description": "Window width in pixels"},
					"height": map[string]string{"type": "number", "description": "Window height in pixels"},
				},
				"required": []string{"width", "height"},
			},
		},
		{
			Name:        "browser_drag",
			Description: "Drag and drop from one element to another",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"from": map[string]string{"type": "string", "description": "CSS selector of element to drag from"},
					"to":   map[string]string{"type": "string", "description": "CSS selector of element to drop to"},
				},
				"required": []string{"from", "to"},
			},
		},
		{
			Name:        "browser_file_upload",
			Description: "Upload file to input element",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector of file input"},
					"filepath": map[string]string{"type": "string", "description": "Path to file to upload"},
				},
				"required": []string{"selector", "filepath"},
			},
		},
		{
			Name:        "browser_fill_form",
			Description: "Fill multiple form fields at once",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"fields": map[string]interface{}{
						"type":        "array",
						"description": "Array of {selector, value} objects",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"selector": map[string]string{"type": "string"},
								"value":    map[string]string{"type": "string"},
							},
						},
					},
				},
				"required": []string{"fields"},
			},
		},
		{
			Name:        "browser_close",
			Description: "Close the current page/tab",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "browser_snapshot",
			Description: "Take accessibility snapshot of the page",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"selector": map[string]string{"type": "string", "description": "CSS selector (optional)"},
				},
			},
		},
		{
			Name:        "browser_handle_dialog",
			Description: "Handle browser dialogs (alert, confirm, prompt)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"accept": map[string]string{"type": "boolean", "description": "Accept or dismiss dialog"},
					"text":   map[string]string{"type": "string", "description": "Text to enter for prompt (optional)"},
				},
				"required": []string{"accept"},
			},
		},
		{
			Name:        "browser_network_requests",
			Description: "Get list of network requests made by the page",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "browser_wait_for",
			Description: "Wait for text to appear/disappear or wait for time",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text":    map[string]string{"type": "string", "description": "Text to wait for"},
					"timeout": map[string]string{"type": "number", "description": "Timeout in milliseconds"},
					"state":   map[string]string{"type": "string", "description": "visible or hidden"},
				},
			},
		},
		{
			Name:        "browser_tabs",
			Description: "Manage browser tabs (list, new, close, select)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]string{"type": "string", "description": "Action: list, new, close, select"},
					"url":    map[string]string{"type": "string", "description": "URL for new tab"},
					"index":  map[string]string{"type": "number", "description": "Tab index for close/select"},
				},
				"required": []string{"action"},
			},
		},
		{
			Name:        "browser_install",
			Description: "Check browser installation and provide install instructions",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

func ensureBrowser() error {
	if browserCtx == nil {
		// Try to find Edge or Chrome executable
		execPath := findBrowserExecutable()

		var opts []chromedp.ExecAllocatorOption
		if execPath != "" {
			opts = []chromedp.ExecAllocatorOption{
				chromedp.ExecPath(execPath),
			}
		}

		opts = append(opts,
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)

		allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		browserCtx = &BrowserContext{ctx: ctx, cancel: cancel}
	}
	return nil
}

func findBrowserExecutable() string {
	var paths []string

	switch runtime.GOOS {
	case "windows":
		// Microsoft Edge paths (prioritize Edge on Windows)
		paths = []string{
			os.Getenv("PROGRAMFILES") + "\\Microsoft\\Edge\\Application\\msedge.exe",
			os.Getenv("PROGRAMFILES(X86)") + "\\Microsoft\\Edge\\Application\\msedge.exe",
			os.Getenv("LOCALAPPDATA") + "\\Microsoft\\Edge\\Application\\msedge.exe",
			// Chrome paths
			os.Getenv("PROGRAMFILES") + "\\Google\\Chrome\\Application\\chrome.exe",
			os.Getenv("PROGRAMFILES(X86)") + "\\Google\\Chrome\\Application\\chrome.exe",
			os.Getenv("LOCALAPPDATA") + "\\Google\\Chrome\\Application\\chrome.exe",
		}
	case "darwin":
		paths = []string{
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		}
	case "linux":
		// Try to find in PATH
		if path, err := exec.LookPath("microsoft-edge"); err == nil {
			return path
		}
		if path, err := exec.LookPath("microsoft-edge-stable"); err == nil {
			return path
		}
		if path, err := exec.LookPath("google-chrome"); err == nil {
			return path
		}
		if path, err := exec.LookPath("chromium-browser"); err == nil {
			return path
		}
		if path, err := exec.LookPath("chromium"); err == nil {
			return path
		}

		paths = []string{
			"/usr/bin/microsoft-edge",
			"/usr/bin/microsoft-edge-stable",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium-browser",
			"/usr/bin/chromium",
		}
	}

	// Check each path
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("Found browser at: %s", path)
			return path
		}
	}

	// Return empty string to use default
	return ""
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
		// Return collected console messages
		if browserCtx != nil && len(browserCtx.consoleMessages) > 0 {
			messages := strings.Join(browserCtx.consoleMessages, "\n")
			return messages, nil
		}
		return "No console messages captured", nil

	// New Browser Tools
	case "browser_hover":
		selector := args["selector"].(string)
		var nodes []*cdp.Node
		err := chromedp.Run(ctx,
			chromedp.Nodes(selector, &nodes, chromedp.ByQuery),
		)
		if err != nil {
			return "", err
		}
		if len(nodes) == 0 {
			return "", fmt.Errorf("element not found: %s", selector)
		}

		err = chromedp.Run(ctx, chromedp.MouseClickNode(nodes[0]))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Hovered over %s", selector), nil

	case "browser_resize":
		width := int64(args["width"].(float64))
		height := int64(args["height"].(float64))
		err := chromedp.Run(ctx, chromedp.EmulateViewport(width, height))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Resized window to %dx%d", width, height), nil

	case "browser_drag":
		from := args["from"].(string)
		to := args["to"].(string)

		// Get coordinates of both elements
		var fromX, fromY, toX, toY float64
		err := chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(`
				const fromEl = document.querySelector('%s');
				const toEl = document.querySelector('%s');
				const fromRect = fromEl.getBoundingClientRect();
				const toRect = toEl.getBoundingClientRect();
				({
					fromX: fromRect.left + fromRect.width/2,
					fromY: fromRect.top + fromRect.height/2,
					toX: toRect.left + toRect.width/2,
					toY: toRect.top + toRect.height/2
				});
			`, from, to), &map[string]interface{}{
				"fromX": &fromX,
				"fromY": &fromY,
				"toX":   &toX,
				"toY":   &toY,
			}),
		)
		if err != nil {
			return "", err
		}

		// Perform drag and drop
		err = chromedp.Run(ctx,
			chromedp.MouseClickXY(fromX, fromY),
			chromedp.Sleep(100*time.Millisecond),
			chromedp.MouseClickXY(toX, toY),
		)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Dragged from %s to %s", from, to), nil

	case "browser_file_upload":
		selector := args["selector"].(string)
		filepath := args["filepath"].(string)
		err := chromedp.Run(ctx, chromedp.SendKeys(selector, filepath, chromedp.NodeVisible))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Uploaded file %s to %s", filepath, selector), nil

	case "browser_fill_form":
		fieldsRaw := args["fields"].([]interface{})
		var actions []chromedp.Action

		for _, fieldRaw := range fieldsRaw {
			field := fieldRaw.(map[string]interface{})
			selector := field["selector"].(string)
			value := field["value"].(string)
			actions = append(actions,
				chromedp.Clear(selector),
				chromedp.SendKeys(selector, value, chromedp.NodeVisible),
			)
		}

		err := chromedp.Run(ctx, actions...)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Filled %d form fields", len(fieldsRaw)), nil

	case "browser_close":
		if browserCtx != nil {
			browserCtx.cancel()
			browserCtx = nil
		}
		return "Browser closed", nil

	case "browser_snapshot":
		var snapshot string
		selector, hasSelector := args["selector"].(string)

		script := `
			function getAccessibilityTree(el) {
				const role = el.getAttribute('role') || el.tagName.toLowerCase();
				const name = el.getAttribute('aria-label') || el.getAttribute('alt') || el.textContent?.slice(0, 50) || '';
				const children = Array.from(el.children).map(child => getAccessibilityTree(child));
				return { role, name: name.trim(), children: children.length > 0 ? children : undefined };
			}
			return JSON.stringify(getAccessibilityTree(document.body), null, 2);
		`

		if hasSelector && selector != "" {
			script = fmt.Sprintf(`
				const el = document.querySelector('%s');
				function getAccessibilityTree(el) {
					const role = el.getAttribute('role') || el.tagName.toLowerCase();
					const name = el.getAttribute('aria-label') || el.getAttribute('alt') || el.textContent?.slice(0, 50) || '';
					const children = Array.from(el.children).map(child => getAccessibilityTree(child));
					return { role, name: name.trim(), children: children.length > 0 ? children : undefined };
				}
				return JSON.stringify(getAccessibilityTree(el), null, 2);
			`, selector)
		}

		err := chromedp.Run(ctx, chromedp.Evaluate(script, &snapshot))
		if err != nil {
			return "", err
		}
		return snapshot, nil

	case "browser_handle_dialog":
		accept := args["accept"].(bool)
		text, hasText := args["text"].(string)

		err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			if accept {
				if hasText && text != "" {
					return page.HandleJavaScriptDialog(true).WithPromptText(text).Do(ctx)
				}
				return page.HandleJavaScriptDialog(true).Do(ctx)
			}
			return page.HandleJavaScriptDialog(false).Do(ctx)
		}))
		if err != nil {
			return "", err
		}
		return "Dialog handled", nil

	case "browser_network_requests":
		if browserCtx != nil && len(browserCtx.networkRequests) > 0 {
			requests, _ := json.MarshalIndent(browserCtx.networkRequests, "", "  ")
			return string(requests), nil
		}
		return "No network requests captured", nil

	case "browser_wait_for":
		text, hasText := args["text"].(string)
		timeout := 30000
		if t, ok := args["timeout"].(float64); ok {
			timeout = int(t)
		}
		state, hasState := args["state"].(string)
		if !hasState {
			state = "visible"
		}

		waitCtx, waitCancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		defer waitCancel()

		if hasText && text != "" {
			script := fmt.Sprintf(`
				Array.from(document.querySelectorAll('*')).some(el => 
					el.textContent && el.textContent.includes('%s')
				)
			`, text)

			if state == "hidden" {
				script = fmt.Sprintf(`
					!Array.from(document.querySelectorAll('*')).some(el => 
						el.textContent && el.textContent.includes('%s')
					)
				`, text)
			}

			err := chromedp.Run(waitCtx, chromedp.WaitVisible(`body`))
			if err != nil {
				return "", err
			}

			// Poll for text
			for i := 0; i < timeout/100; i++ {
				var found bool
				chromedp.Run(ctx, chromedp.Evaluate(script, &found))
				if found {
					return fmt.Sprintf("Text '%s' is %s", text, state), nil
				}
				time.Sleep(100 * time.Millisecond)
			}
			return "", fmt.Errorf("timeout waiting for text '%s'", text)
		} else {
			// Just wait for the timeout
			time.Sleep(time.Duration(timeout) * time.Millisecond)
			return fmt.Sprintf("Waited for %dms", timeout), nil
		}

	case "browser_tabs":
		action := args["action"].(string)

		switch action {
		case "list":
			var targets []*target.Info
			err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
				var err error
				targets, err = target.GetTargets().Do(ctx)
				return err
			}))
			if err != nil {
				return "", err
			}

			var tabs []string
			for i, t := range targets {
				if t.Type == "page" {
					tabs = append(tabs, fmt.Sprintf("%d: %s - %s", i, t.Title, t.URL))
				}
			}
			return strings.Join(tabs, "\n"), nil

		case "new":
			url := args["url"].(string)
			err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
				_, err := target.CreateTarget(url).Do(ctx)
				return err
			}))
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Opened new tab: %s", url), nil

		case "close":
			return "Tab close not fully implemented", nil

		case "select":
			return "Tab select not fully implemented", nil

		default:
			return "", fmt.Errorf("unknown tab action: %s", action)
		}

	case "browser_install":
		instructions := getBrowserInstallInstructions()
		return instructions, nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func getBrowserInstallInstructions() string {
	// Check if Edge or Chrome/Chromium is installed
	foundBrowsers := []string{}

	switch runtime.GOOS {
	case "windows":
		// Check Edge
		edgePaths := []string{
			os.Getenv("PROGRAMFILES") + "\\Microsoft\\Edge\\Application\\msedge.exe",
			os.Getenv("PROGRAMFILES(X86)") + "\\Microsoft\\Edge\\Application\\msedge.exe",
			os.Getenv("LOCALAPPDATA") + "\\Microsoft\\Edge\\Application\\msedge.exe",
		}
		for _, path := range edgePaths {
			if _, err := os.Stat(path); err == nil {
				foundBrowsers = append(foundBrowsers, "Microsoft Edge at: "+path)
				break
			}
		}

		// Check Chrome
		chromePaths := []string{
			os.Getenv("PROGRAMFILES") + "\\Google\\Chrome\\Application\\chrome.exe",
			os.Getenv("PROGRAMFILES(X86)") + "\\Google\\Chrome\\Application\\chrome.exe",
			os.Getenv("LOCALAPPDATA") + "\\Google\\Chrome\\Application\\chrome.exe",
		}
		for _, path := range chromePaths {
			if _, err := os.Stat(path); err == nil {
				foundBrowsers = append(foundBrowsers, "Google Chrome at: "+path)
				break
			}
		}

	case "darwin":
		// Check Edge
		edgePath := "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"
		if _, err := os.Stat(edgePath); err == nil {
			foundBrowsers = append(foundBrowsers, "Microsoft Edge at: "+edgePath)
		}

		// Check Chrome
		chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
		if _, err := os.Stat(chromePath); err == nil {
			foundBrowsers = append(foundBrowsers, "Google Chrome at: "+chromePath)
		}

	case "linux":
		// Check Edge
		if path, err := exec.LookPath("microsoft-edge"); err == nil {
			foundBrowsers = append(foundBrowsers, "Microsoft Edge at: "+path)
		} else if path, err := exec.LookPath("microsoft-edge-stable"); err == nil {
			foundBrowsers = append(foundBrowsers, "Microsoft Edge at: "+path)
		}

		// Check Chrome/Chromium
		if path, err := exec.LookPath("google-chrome"); err == nil {
			foundBrowsers = append(foundBrowsers, "Google Chrome at: "+path)
		} else if path, err := exec.LookPath("chromium-browser"); err == nil {
			foundBrowsers = append(foundBrowsers, "Chromium at: "+path)
		} else if path, err := exec.LookPath("chromium"); err == nil {
			foundBrowsers = append(foundBrowsers, "Chromium at: "+path)
		}
	}

	if len(foundBrowsers) > 0 {
		result := "✓ Found browser(s):\n"
		for _, browser := range foundBrowsers {
			result += "  • " + browser + "\n"
		}
		return result
	}

	// Provide installation instructions
	instructions := "⚠ No compatible browser found (Edge or Chrome/Chromium required)\n\n"
	instructions += "Installation instructions:\n\n"

	switch runtime.GOOS {
	case "windows":
		instructions += "Windows:\n\n"
		instructions += "Microsoft Edge (Recommended - pre-installed on Windows 10/11):\n"
		instructions += "  • Already installed on Windows 10/11\n"
		instructions += "  • Or download: https://www.microsoft.com/edge\n"
		instructions += "  • Or use winget: winget install Microsoft.Edge\n\n"
		instructions += "Google Chrome:\n"
		instructions += "  • Download: https://www.google.com/chrome/\n"
		instructions += "  • Or use winget: winget install Google.Chrome\n"
		instructions += "  • Or use chocolatey: choco install googlechrome\n"

	case "darwin":
		instructions += "macOS:\n\n"
		instructions += "Microsoft Edge:\n"
		instructions += "  • Download: https://www.microsoft.com/edge\n"
		instructions += "  • Or use homebrew: brew install --cask microsoft-edge\n\n"
		instructions += "Google Chrome:\n"
		instructions += "  • Download: https://www.google.com/chrome/\n"
		instructions += "  • Or use homebrew: brew install --cask google-chrome\n"

	case "linux":
		instructions += "Linux:\n\n"
		instructions += "Microsoft Edge:\n"
		instructions += "  Ubuntu/Debian:\n"
		instructions += "    curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.gpg\n"
		instructions += "    sudo install -o root -g root -m 644 microsoft.gpg /etc/apt/trusted.gpg.d/\n"
		instructions += "    sudo sh -c 'echo \"deb [arch=amd64] https://packages.microsoft.com/repos/edge stable main\" > /etc/apt/sources.list.d/microsoft-edge.list'\n"
		instructions += "    sudo apt update && sudo apt install microsoft-edge-stable\n\n"
		instructions += "Google Chrome/Chromium:\n"
		instructions += "  Ubuntu/Debian:\n"
		instructions += "    sudo apt-get update\n"
		instructions += "    sudo apt-get install chromium-browser\n"
		instructions += "  Fedora:\n"
		instructions += "    sudo dnf install chromium\n"
		instructions += "  Arch:\n"
		instructions += "    sudo pacman -S chromium\n"

	default:
		instructions += "Unknown operating system. Please install Microsoft Edge or Chrome/Chromium manually.\n"
	}

	return instructions
}

# MCP ChromeDP Server - Example Configuration

Example for Claude Desktop configuration:

## Windows
Path: %APPDATA%\Claude\claude_desktop_config.json

```json
{
  "mcpServers": {
    "chromedp-browser": {
      "command": "E:\\Workspace\\mcp_auto_test\\workspaceGo\\mcp-chromedp-server.exe",
      "args": []
    }
  }
}
```

## macOS
Path: ~/Library/Application Support/Claude/claude_desktop_config.json

```json
{
  "mcpServers": {
    "chromedp-browser": {
      "command": "/path/to/workspaceGo/mcp-chromedp-server",
      "args": []
    }
  }
}
```

## Linux
Path: ~/.config/Claude/claude_desktop_config.json

```json
{
  "mcpServers": {
    "chromedp-browser": {
      "command": "/path/to/workspaceGo/mcp-chromedp-server",
      "args": []
    }
  }
}
```

## Using with go run (development)

```json
{
  "mcpServers": {
    "chromedp-browser": {
      "command": "go",
      "args": ["run", "E:\\Workspace\\mcp_auto_test\\workspaceGo\\main.go"],
      "cwd": "E:\\Workspace\\mcp_auto_test\\workspaceGo"
    }
  }
}
```

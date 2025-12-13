# Test MCP Server
# This script tests if the MCP server responds correctly

Write-Host "Testing MCP ChromeDP Server..." -ForegroundColor Cyan

# Test 1: Initialize
Write-Host "`nTest 1: Initialize" -ForegroundColor Yellow
$initRequest = '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}'
$initResponse = $initRequest | & ".\mcp-chromedp-server.exe"
Write-Host "Response: $initResponse"

# Test 2: List tools
Write-Host "`nTest 2: List Tools" -ForegroundColor Yellow
$listRequest = '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
$listResponse = $listRequest | & ".\mcp-chromedp-server.exe"
Write-Host "Response: $listResponse"

Write-Host "`nTests completed!" -ForegroundColor Green

# MCP ChromeDP Server

MCP Server tương tự Playwright-MCP sử dụng ChromeDP (Chrome DevTools Protocol) cho browser automation. Server này được viết bằng Go, tối ưu cho deployment với không cần dependencies ngoài ngoại trừ Chrome/Chromium.

## Đặc điểm

- ✅ **Pure Go**: Không cần Selenium server hay Node.js
- ✅ **Nhẹ**: Chỉ cần Chrome/Chromium hoặc Microsoft Edge đã cài
- ✅ **Headless**: Chạy browser ẩn, tối ưu cho server
- ✅ **MCP Protocol**: Tương thích với Model Context Protocol
- ✅ **Đầy đủ chức năng**: 30 tools toàn diện
- ✅ **Multi-browser**: Tự động phát hiện Edge hoặc Chrome

## Cài đặt

### Yêu cầu

1. Go 1.21 trở lên
2. Một trong các trình duyệt sau:
   - **Microsoft Edge** (khuyến nghị cho Windows - thường đã cài sẵn)
   - Google Chrome
   - Chromium

### Build

```bash
cd workspaceGo
go mod download
go build -o mcp-chromedp-server main.go
```

### Chạy trực tiếp

```bash
go run main.go
```

## Cấu hình với Claude Desktop

Thêm vào file cấu hình Claude Desktop:

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "chromedp": {
      "command": "E:\\Workspace\\mcp_auto_test\\workspaceGo\\mcp-chromedp-server.exe",
      "args": []
    }
  }
}
```

Hoặc nếu dùng `go run`:

```json
{
  "mcpServers": {
    "chromedp": {
      "command": "go",
      "args": ["run", "E:\\Workspace\\mcp_auto_test\\workspaceGo\\main.go"]
    }
  }
}
```

## Tools có sẵn

### Navigation & Browser Control

- **playwright_navigate** / **playwright_goto** - Điều hướng đến URL
  ```json
  { "url": "https://example.com" }
  ```

- **playwright_go_back** - Quay lại trang trước
  ```json
  {}
  ```

- **playwright_go_forward** - Tiến đến trang sau
  ```json
  {}
  ```

- **playwright_reload** - Tải lại trang
  ```json
  {}
  ```

### Screenshots & PDF

- **playwright_screenshot** - Chụp ảnh màn hình
  ```json
  {
    "selector": "#element",  // optional, full page nếu bỏ qua
    "fullPage": true         // optional
  }
  ```

- **playwright_pdf** - Tạo PDF từ trang
  ```json
  {}
  ```

### Interaction

- **playwright_click** - Click vào element
  ```json
  { "selector": "button.submit" }
  ```

- **playwright_fill** - Điền vào input field
  ```json
  {
    "selector": "input[name='email']",
    "value": "user@example.com"
  }
  ```

- **playwright_type** - Gõ text vào element
  ```json
  {
    "selector": "textarea",
    "text": "Hello world"
  }
  ```

- **playwright_press** - Nhấn phím
  ```json
  {
    "selector": "input",  // optional
    "key": "Enter"
  }
  ```

- **playwright_select_option** - Chọn option trong dropdown
  ```json
  {
    "selector": "select#country",
    "value": "vietnam"
  }
  ```

### Element Operations

- **playwright_get_text** - Lấy text của element
  ```json
  { "selector": "h1" }
  ```

- **playwright_get_attribute** - Lấy attribute của element
  ```json
  {
    "selector": "a.link",
    "attribute": "href"
  }
  ```

- **playwright_wait_for_selector** - Chờ element xuất hiện
  ```json
  {
    "selector": ".loading",
    "timeout": 5000  // optional, mặc định 30000ms
  }
  ```

- **playwright_query_selector** - Kiểm tra element có tồn tại
  ```json
  { "selector": ".error-message" }
  ```

### Evaluation

- **playwright_evaluate** - Thực thi JavaScript
  ```json
  { "script": "document.title" }
  ```

- **playwright_console** - Lấy console messages
  ```json
  {}
  ```

### Advanced Browser Tools

- **browser_hover** - Hover vào element
  ```json
  { "selector": ".menu-item" }
  ```

- **browser_resize** - Thay đổi kích thước cửa sổ
  ```json
  { "width": 1920, "height": 1080 }
  ```

- **browser_drag** - Kéo thả element
  ```json
  {
    "from": "#draggable",
    "to": "#droptarget"
  }
  ```

- **browser_file_upload** - Upload file
  ```json
  {
    "selector": "input[type='file']",
    "filepath": "C:\\path\\to\\file.pdf"
  }
  ```

- **browser_fill_form** - Điền nhiều trường form cùng lúc
  ```json
  {
    "fields": [
      {"selector": "#name", "value": "John Doe"},
      {"selector": "#email", "value": "john@example.com"},
      {"selector": "#phone", "value": "123456789"}
    ]
  }
  ```

- **browser_close** - Đóng browser
  ```json
  {}
  ```

- **browser_snapshot** - Chụp accessibility tree (tốt hơn screenshot để phân tích cấu trúc)
  ```json
  { "selector": "#content" }
  ```

- **browser_handle_dialog** - Xử lý dialog (alert/confirm/prompt)
  ```json
  {
    "accept": true,
    "text": "response text"  // cho prompt
  }
  ```

- **browser_network_requests** - Lấy danh sách network requests
  ```json
  {}
  ```

- **browser_wait_for** - Chờ text xuất hiện/biến mất
  ```json
  {
    "text": "Loading complete",
    "state": "visible",
    "timeout": 5000
  }
  ```

- **browser_tabs** - Quản lý tabs
  ```json
  { "action": "list" }
  { "action": "new", "url": "https://example.com" }
  ```

- **browser_install** - Kiểm tra và hướng dẫn cài Chrome
  ```json
  {}
  ```

## Ví dụ sử dụng

### Tìm kiếm Google

```
Hãy dùng browser automation để:
1. Mở https://google.com
2. Tìm kiếm "ChromeDP golang"
3. Chụp ảnh kết quả
```

### Điền form

```
Hãy:
1. Mở https://example.com/form
2. Điền email: test@example.com
3. Điền password: secretpass
4. Click nút submit
```

### Scraping data

```
Hãy:
1. Mở https://news.ycombinator.com
2. Lấy text của tiêu đề bài đầu tiên
3. Lấy link của bài đó
```

## Ưu điểm so với Selenium

1. **Không cần WebDriver**: ChromeDP sử dụng Chrome DevTools Protocol trực tiếp
2. **Native Go**: Build thành binary duy nhất, dễ deploy
3. **Hiệu năng cao**: Ít overhead hơn Selenium
4. **Dễ cài đặt**: Chỉ cần Chrome và Go binary

## Troubleshooting

### Browser không tìm thấy

Server tự động tìm kiếm theo thứ tự: **Edge → Chrome → Chromium**

Nếu gặp lỗi không tìm thấy browser:

**Windows**: 
- Edge thường đã cài sẵn trên Windows 10/11
- Hoặc cài Chrome: https://www.google.com/chrome/
- Hoặc dùng winget: `winget install Microsoft.Edge` hoặc `winget install Google.Chrome`

**Linux**: Cài Edge hoặc Chromium
```bash
# Edge
curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.gpg
sudo install -o root -g root -m 644 microsoft.gpg /etc/apt/trusted.gpg.d/
sudo sh -c 'echo "deb [arch=amd64] https://packages.microsoft.com/repos/edge stable main" > /etc/apt/sources.list.d/microsoft-edge.list'
sudo apt update && sudo apt install microsoft-edge-stable

# Hoặc Chromium
sudo apt-get install chromium-browser
```

**macOS**: Cài Edge hoặc Chrome
```bash
# Edge
brew install --cask microsoft-edge

# Hoặc Chrome
brew install --cask google-chrome
```

### Timeout errors

Tăng timeout trong code hoặc arguments của tool:
```json
{ "timeout": 60000 }  // 60 seconds
```

## Development

### Structure

```
workspaceGo/
├── main.go        # MCP server implementation
├── go.mod         # Go dependencies
└── README.md      # Documentation
```

### Testing

```bash
# Test build
go build

# Test với input thủ công
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | go run main.go

# Test tool
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | go run main.go
```

## License

MIT

## Credits

- [ChromeDP](https://github.com/chromedp/chromedp) - Chrome DevTools Protocol for Go
- [Model Context Protocol](https://modelcontextprotocol.io/) - MCP Specification

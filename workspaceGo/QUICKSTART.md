# MCP ChromeDP Server - Quick Start Guide

## 🚀 Cài đặt nhanh

### Bước 1: Build server
```bash
cd E:\Workspace\mcp_auto_test\workspaceGo
go build -o mcp-chromedp-server.exe main.go
```

### Bước 2: Cấu hình Claude Desktop

Mở file: `%APPDATA%\Claude\claude_desktop_config.json`

Thêm:
```json
{
  "mcpServers": {
    "browser": {
      "command": "E:\\Workspace\\mcp_auto_test\\workspaceGo\\mcp-chromedp-server.exe"
    }
  }
}
```

### Bước 3: Restart Claude Desktop

Sau khi restart, bạn sẽ thấy 18 tools mới có sẵn!

## 🎯 Ví dụ sử dụng

### Ví dụ 1: Tìm kiếm Google
```
Prompt cho Claude:

"Hãy mở Google và tìm kiếm 'ChromeDP golang', sau đó chụp ảnh kết quả"
```

Claude sẽ sử dụng:
1. `playwright_navigate` - Mở google.com
2. `playwright_fill` - Điền từ khóa vào search box
3. `playwright_press` - Nhấn Enter
4. `playwright_screenshot` - Chụp ảnh

### Ví dụ 2: Scrape dữ liệu
```
Prompt cho Claude:

"Hãy vào Hacker News và lấy 5 tiêu đề hàng đầu"
```

Claude sẽ sử dụng:
1. `playwright_navigate` - Mở news.ycombinator.com
2. `playwright_query_selector` - Tìm elements
3. `playwright_get_text` - Lấy text từ mỗi tiêu đề

### Ví dụ 3: Điền form
```
Prompt cho Claude:

"Vào example.com/login, điền username 'admin' và password 'secret123', rồi submit"
```

Claude sẽ sử dụng:
1. `playwright_navigate` - Mở trang login
2. `playwright_fill` - Điền username
3. `playwright_fill` - Điền password  
4. `playwright_click` - Click nút submit

### Ví dụ 4: Kiểm tra website
```
Prompt cho Claude:

"Kiểm tra xem example.com có hiển thị thông báo lỗi không"
```

Claude sẽ sử dụng:
1. `playwright_navigate` - Mở website
2. `playwright_query_selector` - Tìm error message
3. `playwright_get_text` - Đọc nội dung lỗi (nếu có)

## 📋 Danh sách đầy đủ 30 tools

### Navigation (5 tools)
- ✅ playwright_navigate
- ✅ playwright_goto  
- ✅ playwright_go_back
- ✅ playwright_go_forward
- ✅ playwright_reload

### Screenshots & PDF (2 tools)
- ✅ playwright_screenshot
- ✅ playwright_pdf

### Interaction (5 tools)
- ✅ playwright_click
- ✅ playwright_fill
- ✅ playwright_type
- ✅ playwright_press
- ✅ playwright_select_option

### Element Operations (4 tools)
- ✅ playwright_get_text
- ✅ playwright_get_attribute
- ✅ playwright_wait_for_selector
- ✅ playwright_query_selector

### Evaluation (2 tools)
- ✅ playwright_evaluate
- ✅ playwright_console

### Advanced Browser Tools (12 tools)
- ✅ browser_hover - Hover vào element
- ✅ browser_resize - Thay đổi kích thước viewport
- ✅ browser_drag - Kéo thả element
- ✅ browser_file_upload - Upload file
- ✅ browser_fill_form - Điền nhiều trường form
- ✅ browser_close - Đóng browser
- ✅ browser_snapshot - Accessibility tree snapshot
- ✅ browser_handle_dialog - Xử lý alert/confirm/prompt
- ✅ browser_network_requests - Lấy network requests
- ✅ browser_wait_for - Chờ text xuất hiện
- ✅ browser_tabs - Quản lý tabs
- ✅ browser_install - Kiểm tra cài đặt Chrome

## 🔧 Troubleshooting

### Lỗi "Browser not found"
Server tự động tìm Edge hoặc Chrome. Trên Windows 10/11, Edge đã cài sẵn.
Nếu chưa có: 
- Edge: https://www.microsoft.com/edge
- Chrome: https://www.google.com/chrome/

### Lỗi "timeout"
Element có thể load chậm, thử tăng timeout:
```json
{ "selector": ".slow-element", "timeout": 60000 }
```

### Server không response
1. Kiểm tra file exe đã build chưa
2. Kiểm tra đường dẫn trong config
3. Xem logs tại stderr của Claude Desktop

## 💡 Tips

1. **Selector CSS**: Dùng DevTools (F12) để tìm selector chính xác
2. **Wait before interact**: Dùng `wait_for_selector` trước khi click/fill
3. **Screenshots for debugging**: Chụp ảnh để debug khi automation không hoạt động
4. **Evaluate for complex tasks**: Dùng JavaScript evaluation cho logic phức tạp

## 🎓 Chi tiết kỹ thuật

### Tại sao dùng ChromeDP thay vì Selenium?

| Feature | ChromeDP | Selenium |
|---------|----------|----------|
| Dependencies | Chỉ cần Chrome | Cần WebDriver + Server |
| Language | Pure Go | Java/Python/Node wrapper |
| Binary size | ~14MB | Phụ thuộc runtime |
| Performance | Native CDP | HTTP overhead |
| Deploy | Copy 1 file exe | Nhiều dependencies |

### Architecture

```
┌─────────────────┐
│ Claude Desktop  │
└────────┬────────┘
         │ stdio (JSON-RPC)
         │
┌────────▼────────┐
│  MCP Server     │
│  (Go)           │
└────────┬────────┘
         │ Chrome DevTools Protocol
         │
┌────────▼────────┐
│  Chrome/        │
│  Chromium       │
└─────────────────┘
```

### Protocol Flow

1. Claude gửi JSON-RPC qua stdin
2. Server parse request và gọi ChromeDP
3. ChromeDP điều khiển Chrome qua CDP
4. Kết quả trả về qua stdout

## 📚 Tài nguyên

- [ChromeDP Documentation](https://github.com/chromedp/chromedp)
- [MCP Specification](https://modelcontextprotocol.io/)
- [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/)

## 🤝 Support

Nếu gặp vấn đề, kiểm tra:
1. Chrome đã cài chưa
2. Go version >= 1.21
3. Build thành công chưa
4. Config path đúng chưa

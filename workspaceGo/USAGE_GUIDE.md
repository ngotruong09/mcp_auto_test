# Hướng dẫn Sử dụng MCP ChromeDP Server

## Mục lục
- [Giới thiệu](#giới-thiệu)
- [Cài đặt và Cấu hình](#cài-đặt-và-cấu-hình)
- [Các Tools có sẵn](#các-tools-có-sẵn)
- [Hướng dẫn chi tiết từng Tool](#hướng-dẫn-chi-tiết-từng-tool)
- [10+ Ví dụ Sample Prompts](#10-ví-dụ-sample-prompts)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Giới thiệu

MCP ChromeDP Server là một browser automation server sử dụng Chrome DevTools Protocol, cho phép bạn điều khiển trình duyệt thông qua Claude Desktop với 30 tools toàn diện.

### Tại sao sử dụng MCP ChromeDP Server?

- ✅ Tự động hóa các tác vụ web browsing
- ✅ Scrape dữ liệu từ websites
- ✅ Test websites tự động
- ✅ Chụp screenshots và tạo PDF
- ✅ Điền forms tự động
- ✅ Tương tác với web apps phức tạp

## Cài đặt và Cấu hình

### Bước 1: Build Server

```bash
cd E:\Workspace\mcp_auto_test\workspaceGo
go build -o mcp-chromedp-server.exe main.go
```

### Bước 2: Cấu hình Claude Desktop

Mở file: `%APPDATA%\Claude\claude_desktop_config.json`

Thêm cấu hình:
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

Sau khi restart, bạn sẽ thấy biểu tượng 🔌 với 30 tools available.

## Các Tools có sẵn

### 📍 Navigation Tools (5)
| Tool | Mô tả |
|------|-------|
| `playwright_navigate` | Điều hướng đến URL |
| `playwright_goto` | Alias của navigate |
| `playwright_go_back` | Quay lại trang trước |
| `playwright_go_forward` | Tiến tới trang sau |
| `playwright_reload` | Reload trang hiện tại |

### 📸 Capture Tools (2)
| Tool | Mô tả |
|------|-------|
| `playwright_screenshot` | Chụp ảnh màn hình |
| `playwright_pdf` | Tạo PDF từ page |

### 🖱️ Interaction Tools (5)
| Tool | Mô tả |
|------|-------|
| `playwright_click` | Click vào element |
| `playwright_fill` | Điền giá trị vào input |
| `playwright_type` | Gõ text vào element |
| `playwright_press` | Nhấn phím keyboard |
| `playwright_select_option` | Chọn option trong dropdown |

### 🔍 Element Query Tools (4)
| Tool | Mô tả |
|------|-------|
| `playwright_get_text` | Lấy text content |
| `playwright_get_attribute` | Lấy giá trị attribute |
| `playwright_wait_for_selector` | Chờ element xuất hiện |
| `playwright_query_selector` | Kiểm tra element tồn tại |

### ⚙️ JavaScript & Console (2)
| Tool | Mô tả |
|------|-------|
| `playwright_evaluate` | Thực thi JavaScript |
| `playwright_console` | Lấy console messages |

### 🚀 Advanced Browser Tools (12)
| Tool | Mô tả |
|------|-------|
| `browser_hover` | Hover vào element |
| `browser_resize` | Thay đổi viewport size |
| `browser_drag` | Drag & drop element |
| `browser_file_upload` | Upload file |
| `browser_fill_form` | Điền nhiều fields cùng lúc |
| `browser_close` | Đóng browser |
| `browser_snapshot` | Accessibility tree snapshot |
| `browser_handle_dialog` | Xử lý alert/confirm/prompt |
| `browser_network_requests` | Lấy network requests |
| `browser_wait_for` | Chờ text xuất hiện |
| `browser_tabs` | Quản lý tabs |
| `browser_install` | Kiểm tra Chrome installation |

## Hướng dẫn chi tiết từng Tool

### Navigation Tools

#### 1. playwright_navigate / playwright_goto
**Công dụng:** Mở một URL trong browser

**Parameters:**
- `url` (string, required): URL cần mở

**Example:**
```
Claude, hãy mở trang https://github.com
```

#### 2. playwright_go_back
**Công dụng:** Quay lại trang trước trong history

**Example:**
```
Hãy quay lại trang trước
```

#### 3. playwright_reload
**Công dụng:** Reload trang hiện tại

**Example:**
```
Reload trang này để lấy dữ liệu mới nhất
```

### Capture Tools

#### 4. playwright_screenshot
**Công dụng:** Chụp ảnh màn hình

**Parameters:**
- `selector` (string, optional): CSS selector của element cần chụp
- `fullPage` (boolean, optional): Chụp toàn bộ trang có scroll

**Example:**
```
Chụp ảnh toàn bộ trang web này
Chụp ảnh phần header của trang
```

#### 5. playwright_pdf
**Công dụng:** Xuất trang web sang PDF

**Example:**
```
Tạo file PDF từ trang này
```

### Interaction Tools

#### 6. playwright_click
**Công dụng:** Click vào một element

**Parameters:**
- `selector` (string, required): CSS selector

**Example:**
```
Click vào nút "Sign In"
Click vào link đầu tiên trong danh sách
```

#### 7. playwright_fill
**Công dụng:** Clear và điền giá trị mới vào input field

**Parameters:**
- `selector` (string, required): CSS selector
- `value` (string, required): Giá trị cần điền

**Example:**
```
Điền "john@example.com" vào ô email
```

#### 8. playwright_type
**Công dụng:** Gõ text vào element (không clear trước)

**Parameters:**
- `selector` (string, required): CSS selector
- `text` (string, required): Text cần gõ

**Example:**
```
Gõ "Hello World" vào textarea
```

#### 9. playwright_press
**Công dụng:** Nhấn một phím trên keyboard

**Parameters:**
- `key` (string, required): Tên phím (Enter, Escape, Tab, etc.)
- `selector` (string, optional): Element để focus trước khi nhấn

**Example:**
```
Nhấn Enter để submit form
Nhấn Escape để đóng dialog
```

#### 10. playwright_select_option
**Công dụng:** Chọn option trong dropdown

**Parameters:**
- `selector` (string, required): CSS selector của select element
- `value` (string, required): Value của option

**Example:**
```
Chọn "Vietnam" trong dropdown quốc gia
```

### Element Query Tools

#### 11. playwright_get_text
**Công dụng:** Lấy text content của element

**Parameters:**
- `selector` (string, required): CSS selector

**Example:**
```
Lấy text của tiêu đề chính
Cho tôi biết giá sản phẩm hiển thị trên trang
```

#### 12. playwright_get_attribute
**Công dụng:** Lấy giá trị attribute của element

**Parameters:**
- `selector` (string, required): CSS selector
- `attribute` (string, required): Tên attribute

**Example:**
```
Lấy href của link đầu tiên
Lấy src của hình ảnh chính
```

#### 13. playwright_wait_for_selector
**Công dụng:** Chờ element xuất hiện (visible)

**Parameters:**
- `selector` (string, required): CSS selector
- `timeout` (number, optional): Timeout ms (default: 30000)

**Example:**
```
Chờ cho loading spinner biến mất
Chờ kết quả tìm kiếm xuất hiện
```

#### 14. playwright_query_selector
**Công dụng:** Kiểm tra element có tồn tại không

**Parameters:**
- `selector` (string, required): CSS selector

**Example:**
```
Kiểm tra xem có thông báo lỗi không
```

### JavaScript & Console

#### 15. playwright_evaluate
**Công dụng:** Thực thi JavaScript code trong page context

**Parameters:**
- `script` (string, required): JavaScript code

**Example:**
```
Chạy JavaScript để lấy document.title
Đếm số lượng links trên trang bằng JavaScript
```

#### 16. playwright_console
**Công dụng:** Lấy console messages từ page

**Example:**
```
Cho tôi xem console messages
```

### Advanced Browser Tools

#### 17. browser_hover
**Công dụng:** Hover chuột vào element (hiển thị dropdown, tooltip)

**Parameters:**
- `selector` (string, required): CSS selector

**Example:**
```
Hover vào menu "Products" để xem dropdown
```

#### 18. browser_resize
**Công dụng:** Thay đổi kích thước viewport

**Parameters:**
- `width` (number, required): Chiều rộng pixels
- `height` (number, required): Chiều cao pixels

**Example:**
```
Resize cửa sổ về kích thước mobile 375x667
Test trang ở kích thước tablet 768x1024
```

#### 19. browser_drag
**Công dụng:** Kéo element từ vị trí này sang vị trí khác

**Parameters:**
- `from` (string, required): CSS selector element nguồn
- `to` (string, required): CSS selector element đích

**Example:**
```
Kéo task từ "To Do" sang "In Progress"
```

#### 20. browser_file_upload
**Công dụng:** Upload file vào input[type=file]

**Parameters:**
- `selector` (string, required): CSS selector của file input
- `filepath` (string, required): Đường dẫn file

**Example:**
```
Upload file C:\Documents\resume.pdf vào form
```

#### 21. browser_fill_form
**Công dụng:** Điền nhiều fields cùng lúc (nhanh hơn fill từng field)

**Parameters:**
- `fields` (array, required): Mảng các {selector, value}

**Example:**
```
Điền form với: name=John Doe, email=john@test.com, phone=123456789
```

#### 22. browser_close
**Công dụng:** Đóng browser/page hiện tại

**Example:**
```
Đóng browser
```

#### 23. browser_snapshot
**Công dụng:** Lấy accessibility tree của page (phân tích cấu trúc)

**Parameters:**
- `selector` (string, optional): CSS selector (optional, whole page nếu bỏ qua)

**Example:**
```
Phân tích accessibility structure của trang
Lấy snapshot của phần main content
```

#### 24. browser_handle_dialog
**Công dụng:** Xử lý JavaScript dialogs (alert, confirm, prompt)

**Parameters:**
- `accept` (boolean, required): true = accept, false = dismiss
- `text` (string, optional): Text cho prompt dialog

**Example:**
```
Accept alert dialog
Dismiss confirmation
```

#### 25. browser_network_requests
**Công dụng:** Lấy danh sách network requests đã thực hiện

**Example:**
```
Cho tôi xem các API requests trang này đã gọi
```

#### 26. browser_wait_for
**Công dụng:** Chờ text xuất hiện/biến mất hoặc chờ một khoảng thời gian

**Parameters:**
- `text` (string, optional): Text cần chờ
- `state` (string, optional): "visible" hoặc "hidden"
- `timeout` (number, optional): Timeout ms

**Example:**
```
Chờ text "Loading complete" xuất hiện
Chờ 3 giây
```

#### 27. browser_tabs
**Công dụng:** Quản lý tabs

**Parameters:**
- `action` (string, required): "list", "new", "close", "select"
- `url` (string, optional): URL cho action "new"
- `index` (number, optional): Index cho "close"/"select"

**Example:**
```
List tất cả tabs đang mở
Mở tab mới với URL https://example.com
```

#### 28. browser_install
**Công dụng:** Kiểm tra Chrome đã cài chưa và hiển thị hướng dẫn

**Example:**
```
Kiểm tra Chrome đã cài đặt chưa
```

## 10+ Ví dụ Sample Prompts

### 1. 🔍 Tìm kiếm Google và Phân tích Kết quả

**Prompt:**
```
Hãy mở Google, tìm kiếm "ChromeDP golang tutorial", 
chờ kết quả xuất hiện, sau đó lấy text của 5 kết quả đầu tiên 
và chụp ảnh màn hình.
```

**Tools sử dụng:**
- `playwright_navigate` → Mở google.com
- `playwright_fill` → Điền từ khóa vào search box
- `playwright_press` → Nhấn Enter
- `playwright_wait_for_selector` → Chờ kết quả
- `playwright_get_text` → Lấy text từng kết quả (x5)
- `playwright_screenshot` → Chụp ảnh

---

### 2. 📝 Điền Form Đăng ký Tự động

**Prompt:**
```
Vào trang https://example.com/register và điền form đăng ký với:
- Full Name: John Doe
- Email: johndoe@example.com
- Phone: +84 123 456 789
- Password: SecurePass123
Sau đó click nút Submit và chụp ảnh xác nhận.
```

**Tools sử dụng:**
- `playwright_navigate` → Mở trang register
- `browser_fill_form` → Điền tất cả fields cùng lúc
- `playwright_click` → Click Submit
- `playwright_screenshot` → Chụp ảnh confirmation

---

### 3. 🛒 Scrape Giá Sản phẩm E-commerce

**Prompt:**
```
Mở trang https://example.com/products, lấy tên và giá của 
10 sản phẩm đầu tiên, sau đó tạo một bảng tổng hợp.
```

**Tools sử dụng:**
- `playwright_navigate` → Mở trang products
- `playwright_wait_for_selector` → Chờ products load
- `playwright_get_text` → Lấy tên sản phẩm (x10)
- `playwright_get_text` → Lấy giá (x10)

**Output:** Claude sẽ tổng hợp thành bảng markdown

---

### 4. 📊 Test Responsive Design

**Prompt:**
```
Mở trang https://example.com, test responsive bằng cách:
1. Xem ở desktop (1920x1080) - chụp ảnh
2. Resize về tablet (768x1024) - chụp ảnh
3. Resize về mobile (375x667) - chụp ảnh
So sánh và cho biết layout có thay đổi đúng không.
```

**Tools sử dụng:**
- `playwright_navigate`
- `browser_resize` → 1920x1080
- `playwright_screenshot`
- `browser_resize` → 768x1024
- `playwright_screenshot`
- `browser_resize` → 375x667
- `playwright_screenshot`

---

### 5. 📰 Scrape Tin tức Hacker News

**Prompt:**
```
Vào Hacker News (news.ycombinator.com), lấy:
- 10 tiêu đề hàng đầu
- Số điểm của mỗi bài
- Link của mỗi bài
Trình bày dưới dạng bảng.
```

**Tools sử dụng:**
- `playwright_navigate`
- `playwright_get_text` → Lấy tiêu đề (x10)
- `playwright_get_text` → Lấy điểm (x10)
- `playwright_get_attribute` → Lấy href (x10)

---

### 6. 🎯 Automation Test Login Flow

**Prompt:**
```
Test login flow của https://example.com/login:
1. Điền username: testuser
2. Điền password: testpass123
3. Click nút Login
4. Kiểm tra xem có chuyển đến dashboard không
5. Kiểm tra có error message không
6. Chụp ảnh kết quả
```

**Tools sử dụng:**
- `playwright_navigate`
- `playwright_fill` → Username
- `playwright_fill` → Password
- `playwright_click` → Login button
- `playwright_wait_for_selector` → Chờ dashboard hoặc error
- `playwright_query_selector` → Check error message
- `playwright_screenshot`

---

### 7. 📄 Tạo PDF Báo cáo

**Prompt:**
```
Mở trang báo cáo tài chính https://example.com/annual-report,
chờ trang load xong, sau đó tạo file PDF và chụp ảnh preview.
```

**Tools sử dụng:**
- `playwright_navigate`
- `playwright_wait_for_selector` → Chờ content load
- `playwright_pdf` → Tạo PDF
- `playwright_screenshot` → Preview

---

### 8. 🎨 Test Menu Dropdown

**Prompt:**
```
Vào trang https://example.com, hover vào menu "Products" 
để xem dropdown menu, lấy text của tất cả items trong dropdown,
sau đó chụp ảnh dropdown đang mở.
```

**Tools sử dụng:**
- `playwright_navigate`
- `browser_hover` → Hover vào Products menu
- `playwright_wait_for_selector` → Chờ dropdown hiện
- `playwright_get_text` → Lấy text các items
- `playwright_screenshot`

---

### 9. 🔐 Upload File và Submit

**Prompt:**
```
Vào trang https://example.com/upload:
1. Upload file C:\Documents\report.pdf
2. Điền description: "Monthly Report December 2025"
3. Chọn category "Financial" từ dropdown
4. Click Submit
5. Chờ success message và chụp ảnh
```

**Tools sử dụng:**
- `playwright_navigate`
- `browser_file_upload`
- `playwright_fill` → Description
- `playwright_select_option` → Category
- `playwright_click` → Submit
- `browser_wait_for` → Chờ success message
- `playwright_screenshot`

---

### 10. 🌐 Kiểm tra Network Requests

**Prompt:**
```
Mở trang https://example.com/dashboard, lấy danh sách 
tất cả API requests mà trang đã gọi, và cho biết:
- Số lượng requests
- Các endpoints được gọi
- HTTP methods sử dụng
```

**Tools sử dụng:**
- `playwright_navigate`
- `playwright_wait_for_selector` → Chờ page load
- `browser_network_requests` → Lấy request list

---

### 11. 🎭 Test Drag & Drop

**Prompt:**
```
Vào trang Kanban board tại https://example.com/board,
kéo task "Complete documentation" từ cột "To Do" 
sang cột "In Progress", sau đó verify task đã di chuyển.
```

**Tools sử dụng:**
- `playwright_navigate`
- `browser_drag` → Kéo task
- `playwright_query_selector` → Verify vị trí mới
- `playwright_screenshot`

---

### 12. 📱 Accessibility Audit

**Prompt:**
```
Phân tích accessibility của trang https://example.com:
1. Lấy accessibility snapshot
2. Kiểm tra heading structure
3. Kiểm tra có alt text cho images không
4. Tạo report
```

**Tools sử dụng:**
- `playwright_navigate`
- `browser_snapshot` → Lấy a11y tree
- `playwright_evaluate` → Check headings và images
- Analysis by Claude

---

### 13. 🔔 Handle Alert Dialog

**Prompt:**
```
Vào trang https://example.com/demo:
1. Click nút "Delete Item"
2. Sẽ có confirm dialog, hãy accept
3. Kiểm tra item đã bị xóa chưa
```

**Tools sử dụng:**
- `playwright_navigate`
- `playwright_click` → Delete button
- `browser_handle_dialog` → Accept confirmation
- `playwright_query_selector` → Verify deletion

---

### 14. ⏱️ Performance Test với Wait

**Prompt:**
```
Test loading speed của https://example.com:
1. Mở trang
2. Đo thời gian từ lúc load đến khi text "Ready" xuất hiện
3. Chụp ảnh khi đã ready
```

**Tools sử dụng:**
- `playwright_navigate`
- `browser_wait_for` → Chờ "Ready" text (ghi nhận timeout)
- `playwright_screenshot`

---

### 15. 🔄 Multi-tab Management

**Prompt:**
```
Mở 3 tabs với các URL:
- https://github.com
- https://stackoverflow.com
- https://news.ycombinator.com
List tất cả tabs, sau đó lấy title của mỗi tab.
```

**Tools sử dụng:**
- `browser_tabs` → new (x3)
- `browser_tabs` → list
- `playwright_evaluate` → document.title (mỗi tab)

---

## Best Practices

### 1. Always Wait Before Interact
```
❌ Bad: Click ngay sau khi navigate
✅ Good: Navigate → Wait for selector → Click
```

### 2. Use Specific Selectors
```
❌ Bad: selector = "button"
✅ Good: selector = "button[type='submit']#login-btn"
```

### 3. Handle Timeouts
```
Khi page load chậm, tăng timeout:
"Chờ kết quả xuất hiện trong 60 giây"
→ Claude sẽ dùng timeout: 60000
```

### 4. Batch Form Fills
```
❌ Bad: Fill từng field riêng lẻ (5 tools calls)
✅ Good: Dùng browser_fill_form (1 tool call)
```

### 5. Screenshot for Verification
```
Luôn chụp ảnh sau các action quan trọng để verify
```

### 6. Use Accessibility Snapshot for Structure Analysis
```
❌ Bad: Screenshot → OCR text
✅ Good: browser_snapshot → Structured data
```

### 7. Combine JavaScript Evaluation
```
Thay vì query nhiều lần, dùng evaluate để chạy JS phức tạp
```

## Troubleshooting

### Lỗi thường gặp

#### 1. "Element not found"
**Nguyên nhân:** Selector sai hoặc element chưa load
**Giải pháp:** 
- Dùng `playwright_wait_for_selector` trước
- Kiểm tra selector bằng DevTools (F12)

#### 2. "Timeout waiting for selector"
**Nguyên nhân:** Element không xuất hiện trong thời gian chờ
**Giải pháp:**
- Tăng timeout: "Chờ 60 giây"
- Kiểm tra selector có đúng không
- Kiểm tra element có bị ẩn bởi CSS không

#### 3. "Cannot click element"
**Nguyên nhân:** Element bị che hoặc không visible
**Giải pháp:**
- Scroll đến element trước
- Đợi overlay/loading biến mất
- Dùng JavaScript click: `playwright_evaluate`

#### 4. "Screenshot is blank"
**Nguyên nhân:** Page chưa render xong
**Giải pháp:**
- Chờ một chút: `browser_wait_for` với timeout
- Wait for key element xuất hiện

#### 5. "File upload failed"
**Nguyên nhân:** Đường dẫn file sai
**Giải pháp:**
- Dùng absolute path
- Kiểm tra file tồn tại
- Windows: dùng double backslash `C:\\path\\file.pdf`

### Debug Tips

1. **Chụp ảnh mỗi bước** để xem browser đang ở trạng thái nào
2. **Lấy console messages** để xem có errors không
3. **Check network requests** để verify API calls
4. **Dùng accessibility snapshot** để phân tích DOM structure
5. **Try evaluate JavaScript** để test selectors

## Tổng kết

MCP ChromeDP Server cung cấp 30 tools mạnh mẽ cho browser automation. Với các sample prompts trên, bạn có thể:

- ✅ Tự động hóa các tác vụ web browsing
- ✅ Scrape dữ liệu hiệu quả
- ✅ Test websites và web apps
- ✅ Tương tác với forms phức tạp
- ✅ Phân tích cấu trúc và performance

**Tips cuối:** Hãy thử kết hợp nhiều tools trong một prompt để Claude tự động hóa workflow phức tạp!

---

📚 **Tài liệu khác:**
- [README.md](README.md) - Giới thiệu và cài đặt
- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
- [TOOLS_COMPARISON.md](TOOLS_COMPARISON.md) - So sánh với Playwright
- [UPDATE_SUMMARY.md](UPDATE_SUMMARY.md) - Tổng kết cập nhật

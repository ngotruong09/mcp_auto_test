# 30 Sample Prompts cho MCP ChromeDP Server

Tổng hợp 30 prompts mẫu từ cơ bản đến nâng cao để sử dụng MCP ChromeDP Server qua Claude Desktop.

## 📚 Mục lục
- [Cơ bản (1-10)](#cơ-bản)
- [Trung cấp (11-20)](#trung-cấp)
- [Nâng cao (21-30)](#nâng-cao)

---

## Cơ bản

### 1. Mở Website Đơn giản
```
Hãy mở trang https://example.com và chụp ảnh màn hình.
```
**Độ phức tạp:** ⭐  
**Tools:** navigate, screenshot

---

### 2. Tìm kiếm Google
```
Mở Google, tìm kiếm "weather in Hanoi", và cho tôi xem kết quả đầu tiên.
```
**Độ phức tạp:** ⭐⭐  
**Tools:** navigate, fill, press, wait_for_selector, get_text

---

### 3. Lấy Thông tin Văn bản
```
Vào trang https://news.ycombinator.com và lấy tiêu đề của 5 bài viết hàng đầu.
```
**Độ phức tạp:** ⭐⭐  
**Tools:** navigate, wait_for_selector, get_text (x5)

---

### 4. Kiểm tra Element Tồn tại
```
Mở trang https://example.com/login và kiểm tra xem có nút "Forgot Password" không.
```
**Độ phức tạp:** ⭐  
**Tools:** navigate, query_selector

---

### 5. Lấy Link từ Page
```
Vào Wikipedia trang Python programming, lấy href của 10 links đầu tiên trong content.
```
**Độ phức tạp:** ⭐⭐  
**Tools:** navigate, wait_for_selector, get_attribute (x10)

---

### 6. Điền Form Đơn giản
```
Mở https://example.com/contact, điền email là "test@example.com" 
và message là "Hello, this is a test", sau đó chụp ảnh.
```
**Độ phức tạp:** ⭐⭐  
**Tools:** navigate, fill (x2), screenshot

---

### 7. Click và Navigate
```
Vào GitHub homepage, click vào link "Explore", chờ trang mới load, 
rồi chụp ảnh.
```
**Độ phức tạp:** ⭐⭐  
**Tools:** navigate, click, wait_for_selector, screenshot

---

### 8. Reload và Compare
```
Mở https://example.com, chụp ảnh, reload trang, chụp ảnh lần 2, 
cho biết có gì thay đổi không.
```
**Độ phức tạp:** ⭐⭐  
**Tools:** navigate, screenshot, reload, screenshot

---

### 9. Nhấn Phím đặc biệt
```
Vào Google, focus vào search box, nhấn phím "Escape" để clear suggestion dropdown.
```
**Độ phức tạp:** ⭐⭐  
**Tools:** navigate, click, press

---

### 10. Kiểm tra Chrome Installation
```
Kiểm tra xem Chrome đã được cài đặt trên máy này chưa.
```
**Độ phức tạp:** ⭐  
**Tools:** browser_install

---

## Trung cấp

### 11. Scrape E-commerce Products
```
Vào trang https://example-shop.com/products, lấy tên, giá, và rating 
của 10 sản phẩm đầu tiên, tổng hợp thành bảng.
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, wait_for_selector, get_text (x30), evaluation

---

### 12. Test Login Flow
```
Test login flow tại https://example.com/login:
1. Điền username: admin
2. Điền password: admin123
3. Click Login
4. Kiểm tra có redirect về dashboard không
5. Nếu có error, chụp ảnh error message
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, fill (x2), click, wait_for_selector, query_selector, screenshot

---

### 13. Multi-step Form
```
Điền form đăng ký tại https://example.com/register với:
- First Name: John
- Last Name: Doe  
- Email: john.doe@example.com
- Password: SecurePass123
- Confirm Password: SecurePass123
- Accept Terms checkbox: checked
Sau đó submit và chụp ảnh confirmation.
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, browser_fill_form, click, wait_for_selector, screenshot

---

### 14. Dropdown Selection
```
Vào https://example.com/booking:
1. Chọn "Economy" từ dropdown Class
2. Chọn "Vietnam" từ dropdown Country
3. Chọn "Hanoi" từ dropdown City
4. Click Search và chụp ảnh kết quả
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, select_option (x3), click, wait_for_selector, screenshot

---

### 15. Hover Menu Navigation
```
Vào trang https://example.com, hover vào menu "Products", 
từ dropdown chọn "Electronics", chờ trang mới load, 
lấy danh sách 5 categories con.
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, browser_hover, wait_for_selector, click, get_text (x5)

---

### 16. Responsive Testing
```
Test responsive design của https://example.com:
- Desktop 1920x1080: chụp ảnh
- Tablet 768x1024: chụp ảnh  
- Mobile 375x667: chụp ảnh
So sánh layout và cho nhận xét.
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, browser_resize (x3), screenshot (x3)

---

### 17. JavaScript Evaluation
```
Vào https://example.com và dùng JavaScript để:
1. Đếm số lượng links trên trang
2. Đếm số lượng images
3. Lấy title của page
4. Lấy meta description
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, evaluate (x4)

---

### 18. Wait for Dynamic Content
```
Vào trang https://example.com/search, tìm kiếm "laptop", 
chờ cho đến khi text "Found X results" xuất hiện, 
sau đó lấy số lượng kết quả.
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, fill, press, browser_wait_for, get_text

---

### 19. PDF Generation
```
Tạo PDF từ trang https://example.com/annual-report-2025:
1. Chờ tất cả charts load xong
2. Tạo PDF
3. Chụp preview trang đầu
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, wait_for_selector, playwright_pdf, screenshot

---

### 20. Network Monitoring
```
Mở trang https://example.com/dashboard, chờ page load xong,
lấy danh sách tất cả API endpoints mà trang đã gọi,
phân loại theo HTTP method (GET, POST, etc.).
```
**Độ phức tạp:** ⭐⭐⭐  
**Tools:** navigate, wait_for_selector, browser_network_requests

---

## Nâng cao

### 21. Complex Scraping Workflow
```
Scrape dữ liệu từ https://news.ycombinator.com:
1. Lấy 30 bài viết đầu tiên (title, points, comments count)
2. Click vào bài có điểm cao nhất
3. Lấy nội dung và top 5 comments
4. Quay lại trang chủ
5. Tổng hợp data thành report
```
**Độ phức tạp:** ⭐⭐⭐⭐  
**Tools:** navigate, get_text (x90+), click, wait, go_back, evaluation

---

### 22. E-commerce Comparison
```
So sánh giá sản phẩm "iPhone 15 Pro" trên 3 trang:
- https://shop1.example.com
- https://shop2.example.com  
- https://shop3.example.com
Tạo bảng so sánh giá, rating, và availability.
```
**Độ phức tạp:** ⭐⭐⭐⭐  
**Tools:** navigate (x3), fill (x3), click (x3), get_text (x9), comparison

---

### 23. File Upload and Processing
```
Vào https://example.com/converter:
1. Upload file C:\Documents\data.csv
2. Chọn output format "JSON"
3. Click Convert
4. Chờ processing complete (có progress bar)
5. Download result
```
**Độ phức tạp:** ⭐⭐⭐⭐  
**Tools:** navigate, browser_file_upload, select_option, click, browser_wait_for

---

### 24. Drag & Drop Kanban
```
Vào Kanban board tại https://example.com/board:
1. Lấy danh sách tasks trong "To Do"
2. Kéo 3 tasks đầu tiên sang "In Progress"
3. Chụp ảnh board sau khi di chuyển
4. Verify tasks đã ở đúng cột
```
**Độ phức tạp:** ⭐⭐⭐⭐  
**Tools:** navigate, get_text, browser_drag (x3), screenshot, query_selector

---

### 25. Multi-tab Research
```
Research về "Claude AI":
- Tab 1: Google search "Claude AI", lấy 5 kết quả đầu
- Tab 2: Mở Wikipedia page, lấy summary
- Tab 3: Mở official website, lấy features
- Tab 4: Mở pricing page, lấy plans
Tổng hợp thành research report.
```
**Độ phức tạp:** ⭐⭐⭐⭐⭐  
**Tools:** browser_tabs (new x4), navigate, fill, get_text, compilation

---

### 26. Accessibility Audit
```
Audit accessibility của https://example.com:
1. Lấy accessibility snapshot
2. Kiểm tra heading hierarchy (h1, h2, h3)
3. Kiểm tra alt text cho images
4. Kiểm tra form labels
5. Kiểm tra color contrast (qua evaluate)
6. Tạo accessibility report với scoring
```
**Độ phức tạp:** ⭐⭐⭐⭐⭐  
**Tools:** navigate, browser_snapshot, evaluate (x5), analysis

---

### 27. Performance Testing
```
Test performance của https://example.com:
1. Measure load time (navigate đến ready state)
2. Count số lượng network requests
3. Tính tổng size của resources
4. Identify slow requests (>1s)
5. Screenshot của Network waterfall (qua DevTools)
```
**Độ phức tạp:** ⭐⭐⭐⭐⭐  
**Tools:** navigate, browser_network_requests, evaluate, timing analysis

---

### 28. Automated Testing Suite
```
Run test suite cho https://example.com/app:

Test 1: Login với valid credentials → Success
Test 2: Login với invalid password → Error message
Test 3: Create new item → Item appears in list
Test 4: Edit item → Changes saved
Test 5: Delete item → Item removed
Test 6: Logout → Redirect to login

Tạo test report với pass/fail cho từng test.
```
**Độ phức tạp:** ⭐⭐⭐⭐⭐  
**Tools:** navigate, fill, click, wait, query_selector, screenshot (x6+)

---

### 29. Dialog Handling Workflow
```
Test dialog handling tại https://example.com/demo:
1. Click "Show Alert" → Accept alert
2. Click "Show Confirm" → Dismiss confirm
3. Click "Show Prompt" → Enter "Test Response" → Accept
4. Verify responses được ghi nhận đúng
5. Chụp ảnh kết quả
```
**Độ phức tạp:** ⭐⭐⭐⭐  
**Tools:** navigate, click (x3), browser_handle_dialog (x3), get_text, screenshot

---

### 30. Complex Form Automation
```
Điền application form tại https://example.com/apply:

Personal Info:
- Full Name, DOB, Gender, Nationality

Contact:  
- Email, Phone, Address, City, Postal Code

Education:
- Degree, University, Year, GPA

Experience:
- Company, Position, Years, Description

Documents:
- Upload Resume (PDF)
- Upload Cover Letter (PDF)

Submit và handle confirmation dialog.
Chụp ảnh confirmation page.
```
**Độ phức tạp:** ⭐⭐⭐⭐⭐  
**Tools:** navigate, browser_fill_form, select_option (x4), browser_file_upload (x2), click, browser_handle_dialog, screenshot

---

## 📊 Phân loại theo Use Case

### Web Scraping
- Prompts: 3, 5, 11, 21, 22

### Testing  
- Prompts: 4, 7, 12, 16, 26, 27, 28

### Form Automation
- Prompts: 6, 13, 14, 23, 30

### Navigation & Interaction
- Prompts: 1, 2, 7, 9, 15, 18, 24

### Content Capture
- Prompts: 8, 16, 19, 20

### Advanced Workflows
- Prompts: 21, 22, 25, 27, 28, 29, 30

---

## 💡 Tips khi viết Prompts

### 1. Rõ ràng và Chi tiết
```
❌ "Vào Google tìm kiếm"
✅ "Vào Google.com, tìm kiếm 'Python tutorial', chờ kết quả, lấy 5 links đầu"
```

### 2. Specify Actions và Expected Results
```
✅ "Click nút Login, chờ redirect đến dashboard, verify có text 'Welcome'"
```

### 3. Handle Edge Cases
```
✅ "Nếu có error message, chụp ảnh error và báo cáo"
```

### 4. Use Timeouts cho Slow Pages
```
✅ "Chờ loading spinner biến mất (timeout 30 giây)"
```

### 5. Combine Multiple Steps
```
✅ Claude sẽ tự động break down thành tools sequence
```

---

## 🎯 Kết luận

30 sample prompts này cover:
- ✅ All 30 tools của MCP ChromeDP Server
- ✅ Use cases từ đơn giản → phức tạp
- ✅ Real-world scenarios
- ✅ Best practices

**Bắt đầu từ prompts cơ bản**, sau đó dần nâng cao! 🚀

# Hỗ trợ Microsoft Edge

## Tổng quan

MCP ChromeDP Server hiện đã **hỗ trợ cả Microsoft Edge và Google Chrome**! Server sẽ tự động phát hiện và sử dụng browser có sẵn trên hệ thống.

## Thứ tự ưu tiên

Server tự động tìm kiếm browser theo thứ tự:

### Windows
1. **Microsoft Edge** (ưu tiên - thường đã cài sẵn trên Windows 10/11)
2. Google Chrome
3. Chromium

### macOS
1. Microsoft Edge
2. Google Chrome

### Linux  
1. Microsoft Edge
2. Google Chrome
3. Chromium

## Tại sao nên dùng Edge?

### ✅ Ưu điểm của Edge trên Windows

1. **Đã cài sẵn** - Windows 10/11 đều có Edge pre-installed
2. **Tích hợp tốt** - Native integration với Windows
3. **Performance** - Tối ưu hóa cho Windows
4. **Same engine** - Dùng Chromium engine giống Chrome
5. **DevTools Protocol** - Hỗ trợ đầy đủ CDP như Chrome

### ✅ Khi nào dùng Chrome?

- Khi bạn đã cài Chrome và muốn consistency
- Khi test cross-browser trên nhiều platforms
- Khi cần Chrome-specific extensions (ít dùng trong headless mode)

## Cách hoạt động

### Auto-detection

```go
// Server tự động check các đường dẫn sau:
Windows:
  %PROGRAMFILES%\Microsoft\Edge\Application\msedge.exe
  %PROGRAMFILES(X86)%\Microsoft\Edge\Application\msedge.exe
  %LOCALAPPDATA%\Microsoft\Edge\Application\msedge.exe
  %PROGRAMFILES%\Google\Chrome\Application\chrome.exe
  %PROGRAMFILES(X86)%\Google\Chrome\Application\chrome.exe
  %LOCALAPPDATA%\Google\Chrome\Application\chrome.exe

macOS:
  /Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge
  /Applications/Google Chrome.app/Contents/MacOS/Google Chrome

Linux:
  /usr/bin/microsoft-edge
  /usr/bin/microsoft-edge-stable
  /usr/bin/google-chrome
  /usr/bin/chromium-browser
  /usr/bin/chromium
```

### Log output

Khi server start, bạn sẽ thấy log:
```
Found browser at: C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe
```

## Cài đặt Edge (nếu chưa có)

### Windows

**Kiểm tra Edge đã có chưa:**
```powershell
Test-Path "${env:ProgramFiles(x86)}\Microsoft\Edge\Application\msedge.exe"
```

**Cài đặt Edge:**
```powershell
# Method 1: WinGet (Windows 10/11)
winget install Microsoft.Edge

# Method 2: Download
# https://www.microsoft.com/edge

# Method 3: Chocolatey
choco install microsoft-edge
```

### macOS

```bash
# Homebrew
brew install --cask microsoft-edge

# Hoặc download trực tiếp
# https://www.microsoft.com/edge
```

### Linux (Ubuntu/Debian)

```bash
# Thêm Microsoft repository
curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.gpg
sudo install -o root -g root -m 644 microsoft.gpg /etc/apt/trusted.gpg.d/
sudo sh -c 'echo "deb [arch=amd64] https://packages.microsoft.com/repos/edge stable main" > /etc/apt/sources.list.d/microsoft-edge.list'

# Cài đặt
sudo apt update
sudo apt install microsoft-edge-stable
```

### Linux (Fedora)

```bash
# Thêm repository
sudo rpm --import https://packages.microsoft.com/keys/microsoft.asc
sudo dnf config-manager --add-repo https://packages.microsoft.com/yumrepos/edge
sudo mv /etc/yum.repos.d/packages.microsoft.com_yumrepos_edge.repo /etc/yum.repos.d/microsoft-edge.repo

# Cài đặt
sudo dnf install microsoft-edge-stable
```

## Kiểm tra Browser đã cài

Sử dụng tool `browser_install`:

```
Prompt cho Claude:
"Kiểm tra xem Edge hoặc Chrome đã cài chưa"
```

**Output mẫu khi có Edge:**
```
✓ Found browser(s):
  • Microsoft Edge at: C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe
```

**Output mẫu khi có cả hai:**
```
✓ Found browser(s):
  • Microsoft Edge at: C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe
  • Google Chrome at: C:\Program Files\Google\Chrome\Application\chrome.exe
```

## So sánh Edge vs Chrome

| Feature | Microsoft Edge | Google Chrome | Chromium |
|---------|---------------|---------------|----------|
| **Pre-installed (Windows)** | ✅ Yes | ❌ No | ❌ No |
| **DevTools Protocol** | ✅ Full | ✅ Full | ✅ Full |
| **Headless Mode** | ✅ Yes | ✅ Yes | ✅ Yes |
| **Screenshot** | ✅ Yes | ✅ Yes | ✅ Yes |
| **PDF Export** | ✅ Yes | ✅ Yes | ✅ Yes |
| **JavaScript Eval** | ✅ Yes | ✅ Yes | ✅ Yes |
| **Network Monitoring** | ✅ Yes | ✅ Yes | ✅ Yes |
| **Performance** | ⚡ Optimized for Windows | ⚡ Cross-platform | ⚡ Lightweight |
| **Updates** | 🔄 Windows Update | 🔄 Auto-update | 📦 Manual |

**Kết luận:** Cả Edge và Chrome đều hoạt động hoàn hảo với MCP ChromeDP Server!

## Testing với Edge

Sau khi build và configure, test với Edge:

### Test 1: Basic Navigation
```
Prompt: "Mở trang https://example.com và chụp ảnh"
```

### Test 2: Check Browser
```
Prompt: "Kiểm tra browser đang dùng là gì"
→ Server log sẽ show: Found browser at: ...msedge.exe
```

### Test 3: Full workflow
```
Prompt: "Vào Google, tìm 'Microsoft Edge features', lấy 5 kết quả đầu"
→ Hoạt động giống hệt với Chrome
```

## FAQ

### Q: Server sẽ dùng Edge hay Chrome?
**A:** Server tự động chọn theo thứ tự ưu tiên. Trên Windows, Edge được ưu tiên vì thường đã cài sẵn.

### Q: Làm sao force dùng Chrome thay vì Edge?
**A:** Hiện tại server auto-detect. Nếu muốn force Chrome, bạn có thể:
1. Uninstall Edge (không khuyến nghị)
2. Hoặc modify code để đảo thứ tự trong `findBrowserExecutable()`

### Q: Edge và Chrome có chạy khác nhau không?
**A:** Không! Cả hai đều dùng Chromium engine và Chrome DevTools Protocol, nên hoạt động giống hệt nhau.

### Q: Headless mode có hoạt động với Edge không?
**A:** Có! Edge hỗ trợ headless mode giống Chrome.

### Q: Performance có khác biệt không?
**A:** Trên Windows, Edge có thể chạy tốt hơn một chút vì được optimize cho Windows. Nhưng sự khác biệt là minimal.

### Q: Tôi có thể dùng Edge Dev/Beta/Canary không?
**A:** Hiện tại server chỉ detect Edge Stable. Nhưng bạn có thể modify `findBrowserExecutable()` để thêm các paths của Dev/Beta channel.

## Technical Details

### Edge DevTools Protocol

Edge sử dụng **Microsoft Edge DevTools Protocol**, là implementation của Chrome DevTools Protocol:

- **Compatible API** - 100% tương thích với Chrome CDP
- **Same commands** - Tất cả commands giống Chrome
- **Same events** - Events và responses giống nhau
- **Same domains** - Page, DOM, Network, Runtime, etc.

### Browser Executable Paths

Server detect Edge tại:

**Windows:**
```
C:\Program Files\Microsoft\Edge\Application\msedge.exe
C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe
%LOCALAPPDATA%\Microsoft\Edge\Application\msedge.exe
```

**macOS:**
```
/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge
```

**Linux:**
```
/usr/bin/microsoft-edge
/usr/bin/microsoft-edge-stable
/usr/bin/microsoft-edge-beta
/usr/bin/microsoft-edge-dev
```

## Migration từ Chrome-only

Nếu bạn đang dùng version cũ (chỉ support Chrome):

### Before (Chrome only)
```go
opts := chromedp.DefaultExecAllocatorOptions[:]
```

### After (Edge + Chrome)
```go
execPath := findBrowserExecutable()  // Auto-detect Edge or Chrome
opts := []chromedp.ExecAllocatorOption{
    chromedp.ExecPath(execPath),
}
```

**No breaking changes!** Server tự động backward compatible.

## Kết luận

✅ **MCP ChromeDP Server giờ hỗ trợ cả Edge và Chrome**  
✅ **Auto-detection** - Không cần config gì thêm  
✅ **Windows-friendly** - Edge đã cài sẵn trên Win 10/11  
✅ **100% compatible** - All 30 tools hoạt động giống nhau  
✅ **No performance difference** - Cả hai đều dùng Chromium  

**Recommendation:** 
- **Windows users:** Dùng Edge (đã có sẵn)
- **Mac/Linux users:** Edge hoặc Chrome đều ok
- **Developers:** Không cần quan tâm, server tự chọn!

🚀 Happy automation với Edge hoặc Chrome!

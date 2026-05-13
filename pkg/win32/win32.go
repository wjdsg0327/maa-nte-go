package win32

import (
	"syscall"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	getWindowRectProc    = user32.NewProc("GetWindowRect")
	getForegroundWindow  = user32.NewProc("GetForegroundWindow")
	isWindowVisibleProc  = user32.NewProc("IsWindowVisible")
	getWindowTextWProc   = user32.NewProc("GetWindowTextW")
	getWindowTextLenProc = user32.NewProc("GetWindowTextLengthW")
	enumWindowsProc      = user32.NewProc("EnumWindows")
)

type Rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func (r *Rect) Width() int32 {
	return r.Right - r.Left
}

func (r *Rect) Height() int32 {
	return r.Bottom - r.Top
}

type HWND uintptr

type WindowInfo struct {
	Handle     HWND
	Title      string
	Rect       Rect
	ClassName  string
}

// GetWindowRect 获取窗口位置和大小
func GetWindowRect(hwnd HWND) (Rect, bool) {
	var rect Rect
	ret, _, _ := getWindowRectProc.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)),
	)
	return rect, ret != 0
}

// IsWindowVisible 检查窗口是否可见
func IsWindowVisible(hwnd HWND) bool {
	ret, _, _ := isWindowVisibleProc.Call(uintptr(hwnd))
	return ret != 0
}

// GetWindowText 获取窗口标题
func GetWindowText(hwnd HWND) string {
	length := getTextLengthW(hwnd)
	if length == 0 {
		return ""
	}
	buf := make([]uint16, length+1)
	getWindowTextWProc.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	return syscall.UTF16ToString(buf)
}

func getTextLengthW(hwnd HWND) int {
	ret, _, _ := getTextLengthWProc.Call(uintptr(hwnd))
	return int(ret)
}

var getTextLengthWProc = user32.NewProc("GetWindowTextLengthW")

// GetAllVisibleWindows 获取所有可见窗口
func GetAllVisibleWindows() []WindowInfo {
	var windows []WindowInfo

	cb := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		hWnd := HWND(hwnd)
		if !IsWindowVisible(hWnd) {
			return 1 // 继续枚举
		}

		title := GetWindowText(hWnd)
		if title == "" {
			return 1
		}

		rect, ok := GetWindowRect(hWnd)
		if !ok || rect.Width() <= 0 || rect.Height() <= 0 {
			return 1
		}

		windows = append(windows, WindowInfo{
			Handle: hWnd,
			Title:  title,
			Rect:   rect,
		})

		return 1 // 继续枚举
	})

	enumWindowsProc.Call(cb, 0)

	return windows
}

// FindWindowByTitle 根据标题查找窗口
func FindWindowByTitle(title string) *WindowInfo {
	windows := GetAllVisibleWindows()
	for _, w := range windows {
		if contains(w.Title, title) {
			return &w
		}
	}
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	moduser32   = syscall.MustLoadDLL("user32.dll")
	modcomctl32 = syscall.MustLoadDLL("comctl32")
	moduxtheme  = syscall.MustLoadDLL("uxtheme")

	procGetDesktopWindow              = moduser32.MustFindProc("GetDesktopWindow")
	procSetWindowText                 = moduser32.MustFindProc("SetWindowTextW")
	procGetWindowText                 = moduser32.MustFindProc("GetWindowTextW")
	procGetWindowTextLength           = moduser32.MustFindProc("GetWindowTextLengthW")
	procGetDialogBaseUnits            = moduser32.MustFindProc("GetDialogBaseUnits")
	procMapWindowPoints               = moduser32.MustFindProc("MapWindowPoints")
	procInitCommonControls            = modcomctl32.MustFindProc("InitCommonControls")
	procGetThemeBackgroundContentRect = moduxtheme.MustFindProc("GetThemeBackgroundContentRect")
)

func init() {
	InitCommonControls()
}

func GetDesktopWindow() win.HWND {
	r1, _, err := syscall.Syscall(procGetDesktopWindow.Addr(), 0, 0, 0, 0)
	if err != 0 {
		panic(err)
	}
	return win.HWND(r1)
}

func GetWindowText(hWnd win.HWND) string {
	r0, _, _ := syscall.Syscall(procGetWindowTextLength.Addr(), 1, uintptr(hWnd), 0, 0)
	if r0 < 80 {
		var buffer [80]uint16
		r0, _, _ := syscall.Syscall(procGetWindowText.Addr(), 3, uintptr(hWnd), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
		return syscall.UTF16ToString(buffer[:r0])
	}
	buffer := make([]uint16, r0)
	r0, _, _ = syscall.Syscall(procGetWindowText.Addr(), 3, uintptr(hWnd), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	return syscall.UTF16ToString(buffer[:r0])
}

func GetWindowTextLength(hWnd win.HWND) int32 {
	r0, _, _ := syscall.Syscall(procGetWindowTextLength.Addr(), 1, uintptr(hWnd), 0, 0)
	return int32(r0)
}

func SetWindowText(hWnd win.HWND, text *uint16) win.BOOL {
	r0, _, _ := syscall.Syscall(procSetWindowText.Addr(), 2, uintptr(hWnd), uintptr(unsafe.Pointer(text)), 0)
	return win.BOOL(r0)
}

func GetDialogBaseUnits() (v uint16, h uint16) {
	r0, _, _ := syscall.Syscall(procGetDialogBaseUnits.Addr(), 0, 0, 0, 0)
	return win.HIWORD(uint32(r0)), win.LOWORD(uint32(r0))
}

func MapWindowPoints(hWndFrom win.HWND, hWndTo win.HWND, lpPoints *win.POINT, cPoints uint32) (dx uint16, dy uint16) {
	r0, _, _ := syscall.Syscall6(procMapWindowPoints.Addr(), 4, uintptr(hWndFrom), uintptr(hWndTo), uintptr(unsafe.Pointer(lpPoints)), uintptr(cPoints), 0, 0)
	return win.HIWORD(uint32(r0)), win.LOWORD(uint32(r0))
}

func InitCommonControls() {
	syscall.Syscall(procInitCommonControls.Addr(), 0, 0, 0, 0)
}

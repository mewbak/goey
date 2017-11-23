package goey

import (
	"github.com/lxn/win"
	"image"
	"syscall"
	"unsafe"
)

var (
	hr struct {
		className *uint16
		atom      win.ATOM
	}
)

func init() {
	var err error
	hr.className, err = syscall.UTF16PtrFromString("goey_hr")
	if err != nil {
		panic(err)
	}
}

func registerHRClass(hInst win.HINSTANCE, wndproc uintptr) (win.ATOM, error) {
	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.HInstance = hInst
	wc.LpfnWndProc = wndproc
	wc.HCursor = win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	wc.HbrBackground = win.GetSysColorBrush(win.COLOR_3DFACE)
	wc.LpszClassName = hr.className

	atom := win.RegisterClassEx(&wc)
	if atom == 0 {
		return 0, syscall.GetLastError()
	}
	return atom, nil
}

func (w *HR) Mount(parent NativeWidget) (MountedWidget, error) {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		return nil, syscall.GetLastError()
	}
	if hr.atom == 0 {
		atom, err := registerHRClass(hInstance, syscall.NewCallback(hrWindowProc))
		if err != nil {
			return nil, err
		}
		hr.atom = atom
	}

	hwnd := win.CreateWindowEx(0, hr.className, nil, win.WS_CHILD|win.WS_VISIBLE,
		10, 10, 100, 100,
		parent.hWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}

	retval := &mountedHR{NativeWidget: NativeWidget{hwnd}}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedHR struct {
	NativeWidget
}

func (w *mountedHR) MeasureWidth() (DP, DP) {
	return 1, 0xffff
}

func (w *mountedHR) MeasureHeight(width DP) (DP, DP) {
	// Same as static text
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13, 13
}

func (w *mountedHR) SetBounds(bounds image.Rectangle) {
	w.NativeWidget.SetBounds(bounds)
}

func (w *mountedHR) UpdateProps(data_ Widget) error {
	return nil
}

func hrWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedHR)(unsafe.Pointer(w))
			ptr.hWnd = 0
		}
		// Defer to the old window proc

	case win.WM_PAINT:
		ps := win.PAINTSTRUCT{}
		rect := win.RECT{}
		hdc := win.BeginPaint(hwnd, &ps)
		win.GetClientRect(hwnd, &rect)
		win.MoveToEx(hdc, int(rect.Left), int(rect.Top+rect.Bottom)/2, nil)
		win.LineTo(hdc, rect.Right, (rect.Top+rect.Bottom)/2)
		win.EndPaint(hwnd, &ps)
	}

	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

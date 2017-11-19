package goey

import (
	"github.com/lxn/win"
	"image"
	"sync/atomic"
	"syscall"
	"unsafe"
)

var (
	mainWindow struct {
		className *uint16
		atom      win.ATOM
	}
	mainWindowCount int32 = 0
	hMessageFont    win.HFONT
	activeWindow    uintptr
)

func init() {
	var err error
	mainWindow.className, err = syscall.UTF16PtrFromString("goey_mainwindow")
	if err != nil {
		panic(err)
	}

	// Determine the mssage font
	var ncm win.NONCLIENTMETRICS
	ncm.CbSize = uint32(unsafe.Sizeof(ncm))
	if rc := win.SystemParametersInfo(win.SPI_GETNONCLIENTMETRICS, ncm.CbSize, unsafe.Pointer(&ncm), 0); rc {
		hMessageFont = win.CreateFontIndirect(&ncm.LfMessageFont)
		if hMessageFont == 0 {
			println("failed CreateFontIndirect")
		}
	} else {
		println("failed SystemParametersInfo")
	}
}

type MainWindow struct {
	mountedVBox

	hWnd               win.HWND
	dpi                image.Point
	clientMinimumWidth int
}

func registerMainWindowClass(hInst win.HINSTANCE, wndproc uintptr) (win.ATOM, error) {
	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.HInstance = hInst
	wc.LpfnWndProc = wndproc
	wc.HCursor = win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	wc.HbrBackground = win.GetSysColorBrush(win.COLOR_3DFACE)
	wc.LpszClassName = mainWindow.className

	atom := win.RegisterClassEx(&wc)
	if atom == 0 {
		return 0, syscall.GetLastError()
	}
	return atom, nil
}

func onSize(hwnd win.HWND, mw *MainWindow) {
	// The recommended margin in device independent pixels.
	const margin = 11

	// Get the client rect for the main window.  This is the layout region.
	rect := win.RECT{}
	win.GetClientRect(hwnd, &rect)

	// Yes it's ugly, the SetBounds method for windows uses the screen DPI to
	// convert device independent pixels into actual pixels, but the DPI can change
	// from window to window when the computer has multiple monitors.  Fortunately,
	// all layout should happen in the GUI thread.
	dpi = mw.dpi

	// We will adjust the margins based on the screen size and preferred width
	// of the content.
	cpw := mw.ClientMinimumWidth()
	width := rect.Right - rect.Left
	if cpw+2*margin <= int(width) {
		mw.SetBounds(image.Rect(int(rect.Left)+margin, int(rect.Top)+margin, int(rect.Right)-margin, int(rect.Bottom)-margin))
	} else if cpw >= int(width) {
		mw.SetBounds(image.Rect(int(rect.Left), int(rect.Top)+margin, int(rect.Right), int(rect.Bottom)-margin))
	} else {
		left := (int(width) - cpw) / 2
		mw.SetBounds(image.Rect(left, int(rect.Top)+margin, left+cpw, int(rect.Bottom)-margin))
	}
}

func wndproc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {

	switch msg {
	case win.WM_CREATE:
		// Maintain count of open windows.
		atomic.AddInt32(&mainWindowCount, 1)
		// Defer to default window proc

	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*MainWindow)(unsafe.Pointer(w))
			ptr.hWnd = 0
		}
		// Make sure we are no longer linked to as the active window
		atomic.CompareAndSwapUintptr(&activeWindow, uintptr(hwnd), 0)
		// If this is the last main window visible, post the quit message so that the
		// message loop terminates.
		if newval := atomic.AddInt32(&mainWindowCount, -1); newval == 0 {
			win.PostQuitMessage(0)
		}
		// Defer to the default window proc

	case win.WM_ACTIVATE:
		if wParam != 0 {
			atomic.StoreUintptr(&activeWindow, uintptr(hwnd))
		}
		// Defer to the default window proc

	case win.WM_SIZE:
		w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
		onSize(hwnd, (*MainWindow)(unsafe.Pointer(w)))
		// Defer to the default window proc

	case win.WM_GETMINMAXINFO:
		mmi := (*win.MINMAXINFO)(unsafe.Pointer(lParam))
		if mmi.PtMinTrackSize.X > 0 {
			mmi.PtMinTrackSize.Y = mmi.PtMinTrackSize.X
		}
		return 0

	case win.WM_COMMAND:
		if n := win.HIWORD(uint32(wParam)); n == win.BN_CLICKED || n == win.EN_UPDATE {
			return win.SendDlgItemMessage(hwnd, int32(win.LOWORD(uint32(wParam))), msg, wParam, lParam)
		}
		// Defer to the default window proc

	}

	// Let the default window proc handle all other messages
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func newMainWindow(title string, children []Widget) (*MainWindow, error) {
	const Width = 640
	const Height = 480

	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		return nil, syscall.GetLastError()
	}
	//GetStartupInfo(&info);
	if win.OleInitialize() != win.S_OK {
		return nil, syscall.GetLastError()
	}
	if mainWindow.atom == 0 {
		atom, err := registerMainWindowClass(hInstance, syscall.NewCallback(wndproc))
		if err != nil {
			return nil, err
		}
		mainWindow.atom = atom
	}

	style := uint32(win.WS_OVERLAPPEDWINDOW)
	//if !settings.Resizable {
	//	style = win.WS_OVERLAPPED | win.WS_CAPTION | win.WS_MINIMIZEBOX | win.WS_SYSMENU
	//}

	rect := win.RECT{0, 0, Width, Height}
	win.AdjustWindowRect(&rect, win.WS_OVERLAPPEDWINDOW, false)

	var clientRect win.RECT
	win.GetClientRect(GetDesktopWindow(), &clientRect)
	left := (clientRect.Right / 2) - ((rect.Right - rect.Left) / 2)
	top := (clientRect.Bottom / 2) - ((rect.Bottom - rect.Top) / 2)
	rect.Right = rect.Right - rect.Left + left
	rect.Left = left
	rect.Bottom = rect.Bottom - rect.Top + top
	rect.Top = top

	windowName_, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return nil, err
	}
	hwnd := win.CreateWindowEx(win.WS_EX_CONTROLPARENT, mainWindow.className, windowName_, style, rect.Left, rect.Top,
		rect.Right-rect.Left, rect.Bottom-rect.Top,
		win.HWND_DESKTOP, 0, hInstance, nil)
	if hwnd == 0 {
		win.OleUninitialize()
		return nil, syscall.GetLastError()
	}

	// Set the font for the window
	if hMessageFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, 0, 0)
	}

	retval := &MainWindow{hWnd: hwnd}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	// Determine the DPI for this window
	hdc := win.GetDC(hwnd)
	retval.dpi.X = int(win.GetDeviceCaps(hdc, win.LOGPIXELSX))
	retval.dpi.Y = int(win.GetDeviceCaps(hdc, win.LOGPIXELSY))
	win.ReleaseDC(hwnd, hdc)

	vbox := VBox{children}
	mounted, err := vbox.Mount(NativeWidget{hwnd})
	if err != nil {
		win.DestroyWindow(hwnd)
		return nil, err
	}
	retval.mountedVBox = *mounted.(*mountedVBox)

	win.ShowWindow(hwnd, win.SW_SHOW /* info.wShowWindow */)
	win.UpdateWindow(hwnd)
	win.SetFocus(hwnd)

	return retval, nil
}

func (w *MainWindow) Close() {
	// Want to be able to close windows in Go, even if they have already been
	// destroyed in the Win32 system
	if w.hWnd != 0 {
		win.DestroyWindow(w.hWnd)
	}
	win.OleUninitialize()
}

func (w *MainWindow) ClientMinimumWidth() int {
	if w.clientMinimumWidth > 0 {
		return w.clientMinimumWidth
	}

	w.clientMinimumWidth = w.mountedVBox.MinimumWidth().ToPixelsX()
	return w.clientMinimumWidth
}

func (w *MainWindow) SetText(value string) error {
	return NativeWidget{w.hWnd}.SetText(value)
}

func (w *MainWindow) SetChildren(children []Widget) error {
	// Defer to the vertical box holding the children.
	vbox := VBox{children}
	err := w.mountedVBox.UpdateProps(&vbox)
	w.clientMinimumWidth = 0
	// Whether or not an error has occured, redo the layout so the children
	// are placed.
	currentHwnd := win.HWND_TOP
	for _, v := range w.children {
		currentHwnd = v.SetOrder(currentHwnd)
	}
	w.clientMinimumWidth = 0
	onSize(w.hWnd, w)
	// ... and we're done
	return err
}

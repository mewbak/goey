package goey

import (
	"image"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
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

const (
	Scale = 1
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
		ncm.LfMessageFont.LfHeight = int32(float64(ncm.LfMessageFont.LfHeight) * Scale)
		ncm.LfMessageFont.LfWidth = int32(float64(ncm.LfMessageFont.LfWidth) * Scale)
		hMessageFont = win.CreateFontIndirect(&ncm.LfMessageFont)
		if hMessageFont == 0 {
			println("failed CreateFontIndirect")
		}
	} else {
		println("failed SystemParametersInfo")
	}
}

type mainWindow struct {
	vbox mountedVBox

	hWnd             win.HWND
	dpi              image.Point
	clientMinWidth   int
	clientMaxWidth   int
	clientMinHeight  int
	clientHeight     int
	scrollbarVisible bool
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

func onSizeCalcMargin(clientMinWidth int, availableWidth int, margin int) int {
	if clientMinWidth+2*margin <= availableWidth {
		return margin
	} else if clientMinWidth >= availableWidth {
		return 0
	} else {
		return (availableWidth - clientMinWidth) / 2
	}
}

func onSize(hwnd win.HWND, mw *MainWindow) {
	// The recommended margin in device independent pixels.
	const margin = DIP(11)

	// Get the client rect for the main window.  This is the layout region.
	rect := win.RECT{}
	win.GetClientRect(hwnd, &rect)
	scrollPos := 0

	// Yes it's ugly, the SetBounds method for windows uses the screen DPI to
	// convert device independent pixels into actual pixels, but the DPI can change
	// from window to window when the computer has multiple monitors.  Fortunately,
	// all layout should happen in the GUI thread.
	mw.updateGlobalDPI()

	// We will adjust the margins based on the screen size and preferred width
	// of the content.
	availableMargin := onSizeCalcMargin(mw.clientMinWidth, int(rect.Right-rect.Left), margin.PixelsX())
	width := int(rect.Right-rect.Left) - 2*availableMargin
	minHeight, _ := mw.vbox.MeasureHeight(ToDIPX(width))
	if minHeight.PixelsY() > int(rect.Bottom-rect.Top)-2*margin.PixelsY() {
		if !mw.scrollbarVisible {
			// Create the scroll bar
			ShowScrollBar(hwnd, win.SB_VERT, win.TRUE)
			mw.scrollbarVisible = true

			// The client rect will have changed.  Need to refresh.
			win.GetClientRect(hwnd, &rect)

			// Update the calculations above, since we now need to account for
			// the width of the scroll bar
			availableMargin := onSizeCalcMargin(mw.clientMinWidth, int(rect.Right-rect.Left), margin.PixelsX())
			width := int(rect.Right-rect.Left) - 2*availableMargin
			minHeight, _ = mw.vbox.MeasureHeight(ToDIPX(width))
		}
		si := win.SCROLLINFO{
			FMask: win.SIF_PAGE | win.SIF_RANGE,
			NMin:  0,
			NMax:  int32(minHeight.PixelsY() + 2*margin.PixelsY()),
			NPage: uint32(rect.Bottom - rect.Top),
		}
		si.CbSize = uint32(unsafe.Sizeof(si))
		win.SetScrollInfo(hwnd, win.SB_VERT, &si, true)
		si.FMask = win.SIF_POS
		win.GetScrollInfo(hwnd, win.SB_CTL, &si)
		scrollPos = int(si.NPos)

		// Perform layout
		mw.vbox.SetBounds(image.Rect(int(rect.Left)+availableMargin, int(rect.Top)+margin.PixelsY()-scrollPos, int(rect.Right)-availableMargin, int(rect.Top)+margin.PixelsY()+minHeight.PixelsY()-scrollPos))
	} else {
		if mw.scrollbarVisible {
			// Remove the scroll bar
			ShowScrollBar(hwnd, win.SB_VERT, win.FALSE)
			mw.scrollbarVisible = false

			// The client rect will have changed.  Need to refresh.
			win.GetClientRect(hwnd, &rect)

			// Update the calculations above, since we now need to account for
			// the width of the scroll bar
			availableMargin = onSizeCalcMargin(mw.clientMinWidth, int(rect.Right-rect.Left), margin.PixelsX())
		}

		// Perform layout
		mw.vbox.SetBounds(image.Rect(int(rect.Left)+availableMargin, int(rect.Top)+margin.PixelsY(), int(rect.Right)-availableMargin, int(rect.Bottom)-margin.PixelsY()))
	}

	// Update the position of all of the children
	win.InvalidateRect(hwnd, &rect, true)
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
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			mw := (*MainWindow)(unsafe.Pointer(w))
			mmi := (*win.MINMAXINFO)(unsafe.Pointer(lParam))
			if mmi.PtMinTrackSize.X < int32(mw.clientMinWidth) {
				mmi.PtMinTrackSize.X = int32(mw.clientMinWidth)
			}
			//if mmi.PtMaxTrackSize.X > int32(mw.clientMaxWidth) {
			//	mmi.PtMaxTrackSize.X = int32(mw.clientMaxWidth)
			//}
			mmi.PtMinTrackSize.Y += int32(mw.clientMinHeight)
			if mmi.PtMinTrackSize.Y > 480 {
				mmi.PtMinTrackSize.Y = 480
			}
		}
		return 0

	case win.WM_VSCROLL:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			mw := (*MainWindow)(unsafe.Pointer(w))

			// Get all the vertial scroll bar information.
			si := win.SCROLLINFO{FMask: win.SIF_ALL}
			si.CbSize = uint32(unsafe.Sizeof(si))
			win.GetScrollInfo(hwnd, win.SB_VERT, &si)

			// Save the position for comparison later on.
			yPos := si.NPos
			switch win.LOWORD(uint32(wParam)) {
			// User clicked the HOME keyboard key.
			case win.SB_TOP:
				si.NPos = si.NMin

			// User clicked the END keyboard key.
			case win.SB_BOTTOM:
				si.NPos = si.NMax

			// User clicked the top arrow.
			case win.SB_LINEUP:
				si.NPos -= int32(DIP(13).PixelsY())

			// User clicked the bottom arrow.
			case win.SB_LINEDOWN:
				si.NPos += int32(DIP(13).PixelsY())

			// User clicked the scroll bar shaft above the scroll box.
			case win.SB_PAGEUP:
				si.NPos -= int32(si.NPage)

			// User clicked the scroll bar shaft below the scroll box.
			case win.SB_PAGEDOWN:
				si.NPos += int32(si.NPage)

			// User dragged the scroll box.
			case win.SB_THUMBTRACK:
				si.NPos = si.NTrackPos
			}

			// Set the position and then retrieve it.  Due to adjustments
			// by Windows it may not be the same as the value set.
			si.FMask = win.SIF_POS
			win.SetScrollInfo(hwnd, win.SB_VERT, &si, true)
			win.GetScrollInfo(hwnd, win.SB_VERT, &si)

			// If the position has changed, scroll window and update it.
			if si.NPos != yPos {
				const margin = DIP(11)

				rect := win.RECT{}
				win.GetClientRect(hwnd, &rect)
				availableMargin := onSizeCalcMargin(mw.clientMinWidth, int(rect.Right-rect.Left), margin.PixelsX())
				width := int(rect.Right-rect.Left) - 2*availableMargin
				minHeight, _ := mw.vbox.MeasureHeight(ToDIPX(width))
				mw.vbox.SetBounds(image.Rect(int(rect.Left)+availableMargin, int(rect.Top)+margin.PixelsY()-int(si.NPos), int(rect.Right)-availableMargin, int(rect.Top)+margin.PixelsY()+minHeight.PixelsY()-int(si.NPos)))

				// TODO:  Use ScrollWindow function to reduce flicker during scrolling
				win.InvalidateRect(hwnd, &rect, true)
			}

			return 0
		}

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
	hwnd := win.CreateWindowEx(win.WS_EX_CONTROLPARENT, mainWindow.className, windowName_, style,
		rect.Left, rect.Top, rect.Right-rect.Left, rect.Bottom-rect.Top,
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

	vbox := VBox{children, MainStart, Stretch}
	mounted, err := vbox.Mount(NativeWidget{hwnd})
	if err != nil {
		win.DestroyWindow(hwnd)
		return nil, err
	}
	retval.vbox = *mounted.(*mountedVBox)
	retval.determineSizeConstraints()

	win.ShowWindow(hwnd, win.SW_SHOW /* info.wShowWindow */)
	win.UpdateWindow(hwnd)
	win.SetFocus(hwnd)

	return retval, nil
}

func (w *MainWindow) determineSizeConstraints() {
	w.updateGlobalDPI()

	clientMinWidth, clientMaxWidth := w.vbox.MeasureWidth()
	w.clientMinWidth, w.clientMaxWidth = (clientMinWidth + 22).PixelsX(), (clientMaxWidth + 22).PixelsX()
	clientMinHeight, _ := w.vbox.MeasureHeight(clientMaxWidth)
	w.clientMinHeight = (clientMinHeight + 22).PixelsY()
	if w.clientMinHeight > 480 {
		w.clientMinHeight = 480
	}
}

func (w *MainWindow) close() {
	// Want to be able to close windows in Go, even if they have already been
	// destroyed in the Win32 system
	if w.hWnd != 0 {
		win.DestroyWindow(w.hWnd)
	}
	win.OleUninitialize()
}

func (w *MainWindow) setAlignment(main MainAxisAlign, cross CrossAxisAlign) error {
	w.vbox.alignMain = main
	w.vbox.alignCross = cross

	onSize(w.hWnd, w)
	return nil
}

func (w *MainWindow) setChildren(children []Widget) error {
	// Defer to the vertical box holding the children.
	vbox := VBox{children, w.vbox.alignMain, w.vbox.alignCross}
	err := w.vbox.UpdateProps(&vbox)
	// Whether or not an error has occured, redo the layout so the children
	// are placed.
	currentHwnd := win.HWND_TOP
	for _, v := range w.vbox.children {
		currentHwnd = v.SetOrder(currentHwnd)
	}
	// Determine the size constraints for the window
	w.determineSizeConstraints()
	onSize(w.hWnd, w)
	// ... and we're done
	return err
}

func (w *MainWindow) setTitle(value string) error {
	return NativeWidget{w.hWnd}.SetText(value)
}

func (w *MainWindow) updateGlobalDPI() {
	DPI = image.Point{int(float32(w.dpi.X) * Scale), int(float32(w.dpi.Y) * Scale)}
}

package goey

import (
	"image"
	"sync/atomic"
	"syscall"
	"unsafe"

	win2 "bitbucket.org/rj/goey/syscall"
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

type windowImpl struct {
	hWnd                    win.HWND
	dpi                     image.Point
	child                   Element
	childSize               Size
	minimumSize             Size
	horizontalScroll        bool
	horizontalScrollVisible bool
	horizontalScrollPos     Length
	verticalScroll          bool
	verticalScrollVisible   bool
	verticalScrollPos       Length
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

func (mw *windowImpl) onSize(hwnd win.HWND) {
	if mw.child == nil {
		return
	}

	// Yes it's ugly, the SetBounds method for windows uses the screen DPI to
	// convert device independent pixels into actual pixels, but the DPI can change
	// from window to window when the computer has multiple monitors.  Fortunately,
	// all layout should happen in the GUI thread.
	mw.updateGlobalDPI()

	// Get the client rect for the main window.  This is the layout region.
	rect := win.RECT{}
	win.GetClientRect(hwnd, &rect)
	size := mw.layoutChild(Size{
		FromPixelsX(int(rect.Right - rect.Left)),
		FromPixelsY(int(rect.Bottom - rect.Top)),
	})
	for mw.showScroll(size, rect) {
		win.GetClientRect(hwnd, &rect)
		size = mw.layoutChild(Size{
			FromPixelsX(int(rect.Right - rect.Left)),
			FromPixelsY(int(rect.Bottom - rect.Top)),
		})
	}
	mw.childSize = size

	// Position the child element.
	mw.child.SetBounds(Rectangle{
		Point{-mw.horizontalScrollPos, -mw.verticalScrollPos},
		Point{size.Width - mw.horizontalScrollPos, size.Height - mw.verticalScrollPos},
	})

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
			ptr := (*windowImpl)(unsafe.Pointer(w))
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

	case win.WM_SETFOCUS:
		// The main window doesn't need focus, we want to delegate to a control
		if hwnd == win.GetFocus() { // Is this always true
			child := win.GetWindow(hwnd, win.GW_CHILD)
			for child != 0 {
				if style := win.GetWindowLong(child, win.GWL_STYLE); (style & win.WS_TABSTOP) != 0 {
					win.SetFocus(child)
					break
				}
				child = win.GetWindow(child, win.GW_HWNDNEXT)
			}
		}
		// Defer to the default window proc

	case win.WM_SIZE:
		w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
		mw := (*windowImpl)(unsafe.Pointer(w))
		mw.onSize(hwnd)
		// Defer to the default window proc

	case win.WM_GETMINMAXINFO:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			mw := (*windowImpl)(unsafe.Pointer(w))
			mw.updateGlobalDPI()
			// We need to identify how much addition width and height are
			// required for the windows border, statusbar, etc.
			// TODO:  This value could be cached
			windowRect := win.RECT{}
			win.GetWindowRect(hwnd, &windowRect)
			clientRect := win.RECT{}
			win.GetClientRect(hwnd, &clientRect)
			dx := (windowRect.Right - windowRect.Left) - (clientRect.Right - clientRect.Left)
			dy := (windowRect.Bottom - windowRect.Top) - (clientRect.Bottom - clientRect.Top)
			if mw.verticalScroll {
				dx += FromPixelsX(int(win.GetSystemMetrics(win.SM_CXVSCROLL)))
			}
			if mw.horizontalScroll {
				dy += FromPixelsY(int(win.GetSystemMetrics(win.SM_CYHSCROLL)))
			}

			// Update tracking information based on our minimum size
			mmi := (*win.MINMAXINFO)(unsafe.Pointer(lParam))
			if limit := int32(mw.minimumSize.Width.PixelsX()) + dx; mmi.PtMinTrackSize.X < limit {
				mmi.PtMinTrackSize.X = limit
			}
			if limit := int32(mw.minimumSize.Height.PixelsY()) + dy; mmi.PtMinTrackSize.Y < limit {
				mmi.PtMinTrackSize.Y = limit
			}
			// If scrolling is enabled for either direction, we can relax the
			// minimum window size.  These limits are fairly arbitrary, but we do need to
			// leave enough space for the scroll bars.
			if limit := int32((120 * DIP).PixelsX()); mw.horizontalScroll && mmi.PtMinTrackSize.X > limit {
				mmi.PtMinTrackSize.X = limit
			}
			if limit := int32((120 * DIP).PixelsY()); mw.verticalScroll && mmi.PtMinTrackSize.Y > limit {
				mmi.PtMinTrackSize.Y = limit
			}
		}
		return 0

	case win.WM_HSCROLL:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			mw := (*windowImpl)(unsafe.Pointer(w))
			mw.setScrollPos(win.SB_HORZ, wParam)
			return 0
		}

	case win.WM_VSCROLL:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			mw := (*windowImpl)(unsafe.Pointer(w))
			mw.setScrollPos(win.SB_VERT, wParam)
			return 0
		}

	case win.WM_COMMAND:
		if n := win.HIWORD(uint32(wParam)); n == win.BN_CLICKED || n == win.EN_UPDATE {
			return win.SendDlgItemMessage(hwnd, int32(win.LOWORD(uint32(wParam))), msg, wParam, lParam)
		}
		// Defer to the default window proc

	case win.WM_NOTIFY:
		if n := (*win.NMHDR)(unsafe.Pointer(lParam)); true {
			return win.SendMessage(n.HwndFrom, win.WM_NOTIFY, wParam, lParam)
		}
		// Defer to the default window proc

	}

	// Let the default window proc handle all other messages
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func newWindow(title string, child Widget) (*Window, error) {
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
		if atom == 0 {
			panic("internal error:  atom==0 although no error returned")
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
	win.GetClientRect(win2.GetDesktopWindow(), &clientRect)
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

	retval := &Window{windowImpl{hWnd: hwnd}}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(&retval.windowImpl)))

	// Determine the DPI for this window
	hdc := win.GetDC(hwnd)
	retval.dpi.X = int(win.GetDeviceCaps(hdc, win.LOGPIXELSX))
	retval.dpi.Y = int(win.GetDeviceCaps(hdc, win.LOGPIXELSY))
	win.ReleaseDC(hwnd, hdc)

	err = retval.setChild(child)
	if err != nil {
		win.DestroyWindow(hwnd)
		return nil, err
	}

	win.ShowWindow(hwnd, win.SW_SHOW /* info.wShowWindow */)
	win.UpdateWindow(hwnd)

	return retval, nil
}

func (w *windowImpl) close() {
	// Want to be able to close windows in Go, even if they have already been
	// destroyed in the Win32 system
	if w.hWnd != 0 {
		win.DestroyWindow(w.hWnd)
		w.hWnd = 0
	}
	win.OleUninitialize()
}

func (w *windowImpl) getChild() Element {
	return w.child
}

// NativeHandle returns the handle to the platform-specific window handle
// (i.e. a HWND on WIN32).
func (w *windowImpl) NativeHandle() win.HWND {
	return w.hWnd
}

func (w *windowImpl) message(m *Message) {
	m.title = win2.GetWindowText(w.hWnd)
	m.handle = uintptr(w.hWnd)
}

func (w *windowImpl) setChild(child Widget) (err error) {
	// Update the child element
	w.child, err = DiffChild(Control{w.hWnd}, w.child, child)
	// Whether or not an error has occured, redo the layout so the children
	// are placed.
	if w.child != nil {
		w.child.SetOrder(win.HWND_TOP)
		w.childSize = Size{}
		w.updateGlobalDPI()
		w.minimumSize = w.child.MinimumSize()

		// Determine the size constraints for the window
		w.onSize(w.hWnd)
	}
	// ... and we're done
	return err
}

func (w *windowImpl) setScroll(hscroll, vscroll bool) {
	// Copy the new parameters for the window into the fields.
	w.horizontalScroll, w.verticalScroll = hscroll, vscroll

	// If either scrollbar is being disabled, make sure to the state for
	// that scrollbar, and to hide it.
	if !w.horizontalScroll {
		w.horizontalScrollPos = 0
		w.horizontalScrollVisible = false
		win2.ShowScrollBar(w.hWnd, win.SB_HORZ, win.FALSE)
	}
	if !w.verticalScroll {
		w.verticalScrollPos = 0
		w.verticalScrollVisible = false
		win2.ShowScrollBar(w.hWnd, win.SB_VERT, win.FALSE)
	}

	// Even when enabled, the scrollbars appear only if required
	// by the size of the windows content.  Run the event handler
	// for resizing to reevaluate the logic for displaying the
	// scroll bars.
	w.onSize(w.hWnd)
}

func (w *windowImpl) setScrollPos(direction int32, wParam uintptr) {
	// Get all the vertial scroll bar information.
	si := win.SCROLLINFO{FMask: win.SIF_ALL}
	si.CbSize = uint32(unsafe.Sizeof(si))
	win.GetScrollInfo(w.hWnd, direction, &si)

	// Save the position for comparison later on.
	currentPos := si.NPos
	switch win.LOWORD(uint32(wParam)) {
	// User clicked the HOME keyboard key.
	case win.SB_TOP:
		si.NPos = si.NMin

	// User clicked the END keyboard key.
	case win.SB_BOTTOM:
		si.NPos = si.NMax

	// User clicked the top arrow.
	case win.SB_LINEUP:
		if direction == win.SB_HORZ {
			si.NPos -= int32((13 * DIP).PixelsX())
		} else {
			si.NPos -= int32((13 * DIP).PixelsY())
		}

	// User clicked the bottom arrow.
	case win.SB_LINEDOWN:
		if direction == win.SB_HORZ {
			si.NPos += int32((13 * DIP).PixelsX())
		} else {
			si.NPos += int32((13 * DIP).PixelsY())
		}

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
	win.SetScrollInfo(w.hWnd, direction, &si, true)
	win.GetScrollInfo(w.hWnd, direction, &si)

	// If the position has changed, scroll window and update it.
	if si.NPos != currentPos {
		if direction == win.SB_HORZ {
			w.horizontalScrollPos = FromPixelsX(int(si.NPos))
		} else {
			w.verticalScrollPos = FromPixelsY(int(si.NPos))
		}
		rect := win.RECT{}
		win.GetClientRect(w.hWnd, &rect)
		w.child.SetBounds(Rectangle{
			Point{-w.horizontalScrollPos, -w.verticalScrollPos},
			Point{w.childSize.Width - w.horizontalScrollPos, w.childSize.Height - w.verticalScrollPos},
		})

		// TODO:  Use ScrollWindow function to reduce flicker during scrolling
		win.InvalidateRect(w.hWnd, &rect, true)
	}
}

func (w *windowImpl) showScroll(size Size, rect win.RECT) bool {
	if w.horizontalScroll {
		if size.Width > FromPixelsY(int(rect.Right-rect.Left)) {
			if !w.horizontalScrollVisible {
				// Create the scroll bar
				win2.ShowScrollBar(w.hWnd, win.SB_HORZ, win.TRUE)
				w.horizontalScrollVisible = true
				return true
			}
			si := win.SCROLLINFO{
				FMask: win.SIF_PAGE | win.SIF_RANGE,
				NMin:  0,
				NMax:  int32(size.Width.PixelsX()),
				NPage: uint32(rect.Right - rect.Left),
			}
			si.CbSize = uint32(unsafe.Sizeof(si))
			win.SetScrollInfo(w.hWnd, win.SB_HORZ, &si, true)
			si.FMask = win.SIF_POS
			win.GetScrollInfo(w.hWnd, win.SB_HORZ, &si)
			w.horizontalScrollPos = FromPixelsX(int(si.NPos))
		} else if w.horizontalScrollVisible {
			// Remove the scroll bar
			win2.ShowScrollBar(w.hWnd, win.SB_HORZ, win.FALSE)
			w.horizontalScrollPos = 0
			w.horizontalScrollVisible = false
			return true
		}
	}

	if w.verticalScroll {
		if size.Height > FromPixelsY(int(rect.Bottom-rect.Top)) {
			if !w.verticalScrollVisible {
				// Create the scroll bar
				win2.ShowScrollBar(w.hWnd, win.SB_VERT, win.TRUE)
				w.verticalScrollVisible = true
				return true
			}
			si := win.SCROLLINFO{
				FMask: win.SIF_PAGE | win.SIF_RANGE,
				NMin:  0,
				NMax:  int32(size.Height.PixelsY()),
				NPage: uint32(rect.Bottom - rect.Top),
			}
			si.CbSize = uint32(unsafe.Sizeof(si))
			win.SetScrollInfo(w.hWnd, win.SB_VERT, &si, true)
			si.FMask = win.SIF_POS
			win.GetScrollInfo(w.hWnd, win.SB_VERT, &si)
			w.verticalScrollPos = FromPixelsY(int(si.NPos))
		} else if w.verticalScrollVisible {
			// Remove the scroll bar
			win2.ShowScrollBar(w.hWnd, win.SB_VERT, win.FALSE)
			w.verticalScrollPos = 0
			w.verticalScrollVisible = false
			return true
		}
	}

	return false
}

func (w *windowImpl) setIcon(img image.Image) error {
	hicon, _, err := imageToIcon(img)
	if err != nil {
		return err
	}
	win2.SetClassLongPtr(w.hWnd, win2.GCLP_HICON, uintptr(hicon))

	return nil
}

func (w *windowImpl) setTitle(value string) error {
	return Control{w.hWnd}.SetText(value)
}

func (w *windowImpl) updateGlobalDPI() {
	DPI = image.Point{int(float32(w.dpi.X) * Scale), int(float32(w.dpi.Y) * Scale)}
}

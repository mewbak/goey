package goey

import (
	win2 "bitbucket.org/rj/goey/syscall"
	"github.com/lxn/win"
	"image"
	"sync/atomic"
	"syscall"
	"unsafe"
)

func init() {
	// If the return of the call to InitCommonControlsEx is checked, we see
	// false, which according to the documentation indicates that it failed.
	// However, there is no error with syscall.GetLastError().
	//
	// Note:  The init function for github.com/lxn/win also calls this
	// function, but does not include ICC_STANDARD_CLASSES.
	initCtrls := win.INITCOMMONCONTROLSEX{}
	initCtrls.DwSize = uint32(unsafe.Sizeof(initCtrls))
	initCtrls.DwICC = win.ICC_STANDARD_CLASSES | win.ICC_DATE_CLASSES
	win.InitCommonControlsEx(&initCtrls)
}

// Control ID

var (
	currentControlID uint32 = 100
)

func nextControlID() uint32 {
	return atomic.AddUint32(&currentControlID, 1)
}

// NativeWidget is an opaque type used as a platform-specific handle to a
// window or widget for WIN32 builds.
type NativeWidget struct {
	hWnd win.HWND
}

// Text copies text of the underlying window
func (w NativeWidget) Text() string {
	return win2.GetWindowText(w.hWnd)
}

func (w NativeWidget) SetDisabled(value bool) {
	win.EnableWindow(w.hWnd, !value)
}

func (w *NativeWidget) SetBounds(bounds image.Rectangle) {
	win.MoveWindow(w.hWnd, int32(bounds.Min.X), int32(bounds.Min.Y), int32(bounds.Dx()), int32(bounds.Dy()), true)
}

func (w *NativeWidget) SetOrder(previous win.HWND) win.HWND {
	// Note, the argument previous may be 0 when setting the first child.
	// Fortunately, this corresponds to HWND_TOP, which sets the window
	// to top of the z-order.
	win.SetWindowPos(w.hWnd, previous, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
	return w.hWnd
}

func (w NativeWidget) SetText(value string) error {
	utf16, err := syscall.UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	rc := win2.SetWindowText(w.hWnd, utf16)
	if rc == 0 {
		return syscall.GetLastError()
	}
	return nil
}

func (w *NativeWidget) Close() {
	if w.hWnd != 0 {
		win.DestroyWindow(w.hWnd)
		w.hWnd = 0
	}
}

// NativeMountedWidget contains platform-specific methods that all widgets
// must support on WIN32
type NativeMountedWidget interface {
	MeasureWidth() (min Length, max Length)
	MeasureHeight(width Length) (min Length, max Length)
	SetBounds(bounds image.Rectangle)
	SetOrder(previous win.HWND) win.HWND
}

func calculateHGap(previous MountedWidget, current MountedWidget) Length {
	// The vertical gap between most controls is 11 relative pixels.  However,
	// there are different rules for between a label and its associated control,
	// or between related controls.  These relationship do not appear in the
	// model provided by this package, so these relationships need to be
	// inferred from the order and type of controls.
	//
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	if _, ok := previous.(*mountedButton); ok {
		if _, ok := current.(*mountedButton); ok {
			// Any pair of successive buttons will be assumed to be in a
			// related group.
			return 7 * DIP
		}
	}

	// The spacing between unrelated controls.
	return 11 * DIP
}

func calculateVGap(previous MountedWidget, current MountedWidget) Length {
	// The vertical gap between most controls is 11 relative pixels.  However,
	// there are different rules for between a label and its associated control,
	// or between related controls.  These relationship do not appear in the
	// model provided by this package, so these relationships need to be
	// inferred from the order and type of controls.
	//
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	if _, ok := previous.(*mountedLabel); ok {
		// Any label immediately preceding any other control will be assumed to
		// be 'associated'.
		return 5 * DIP
	}
	if _, ok := previous.(*mountedCheckbox); ok {
		if _, ok := current.(*mountedCheckbox); ok {
			// Any pair of successive checkboxes will be assumed to be in a
			// related group.
			return 7 * DIP
		}
	}

	// The spacing between unrelated controls.  This is also the default space
	// between paragraphs of text.
	return 11 * DIP
}

func distributeVSpace(align MainAxisAlign, childrenCount int, actualHeight int, minHeight int, maxHeight int) (extraGap int, posY int, scale1 Length, scale2 Length) {
	if actualHeight < minHeight {
		panic("not implemented")
	}

	// If there is more space than necessary, then we need to distribute the extra space.
	if actualHeight >= maxHeight {
		switch align {
		case MainStart:
			// No need to do any adjustment.  The algorithm below will lay out
			// controls aligned to the top.
		case MainCenter:
			// Adjust the starting position to align the contents.
			posY += (actualHeight - maxHeight) / 2

		case MainEnd:
			// Adjust the starting position to align the contents.
			posY += actualHeight - maxHeight

		case SpaceAround:
			extraGap = (actualHeight - maxHeight) / (childrenCount + 1)
			posY += extraGap

		case SpaceBetween:
			if childrenCount > 1 {
				extraGap = (actualHeight - maxHeight) / (childrenCount - 1)
			} else {
				// There are no controls between which to put the extra space.
				// The following essentially convert SpaceBetween to SpaceAround
				extraGap = (actualHeight - maxHeight) / (childrenCount + 1)
				posY += extraGap
			}
		}
	}

	// Calculate scaling to use extra vertical space when available
	scale1, scale2 = 0, 1
	if actualHeight > minHeight && maxHeight > minHeight {
		// We are not doing an actual conversion from pixels to DIPs below.
		// However, the two scale factors are used as a ratio, so any
		// scaling would not affect the final result
		scale1, scale2 = Length(actualHeight-minHeight), Length(maxHeight-minHeight)
	}

	return extraGap, posY, scale1, scale2
}

func setBoundsWithAlign(widget MountedWidget, bounds image.Rectangle, align CrossAxisAlign, scale1, scale2 Length) (moveY int) {
	width := bounds.Dx()
	widthDP := FromPixelsX(width)
	min, max := widget.MeasureHeight(widthDP)
	h := (min + (max-min)*scale1/scale2).PixelsY()

	switch align {
	case CrossStart:
		_, maxX := widget.MeasureWidth()
		if newWidth := maxX.PixelsX(); newWidth < width {
			widget.SetBounds(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+newWidth, bounds.Min.Y+h))
		} else {
			widget.SetBounds(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+h))
		}
	case CrossCenter:
		_, maxX := widget.MeasureWidth()
		if newWidth := maxX.PixelsX(); newWidth < width {
			x1 := (bounds.Min.X + bounds.Max.X - newWidth) / 2
			x2 := (bounds.Min.X + bounds.Max.X + newWidth) / 2
			widget.SetBounds(image.Rect(x1, bounds.Min.Y, x2, bounds.Min.Y+h))
		} else {
			widget.SetBounds(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+h))
		}
	case CrossEnd:
		_, maxX := widget.MeasureWidth()
		if newWidth := maxX.PixelsX(); newWidth < width {
			widget.SetBounds(image.Rect(bounds.Max.X-newWidth, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+h))
		} else {
			widget.SetBounds(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+h))
		}
	case Stretch:
		widget.SetBounds(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Min.Y+h))
	}

	return h
}

func subclassWindowProcedure(hWnd win.HWND, oldWindowProc *uintptr, newWindowProc uintptr) {
	// We need a copy of the address of the old window proc when subclassing.
	// Unhandled messages need to be forwarded.
	if *oldWindowProc == 0 {
		*oldWindowProc = win.GetWindowLongPtr(hWnd, win.GWLP_WNDPROC)
	} else {
		// Paranoia.  Windows creates with the same class should have the same
		// window proc set, but just in case we will double check.
		tmp := win.GetWindowLongPtr(hWnd, win.GWLP_WNDPROC)
		if tmp != *oldWindowProc {
			panic("Window procedure does not match.")
		}
	}

	// Subclass the window by setting a new window proc.
	win.SetWindowLongPtr(hWnd, win.GWLP_WNDPROC, newWindowProc)
}

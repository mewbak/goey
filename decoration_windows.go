package goey

import (
	"github.com/lxn/win"
	"image/color"
	"syscall"
	"unsafe"
)

var (
	decoration struct {
		className *uint16
		atom      win.ATOM
	}
)

func init() {
	var err error
	decoration.className, err = syscall.UTF16PtrFromString("GoeyBackground")
	if err != nil {
		panic(err)
	}
}

func (w *Decoration) mount(parent Control) (Element, error) {
	if decoration.atom == 0 {
		var wc win.WNDCLASSEX
		wc.CbSize = uint32(unsafe.Sizeof(wc))
		wc.HInstance = win.GetModuleHandle(nil)
		wc.LpfnWndProc = syscall.NewCallback(wndprocDecoration)
		wc.HCursor = win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
		wc.HbrBackground = win.GetSysColorBrush(win.COLOR_3DFACE)
		wc.LpszClassName = decoration.className

		atom := win.RegisterClassEx(&wc)
		if atom == 0 {
			return nil, syscall.GetLastError()
		}
		decoration.atom = atom
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE)
	hwnd := win.CreateWindowEx(win.WS_EX_CONTROLPARENT, decoration.className, nil, style,
		10, 10, 100, 100,
		parent.hWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}

	retval := &mountedDecoration{
		Control: Control{hwnd},
		fill:    w.Fill,
		stroke:  w.Stroke,
		padding: w.Padding,
		radius:  w.Radius,
		hBrush:  createBrush(w.Fill),
		hPen:    createPen(w.Stroke),
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	if w.Child != nil {
		child, err := w.Child.Mount(Control{hwnd})
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}
		retval.child = child
	}

	return retval, nil
}

type mountedDecoration struct {
	Control
	fill    color.RGBA
	stroke  color.RGBA
	padding Padding
	radius  Length
	hBrush  win.HBRUSH
	hPen    win.HPEN

	child Element
}

func createBrush(clr color.RGBA) win.HBRUSH {
	// This function create a brush for the requested color.
	//
	// If the color is either white or black, then the stock brush is returned.
	// Note that these can safely be passed to DeleteObject, where they will
	// be ignored.  So we can safely interchange calls to GetStockObject and
	// CreateBrushIndirect.

	if clr.A == 0 {
		return win.HBRUSH(win.GetStockObject(win.NULL_BRUSH))
	} else if clr.R == 0 && clr.G == 0 && clr.B == 0 {
		return win.HBRUSH(win.GetStockObject(win.BLACK_BRUSH))
	} else if clr.R == 0xff && clr.G == 0xff && clr.B == 0xff {
		return win.HBRUSH(win.GetStockObject(win.WHITE_BRUSH))
	} else {
		lb := win.LOGBRUSH{
			LbStyle: win.BS_SOLID,
			LbColor: win.COLORREF(uint32(clr.B)<<16 | uint32(clr.G)<<8 | uint32(clr.R)),
		}
		return win.CreateBrushIndirect(&lb)
	}
}

func createPen(clr color.RGBA) win.HPEN {
	// This function create a brush for the requested color.
	//
	// If the color is either white or black, then the stock brush is returned.
	// Note that these can safely be passed to DeleteObject, where they will
	// be ignored.  So we can safely interchange calls to GetStockObject and
	// CreateBrushIndirect.

	if clr.A == 0 {
		return win.HPEN(win.GetStockObject(win.NULL_PEN))
	} else if clr.R == 0 && clr.G == 0 && clr.B == 0 {
		return win.HPEN(win.GetStockObject(win.BLACK_PEN))
	} else if clr.R == 0xff && clr.G == 0xff && clr.B == 0xff {
		return win.HPEN(win.GetStockObject(win.WHITE_PEN))
	} else {
		lb := win.LOGBRUSH{
			LbStyle: win.BS_SOLID,
			LbColor: win.COLORREF(uint32(clr.B)<<16 | uint32(clr.G)<<8 | uint32(clr.R)),
		}
		return win.ExtCreatePen(win.PS_COSMETIC|win.PS_SOLID, 1, &lb, 0, nil)
	}
}

func (w *mountedDecoration) Close() {
	w.child.Close()
	w.Control.Close()
}

func (w *mountedDecoration) MeasureWidth() (Length, Length) {
	if w.child != nil {
		px := FromPixelsX(1)
		min, max := w.child.MeasureWidth()
		min += 2*px + w.padding.Left + w.padding.Right
		max += 2*px + w.padding.Left + w.padding.Right
		return min, max
	}

	return 13 * DIP, 13 * DIP
}

func (w *mountedDecoration) MeasureHeight(width Length) (Length, Length) {
	if w.child != nil {
		px := FromPixelsX(1)
		py := FromPixelsY(1)
		min, max := w.child.MeasureHeight(width - 2*px)
		min += 2*py + w.padding.Top + w.padding.Bottom
		max += 2*py + w.padding.Top + w.padding.Bottom
		return min, max
	}

	return 13 * DIP, 13 * DIP
}

func (w *mountedDecoration) SetBounds(bounds Rectangle) {
	// Update background control position
	w.Control.SetBounds(bounds)

	px := FromPixelsX(1)
	py := FromPixelsY(1)
	position := bounds.Min
	bounds.Min.X += px + w.padding.Left - position.X
	bounds.Min.Y += py + w.padding.Top - position.Y
	bounds.Max.X -= px + w.padding.Right + position.X
	bounds.Max.Y -= py + w.padding.Bottom + position.Y
	w.child.SetBounds(bounds)
}

func (w *mountedDecoration) SetOrder(previous win.HWND) win.HWND {
	previous = w.Control.SetOrder(previous)
	previous = w.child.SetOrder(previous)
	return previous
}

func (w *mountedDecoration) updateProps(data *Decoration) error {
	if w.fill != data.Fill {
		// Free the old brush
		if w.hBrush != 0 {
			win.DeleteObject(win.HGDIOBJ(w.hBrush))
		}

		// Allocate the new brush
		w.hBrush = createBrush(data.Fill)
		if w.hBrush == 0 {
			return syscall.GetLastError()
		}
		w.fill = data.Fill
	}

	if w.stroke != data.Stroke {
		if w.hPen != 0 {
			win.DeleteObject(win.HGDIOBJ(w.hPen))
		}

		w.hPen = createPen(data.Stroke)
		if w.hPen == 0 {
			return syscall.GetLastError()
		}
		w.stroke = data.Stroke
	}

	w.padding = data.Padding
	w.radius = data.Radius

	child, err := DiffChild(Control{w.hWnd}, w.child, data.Child)
	if err != nil {
		return err
	}
	w.child = child

	return nil
}

func wndprocDecoration(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_PAINT:
		// Fill with the proper background color
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedDecoration)(unsafe.Pointer(w))

			ps := win.PAINTSTRUCT{}
			cr := win.RECT{}
			win.GetClientRect(hwnd, &cr)
			hdc := win.BeginPaint(hwnd, &ps)
			win.SelectObject(hdc, win.HGDIOBJ(ptr.hBrush))
			win.SelectObject(hdc, win.HGDIOBJ(ptr.hPen))
			if ptr.radius > 0 {
				rx := ptr.radius.PixelsX()
				ry := ptr.radius.PixelsY()
				win.RoundRect(hdc, cr.Left, cr.Top, cr.Right, cr.Bottom, int32(rx), int32(ry))
			} else {
				win.Rectangle_(hdc, cr.Left, cr.Top, cr.Right, cr.Bottom)
			}
			win.EndPaint(hwnd, &ps)
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

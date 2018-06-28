package goey

import (
	win2 "bitbucket.org/rj/goey/syscall"
	"github.com/lxn/win"
	"syscall"
	"time"
	"unsafe"
)

var (
	datetimepickClassName     *uint16
	oldDateTimePickWindowProc uintptr
)

func init() {
	var err error
	datetimepickClassName, err = syscall.UTF16PtrFromString("SysDateTimePick32")
	if err != nil {
		panic(err)
	}
}

func (w *DateInput) systemTime() win.SYSTEMTIME {
	return win.SYSTEMTIME{
		WYear:   uint16(w.Value.Year()),
		WMonth:  uint16(w.Value.Month()),
		WDay:    uint16(w.Value.Day()),
		WHour:   uint16(w.Value.Hour()),
		WMinute: uint16(w.Value.Minute()),
		WSecond: uint16(w.Value.Second()),
	}
}

func (w *DateInput) mount(parent Control) (Element, error) {
	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP)
	hwnd := win.CreateWindowEx(0, datetimepickClassName, nil,
		style,
		10, 10, 100, 100,
		parent.hWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}

	// Set the font for the window
	if hMessageFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, uintptr(hMessageFont), 0)
	}

	// Set the properties for the control
	st := w.systemTime()
	win.SendMessage(hwnd, win.DTM_SETSYSTEMTIME, win.GDT_VALID, uintptr(unsafe.Pointer(&st)))
	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &oldDateTimePickWindowProc, syscall.NewCallback(dateinputWindowProc))

	retval := &mountedDateInput{
		Control:  Control{hwnd},
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedDateInput struct {
	Control
	onChange func(value time.Time)
	onFocus  func()
	onBlur   func()
}

func (w *mountedDateInput) Layout(bc Constraint) Size {
	// Determine ideal width.
	width := 75 * DIP
	height := 23 * DIP
	return bc.Constrain(Size{width, height})
}

func (w *mountedDateInput) MinimumSize() Size {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	// Unclear what the correct width should be.  Using button for the moment
	// Height set to match a text box.
	return Size{75 * DIP, 23 * DIP}
}

func (w *mountedDateInput) updateProps(data *DateInput) error {
	st := data.systemTime()
	win.SendMessage(w.hWnd, win.DTM_SETSYSTEMTIME, win.GDT_VALID, uintptr(unsafe.Pointer(&st)))

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	return nil
}

func dateinputWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		dateinputGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := dateinputGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := dateinputGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_NOTIFY:
		switch code := (*win.NMHDR)(unsafe.Pointer(lParam)).Code; code {
		case win.DTN_DATETIMECHANGE:
			if w := dateinputGetPtr(hwnd); w.onChange != nil {
				nmhdr := (*win.NMDATETIMECHANGE)(unsafe.Pointer(lParam))
				st := time.Date(int(nmhdr.St.WYear), time.Month(nmhdr.St.WMonth), int(nmhdr.St.WDay), int(nmhdr.St.WHour), int(nmhdr.St.WMinute), int(nmhdr.St.WSecond), 0, time.Local)
				w.onChange(st)
			}

		case win2.MCN_SELECT:
			nmhdr := (*win2.NMSELCHANGE)(unsafe.Pointer(lParam))
			win.SendMessage(hwnd, win.DTM_SETSYSTEMTIME, win.GDT_VALID, uintptr(unsafe.Pointer(&nmhdr.StSelStart)))
			win.SendMessage(hwnd, win2.DTM_CLOSEMONTHCAL, 0, 0)
		}
		return 0

	}

	return win.CallWindowProc(oldDateTimePickWindowProc, hwnd, msg, wParam, lParam)
}

func dateinputGetPtr(hwnd win.HWND) *mountedDateInput {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*mountedDateInput)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd {
		panic("Internal error.")
	}

	return ptr
}

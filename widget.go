package goey

// Kind identifies the different kind of widgets.  Most widgets have two
// concrete types associated with their behaviour.  First, there is a type with data
// to describe the widget when unmounted.  Second, there is a type with a handle
// to the windowing system when mounted.  Automatic reconciliation of two widget
// trees relies on Kind to match the unmounted and mounted widgets.
type Kind struct {
	name string
}

// Widget is an interface that wraps any type describing a GUI widget or control.
// A Widget can be 'mounted' to instantiate a widget or control in the GUI.
type Widget interface {
	// Kind returns the concrete type's Kind.  The returned value should
	// be constant, and the same for all instances of a concrete type.
	// Users should not need to use this method directly.
	Kind() *Kind
	// Mount creates a widget or control in the GUI.  The newly created widget
	// will be a child of the widget specified by parent.  If non-nil, the returned
	// Element must have a matching kind.
	Mount(parent NativeWidget) (Element, error)
}

// Element is an interface that wrap any type representing an existing
// widget or control in the GUI.
type Element interface {
	// NativeMountedWidget provides platform dependent methods.  These should
	// not be used by client libraries, but exist for the internal implementation
	// of platform dependent code.
	NativeMountedWidget

	// Close removes the widget from the GUI, and frees any associated resources.
	Close()
	// Kind returns the concrete type for the Element.
	// Users should not need to use this method directly.
	Kind() *Kind
	// MeasureWidth returns the minimum and maximum comfortable widths for
	// this widget.
	MeasureWidth() (min Length, max Length)
	// MeasureHeight returns the minimum and maximum comfortable heights for
	// this widget.
	MeasureHeight(width Length) (min Length, max Length)
	// SetBounds updates the position of the widget.
	SetBounds(bounds Rectangle)
	// UpdateProps will update the properties of the widget.  The Kind for
	// the parameter data must match the Kind for the interface.
	UpdateProps(data Widget) error
}

func calculateHGap(previous Element, current Element) Length {
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

func calculateVGap(previous Element, current Element) Length {
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

func distributeVSpace(align MainAxisAlign, childrenCount int, actualHeight Length, minHeight Length, maxHeight Length) (extraGap Length, posY Length, scale1 int, scale2 int) {
	if actualHeight < minHeight {
		println("actualHeight:", actualHeight.String())
		println("minHeight:", minHeight.String())
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
			extraGap = (actualHeight - maxHeight).Scale(1, childrenCount+1)
			posY += extraGap

		case SpaceBetween:
			if childrenCount > 1 {
				extraGap = (actualHeight - maxHeight).Scale(1, childrenCount-1)
			} else {
				// There are no controls between which to put the extra space.
				// The following essentially convert SpaceBetween to SpaceAround
				extraGap = (actualHeight - maxHeight).Scale(1, childrenCount+1)
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
		scale1, scale2 = int(actualHeight-minHeight), int(maxHeight-minHeight)
	}

	return extraGap, posY, scale1, scale2
}

func setBoundsWithAlign(widget Element, bounds Rectangle, align CrossAxisAlign, scale1, scale2 int) (moveY Length) {
	width := bounds.Dx()
	min, max := widget.MeasureHeight(width)
	h := min + (max-min).Scale(scale1, scale2)

	switch align {
	case CrossStart:
		_, maxX := widget.MeasureWidth()
		if maxX < width {
			widget.SetBounds(Rectangle{bounds.Min, Point{bounds.Min.X + maxX, bounds.Min.Y + h}})
		} else {
			widget.SetBounds(Rectangle{bounds.Min, Point{bounds.Max.X, bounds.Min.Y + h}})
		}
	case CrossCenter:
		_, maxX := widget.MeasureWidth()
		if maxX < width {
			x1 := (bounds.Min.X + bounds.Max.X - maxX) / 2
			x2 := (bounds.Min.X + bounds.Max.X + maxX) / 2
			widget.SetBounds(Rectangle{Point{x1, bounds.Min.Y}, Point{x2, bounds.Min.Y + h}})
		} else {
			widget.SetBounds(Rectangle{bounds.Min, Point{bounds.Max.X, bounds.Min.Y + h}})
		}
	case CrossEnd:
		_, maxX := widget.MeasureWidth()
		if maxX < width {
			widget.SetBounds(Rectangle{Point{bounds.Max.X - maxX, bounds.Min.Y}, Point{bounds.Max.X, bounds.Min.Y + h}})
		} else {
			widget.SetBounds(Rectangle{bounds.Min, Point{bounds.Max.X, bounds.Min.Y + h}})
		}
	case Stretch:
		widget.SetBounds(Rectangle{bounds.Min, Point{bounds.Max.X, bounds.Min.Y + h}})
	}

	return h
}

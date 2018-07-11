package goey

// Kind identifies the different kind of widgets.  Most widgets have two
// concrete types associated with their behaviour.  First, there is a type with data
// to describe the widget when unmounted.  Second, there is a type with a handle
// to the windowing system when mounted.  Automatic reconciliation of two widget
// trees relies on Kind to match the unmounted and mounted widgets.
//
// Note that comparison of kinds is done by address, and not done using the value of any fields.
// Any internal state is simply to help with debugging.
type Kind struct {
	name string
}

// NewKind creates a new kind.  The name should identify the type used for the widget,
// but is currently unused.
func NewKind(name string) Kind {
	return Kind{name}
}

// String returns the string with the name of the widget and element kind.
func (k Kind) String() string {
	return k.name
}

// Widget is an interface that wraps any type describing part of a GUI.
// A widget can be 'mounted' to create controls using the platform GUI.
type Widget interface {
	// Kind returns the concrete type's Kind.  The returned value should
	// be constant, and the same for all instances of a concrete type.
	// Users should not need to use this method directly.
	Kind() *Kind
	// Mount creates a widget or control in the GUI.  The newly created widget
	// will be a child of the widget specified by parent.  If non-nil, the returned
	// Element must have a matching kind.
	Mount(parent Control) (Element, error)
}

// Element is an interface that wraps any type representing a control, or group
// of controls, created using the platform GUI.
type Element interface {
	// NativeElement provides platform dependent methods.  These should
	// not be used by client libraries, but exist for the internal implementation
	// of platform dependent code.
	NativeElement

	// Close removes the widget from the GUI, and frees any associated resources.
	Close()
	// Kind returns the concrete type for the Element.
	// Users should not need to use this method directly.
	Kind() *Kind
	// Layout determines the best size for an element that sastisfies the
	// constraints.
	Layout(Constraint) Size
	// MinIntrinsicHeight returns the minimum height that this element requires
	// to be correctly displayed.
	MinIntrinsicHeight(width Length) Length
	// MinIntrinsicWidth returns the minimum width that this element requires
	// to be correctly displayed.
	MinIntrinsicWidth(height Length) Length
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
	if _, ok := previous.(*buttonElement); ok {
		if _, ok := current.(*buttonElement); ok {
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
	if _, ok := previous.(*labelElement); ok {
		// Any label immediately preceding any other control will be assumed to
		// be 'associated'.
		return 5 * DIP
	}
	if _, ok := previous.(*checkboxElement); ok {
		if _, ok := current.(*checkboxElement); ok {
			// Any pair of successive checkboxes will be assumed to be in a
			// related group.
			return 7 * DIP
		}
	}

	// The spacing between unrelated controls.  This is also the default space
	// between paragraphs of text.
	return 11 * DIP
}

package base

const (
	// Inf is a sentinal value indicating an unbounded width or height.
	Inf Length = 0x7fffffff
)

func guardInf(a, b Length) Length {
	if a == Inf {
		return Inf
	}
	return b
}

func max(a, b Length) Length {
	if a > b {
		return a
	}
	return b
}

// Constraints represents box constraints on width and height for the layout of
// rectangular widgets.  For each dimension, the constraints specify the
// minimum and maximum allowed size for a widget.
//
// The constraints in a dimension are called 'tight' if the minimum and
// maximum values are equal, which essential requires the widget to take a
// fixed size. On the other hand, if the minimum allowed value is zero, then
// the constraints on that dimension is 'loose'.
//
// A sentinal value can be used to indicate that the maximum size in a
// dimension is infinite.  The constraints on that dimension are called
// 'unbounded'.  If the constraints for a dimension are both tight and unbounded,
// the dimension is 'expanding'.
//
// (This type is similar to BoxConstraints type in flutter library rendering)
type Constraints struct {
	Min, Max Size
}

// Expand creates box constraints that force elements to expand to as large as
// possible.  The constraints for both width and height will be tight and
// unbounded.
func Expand() Constraints {
	return Constraints{Size{Inf, Inf}, Size{Inf, Inf}}
}

// ExpandHeight creates box constraints with a fixed width and that forces
// elements to expand to as high as possible.
func ExpandHeight(width Length) Constraints {
	return Constraints{Size{width, Inf}, Size{width, Inf}}
}

// ExpandWidth creates box constraints with a fixed height and that forces
// elements to expand to as wide as possible.
func ExpandWidth(height Length) Constraints {
	return Constraints{Size{Inf, height}, Size{Inf, height}}
}

// Loose creates box constraints that forbid sizes larger than the given size.
// The constraints for both width and height will be loose and bounded.
func Loose(size Size) Constraints {
	return Constraints{Size{}, size}
}

// Tight creates a box constraints that is respected only by the given size.
func Tight(size Size) Constraints {
	return Constraints{size, size}
}

// TightWidth creates a box constraints that is respected only by sizes with
// the given width.  The height is unconstrained (i.e. loose and unbounded).
func TightWidth(width Length) Constraints {
	return Constraints{Size{width, 0}, Size{width, Inf}}
}

// TightHeight creates a box constraints that is respected only by sizes with
// the given height.  The width is unconstrained (i.e. loose and unbounded).
func TightHeight(height Length) Constraints {
	return Constraints{Size{0, height}, Size{Inf, height}}
}

// Constrain returns the size that satisfies the constraints while staying as
// close as possible to the passed size.
func (bc Constraints) Constrain(size Size) Size {
	return Size{
		Width:  size.Width.Clamp(bc.Min.Width, bc.Max.Width),
		Height: size.Height.Clamp(bc.Min.Height, bc.Max.Height),
	}
}

// ConstrainAndAttemptToPreserveAspectRatio returns the size that satisfies the
// constraints while staying close to the passed size and maintaining the aspect
// ratio of the passed size.
func (bc Constraints) ConstrainAndAttemptToPreserveAspectRatio(size Size) Size {
	if bc.IsTight() {
		return bc.Min
	}

	width := size.Width
	height := size.Height

	if width > bc.Max.Width {
		width = bc.Max.Width
		height = width.Scale(int(size.Height), int(size.Width))
	}

	if height > bc.Max.Height {
		height = bc.Max.Height
		width = height.Scale(int(size.Width), int(size.Height))
	}

	if width < bc.Min.Width {
		width = bc.Min.Width
		height = width.Scale(int(size.Height), int(size.Width))
	}

	if height < bc.Min.Height {
		height = bc.Min.Height
		width = height.Scale(int(size.Width), int(size.Height))
	}

	return bc.Constrain(Size{width, height})
}

// ConstrainHeight returns the length that satisfies the constraints for height
// while staying as close as possible to the passed height.
func (bc Constraints) ConstrainHeight(height Length) Length {
	return height.Clamp(bc.Min.Height, bc.Max.Height)
}

// ConstrainWidth returns the length that satisfies the constraints for width
// while staying as close as possible to the passed height.
func (bc Constraints) ConstrainWidth(width Length) Length {
	return width.Clamp(bc.Min.Width, bc.Max.Width)
}

// Enforce returns new box constraints that respect the constraints,
// while respecting the constraints of the method receiver as closely as possible.
/*func (bc Constraints) Enforce(constraints Constraints) Constraints {
	minWidth := bc.Min.Width.Clamp(constraints.Min.Width, constraints.Max.Width)
	maxWidth := bc.Max.Width.Clamp(constraints.Min.Width, constraints.Max.Width)
	minHeight := bc.Min.Height.Clamp(constraints.Min.Height, constraints.Max.Height)
	maxHeight := bc.Max.Height.Clamp(constraints.Min.Height, constraints.Max.Height)
	return Constraints{Size{minWidth, minHeight}, Size{maxWidth, maxHeight}}
}*/

// HasBoundedHeight is true if the maximum height is bounded.
func (bc Constraints) HasBoundedHeight() bool {
	return bc.Max.Height < Inf
}

// HasBoundedWidth is true if the maximum width is bounded.
func (bc Constraints) HasBoundedWidth() bool {
	return bc.Max.Width < Inf
}

// HasTightWidth is true if the width is tight (only one value of width
// satisfies the constraints).
func (bc Constraints) HasTightWidth() bool {
	return bc.Min.Width >= bc.Max.Width
}

// HasTightHeight is true if the height is tight (only one value of height
// satisfies the constraints).
func (bc Constraints) HasTightHeight() bool {
	return bc.Min.Height >= bc.Max.Height
}

// Inset returns a new set of box constraints such that a size that satisfies
// those new constraints can be increased by width and height and will satisfy
// the original constrains.
func (bc Constraints) Inset(width Length, height Length) Constraints {
	deflatedMinWidth := guardInf(bc.Min.Width, max(0, bc.Min.Width-width))
	deflatedMinHeight := guardInf(bc.Min.Height, max(0, bc.Min.Height-height))

	return Constraints{
		Size{deflatedMinWidth, deflatedMinHeight},
		Size{
			max(deflatedMinWidth, guardInf(bc.Max.Width, bc.Max.Width-width)),
			max(deflatedMinHeight, guardInf(bc.Max.Height, bc.Max.Height-height)),
		},
	}
}

// IsBounded is true if both the width and height are bounded.
func (bc Constraints) IsBounded() bool {
	return bc.HasBoundedWidth() && bc.HasBoundedHeight()
}

// IsNormalized is true if both the width and height constraints are normalized.
// A set of constraints are normalized if 0 < Min < Max.
func (bc Constraints) IsNormalized() bool {
	return bc.Min.Width >= 0.0 &&
		bc.Min.Width <= bc.Max.Width &&
		bc.Min.Height >= 0.0 &&
		bc.Min.Height <= bc.Max.Height
}

// IsSatisfiedBy returns true if the passed size satisfies the both the width
// and height constraints.
func (bc Constraints) IsSatisfiedBy(size Size) bool {
	return bc.Min.Width <= size.Width &&
		size.Width <= bc.Max.Width &&
		bc.Min.Height <= size.Height &&
		size.Height <= bc.Max.Height
}

// IsTight returns true if both the width and height are tightly constrained.
func (bc Constraints) IsTight() bool {
	return bc.HasTightWidth() && bc.HasTightHeight()
}

// IsZero returns true if the constrain is the zero value.
func (bc Constraints) IsZero() bool {
	return bc.Min.Width == 0 && bc.Min.Height == 0 && bc.Max.Width == 0 && bc.Max.Height == 0
}

// Loosen creates a new box constraint with the minimum width and height
// requirements removed.
func (bc Constraints) Loosen() Constraints {
	return Constraints{Size{}, bc.Max}
}

// LoosenHeight creates a new box constraint with the minimum height
// requirements removed.
func (bc Constraints) LoosenHeight() Constraints {
	return Constraints{Size{bc.Min.Width, 0}, bc.Max}
}

// LoosenWidth creates a new box constraint with the minimum width
// requirements removed.
func (bc Constraints) LoosenWidth() Constraints {
	return Constraints{Size{0, bc.Min.Height}, bc.Max}
}

// Tighten creates a new box constraint with tight width and height
// requirements matching as closely as possible the passed size.
// The new constrains will be tight, but will only match the requested size if
// the size satisfies the original constraints.
func (bc Constraints) Tighten(size Size) Constraints {
	bc.Min.Width = size.Width.Clamp(bc.Min.Width, bc.Max.Width)
	bc.Max.Width = bc.Min.Width
	bc.Min.Height = size.Height.Clamp(bc.Min.Height, bc.Max.Height)
	bc.Max.Height = bc.Min.Height
	return bc
}

// TightenHeight creates a new box constraint with a tight height
// requirements matching as closely as possible the length.
// The new height constraints will be tight, but will only match the requested
// height if the height statisfies the original constraints.
func (bc Constraints) TightenHeight(height Length) Constraints {
	bc.Min.Height = height.Clamp(bc.Min.Height, bc.Max.Height)
	bc.Max.Height = bc.Min.Height
	return bc
}

// TightenWidth creates a new box constraint with a tight width
// requirements matching as closely as possible the length.
// The new width constraints will be tight, but will only match the requested
// width if the width statisfies the original constraints.
func (bc Constraints) TightenWidth(width Length) Constraints {
	bc.Min.Width = width.Clamp(bc.Min.Width, bc.Max.Width)
	bc.Max.Width = bc.Min.Width
	return bc
}

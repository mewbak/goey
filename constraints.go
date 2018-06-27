package goey

const (
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

// Size represents the size of a rectangular element.
type Size struct {
	Width, Height Length
}

// IsZero returns true if the constrain is the zero value.
func (s *Size) IsZero() bool {
	return s.Width == 0 && s.Height == 0
}

// String returns a string representation of the size.
func (s *Size) String() string {
	return "(" + s.Width.String() + "x" + s.Height.String() + ")"
}

// Constraint represents box constraints on width and height for the layout of
// rectangular widgets.
type Constraint struct {
	Min, Max Size
}

// Expand creates box constraints that force elements to expand to as large as
// possible.
func Expand() Constraint {
	return Constraint{Size{Inf, Inf}, Size{Inf, Inf}}
}

// ExpandHeight creates box constraints with a fixed width and that force
// elements to expand to as high as possible.
func ExpandHeight(width Length) Constraint {
	return Constraint{Size{width, Inf}, Size{width, Inf}}
}

// ExpandWidth creates box constraints with a fixed height and that force
// elements to expand to as wide as possible.
func ExpandWidth(height Length) Constraint {
	return Constraint{Size{Inf, height}, Size{Inf, height}}
}

// Loose creates box constraints that forbid sizes larger than the given size.
func Loose(size Size) Constraint {
	return Constraint{Size{}, size}
}

// Tight creates a box constraints that is respected only by the given size.
func Tight(size Size) Constraint {
	return Constraint{size, size}
}

// TightWidth creates a box constraints that is respected only by sizes with
// the given width.  The height is unconstrained.
func TightWidth(width Length) Constraint {
	return Constraint{Size{width, 0}, Size{width, Inf}}
}

// TightHeight creates a box constraints that is respected only by sizes with
// the given height.  The width is unconstrained.
func TightHeight(height Length) Constraint {
	return Constraint{Size{0, height}, Size{Inf, height}}
}

// Constrain returns the size that satisfies the constraints while staying as
// close as possible to the passed size.
func (bc Constraint) Constrain(size Size) Size {
	return Size{
		Width:  size.Width.Clamp(bc.Min.Width, bc.Max.Width),
		Height: size.Height.Clamp(bc.Min.Height, bc.Max.Height),
	}
}

// ConstrainAndAttemptToPreserveAspectRatio returns the size that satisfies the
// constraints while staying close to the passed size and maintaining the aspect
// ratio of the passed size.
func (bc Constraint) ConstrainAndAttemptToPreserveAspectRatio(size Size) Size {
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
func (bc Constraint) ConstrainHeight(height Length) Length {
	return height.Clamp(bc.Min.Height, bc.Max.Height)
}

// ConstrainWidth returns the length that satisfies the constraints for width
// while staying as close as possible to the passed height.
func (bc Constraint) ConstrainWidth(width Length) Length {
	return width.Clamp(bc.Min.Width, bc.Max.Width)
}

// Enforce returns new box constraints that respect the constraints,
// while respecting the constraints of the method receiver as closely as possible.
func (bc Constraint) Enforce(constraints Constraint) Constraint {
	minWidth := bc.Min.Width.Clamp(constraints.Min.Width, constraints.Max.Width)
	maxWidth := bc.Max.Width.Clamp(constraints.Min.Width, constraints.Max.Width)
	minHeight := bc.Min.Height.Clamp(constraints.Min.Height, constraints.Max.Height)
	maxHeight := bc.Max.Height.Clamp(constraints.Min.Height, constraints.Max.Height)
	return Constraint{Size{minWidth, minHeight}, Size{maxWidth, maxHeight}}
}

// HasBoundedHeight is true if the maximum height is bounded.
func (bc Constraint) HasBoundedHeight() bool {
	return bc.Max.Height < Inf
}

// HasBoundedWidth is true if the maximum width is bounded.
func (bc Constraint) HasBoundedWidth() bool {
	return bc.Max.Width < Inf
}

// HasTightWidth is true if the width is tight (only one value of width
// satisfies the constraints).
func (bc Constraint) HasTightWidth() bool {
	return bc.Min.Width >= bc.Max.Width
}

// HasTightHeight is true if the height is tight (only one value of height
// satisfies the constraints).
func (bc Constraint) HasTightHeight() bool {
	return bc.Min.Height >= bc.Max.Height
}

// Inset returns a new set of box constraints such that a size that satisfies
// those new constraints can be increased by width and height and will satisfy
// the original constrains.
func (bc Constraint) Inset(width Length, height Length) Constraint {
	deflatedMinWidth := guardInf(bc.Min.Width, max(0, bc.Min.Width-width))
	deflatedMinHeight := guardInf(bc.Min.Height, max(0, bc.Min.Height-height))

	return Constraint{
		Size{deflatedMinWidth, deflatedMinHeight},
		Size{
			max(deflatedMinWidth, guardInf(bc.Max.Width, bc.Max.Width-width)),
			max(deflatedMinHeight, guardInf(bc.Max.Height, bc.Max.Height-height)),
		},
	}
}

// IsBounded is true if both the width and height are bounded.
func (bc Constraint) IsBounded() bool {
	return bc.HasBoundedWidth() && bc.HasBoundedHeight()
}

// IsNormalized is true if both the width and height constraints are normalized.
// A set of constraints are normalized if 0 < Min < Max.
func (bc Constraint) IsNormalized() bool {
	return bc.Min.Width >= 0.0 &&
		bc.Min.Width <= bc.Max.Width &&
		bc.Min.Height >= 0.0 &&
		bc.Min.Height <= bc.Max.Height
}

// IsSatisfiedBy returns true if the passed size satisfies the both the width
// and height constraints.
func (bc Constraint) IsSatisfiedBy(size Size) bool {
	return bc.Min.Width <= size.Width &&
		size.Width <= bc.Max.Width &&
		bc.Min.Height <= size.Height &&
		size.Height <= bc.Max.Height
}

// IsTight returns true if both the width and height are tightly constrained.
func (bc Constraint) IsTight() bool {
	return bc.HasTightWidth() && bc.HasTightHeight()
}

// IsZero returns true if the constrain is the zero value.
func (bc Constraint) IsZero() bool {
	return bc.Min.Width == 0 && bc.Min.Height == 0 && bc.Max.Width == 0 && bc.Max.Height == 0
}

// Loosen creates a new box constraint with the minimum width and height
// requirements removed.
func (bc Constraint) Loosen() Constraint {
	return Constraint{Size{}, bc.Max}
}

// LoosenHeight creates a new box constraint with the minimum height
// requirements removed.
func (bc Constraint) LoosenHeight() Constraint {
	return Constraint{Size{bc.Min.Width, 0}, bc.Max}
}

// LoosenWidth creates a new box constraint with the minimum width
// requirements removed.
func (bc Constraint) LoosenWidth() Constraint {
	return Constraint{Size{0, bc.Min.Height}, bc.Max}
}

// Tighten creates a new box constraint with tight width and height
// requirements matching as closely as possible the passed size.
func (bc Constraint) Tighten(size Size) Constraint {
	bc.Min.Width = size.Width.Clamp(bc.Min.Width, bc.Max.Width)
	bc.Max.Width = bc.Min.Width
	bc.Min.Height = size.Height.Clamp(bc.Min.Height, bc.Max.Height)
	bc.Max.Height = bc.Min.Height
	return bc
}

// TightenHeight creates a new box constraint with a tight height
// requirements matching as closely as possible the length.
func (bc Constraint) TightenHeight(height Length) Constraint {
	bc.Min.Height = height.Clamp(bc.Min.Height, bc.Max.Height)
	bc.Max.Height = bc.Min.Height
	return bc
}

// TightenWidth creates a new box constraint with a tight width
// requirements matching as closely as possible the length.
func (bc Constraint) TightenWidth(width Length) Constraint {
	bc.Min.Width = width.Clamp(bc.Min.Width, bc.Max.Width)
	bc.Max.Width = bc.Min.Width
	return bc
}

package goey

const (
	Inf Length = 0x7fffffff
)

func clamp(min, value, max Length) Length {
	if value > max {
		value = max
	}
	if value < min {
		value = min
	}
	return value
}

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

type Size struct {
	Width, Height Length
}

func (s *Size) IsZero() bool {
	return s.Width == 0 && s.Height == 0
}

func (s *Size) String() string {
	return "(" + s.Width.String() + "x" + s.Height.String() + ")"
}

type Box struct {
	Min, Max Size
}

func Expand() Box {
	return Box{Size{Inf, Inf}, Size{Inf, Inf}}
}

func ExpandWidth(width Length) Box {
	return Box{Size{width, Inf}, Size{width, Inf}}
}

func ExpandHeight(height Length) Box {
	return Box{Size{Inf, height}, Size{Inf, height}}
}

func Loose(size Size) Box {
	return Box{Size{}, size}
}

func Tight(size Size) Box {
	return Box{size, size}
}

func TightWidth(width Length) Box {
	return Box{Size{width, 0}, Size{width, Inf}}
}

func TightHeight(height Length) Box {
	return Box{Size{0, height}, Size{Inf, height}}
}

func (bc Box) Constrain(size Size) Size {
	return Size{
		Width:  clamp(bc.Min.Width, size.Width, bc.Max.Width),
		Height: clamp(bc.Min.Height, size.Height, bc.Max.Height),
	}
}

func (bc Box) ConstrainAndAttemptToPreserveAspectRatio(size Size) Size {
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

func (bc Box) ConstrainHeight(height Length) Length {
	return clamp(bc.Min.Height, height, bc.Max.Height)
}

func (bc Box) ConstrainWidth(width Length) Length {
	return clamp(bc.Min.Width, width, bc.Max.Width)
}

func (bc Box) Deflate(width Length, height Length) Box {
	deflatedMinWidth := guardInf(bc.Min.Width, max(0, bc.Min.Width-width))
	deflatedMinHeight := guardInf(bc.Min.Height, max(0, bc.Min.Height-height))

	return Box{
		Size{deflatedMinWidth, deflatedMinHeight},
		Size{
			max(deflatedMinWidth, guardInf(bc.Max.Width, bc.Max.Width-width)),
			max(deflatedMinHeight, guardInf(bc.Max.Height, bc.Max.Height-height)),
		},
	}
}

func (bc Box) Enforce(constraints Box) Box {
	minWidth := clamp(constraints.Min.Width, bc.Min.Width, constraints.Max.Width)
	maxWidth := clamp(constraints.Min.Width, bc.Max.Width, constraints.Max.Width)
	minHeight := clamp(constraints.Min.Height, bc.Min.Height, constraints.Max.Height)
	maxHeight := clamp(constraints.Min.Height, bc.Max.Height, constraints.Max.Height)
	return Box{Size{minWidth, minHeight}, Size{maxWidth, maxHeight}}
}

func (bc Box) HasBoundedHeight() bool {
	return bc.Max.Height < Inf
}

func (bc Box) HasBoundedWidth() bool {
	return bc.Max.Width < Inf
}

func (bc Box) HasTightWidth() bool {
	return bc.Min.Width >= bc.Max.Width
}

func (bc Box) HasTightHeight() bool {
	return bc.Min.Height >= bc.Max.Height
}

func (bc Box) Inset(padding Length) Box {
	minWidth := bc.Min.Width - 2*padding
	if minWidth < 0 {
		minWidth = 0
	}
	minHeight := bc.Min.Height - 2*padding
	if minHeight < 0 {
		minHeight = 0
	}
	maxWidth := bc.Max.Width
	if maxWidth != Inf {
		maxWidth = bc.Max.Width - 2*padding
		if maxWidth < minWidth {
			maxWidth = minWidth
		}
	}
	maxHeight := bc.Max.Height
	if maxHeight != Inf {
		maxHeight = bc.Max.Height - 2*padding
		if maxHeight < minHeight {
			maxHeight = minHeight
		}
	}

	return Box{Size{minWidth, minHeight}, Size{maxWidth, maxHeight}}
}

func (bc Box) IsBounded() bool {
	return bc.HasBoundedWidth() && bc.HasBoundedHeight()
}

func (bc Box) IsNormalized() bool {
	return bc.Min.Width >= 0.0 &&
		bc.Min.Width <= bc.Max.Width &&
		bc.Min.Height >= 0.0 &&
		bc.Min.Height <= bc.Max.Height
}

func (bc Box) IsSatisfiedBy(size Size) bool {
	return bc.Min.Width <= size.Width &&
		size.Width <= bc.Max.Width &&
		bc.Min.Height <= size.Height &&
		size.Height <= bc.Max.Height
}

func (bc Box) IsTight() bool {
	return bc.HasTightWidth() && bc.HasTightHeight()
}

func (bc Box) IsZero() bool {
	return bc.Min.Width == 0 && bc.Min.Height == 0 && bc.Max.Width == 0 && bc.Max.Height == 0
}

func (bc Box) Loosen() Box {
	return Box{Size{}, bc.Max}
}

func (bc Box) LoosenHeight() Box {
	return Box{Size{bc.Min.Width, 0}, bc.Max}
}

func (bc Box) LoosenWidth() Box {
	return Box{Size{0, bc.Min.Height}, bc.Max}
}

func (bc Box) Tighten(size Size) Box {
	bc.Min.Width = clamp(bc.Min.Width, size.Width, bc.Max.Width)
	bc.Max.Width = bc.Min.Width
	bc.Min.Height = clamp(bc.Min.Height, size.Height, bc.Max.Height)
	bc.Max.Height = bc.Min.Height
	return bc
}

func (bc Box) TightenHeight(height Length) Box {
	bc.Min.Height = clamp(bc.Min.Height, height, bc.Max.Height)
	bc.Max.Height = bc.Min.Height
	return bc
}

func (bc Box) TightenWidth(width Length) Box {
	bc.Min.Width = clamp(bc.Min.Width, width, bc.Max.Width)
	bc.Max.Width = bc.Min.Width
	return bc
}

func (bc Box) OrValue(other Box) Box {
	if bc.IsZero() {
		return other
	}
	return bc
}

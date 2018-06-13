package goey

import (
	"image"

	"golang.org/x/image/math/fixed"
)

var (
	// DPI contains the current DPI (dots per inch) of the monitor.
	// User code should not need to set this directly, as drivers will update
	// this variable as necessary.
	DPI image.Point
)

// Common lengths used when describing GUIs.  Note that the DIP is the natural
// unit for this package.  Because of precision, the PT listed here is somewhat
// smaller than its correct value.
const (
	DIP = Length(1 << 6)
	PT  = Length((96 << 6) / 72)
)

// Length is a distance measured in device-independent pixels.  There are nominally
// 96 DIPs per inch.  This definition corresponds with the definition of a
// pixel for both CSS and on Windows.
type Length fixed.Int26_6

func (dp Length) DIP() float64 {
	return float64(dp) / (1 << 6)
}

func (dp Length) PT() float64 {
	return float64(dp) / ((96 << 6) / 72)
}

// PixelsX converts the distance measurement in DIPs to physical pixels, based
// on the current DPI settings for horizontal scaling.
func (dp Length) PixelsX() int {
	return fixed.Int26_6(dp.Scale(DPI.X, 96)).Round()
}

// PixelsY converts the distance measurement in DIPs to physical pixels, based
// on the current DPI settings for vertical scaling.
func (dp Length) PixelsY() int {
	return fixed.Int26_6(dp.Scale(DPI.Y, 96)).Round()
}

// Scale scales the distance by the ratio of num:den.
func (dp Length) Scale(num, den int) Length {
	return Length(int64(dp) * int64(num) / int64(den))
}

// String returns a human readable distance.
func (dp Length) String() string {
	return fixed.Int26_6(dp).String()
}

// FromPixelsX converts a distance measurement in physical pixels to DIPs, based on
// the current DPI settings for horizontal scaling.
func FromPixelsX(pixels int) Length {
	return Length(pixels<<6).Scale(96, DPI.X)
}

// FromPixelsY converts a distance measurement in physical pixels to DIPs, based on
// the current DPI settings for vertical scaling.
func FromPixelsY(pixels int) Length {
	return Length(pixels<<6).Scale(96, DPI.Y)
}

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point struct {
	X, Y Length
}

// String returns a string representation of p like "(3,4)".
func (p Point) String() string {
	return "(" + p.X.String() + "," + p.Y.String() + ")"
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Pixels returns the vector with the X and Y coordinates measured in pixels.
func (p Point) Pixels() image.Point {
	return image.Point{p.X.PixelsX(), p.Y.PixelsY()}
}

// A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
//
// A Rectangle is also an Image whose bounds are the rectangle itself. At
// returns color.Opaque for points in the rectangle and color.Transparent
// otherwise.
type Rectangle struct {
	Min, Max Point
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Dx returns r's width.
func (r Rectangle) Dx() Length {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle) Dy() Length {
	return r.Max.Y - r.Min.Y
}

// Pixels returns the rectangle with the X and Y coordinates measured in pixels.
func (r Rectangle) Pixels() image.Rectangle {
	return image.Rectangle{r.Min.Pixels(), r.Max.Pixels()}
}

// Size returns r's width and height.
func (r Rectangle) Size() Point {
	return Point{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

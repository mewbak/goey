package goey

import (
	"golang.org/x/image/math/fixed"
	"image"
)

var (
	// DPI contains the current DPI (dots per inch) of the monitor.
	// User code should not need to set this directly, as drivers will update
	// this variable as necessary.
	DPI image.Point
)

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
	return float64(dp) * ((96 / 72) / 1 << 6)
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

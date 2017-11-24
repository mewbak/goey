package goey

import (
	"image"
)

var (
	// DPI contains the current DPI (dots per inch) of the monitor.
	// User code should not need to set this directly, as drivers will update
	// this variable as necessary.
	DPI image.Point
)

// DIP is a distance measured in device-independent pixels.  There are nominally
// 96 DIPs per inch.
type DIP int

// PixelsX converts the distance measurement in DIPs to physical pixels, based
// on the current DPI settings for horizontal scaling.
func (dp DIP) PixelsX() int {
	return int(dp) * DPI.X / 96
}

// PixelsY converts the distance measurement in DIPs to physical pixels, based
// on the current DPI settings for vertical scaling.
func (dp DIP) PixelsY() int {
	return int(dp) * DPI.Y / 96
}

// ToDIPX converts a distance measurement in physical pixels to DIPs, based on
// the current DPI settings for horizontal scaling.
func ToDIPX(pixels int) DIP {
	return DIP(pixels * 96 / DPI.X)
}

// ToDIPY converts a distance measurement in physical pixels to DIPs, based on
// the current DPI settings for vertical scaling.
func ToDIPY(pixels int) DIP {
	return DIP(pixels * 96 / DPI.Y)
}

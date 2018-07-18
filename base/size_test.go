package base

import (
	"fmt"
	"image"
)

func ExampleSize_FromPixels() {
	// Most code should not need to worry about setting the DPI.  Windows will
	// ensure that the DPI is set.
	DPI = image.Point{96, 96}

	size := FromPixels(48, 96+96)
	fmt.Printf("The size is %s.\n", size.String())

	// Output:
	// The size is (48:00x192:00).
}

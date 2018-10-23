// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/cocoa"

	"bitbucket.org/rj/goey/base"
)

type imgElement struct {
	control *cocoa.ImageView
	buffer  []byte
	width   base.Length
	height  base.Length
}

func (w *Img) mount(parent base.Control) (base.Element, error) {
	// Convert the image to an NSImage
	control, buffer, err := cocoa.NewImageView(parent.Handle, w.Image)
	if err != nil {
		return nil, err
	}

	retval := &imgElement{
		control: control,
		buffer:  buffer,
		width:   w.Width,
		height:  w.Height,
	}

	return retval, nil
}

func (w *imgElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *imgElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *imgElement) updateProps(data *Img) error {
	buffer, err := w.control.SetImage(data.Image)
	if err != nil {
		return err
	}
	w.buffer = buffer
	w.width, w.height = data.Width, data.Height
	return nil
}

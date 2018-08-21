package icons

import (
	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
	"errors"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
)

// Icon describes a widget that shows an icon as an image.
type Icon rune

var (
	kind   = base.NewKind("bitbucket.org/rj/goey/icons.Icon")
	assets struct {
		font *truetype.Font
		face font.Face
	}
)

func init() {
	err := error(nil)
	assets.font, err = truetype.Parse(file0[:])
	if err != nil {
		panic("internal error: failed to parse embedded truetype file")
	}

	assets.face = truetype.NewFace(assets.font, &truetype.Options{Size: 32})
}

// New returns a new widget description an image showing the icon with the
// specified rune.
func New(r rune) Icon {
	return Icon(r)
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (i Icon) Kind() *base.Kind {
	return &kind
}

// Mount creates a control in the GUI to display the icon.
// The newly created widget will be a child of the widget specified by parent.
func (i Icon) Mount(parent base.Control) (base.Element, error) {
	ndx := assets.font.Index(rune(i))
	if ndx == 0 {
		return nil, errors.New("rune not available")
	}

	dr, _, _, _, _ := assets.face.Glyph(fixed.P(0, 0), rune(i))

	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	dot := fixed.P(16-dr.Dx()/2-dr.Min.X, 16+dr.Dy()/2-dr.Max.Y)

	dr, mask, maskp, _, _ := assets.face.Glyph(dot, rune(i))

	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Over)
	draw.DrawMask(img, dr, image.Black, image.Point{}, mask, maskp, draw.Over)

	widget := goey.Img{Image: img}
	elem, err := widget.Mount(parent)
	if err != nil {
		return nil, err
	}

	return &iconElement{elem, rune(i)}, nil
}

package goey

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestImg(t *testing.T) {
	bounds := image.Rect(0, 0, 92, 92)
	images := []*image.RGBA{image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds)}
	draw.Draw(images[0], bounds, image.NewUniform(color.RGBA{255, 255, 0, 255}), image.Point{}, draw.Src)
	draw.Draw(images[1], bounds, image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(images[2], bounds, image.NewUniform(color.RGBA{0, 255, 255, 255}), image.Point{}, draw.Src)

	testingRenderWidgets(t, []Widget{
		&Img{Image: images[0]},
		&Img{Image: images[1]},
		&Img{Image: images[2]},
	})
}

func TestImgClose(t *testing.T) {
	bounds := image.Rect(0, 0, 92, 92)
	images := []*image.RGBA{image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds)}
	draw.Draw(images[0], bounds, image.NewUniform(color.RGBA{255, 255, 0, 255}), image.Point{}, draw.Src)
	draw.Draw(images[1], bounds, image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(images[2], bounds, image.NewUniform(color.RGBA{0, 255, 255, 255}), image.Point{}, draw.Src)

	testingRenderWidgets(t, []Widget{
		&Img{Image: images[0]},
		&Img{Image: images[1]},
		&Img{Image: images[2]},
	})
}

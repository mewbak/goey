package main

import (
	"fmt"
	"goey"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"
)

var (
	mainWindow *goey.Window
	clickCount int

	colors = [3]color.RGBA{
		color.RGBA{0xff, 0, 0, 0xff},
		color.RGBA{0, 0xff, 0, 0xff},
		color.RGBA{0, 0, 0xff, 0xff},
	}
	colorNames = [3]string{
		"Red", "Green", "Blue",
	}

	gopher image.Image
)

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func main() {
	err := error(nil)
	gopher, err = loadImage("gopher.png")
	if err != nil {
		println(err.Error())
		return
	}

	mw, err := goey.NewWindow("Colour", render())
	if err != nil {
		println(err.Error())
		return
	}
	defer mw.Close()
	mw.SetAlignment(goey.MainCenter, goey.CrossCenter)
	mainWindow = mw

	goey.Run()
}

func update() {
	err := mainWindow.SetChildren(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func render() []goey.Widget {
	img := image.Image(nil)
	description := ""
	if clickCount%4 == 3 {
		img = gopher
		description = "Image of the Go gopher."
	} else {
		rgbimg := image.NewRGBA(image.Rect(0, 0, 100, 100))
		draw.Draw(rgbimg, rgbimg.Bounds(), image.NewUniform(colors[clickCount%4]), image.Point{}, draw.Src)
		img = rgbimg
		description = colorNames[clickCount%4]
	}

	return []goey.Widget{
		&goey.Button{Text: "Change the colour", OnClick: func() {
			clickCount++
			update()
		}},
		&goey.Img{
			Image:  img,
			Width:  goey.DIP(img.Bounds().Dx()),
			Height: goey.DIP(img.Bounds().Dy()),
		},
		&goey.P{Text:description},
	}
}

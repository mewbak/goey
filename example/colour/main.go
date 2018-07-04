package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"

	"bitbucket.org/rj/goey"
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

func selectImage(index int) (image.Image, string) {
	if clickCount%4 == 3 {
		return gopher, "Image of the Go gopher."
	} else {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		draw.Draw(img, img.Bounds(), image.NewUniform(colors[index%4]), image.Point{}, draw.Src)
		return img, colorNames[index%4]
	}
}

func main() {
	err := error(nil)
	gopher, err = loadImage("gopher.png")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	err = goey.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	// Add the controls
	mw, err := goey.NewWindow("Colour", render())
	if err != nil {
		return err
	}
	mw.SetScroll(false, true)
	mainWindow = mw

	// Set the icon
	img, _ := selectImage(clickCount)
	mw.SetIcon(img)

	return nil
}

func update() {
	err := mainWindow.SetChild(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	img, _ := selectImage(clickCount)
	mainWindow.SetIcon(img)
}

func render() goey.Widget {
	img, description := selectImage(clickCount)

	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			AlignMain:  goey.MainCenter,
			AlignCross: goey.CrossCenter,
			Children: []goey.Widget{
				&goey.Button{Text: "Change the colour", OnClick: func() {
					clickCount++
					update()
				}},
				&goey.Img{
					Image:  img,
					Width:  (1 * goey.DIP).Scale(img.Bounds().Dx(), 1),
					Height: (1 * goey.DIP).Scale(img.Bounds().Dy(), 1),
				},
				&goey.P{Text: description, Align: goey.JustifyCenter},
			},
		},
	}
}

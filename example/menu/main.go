package main

import (
	"bitbucket.org/rj/goey"
	"fmt"
	"image"
	_ "image/png"
	"os"
)

var (
	gopher image.Image
	window *goey.Window
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
	img, err := loadImage("gopher.png")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	gopher = img

	err = goey.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	w, err := goey.NewWindow("Menu", renderWindow())
	if err != nil {
		return err
	}
	w.SetAlignment(goey.MainCenter, goey.CrossCenter)
	window = w
	return nil
}

func updateWindow() {
	window.SetChildren(renderWindow())
}

func renderSidebar() goey.Widget {
	return &goey.VBox{goey.MainCenter, goey.CrossCenter, []goey.Widget{
		&goey.Label{Text: "Example Menu"},
		&goey.Img{Image: gopher},
	}}
}

func renderMainbar() goey.Widget {
	return &goey.VBox{goey.MainCenter, goey.CrossCenter, []goey.Widget{
		&goey.Column{[][]goey.Widget{
			{&goey.Button{Text: "A1"}, &goey.Button{Text: "A2"}, &goey.Button{Text: "A3"}, &goey.Button{Text: "A4"}},
			{&goey.Button{Text: "B1"}, &goey.Button{Text: "B2"}, &goey.Button{Text: "B3"}, &goey.Button{Text: "B4"}},
			{&goey.Button{Text: "C1"}, &goey.Button{Text: "C2"}, &goey.Button{Text: "C3"}, &goey.Button{Text: "C4"}},
			{&goey.Button{Text: "D1"}, &goey.Button{Text: "D2"}, &goey.Button{Text: "D3"}, &goey.Button{Text: "D4"}},
		}},
		&goey.HR{},
		&goey.Button{Text: "Help"},
	}}
}

func renderWindow() []goey.Widget {
	ret := []goey.Widget{&goey.HBox{
		Children: []goey.Widget{
			renderSidebar(),
			renderMainbar(),
		},
	}}

	return ret
}

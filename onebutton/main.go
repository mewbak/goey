package main

import (
	"fmt"
	"goey"
	"strconv"
)

var (
	mainWindow *goey.MainWindow
	clickCount int
)

func main() {
	mw, err := goey.NewMainWindow("One Button", render())
	if err != nil {
		println(err.Error())
		return
	}
	defer mw.Close()
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
	text := "Click me!"
	if clickCount > 0 {
		text = text + "  (" + strconv.Itoa(clickCount) + ")"
	}
	return []goey.Widget{
		&goey.Button{Text: text, OnClick: func() {
			clickCount++
			update()
		}},
	}
}

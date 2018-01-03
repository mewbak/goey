package main

import (
	"fmt"
	"goey"
	"strconv"
)

var (
	mainWindow *goey.Window
	clickCount int
)

func main() {
	err := goey.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	mw, err := goey.NewWindow("One Button", render())
	if err != nil {
		return err
	}
	mw.SetAlignment(goey.MainCenter, goey.CrossCenter)
	mainWindow = mw

	return nil
}

func updateWindow() {
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
			updateWindow()
		}},
	}
}

// This package provides an example application built using the goey package
// that shows a single button.  The button is centered in the window, and, when
// the button is clicked, the button's caption is changed to keep a running
// total.
package main

import (
	"fmt"
	"strconv"

	"bitbucket.org/rj/goey"
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
	mainWindow = mw

	return nil
}

func updateWindow() {
	err := mainWindow.SetChild(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func render() goey.Widget {
	text := "Click me!"
	if clickCount > 0 {
		text = text + "  (" + strconv.Itoa(clickCount) + ")"
	}
	return &goey.Padding{
		Insets: goey.UniformInsets(36 * goey.DIP),
		Child: &goey.Align{
			Child: &goey.Button{Text: text, OnClick: func() {
				clickCount++
				updateWindow()
			}},
		},
	}
}

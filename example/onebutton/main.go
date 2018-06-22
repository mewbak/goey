// Package main for an example application using the goey package to show a
// single button.  The button is centered in the window, and, when the button
// is clicked, the button's caption is changed to keep a running total.
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
		Insets: goey.UniformInset(36 * goey.DIP),
		Child: &goey.Center{
			Child: &goey.Button{Text: text, OnClick: func() {
				clickCount++
				updateWindow()
			}},
		},
	}
}

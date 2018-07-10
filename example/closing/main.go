package main

import (
	"fmt"

	"bitbucket.org/rj/goey"
)

var (
	mainWindow *goey.Window
)

func main() {
	err := goey.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	mw, err := goey.NewWindow("Closing", render())
	if err != nil {
		return err
	}
	mw.SetOnClosing(func() bool {
		// Block closing of the window
		return true
	})
	mainWindow = mw

	return nil
}

func render() goey.Widget {
	return &goey.Padding{
		Insets: goey.UniformInsets(36 * goey.DIP),
		Child: &goey.Align{
			Child: &goey.Button{Text: "Close app", OnClick: func() {
				mainWindow.Close()
			}},
		},
	}
}

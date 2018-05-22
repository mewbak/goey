// Package main for an example application using the goey package to show
// two editable multiline text fields.  As the fields are edited, a
// status line with the total count of characters is maintained.
package main

import (
	"bitbucket.org/rj/goey"
	"fmt"
)

var (
	mainWindow *goey.Window
)

const (
	fillerText = "*** Text Text Text Text ***"
)

func main() {
	err := goey.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	mw, err := goey.NewWindow("Two Fields", render())
	if err != nil {
		return err
	}
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
	return []goey.Widget{
		&goey.Column{
			Children: [][]goey.Widget{
				{
					&goey.Label{Text: fillerText},
				},
				{
					&goey.Label{Text: fillerText},
				},
				{
					&goey.Label{Text: fillerText},
				},
				{
					&goey.Label{Text: fillerText},
				},
				{
					&goey.Label{Text: fillerText},
				},
				{
					&goey.Label{Text: fillerText},
				},
				{
					&goey.Label{Text: fillerText},
				},
				{
					&goey.Label{Text: fillerText},
				},
			},
		},
		&goey.HR{},
		&goey.Label{Text: "The total character count is...  "},
	}
}

func onfocus(ndx int) func() {
	return func() {
		fmt.Println("focus", ndx)
	}
}

func onblur(ndx int) func() {
	return func() {
		fmt.Println("blur", ndx)
	}
}

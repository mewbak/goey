// Package main for an example application using the goey package to show
// two editable multiline text fields.  As the fields are edited, a
// status line with the total count of characters is maintained.
package main

import (
	"bitbucket.org/rj/goey"
	"fmt"
	"strconv"
)

var (
	mainWindow     *goey.Window
	text           [2]string
	characterCount [2]int
	wordCount      [2]int
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
		&goey.VBox{
			Children: []goey.Widget{
				&goey.Label{Text: "This is the most important field:"},
				&goey.TextArea{Value: text[0], Placeholder: "You should type something here.",
					OnChange: func(value string) {
						text[0] = value
						characterCount[0] = len(value)
						updateWindow()
					},
					OnFocus: onfocus(1),
					OnBlur:  onblur(1),
				},
				&goey.Label{Text: "This is a secondary field:"},
				&goey.TextArea{Value: text[1], Placeholder: "...and here.",
					OnChange: func(value string) {
						text[1] = value
						characterCount[1] = len(value)
						updateWindow()
					},
					OnFocus: onfocus(2),
					OnBlur:  onblur(2),
				},
			},
		},
		&goey.HR{},
		&goey.Label{Text: "The total character count is:  " + strconv.Itoa(characterCount[0]+characterCount[1])},
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

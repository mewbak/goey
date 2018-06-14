// Package main for an example application using the goey package to display
// various images.  The images can be cycled by clicking a button, and each
// image has an associated description.
package main

import (
	"bitbucket.org/rj/goey"
	"fmt"
	"strconv"
)

var (
	mainWindow *goey.Window

	feetValue  string
	meterValue string
)

func main() {
	err := goey.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func createWindow() error {
	// Add the controls
	mw, err := goey.NewWindow("Feet to Meters", render())
	if err != nil {
		return err
	}
	mw.SetAlignment(goey.MainCenter, goey.CrossCenter)
	mainWindow = mw

	return nil
}

func update() {
	err := mainWindow.SetChildren(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func render() []goey.Widget {
	return []goey.Widget{
		&Grid{
			Children: [3][3]goey.Widget{{
				&goey.Empty{},
				&goey.TextInput{Value: feetValue, OnChange: func(v string) { feetValue = v }, OnEnterKey: func(v string) { feetValue = v; calculate() }},
				&goey.Label{Text: "feet"},
			}, {
				&goey.Label{Text: "is equivalent to"},
				&goey.Label{Text: meterValue},
				&goey.Label{Text: "meters"},
			}, {
				&goey.Empty{},
				&goey.Empty{},
				&goey.Button{Text: "Calculate", Default: true, OnClick: calculate},
			}},
		},
	}
}

func calculate() {
	feet, err := strconv.ParseFloat(feetValue, 64)
	if err != nil {
		meterValue = "(error)"
	} else {
		meterValue = fmt.Sprintf("%f", feet*0.3048)
	}
	update()
}

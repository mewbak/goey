// This package provides an example application built using the goey package
// that demonstrates most of the controls that are available.
package main

import (
	"fmt"
	"time"

	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
)

var (
	currentCD  string
	s1, s2, s3 bool
	showLorem  bool
	window     *goey.Window
)

func main() {
	err := goey.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	w, err := goey.NewWindow("Example", renderWindow())
	if err != nil {
		return err
	}
	w.SetScroll(false, true)
	window = w
	return nil
}

func updateWindow() {
	window.SetChild(renderWindow())
}

func renderWindow() base.Widget {
	// Fixed part at top
	widgets := []base.Widget{
		&goey.Checkbox{Text: "Show text of lorem ipsum.", Value: showLorem, OnChange: func(ok bool) {
			showLorem = ok
			updateWindow()
		}},
		&goey.HR{},
	}

	// Depends
	if showLorem {
		widgets = append(widgets,
			&goey.P{Text: lorem},
		)
	} else {
		widgets = append(widgets,
			&goey.P{Text: "This is a paragraph, but without much text.", Align: goey.JustifyCenter},
			&goey.Label{Text: "Text input:"},
			&goey.TextInput{Value: "Some input...", Placeholder: "Type some text here.  And some more.  And something really long.",
				OnChange: func(v string) { println("text input ", v) }, OnEnterKey: func(v string) { println("t1* ", v) }},
			&goey.Label{Text: "Password input:"},
			&goey.TextInput{Value: "", Placeholder: "Don't share", Password: true,
				OnChange: func(v string) { println("password input ", v) }},
			&goey.Label{Text: "Integer input:"},
			&goey.IntInput{Value: 3, Placeholder: "Please enter a number",
				OnChange: func(v int64) { println("int input ", v) }},
			&goey.Label{Text: "Date input:"},
			&goey.DateInput{Value: time.Now().Add(24 * time.Hour),
				OnChange: func(v time.Time) { println("date input: ", v.String()) }},
			&goey.HR{},
			&goey.HBox{Children: []base.Widget{
				&goey.Button{Text: "C1", Default: true},
				&goey.Button{Text: "C2"},
			}},
			&goey.HBox{Children: []base.Widget{
				&goey.Button{Text: "D1"},
				&goey.Button{Text: "D2", Disabled: true},
				&goey.Button{Text: "D3"},
			},
				AlignMain: goey.MainEnd,
			},
			&goey.HR{},
			&goey.SelectInput{Items: []string{"Choice 1", "Choice 2", "Choice 3"},
				OnChange: func(v int) { println("select input: ", v) }},
			&goey.TextArea{Value: "", Placeholder: "Room to write",
				OnChange: func(v string) { println("text area: ", v) }},
		)
	}

	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child:  &goey.VBox{Children: widgets},
	}
}

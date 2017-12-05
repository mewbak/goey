package main

import (
	"goey"
)

var (
	currentCD  string
	s1, s2, s3 bool
)

func main() {
	mw, err := goey.NewMainWindow("Example",
		[]goey.Widget{
			&goey.P{Text: lorem},
			&goey.P{Text: "This is a label.", Align: goey.Center},
			&goey.Label{Text: "Input 1:"},
			&goey.TextInput{Value: "Some input...", OnChange: func(v string) { println("t1 ", v) }, OnEnterKey: func(v string) { println("t1* ", v) }},
			&goey.HR{},
			&goey.Button{Text: "This is button 1", OnClick: func() { println("b1") }},
			&goey.Label{Text: "Input 2:"},
			&goey.TextInput{Value: "", Placeholder: "Type some text here.  And some more.  And something really long.", OnChange: func(v string) { println("t2 ", v) }},
			&goey.Checkbox{Text: "Check this box", Value: true, OnChange: func(b bool) { println("c1", b) }},
			&goey.Button{Text: "This is second button", Default: true},
			&goey.HBox{Children: []goey.Widget{
				&goey.Button{Text: "C1"},
				&goey.Button{Text: "C2"},
			}},
			&goey.HBox{Children: []goey.Widget{
				&goey.Button{Text: "D1"},
				&goey.Button{Text: "D2", Disabled: true},
				&goey.Button{Text: "D3"},
			},
				AlignMain: goey.MainEnd,
			},
			&goey.SelectInput{Items: []string{"Choice 1", "Choice 2", "Choice 3"}, OnChange: func(v int) { println("cb1 ", v) }},
			&goey.TextArea{},
		},
	)
	if err != nil {
		println(err.Error())
		return
	}
	defer mw.Close()

	goey.Run()
}

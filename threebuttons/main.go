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
	mw.SetAlignment(goey.SpaceAround, goey.Stretch)
	mainWindow = mw

	goey.Run()
}

func update() {
	err := mainWindow.SetChildren(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func onclick(ndx int) func() {
	return func() {
		fmt.Println("click", ndx)
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

func render() []goey.Widget {
	text := "Click me!"
	if clickCount > 0 {
		text = text + "  (" + strconv.Itoa(clickCount) + ")"
	}
	return []goey.Widget{
		&goey.Button{Text: text,
			Default: true,
			OnClick: func() {
				clickCount++
				update()
			},
			OnFocus: onfocus(1),
			OnBlur:  onblur(1),
		},
		&goey.Button{Text: "Extra button",
			OnClick: onclick(2),
			OnFocus: onfocus(2),
			OnBlur:  onblur(2),
		},
		&goey.Button{Text: "Supplefluous button",
			OnClick: onclick(3),
			OnFocus: onfocus(3),
			OnBlur:  onblur(3),
		},
	}
}

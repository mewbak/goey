package main

import (
	"fmt"
	"goey"
	"strconv"
)

var (
	mainWindow *goey.Window
	clickCount int
	alignMain  = goey.SpaceAround
	alignCross = goey.Stretch
)

func main() {
	mw, err := goey.NewWindow("Three Buttons", render())
	if err != nil {
		println(err.Error())
		return
	}
	defer mw.Close()
	mw.SetAlignment(alignMain, alignCross)
	mainWindow = mw

	goey.Run()
}

func update() {
	err := mainWindow.SetChildren(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func cycleMainAxisAlign() {
	alignMain++
	if alignMain > goey.SpaceBetween {
		alignMain = goey.MainStart
	}
	mainWindow.SetAlignment(alignMain, alignCross)
	update()
}

func cycleCrossAxisAlign() {
	alignCross++
	if alignCross > goey.CrossEnd {
		alignCross = goey.Stretch
	}
	mainWindow.SetAlignment(goey.SpaceAround, alignCross)
	update()
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
			OnClick: cycleMainAxisAlign,
			OnFocus: onfocus(2),
			OnBlur:  onblur(2),
		},
		&goey.Button{Text: "Cycle cross axis align",
			OnClick: cycleCrossAxisAlign,
			OnFocus: onfocus(3),
			OnBlur:  onblur(3),
		},
	}
}

package icons

import (
	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
	"fmt"
)

func ExampleIcon() {
	render := func() base.Widget {
		return &goey.Padding{
			Insets: goey.DefaultInsets(),
			Child: &goey.Align{
				Child: Icon(0xe869),
			},
		}
	}

	createWindow := func() error {
		// Add the controls
		_, err := goey.NewWindow("Icons", render())
		return err
	}

	err := goey.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("OK")

	// Output:
	// OK
}

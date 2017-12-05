package main

import (
	"fmt"
	"goey"
	"time"
)

func main() {
	mw, err := goey.NewMainWindow("Flash", []goey.Widget{
		&goey.Button{Text: "Click me!"},
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	go func() {
		fmt.Println("Up")
		time.Sleep(1 * time.Second)
		goey.Do(func() error {
			fmt.Println("Down")
			mw.Close()
			return nil
		})
	}()

	goey.Run()
}

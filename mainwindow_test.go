package goey

import (
	"fmt"
	"testing"
	"time"
)

func ExampleNewMainWindow() {
	mw, err := NewMainWindow("Test", []Widget{
		&Button{Text: "Click me!"},
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	go func() {
		fmt.Println("Up")
		time.Sleep(1 * time.Second)
		Do(func() error {
			mw.Close()
			return nil
		})
		fmt.Println("Down")
	}()

	Run()
}

func TestNewMainWindow_SetAlignment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	mw, err := NewMainWindow("Test", []Widget{
		&Button{Text: "Click me!"},
		&Button{Text: "Or me!"},
		&Button{Text: "But not me."},
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	go func() {
		t.Logf("Startin alignment tests")
		for i := MainStart; i <= SpaceBetween; i++ {
			for j := Stretch; j <= CrossEnd; j++ {
				time.Sleep(100 * time.Millisecond)
				Do(func() error {
					mw.SetAlignment(i, j)
					return nil
				})
			}
		}
		time.Sleep(100 * time.Millisecond)
		Do(func() error {
			mw.Close()
			return nil
		})
		t.Logf("Stopping alignment tests")
	}()

	Run()
}

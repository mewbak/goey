package goey

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func ExampleNewWindow() {
	// All calls that modify GUI objects need to be schedule ont he GUI thread.
	// This callback will be used to create the top-level window.
	createWindow := func() error {
		// Create a top-level window.
		mw, err := NewWindow("Test", []Widget{
			&Button{Text: "Click me!"},
		})
		if err != nil {
			// This error will be reported back up through the call to
			// Run below.  No need to print or log it here.
			return err
		}

		// We can start a goroutine, but note that we can't modify GUI objects
		// directly.
		go func() {
			fmt.Println("Up")
			time.Sleep(50 * time.Millisecond)
			fmt.Println("Down")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	// Start the GUI thread.
	err := Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Output:
	// Up
	// Down
}

func ExampleWindow_Message() {
	// All calls that modify GUI objects need to be schedule ont he GUI thread.
	// This callback will be used to create the top-level window.
	createWindow := func() error {
		// Create a top-level window.
		mw, err := NewWindow("Test", []Widget{
			&Button{Text: "Click me!"},
		})
		if err != nil {
			// This error will be reported back up through the call to
			// Run below.  No need to print or log it here.
			return err
		}

		// We can start a goroutine, but note that we can't modify GUI objects
		// directly.
		go func() {
			// Show the error message.
			Do(func() error {
				return mw.Message("This is an example message.").WithInfo().Show()
			})

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	// Start the GUI thread.
	err := Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	//Output:
}

func TestWindow_SetAlignment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	createWindow := func() error {

		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		mw, err := NewWindow("Test", []Widget{
			&Button{Text: "Click me!"},
			&Button{Text: "Or me!"},
			&Button{Text: "But not me."},
		})
		if err != nil {
			t.Fatalf("Failed to create window, %s", err)
		}
		if mw == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}

		go func() {
			t.Logf("Starting alignment tests")
			for i := MainStart; i <= SpaceBetween; i++ {
				for j := Stretch; j <= CrossEnd; j++ {
					time.Sleep(50 * time.Millisecond)
					Do(func() error {
						mw.SetAlignment(i, j)
						return nil
					})
					itmp, jtmp := mw.Alignment()
					if itmp != i {
						t.Errorf("Expected main alignment itmp==i, got %d and %d", itmp, i)
					}
					if jtmp != j {
						t.Errorf("Expected cross alignment jtmp==j, got %d and %d", jtmp, j)
					}
				}
			}
			time.Sleep(50 * time.Millisecond)
			t.Logf("Stopping alignment tests")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	err := Run(createWindow)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func TestNewWindow_SetChildren(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	createWindow := func() error {
		widgets := []Widget{}

		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		mw, err := NewWindow("Test", widgets)
		if err != nil {
			t.Fatalf("Failed to create window, %s", err)
		}
		if mw == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}
		mw.SetAlignment(SpaceBetween, CrossCenter)

		go func() {
			t.Logf("Starting set children tests")
			for i := 1; i < 10; i++ {
				time.Sleep(50 * time.Millisecond)
				widgets = append(widgets, &Button{Text: "Button " + strconv.Itoa(i)})
				err := Do(func() error {
					return mw.SetChildren(widgets)
				})
				if err != nil {
					t.Logf("Error setting children, %s", err)
				}
			}
			for i := len(widgets); i > 0; i-- {
				time.Sleep(50 * time.Millisecond)
				widgets = widgets[:i-1]
				err := Do(func() error {
					return mw.SetChildren(widgets)
				})
				if err != nil {
					t.Logf("Error setting children, %s", err)
				}
			}
			time.Sleep(50 * time.Millisecond)
			t.Logf("Stopping set children tests")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	err := Run(createWindow)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func makeImage(t *testing.T, index int) image.Image {
	colors := [3]color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	bounds := image.Rect(0, 0, 32, 32)
	img := image.NewRGBA(bounds)
	draw.Draw(img, image.Rect(0, 0, 11, 32), image.NewUniform(colors[index%3]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(11, 0, 22, 32), image.NewUniform(colors[(index+1)%3]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(22, 0, 32, 32), image.NewUniform(colors[(index+2)%3]), image.Point{}, draw.Src)
	return img
}

func TestNewWindow_SetIcon(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	createWindow := func() error {
		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		mw, err := NewWindow("Test", []Widget{
			&Button{Text: "Click me!"},
		})
		if err != nil {
			t.Fatalf("Failed to create window, %s", err)
		}
		if mw == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}

		go func() {
			t.Logf("Starting set icon tests")
			for i := 0; i < 6; i++ {
				img := makeImage(t, i)

				err := Do(func() error {
					return mw.SetIcon(img)
				})
				if err != nil {
					t.Errorf("Error calling SetTitle, %", err)
				}
				time.Sleep(50 * time.Millisecond)
			}
			t.Logf("Stopping icon tests")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			err = Do(func() error {
				mw.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error calling Close, %", err)
			}
		}()

		return nil
	}

	err := Run(createWindow)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func TestNewWindow_SetTitle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	createWindow := func() error {

		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		mw, err := NewWindow("Test", []Widget{
			&Button{Text: "Click me!"},
			&Button{Text: "Or me!"},
			&Button{Text: "But not me."},
		})
		if err != nil {
			t.Fatalf("Failed to create window, %s", err)
		}
		if mw == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}

		go func() {
			t.Logf("Starting set title tests")
			time.Sleep(50 * time.Millisecond)
			err := Do(func() error {
				return mw.SetTitle("Flash!")
			})
			if err != nil {
				t.Errorf("Error calling SetTitle, %", err)
			}
			t.Logf("Stopping alignment tests")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			err = Do(func() error {
				mw.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error calling Close, %", err)
			}
		}()

		return nil
	}

	err := Run(createWindow)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

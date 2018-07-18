package goey

import (
	"bitbucket.org/rj/goey/base"
)

var (
	mockKind = base.NewKind("bitbucket.org/rj/goey.Mock")
)

func mock(width, height base.Length) base.Element {
	return &mockElement{
		Width:  width,
		Height: height,
	}
}

type mockElement struct {
	Width, Height base.Length
	Bounds        base.Rectangle
}

func (*mockElement) Close() {
}

func (*mockElement) Kind() *base.Kind {
	return &mockKind
}

func (m *mockElement) Layout(bc base.Constraints) base.Size {
	return bc.Constrain(base.Size{m.Width, m.Height})
}

func (m *mockElement) MinIntrinsicHeight(width base.Length) base.Length {
	return m.Height
}

func (m *mockElement) MinIntrinsicWidth(height base.Length) base.Length {
	return m.Width
}

func (m *mockElement) SetBounds(bounds base.Rectangle) {
	m.Bounds = bounds
}

func (*mockElement) UpdateProps(data base.Widget) error {
	return nil
}

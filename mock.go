package goey

var (
	mockKind = Kind{"bitbucket.org/rj/goey.Mock"}
)

type mockElement struct {
	Width, Height Length
}

func (*mockElement) Close() {
}

func (*mockElement) Kind() *Kind {
	return &mockKind
}

func (m *mockElement) Layout(bc Constraint) Size {
	return bc.Constrain(Size{m.Width, m.Height})
}

func (m *mockElement) MinIntrinsicHeight(width Length) Length {
	return m.Height
}

func (m *mockElement) MinIntrinsicWidth(height Length) Length {
	return m.Width
}

func (*mockElement) SetBounds(bounds Rectangle) {
}

func (*mockElement) UpdateProps(data Widget) error {
	return nil
}

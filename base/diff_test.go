package base

import (
	"reflect"
	"testing"
)

type mock struct {
	kind *Kind
	Prop int
}

func (m *mock) Kind() *Kind {
	return m.kind
}

func (m *mock) Mount(parent Control) (Element, error) {
	return &mockElement{
		kind: m.kind,
		Prop: m.Prop,
	}, nil
}

type mockElement struct {
	kind   *Kind
	Closed bool
	Prop   int
}

func (m *mockElement) Close() {
	m.Closed = true
}

func (m *mockElement) Kind() *Kind {
	return m.kind
}

func (m *mockElement) Layout(Constraints) Size {
	return Size{}
}
func (m *mockElement) MinIntrinsicHeight(width Length) Length {
	return 0
}

func (m *mockElement) MinIntrinsicWidth(height Length) Length {
	return 0
}

func (m *mockElement) SetBounds(bounds Rectangle) {

}

func (m *mockElement) updateProps(data *mock) error {
	if m.kind != data.kind {
		panic("Mismatched kinds")
	}
	m.Prop = data.Prop
	return nil
}

func (m *mockElement) UpdateProps(data Widget) error {
	return m.updateProps(data.(*mock))
}

func TestCloseElements(t *testing.T) {
	kind := NewKind("bitbucket.org/rj/goey/base.Mock")

	// Check for no panic on nil or zero-length list
	CloseElements(nil)
	CloseElements([]Element{})

	for _, v := range []int{1, 2, 3, 4, 8, 16} {
		elem := make([]Element, 0, v)
		for i := 0; i < v; i++ {
			elem = append(elem, &mockElement{kind: &kind})
		}

		CloseElements(elem)

		for _, v := range elem {
			if !v.(*mockElement).Closed {
				t.Errorf("Failed to close element")
			}
		}
	}
}

func TestDiffChild(t *testing.T) {
	kind1 := NewKind("bitbucket.org/rj/goey/base.Mock1")
	kind2 := NewKind("bitbucket.org/rj/goey/base.Mock2")

	cases := []struct {
		lhs       Element
		rhs       Widget
		out       Element
		lhsClosed bool
	}{
		{nil, nil, nil, false},
		{nil, &mock{kind: &kind1}, &mockElement{kind: &kind1}, false},
		{nil, &mock{kind: &kind1, Prop: 3}, &mockElement{kind: &kind1, Prop: 3}, false},
		{nil, &mock{kind: &kind2}, &mockElement{kind: &kind2}, false},
		{nil, &mock{kind: &kind2, Prop: 13}, &mockElement{kind: &kind2, Prop: 13}, false},
		{&mockElement{kind: &kind1}, nil, nil, true},
		{&mockElement{kind: &kind1, Prop: 3}, &mock{kind: &kind2, Prop: 13}, &mockElement{kind: &kind2, Prop: 13}, true},
		{&mockElement{kind: &kind1, Prop: 3}, &mock{kind: &kind1, Prop: 13}, &mockElement{kind: &kind1, Prop: 13}, false},
	}

	for i, v := range cases {
		out, err := DiffChild(Control{}, v.lhs, v.rhs)
		if err != nil {
			t.Errorf("Case %d: Unexpected error during DiffChild, %s", i, err)
		}
		if !reflect.DeepEqual(out, v.out) {
			t.Errorf("Case %d: Returned element does not match, got %v, want %v", i, out, v.out)
		}
		if v.lhsClosed && !v.lhs.(*mockElement).Closed {
			t.Errorf("Case %d: Failed to close lhs", i)
		}
	}
}

func TestDiffChildren(t *testing.T) {
	kind1 := NewKind("bitbucket.org/rj/goey/base.Mock1")
	kind2 := NewKind("bitbucket.org/rj/goey/base.Mock2")

	cases := []struct {
		lhs       []Element
		rhs       []Widget
		out       []Element
		lhsClosed bool
	}{
		{nil, nil, nil, false},
		{nil, []Widget{&mock{kind: &kind1}}, []Element{&mockElement{kind: &kind1}}, false},
		{nil, []Widget{&mock{kind: &kind1, Prop: 3}}, []Element{&mockElement{kind: &kind1, Prop: 3}}, false},
		{nil, []Widget{&mock{kind: &kind2}}, []Element{&mockElement{kind: &kind2}}, false},
		{nil, []Widget{&mock{kind: &kind2, Prop: 13}}, []Element{&mockElement{kind: &kind2, Prop: 13}}, false},
		{[]Element{}, nil, nil, true},
		{[]Element{&mockElement{kind: &kind1}}, nil, nil, true},
		{[]Element{&mockElement{kind: &kind2}}, nil, nil, true},
		{
			[]Element{&mockElement{kind: &kind1}},
			[]Widget{&mock{kind: &kind1, Prop: 15}},
			[]Element{&mockElement{kind: &kind1, Prop: 15}},
			false,
		},
		{
			[]Element{&mockElement{kind: &kind2}},
			[]Widget{&mock{kind: &kind2, Prop: 16}},
			[]Element{&mockElement{kind: &kind2, Prop: 16}},
			false,
		},
		{
			[]Element{&mockElement{kind: &kind1}, &mockElement{kind: &kind2, Prop: 17}},
			[]Widget{&mock{kind: &kind1, Prop: 15}},
			[]Element{&mockElement{kind: &kind1, Prop: 15}},
			false,
		},
		{
			[]Element{&mockElement{kind: &kind1}},
			[]Widget{&mock{kind: &kind1, Prop: 15}, &mock{kind: &kind2, Prop: 17}},
			[]Element{&mockElement{kind: &kind1, Prop: 15}, &mockElement{kind: &kind2, Prop: 17}},
			false,
		},
		{
			[]Element{&mockElement{kind: &kind1, Prop: 123}},
			[]Widget{&mock{kind: &kind2}},
			[]Element{&mockElement{kind: &kind2}},
			true,
		},
	}

	for i, v := range cases {
		out, err := DiffChildren(Control{}, append([]Element(nil), v.lhs...), v.rhs)
		if err != nil {
			t.Errorf("Case %d: Unexpected error during DiffChildren, %s", i, err)
		}
		if !reflect.DeepEqual(out, v.out) {
			t.Errorf("Case %d: Returned element does not match, got %v, want %v", i, out, v.out)
		}
		if len(out) < len(v.lhs) {
			for j, v := range v.lhs[len(out):] {
				if !v.(*mockElement).Closed {
					t.Errorf("Case %d: Failed to close lhs[%d]", i, len(out)+j)
				}
			}
		}
		if v.lhsClosed {
			for j, v := range v.lhs {
				if !v.(*mockElement).Closed {
					t.Errorf("Case %d: Failed to close lhs[%d]", i, j)
				}
			}
		}
	}
}

package goey

import (
	"bitbucket.org/rj/goey/base"
	"testing"
)

func TestTabsCreate(t *testing.T) {
	items := []TabItem{
		{"Tab 1", &Button{Text: "Click me!"}},
		{"Tab 2", &Button{Text: "Click me!"}},
	}

	testingRenderWidgets(t,
		&Tabs{Children: items},
		&Tabs{Value: 1, Children: items},
	)
}

func TestTabsClose(t *testing.T) {
	items := []TabItem{
		{"Tab 1", &Button{Text: "Click me!"}},
		{"Tab 2", &Button{Text: "Click me!"}},
	}

	testingRenderWidgets(t,
		&Tabs{Children: items},
		&Tabs{Value: 1, Children: items},
	)
}

func TestTabsUpdate(t *testing.T) {
	items1 := []TabItem{
		{"Tab 1", &Button{Text: "Click me!"}},
		{"Tab 2", &Button{Text: "Don't click me!"}},
	}
	items2 := []TabItem{
		{"Tab A", &Button{Text: "Don't click me!"}},
		{"Tab B", &Button{Text: "Click me!"}},
		{"Tab C", &Button{Text: "..."}},
	}

	testingUpdateWidgets(t, []base.Widget{
		&Tabs{Children: items1},
		&Tabs{Value: 1, Children: items2},
	}, []base.Widget{
		&Tabs{Value: 1, Children: items2},
		&Tabs{Children: items1},
	})
}

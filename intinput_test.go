package goey

import (
	"reflect"
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestIntInputMount(t *testing.T) {
	testingMountWidgets(t,
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
	)
}

func TestIntInputClose(t *testing.T) {
	testingCloseWidgets(t,
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
	)
}

func TestIntInputOnFocus(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&IntInput{},
		&IntInput{},
		&IntInput{},
	)
}

func TestIntInputOnChange(t *testing.T) {
	log := make([]int64, 0)

	testingTypeKeys(t, "1234",
		&IntInput{OnChange: func(v int64) {
			log = append(log, v)
		}})

	want := []int64{1, 12, 123, 1234}
	if !reflect.DeepEqual(want, log) {
		t.Errorf("Wanted %v, got %v", want, log)
	}
}

func TestIntInputOnEnterKey(t *testing.T) {
	got := int64(0)

	testingTypeKeys(t, "1234\n",
		&IntInput{OnEnterKey: func(v int64) {
			got = v
		}})

	const want = 1234
	if got != want {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func TestIntInputUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
	}, []base.Widget{
		&IntInput{Value: 1},
		&IntInput{Value: 4, Disabled: true},
		&IntInput{Value: 5, Placeholder: "***"},
	})
}

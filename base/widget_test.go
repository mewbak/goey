package base

import (
	"fmt"
)

func ExampleKind_String() {
	kind := NewKind("bitbucket.org/rj/goey/base.Example")

	fmt.Println("Kind is", kind.String())

	// Output:
	// Kind is bitbucket.org/rj/goey/base.Example
}

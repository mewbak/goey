# Goey

Package goey provides a declarative, cross-platform GUI for the
[Go](https://golang.org/) language. The range of controls, their supported
properties and events, should roughly match what is available in HTML. However,
properties and events may be limited to support portability. Additionally,
styling of the controls will be limited, with the look of controls matching the
native platform.

## Install

The package can be installed from the command line using the
[go](https://golang.org/cmd/go/) tool.

    go get bitbucket.org/rj/goey

### Linux

Although this package does not use CGO, some of its dependencies do. The build
machine also requires that GTK+ 3 is installed.  This should be installed before
issuing `go get` or you will have error messages during the building of some
of the dependencies.

On Ubuntu:

    sudo apt-get install libgtk-3-dev

## Getting Started

* Package documentation and examples are on [godoc](https://godoc.org/bitbucket.org/rj/goey).
* The minimal GUI example application is [onebutton](https://godoc.org/bitbucket.org/rj/goey/example/onebutton),
  and additional example applications are in the example folder.
* A mock widget is provided in the `mock` package
  ([documentation](https://godoc.org/bitbucket.org/rj/goey/mock)).

### Windows

To get properly themed controls, a manifest is required. Please look at the
source code for the example applications for an example. The manifest needs to
be compiled with `github.com/akavel/rsrc` to create a .syso that will be
recognize by the go build program. Additionally, you could use build flags
(`-ldflags="-H windowsgui"`) to change the type of application built.

## Contribute

Feedback and PRs welcome.

In particular, if anyone has the expertise to provide a port for MacOS, that
would provide support for all major desktop operating systems.

## License

BSD © Robert Johnstone

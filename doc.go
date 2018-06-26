// Package goey provides a declarative, cross-platform GUI.
// The range of controls, their supported properties and events, should roughly
// match what is available in HTML.  However, properties and events may be
// limited to support portability.  Additionally, styling of the controls
// will be limited, with the look of controls matching the native platform.
//
// GUI state is managed using three groups of types.  There are 'widgets', which
// are data-only representations of a desired GUI.  These widgets can be mounted
// to create 'elements', which manage actual, visible GUI resources.  The
// elements manage child elements, some number of platform-specific resources,
// or both.
//
// Windows:  To get properly themed controls, a manifest is required.  Please
// look at the source for the example applications for an example.  This file
// needs to be compiled with github.com/akavel/rsrc to create a .syso that will
// be recognize by the go build program.  Additionally, you could use build flags
// (-ldflags="-H windowsgui") to change the type of application built.
//
// Linux:  Although this package does not use CGO, some of its dependencies
// do.  The build machine also requires that GTK+ 3 is installed.
package goey

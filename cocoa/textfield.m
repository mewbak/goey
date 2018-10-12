#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void* textfieldNew( void* window, char const* text ) {
	NSString* title = [[NSString alloc] initWithUTF8String:text];

	// Create the button
	NSTextField* control = [[NSTextField alloc] init];
	[control setTitleWithMnemonic:title];

	// Add the button as the view for the window
	NSView* cv = [(NSWindow*)window contentView];
	[cv addSubview:control];

	// Return handle to the control
	return control;
}

void textfieldSetTitle( void* handle, char const* text ) {
	NSString* title = [[NSString alloc] initWithUTF8String:text];
	[(NSTextField*)handle setTitleWithMnemonic:title];
}

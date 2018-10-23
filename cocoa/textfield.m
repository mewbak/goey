#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void* textfieldNew( void* superview, char const* text ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( text );

	// Create the button
	NSTextField* control = [[NSTextField alloc] init];
	textfieldSetTitle( control, text );

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

void textfieldSetTitle( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[NSTextField class]] );
	assert( text );

	NSString* title = [[NSString alloc] initWithUTF8String:text];
	[(NSTextField*)handle setTitleWithMnemonic:title];
	[title release];
}

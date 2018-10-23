#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void* textNew( void* superview, char const* text ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( text );

	// Create the text view
	NSText* control = [[NSText alloc] init];
	[control setDrawsBackground:NO];
	textSetText( control, text );

	// Add the control as the view for the window
	[(NSView*)superview addSubview:control];

	return control;
}

void textSetText( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );
	assert( text );

	NSString* nsText = [[NSString alloc] initWithUTF8String:text];
	[(NSText*)handle setText:nsText];
	[nsText release];
}

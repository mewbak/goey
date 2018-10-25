#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void* textNew( void* superview, char const* text ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( text );

	// Create the text view
	NSText* control = [[NSText alloc] init];
	[control setDrawsBackground:NO];
	textSetText( control, text );
	[control setEditable:NO];

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

void textSetAlignment( void* handle, int align ) {
	assert( handle && [(id)handle isKindOfClass:[NSText class]] );

	switch ( align ) {
	default:
	case 0:
		[(NSText*)handle setAlignment:NSTextAlignmentLeft];
		break;
	case 1:
		[(NSText*)handle setAlignment:NSTextAlignmentCenter];
		break;
	case 2:
		[(NSText*)handle setAlignment:NSTextAlignmentRight];
		break;
	case 3:
		[(NSText*)handle setAlignment:NSTextAlignmentJustified];
		break;
	}
}

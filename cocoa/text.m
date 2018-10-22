#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void* textNew( void* window, char const* text ) {
	NSString* nsText = [[NSString alloc] initWithUTF8String:text];

	// Create the text view
	NSText* control = [[NSText alloc] init];
	[control setText:nsText];
	[control setDrawsBackground:NO];

	// Add the control as the view for the window
	NSView* cv = [(NSWindow*)window contentView];
	[cv addSubview:control];

	return control;
}

void textClose( void* handle ) {
	NSText* control = handle;
	[control release];
}

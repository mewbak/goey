#import <Cocoa/Cocoa.h>
#include "cocoa.h"

void* buttonNew(void* window, char const* title) {
	NSString* nsTitle = [[NSString alloc] initWithUTF8String:title];

	// Create the button
	NSButton* control = [[NSButton alloc] init];
	[control setTitle:nsTitle];

	// Add the button as the view for the window
	[(NSWindow*)window setContentView:control];

	// Return handle to the control
	return control;
}

void buttonClose(void* handle) {
	printf("closeWindow\n");
	NSButton* control = handle;
	[control release];
}

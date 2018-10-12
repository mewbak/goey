#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@implementation NSButton ( Goey )

- (void)onclick {
	buttonOnClick( self );
}

- (BOOL)becomeFirstResponder {
	buttonOnFocus( self );
	return YES;
}

@end

void* buttonNew( void* window, char const* title ) {
	NSString* nsTitle = [[NSString alloc] initWithUTF8String:title];

	// Create the button
	NSButton* control = [[NSButton alloc] init];
	[control setTitle:nsTitle];
	[control setTarget:control];
	[control setAction:@selector( onclick )];

	// Add the button as the view for the window
	NSView* cv = [(NSWindow*)window contentView];
	[cv addSubview:control];

	// Return handle to the control
	return control;
}

void buttonClose( void* handle ) {
	printf( "closeWindow\n" );
	NSButton* control = handle;
	[control release];
}

void buttonSetTitle( void* handle, char const* title ) {
	NSString* nsTitle = [[NSString alloc] initWithUTF8String:title];
	[(NSButton*)handle setTitle:nsTitle];
}

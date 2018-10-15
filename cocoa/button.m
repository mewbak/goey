#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GButton : NSButton
- (BOOL)becomeFirstResponder;
- (BOOL)resignFirstResponder;
- (void)onclick;
@end

@implementation GButton

- (void)onclick {
	buttonOnClick( self );
}

- (BOOL)becomeFirstResponder {
	BOOL rc = [super becomeFirstResponder];
	if ( rc ) {
		buttonOnFocus( self );
	}
	return rc;
}

- (BOOL)resignFirstResponder {
	BOOL rc = [super resignFirstResponder];
	if ( rc ) {
		buttonOnBlur( self );
	}
	return rc;
}

@end

void* buttonNew( void* window, char const* title ) {
	NSString* nsTitle = [[NSString alloc] initWithUTF8String:title];

	// Create the button
	GButton* control = [[GButton alloc] init];
	[control setTitle:nsTitle];
	[control setTarget:control];
	[control setAction:@selector( onclick )];

	// Add the button as the view for the window
	NSView* cv = [(NSWindow*)window contentView];
	[cv addSubview:control];

	// Return handle to the control
	return control;
}

void buttonPerformClick( void* handle ) {
	[[(NSButton*)handle cell] performClick:nil];
}

char const* buttonTitle( void* handle ) {
	NSString* title = [(GButton*)handle title];
	return [title cString];
}

void buttonSetTitle( void* handle, char const* title ) {
	NSString* nsTitle = [[NSString alloc] initWithUTF8String:title];
	[(GButton*)handle setTitle:nsTitle];
}

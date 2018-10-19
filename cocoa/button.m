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

- (void)onchange {
	NSInteger s = [self state];
	buttonOnChange( self, s == NSOnState );
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
	assert( [(id)window isKindOfClass:[NSWindow class]] );

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

void* buttonNewCheck( void* window, char const* title, bool_t value ) {
	assert( [(id)window isKindOfClass:[NSWindow class]] );

	NSString* nsTitle = [[NSString alloc] initWithUTF8String:title];

	// Create the button
	GButton* control = [[GButton alloc] init];
	[control setButtonType:NSSwitchButton];
	[control setTitle:nsTitle];
	[control setTarget:control];
    [control setState:value];
	[control setAction:@selector( onchange )];

	// Add the button as the view for the window
	NSView* cv = [(NSWindow*)window contentView];
	[cv addSubview:control];

	// Return handle to the control
	return control;
}

void buttonPerformClick( void* handle ) {
	[[(GButton*)handle cell] performClick:nil];
}

bool_t buttonState( void* handle ) {
	NSInteger s = [(GButton*)handle state];
	return s == NSOnState;
}

void buttonSetState( void* handle, bool_t state ) {
	NSInteger s = state ? NSOnState : NSOffState;
	[(GButton*)handle setState:s];
}

char const* buttonTitle( void* handle ) {
	NSString* title = [(GButton*)handle title];
	return [title cString];
}

void buttonSetTitle( void* handle, char const* title ) {
	NSString* nsTitle = [[NSString alloc] initWithUTF8String:title];
	[(GButton*)handle setTitle:nsTitle];
}

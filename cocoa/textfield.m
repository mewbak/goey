#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GTextField : NSTextField <NSTextFieldDelegate>
- (BOOL)becomeFirstResponder;
- (BOOL)resignFirstResponder;
- (void)controlTextDidChange:(NSNotification*)obj;
@end

@implementation GTextField

- (void)controlTextDidChange:(NSNotification*)obj {
	NSString* v = [self stringValue];
	// Drop const, not representable in Go type system
	textfieldOnChange( self,
	                   (char*)[v cStringUsingEncoding:NSUTF8StringEncoding] );
}

- (BOOL)becomeFirstResponder {
	BOOL rc = [super becomeFirstResponder];
	if ( rc ) {
		textfieldOnFocus( self );
	}
	return rc;
}

- (BOOL)resignFirstResponder {
	BOOL rc = [super resignFirstResponder];
	if ( rc ) {
		textfieldOnBlur( self );
	}
	return rc;
}

@end

void* textfieldNew( void* superview, char const* text ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( text );

	// Create the button
	GTextField* control = [[GTextField alloc] init];
	textfieldSetValue( control, text );
	[control setEditable:YES];
	//[control setUsesSingleLineMode:YES];
	[control setDelegate:control];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

void textfieldSetValue( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[GTextField class]] );
	assert( text );

	NSString* title = [[NSString alloc] initWithUTF8String:text];
	[(GTextField*)handle setStringValue:title];
	[title release];
}

void textfieldSetPlaceholder( void* handle, char const* text ) {
	assert( handle && [(id)handle isKindOfClass:[GTextField class]] );
	assert( text );

	NSString* title = [[NSString alloc] initWithUTF8String:text];
	[[(GTextField*)handle cell] setPlaceholderString:title];
	[title release];
}

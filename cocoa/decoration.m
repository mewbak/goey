#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface GDecoration : NSView
- (void)drawRect:(NSRect)dirtyRect;
- (BOOL)isOpaque;
@end

@implementation GDecoration

- (void)drawRect:(NSRect)dirtyRect {
	printf( "%s\t%p\n", __func__, self );

	NSRect frame = [self frame];
	[[NSColor whiteColor] set];
	[NSBezierPath fillRect:frame];
}

- (BOOL)isOpaque {
	return YES;
}

@end

void* decorationNew( void* superview ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );

	// Create the button
	GDecoration* control = [[GDecoration alloc] init];
	printf( "new dec\t%p\n", control );

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

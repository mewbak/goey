#include "cocoa.h"
#import <Cocoa/Cocoa.h>

#define STYLE_MASK                                                             \
	( NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | \
	  NSResizableWindowMask )

void* windowNew( char const* title, unsigned width, unsigned height ) {
	NSString* appName = [[NSString alloc] initWithUTF8String:title];

	NSWindow* window =
	    [[[NSWindow alloc] initWithContentRect:NSMakeRect( 0, 0, width, height )
	                                 styleMask:STYLE_MASK
	                                   backing:NSBackingStoreBuffered
	                                     defer:NO] autorelease];
	[window cascadeTopLeftFromPoint:NSMakePoint( 20, 20 )];
	[window setTitle:appName];
	[window makeKeyAndOrderFront:nil];
	return window;
}

void windowClose( void* handle ) {
	NSWindow* window = handle;
	[window close];
}

int windowContentSize( void* handle, int* h ) {
	NSUInteger style = [(NSWindow*)handle styleMask];
	NSRect frame = [(NSWindow*)handle frame];
	frame = [NSWindow contentRectForFrameRect:frame styleMask:style];
	*h = frame.size.height;
	return frame.size.width;
}

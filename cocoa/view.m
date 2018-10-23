#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void viewClose( void* handle ) {
	assert( handle );
	assert( [(id)handle isKindOfClass:[NSView class]] );

	[(NSView*)handle removeFromSuperview];
	[(NSView*)handle release];
}

void viewSetFrame( void* handle, int x, int y, int dx, int dy ) {
	assert( [(id)handle isKindOfClass:[NSView class]] );

	NSRect frame = NSMakeRect( x, y, dx, dy );
	[(NSView*)handle setFrame:frame];
	[(NSView*)handle display];
}

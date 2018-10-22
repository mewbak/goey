#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void viewSetFrame( void* handle, int x, int y, int dx, int dy ) {
	assert( [(id)handle isKindOfClass:[NSView class]] );

	NSRect frame = NSMakeRect( x, y, dx, dy );
	[(NSView*)handle setFrame:frame];
	[(NSView*)handle display];
}

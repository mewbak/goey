#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void controlSetEnabled( void* handle, BOOL value ) {
	[(NSControl*)handle setEnabled:value];
}

void controlSetBounds( void* handle, int x, int y, int dx, int dy ) {
	NSRect frame = NSMakeRect( x, y, dx, dy );
	[(NSControl*)handle setFrame:frame];
	[(NSControl*)handle display];
}

int controlIntrinsicContentSize( void* handle, int* h ) {
	// Note that accessing the cell is deprecated, but GNUstep does not have
	// the newer methods needed to gather this information.
	NSCell* cell = [(NSControl*)handle cell];
	NSSize size = [cell cellSize];

	// Return the values
	*h = size.height;
	return size.width;
}

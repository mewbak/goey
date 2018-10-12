#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void viewSetFrame( void* handle, int x, int y, int dx, int dy ) {
	NSRect frame = NSMakeRect( x, y, dx, dy );
	[(NSControl*)handle setFrame:frame];
	[(NSControl*)handle display];
}

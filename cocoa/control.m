#import <Cocoa/Cocoa.h>
#include "cocoa.h"

void controlSetEnabled(void* handle, BOOL value) {
	[(NSControl*)handle setEnabled:value];
}

void controlSetBounds(void* handle, int x, int y, int dx, int dy) {
	NSRect frame = NSMakeRect( x, y, dx, dy );
	[(NSControl*)handle setFrame:frame];
	[(NSControl*)handle display];
}

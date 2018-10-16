#include "cocoa.h"
#import <Cocoa/Cocoa.h>

BOOL controlIsEnabled( void* handle ) {
	return [(NSControl*)handle isEnabled];
}

void controlSetEnabled( void* handle, BOOL value ) {
	[(NSControl*)handle setEnabled:value];
}

void controlClose( void* handle ) {
	assert( [(id)handle isKindOfClass:[NSControl class]] );
	[(NSControl*)handle release];
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

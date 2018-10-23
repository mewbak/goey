#include "cocoa.h"
#import <Cocoa/Cocoa.h>

bool_t controlIsEnabled( void* handle ) {
	return [(NSControl*)handle isEnabled];
}

void controlSetEnabled( void* handle, bool_t value ) {
	[(NSControl*)handle setEnabled:value];
}

nssize_t controlIntrinsicContentSize( void* handle ) {
	// Note that accessing the cell is deprecated, but GNUstep does not have
	// the newer methods needed to gather this information.
	NSCell* cell = [(NSControl*)handle cell];
	NSSize size = [cell cellSize];

	// Return the values
	nssize_t ret = {size.width, size.height};
	return ret;
}

#include "cocoa.h"
#import <Cocoa/Cocoa.h>

bool_t controlIsEnabled( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSControl class]] );

	return [(NSControl*)handle isEnabled];
}

void controlSetEnabled( void* handle, bool_t value ) {
	assert( handle && [(id)handle isKindOfClass:[NSControl class]] );

	[(NSControl*)handle setEnabled:value];
}

nssize_t controlIntrinsicContentSize( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSControl class]] );

	// Note that accessing the cell is deprecated, but GNUstep does not have
	// the newer methods needed to gather this information.
	NSCell* cell = [(NSControl*)handle cell];
	NSSize size = [cell cellSize];

	// Return the values
	nssize_t ret = {size.width, size.height};
	return ret;
}

bool_t controlMakeFirstResponder( void* handle ) {
	assert( handle && [(id)handle isKindOfClass:[NSControl class]] );

	NSWindow* window = [(NSControl*)handle window];
	assert( window );
	return [window makeFirstResponder:(NSControl*)handle];
}

#import <Cocoa/Cocoa.h>
#include "cocoa.h"

void controlSetEnabled(void* handle, BOOL value) {
	[(NSControl*)handle setEnabled:value];
}

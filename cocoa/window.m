#import <Cocoa/Cocoa.h>
#include "cocoa.h"

void* windowNew(char const* title, unsigned width, unsigned height) {
	NSString* appName = [[NSString alloc] initWithUTF8String:title];

	NSWindow* window = [[[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, width, height)
		styleMask:NSTitledWindowMask backing:NSBackingStoreBuffered defer:NO]
		autorelease];
	[window cascadeTopLeftFromPoint:NSMakePoint(20,20)];
	[window setTitle:appName];
	[window makeKeyAndOrderFront:nil];
	return window;
}

void windowClose(void* handle) {
	printf("closeWindow\n");
	NSWindow* window = handle;
	[window close];
}

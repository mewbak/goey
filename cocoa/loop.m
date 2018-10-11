#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void init() {
	NSString* quitString = [[NSString alloc] initWithUTF8String:"Quit "];
	NSString* qString = [[NSString alloc] initWithUTF8String:"q"];

	[NSAutoreleasePool new];
	[NSApplication sharedApplication];
	//[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

	id menubar = [[NSMenu new] autorelease];
	id appMenuItem = [[NSMenuItem new] autorelease];
	[menubar addItem:appMenuItem];
	[NSApp setMainMenu:menubar];
	id appMenu = [[NSMenu new] autorelease];
	id appName = [[NSProcessInfo processInfo] processName];
	id quitTitle = [quitString stringByAppendingString:appName];
	id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
	                                              action:@selector( terminate: )
	                                       keyEquivalent:qString] autorelease];
	[appMenu addItem:quitMenuItem];
	[appMenuItem setSubmenu:appMenu];
}

void run() {
	[NSApp activateIgnoringOtherApps:YES];
	[NSApp run];
}

void stop() {
	[NSApp stop:nil];
}

@interface DoActionOperation : NSOperation {
	void* action;
	void* err;
}
@end

@implementation DoActionOperation

- (id)init {
	self = [super init];
	return self;
}

- (id)main {
	printf( "Ta da...\n" );
	callbackDo();
}

@end

void thunkDo() {
	[NSAutoreleasePool new];

	// DoActionEvent* event = [[DoActionEvent alloc] init:action Err:err];
	// SEL sel = @selector(call);
	//[event performSelectorOnMainThread:sel];

	NSOperation* operation = [[DoActionOperation alloc] init];
	NSOperationQueue* targetQueue = [NSOperationQueue mainQueue];
	[targetQueue addOperation:operation];
}

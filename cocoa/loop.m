#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>
#include <assert.h>

@interface GNOPThread : NSThread
- (void)main;
@end

@implementation GNOPThread

- (void)main {
	// Do nothing.  This is a NOP thread.
	return;
}

@end

static void detachAThread() {
	// We need to make sure that Cocoa is running multithreaded.  Otherwise,
	// use of autopool from other threads will not work propertly.  The notes
	// for NSAutoreleasePool indicate that we need to detach a thread to
	// cause this transition.
	NSThread* thread = [[GNOPThread alloc] init];
	[thread start];
	[thread release];
}

static void initApplication() {
	NSString* quitString = [[NSString alloc] initWithUTF8String:"Quit "];
	NSString* qString = [[NSString alloc] initWithUTF8String:"q"];

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

void init() {
	detachAThread();

	// This is a global release pool that we will keep around.  This will still
	// cause leaks, but until we identify where the autoreleasepool is required,
	// this will get us running.
	NSAutoreleasePool* pool = [NSAutoreleasePool new];
	initApplication();
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
	callbackDo();
}

@end

void thunkDo() {
	assert( [NSThread isMultiThreaded] );
	NSAutoreleasePool* pool = [NSAutoreleasePool new];

	NSOperation* operation = [[DoActionOperation alloc] init];
	NSOperationQueue* targetQueue = [NSOperationQueue mainQueue];
	[targetQueue addOperation:operation];
	[operation release];

	[pool release];
}

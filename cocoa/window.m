#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

#define STYLE_MASK                                                             \
	( NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | \
	  NSResizableWindowMask )

@implementation NSWindow ( Goey )

- (BOOL)windowShouldClose:(id)sender {
	if ( windowShouldClose( sender ) != 0 ) {
		return YES;
	}
	return NO;
}

@end

static void adjustContentOrigin( NSWindow* window ) {
	NSView* cv = window.contentView;
	CGFloat h = cv.frame.size.height;
	[cv setBoundsOrigin:NSMakePoint( 0, -h )];
}

@interface MyWindowDelegate : NSObject <NSWindowDelegate>
- (void)windowWillClose:(NSNotification*)aNotification;
- (void)windowDidResize:(NSNotification*)aNotification;
@end

@implementation MyWindowDelegate

- (void)windowWillClose:(NSNotification*)notification {
	NSWindow* window = [notification object];
	windowWillClose( window );
}

- (void)windowDidResize:(NSNotification*)notification {
	NSWindow* window = [notification object];
	adjustContentOrigin( window );
	windowDidResize( window );
}

@end

void* windowNew( char const* title, unsigned width, unsigned height ) {
	TRACE();

	assert( [NSThread isMainThread] );
	assert( title );

	// Make sure that we have a delegate
	static MyWindowDelegate* delegate = 0;
	if ( !delegate ) {
		delegate = [[MyWindowDelegate alloc] init];
	}

	NSWindow* window =
	    [[NSWindow alloc] initWithContentRect:NSMakeRect( 0, 0, width, height )
	                                styleMask:STYLE_MASK
	                                  backing:NSBackingStoreBuffered
	                                    defer:NO];
	[window cascadeTopLeftFromPoint:NSMakePoint( 20, 20 )];
	[window setDelegate:delegate];
	windowSetTitle( window, title );
	[window makeKeyAndOrderFront:nil];
	adjustContentOrigin( window );
	return window;
}

void windowClose( void* handle ) {
	TRACE();

	assert( [NSThread isMainThread] );
	assert( handle );

	// This will send a blur message to any controls that currently have focus.
	[(NSWindow*)handle makeFirstResponder:NULL];

	// This call to close the window should also release.
	[(NSWindow*)handle close];
}

nssize_t windowContentSize( void* handle ) {
	assert( [NSThread isMainThread] );
	assert( handle );

	NSUInteger style = [(NSWindow*)handle styleMask];
	NSRect frame = [(NSWindow*)handle frame];
	frame = [NSWindow contentRectForFrameRect:frame styleMask:style];

	nssize_t ret = {frame.size.width, frame.size.height};
	return ret;
}

void* windowContentView( void* handle ) {
	assert( [NSThread isMainThread] );
	assert( handle );
	assert( [(id)handle isKindOfClass:[NSWindow class]] );

	return [(NSWindow*)handle contentView];
}

void windowMakeFirstResponder( void* window, void* handle2 ) {
	assert( [NSThread isMainThread] );
	assert( window );
	assert( [(id)window isKindOfClass:[NSWindow class]] );

	NSWindow* w = (NSWindow*)window;
	NSControl* c = (NSControl*)handle2;

	[w makeFirstResponder:c];
}

void windowSetMinSize( void* handle, int width, int height ) {
	assert( [NSThread isMainThread] );
	assert( handle );

	NSWindow* w = (NSWindow*)handle;

	// Adjust size from content to outer frame
	NSRect frame = NSMakeRect( 0, 0, width, height );
	frame = [NSWindow frameRectForContentRect:frame styleMask:[w styleMask]];
	[w setMinSize:NSMakeSize( NSWidth( frame ), NSHeight( frame ) )];
}

void windowSetIconImage( void* handle, void* nsimage ) {
	assert( [NSThread isMainThread] );
	assert( handle );
	assert( nsimage );

	[NSApp setApplicationIconImage:(NSImage*)nsimage];
}

void windowSetTitle( void* handle, char const* title ) {
	assert( [NSThread isMainThread] );
	assert( handle );

	NSString* wtitle = [[NSString alloc] initWithUTF8String:title];
	[(NSWindow*)handle setTitle:wtitle];
	[wtitle release];
}

char const* windowTitle( void* handle ) {
	assert( [NSThread isMainThread] );
	assert( handle );

	char const* cstring =
	    [[(NSWindow*)handle title] cStringUsingEncoding:NSUTF8StringEncoding];
	assert( cstring );
	return cstring;
}

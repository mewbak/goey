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
	windowDidResize( window );
}

@end

void* windowNew( char const* title, unsigned width, unsigned height ) {
	TRACE();

	assert( [NSThread isMainThread] );
	assert( title );

	// Make sure that we have a delegate.  This is required to respond to
	// window events.
	static MyWindowDelegate* delegate = 0;
	if ( !delegate ) {
		delegate = [[MyWindowDelegate alloc] init];
	}

	// Allocate and initialize the window.
	NSWindow* window =
	    [[NSWindow alloc] initWithContentRect:NSMakeRect( 0, 0, width, height )
	                                styleMask:STYLE_MASK
	                                  backing:NSBackingStoreBuffered
	                                    defer:NO];
	[window cascadeTopLeftFromPoint:NSMakePoint( 20, 20 )];
	[window setDelegate:delegate];
	windowSetTitle( window, title );

	// Replace the content view of the windows with a scrolled view.  This will
	// provide scrollbars and scrolling when necessary.
	NSScrollView* sw = [[NSScrollView alloc] init];
	NSView* swv = [[NSView alloc] init];
	[sw setDocumentView:swv];
	[window setContentView:sw];
	[swv release];
	[sw release];

	[window makeKeyAndOrderFront:nil];
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

	NSView* sw = [(NSWindow*)handle contentView];
	assert( sw && [sw isKindOfClass:[NSScrollView class]] );
	return [(NSScrollView*)sw documentView];
}

void windowMakeFirstResponder( void* window, void* handle2 ) {
	assert( [NSThread isMainThread] );
	assert( window && [(id)window isKindOfClass:[NSWindow class]] );

	NSWindow* w = (NSWindow*)window;
	NSControl* c = (NSControl*)handle2;

	[w makeFirstResponder:c];
}

void windowSetContentSize( void* handle, int width, int height ) {
	assert( [NSThread isMainThread] );
	assert( handle && [(id)handle isKindOfClass:[NSWindow class]] );

	NSView* sw = [(NSWindow*)handle contentView];
	assert( sw && [sw isKindOfClass:[NSScrollView class]] );
	[[(NSScrollView*)sw documentView]
	    setFrame:NSMakeRect( 0, 0, width, height )];
}

void windowSetMinSize( void* handle, int width, int height ) {
	assert( [NSThread isMainThread] );
	assert( handle && [(id)handle isKindOfClass:[NSWindow class]] );

	NSWindow* w = (NSWindow*)handle;

	// Adjust size from content to outer frame
	NSRect frame = NSMakeRect( 0, 0, width, height );
	frame = [NSWindow frameRectForContentRect:frame styleMask:[w styleMask]];
	[w setMinSize:NSMakeSize( NSWidth( frame ), NSHeight( frame ) )];
}

void windowSetIconImage( void* handle, void* nsimage ) {
	assert( [NSThread isMainThread] );
	assert( handle && [(id)handle isKindOfClass:[NSWindow class]] );
	assert( nsimage && [(id)nsimage isKindOfClass:[NSImage class]] );

	[NSApp setApplicationIconImage:(NSImage*)nsimage];
}

void windowSetScrollVisible( void* handle, bool_t horz, bool_t vert ) {
	assert( [NSThread isMainThread] );
	assert( handle && [(id)handle isKindOfClass:[NSWindow class]] );

	NSView* sw = [(NSWindow*)handle contentView];
	assert( sw && [sw isKindOfClass:[NSScrollView class]] );
	[(NSScrollView*)sw setHasHorizontalScroller:horz];
	[(NSScrollView*)sw setHasVerticalScroller:vert];
}

void windowSetTitle( void* handle, char const* title ) {
	assert( [NSThread isMainThread] );
	assert( handle && [(id)handle isKindOfClass:[NSWindow class]] );

	NSString* wtitle = [[NSString alloc] initWithUTF8String:title];
	[(NSWindow*)handle setTitle:wtitle];
	[wtitle release];
}

char const* windowTitle( void* handle ) {
	assert( [NSThread isMainThread] );
	assert( handle && [(id)handle isKindOfClass:[NSWindow class]] );

	char const* cstring =
	    [[(NSWindow*)handle title] cStringUsingEncoding:NSUTF8StringEncoding];
	assert( cstring );
	return cstring;
}

#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

@interface MyTabViewDelegate : NSObject /*<NSTabViewDelegate>*/
- (void)tabView:(NSTabView*)tabView
    didSelectTabViewItem:(NSTabViewItem*)tabViewItem;
@end

@implementation MyTabViewDelegate

- (void)tabView:(NSTabView*)tabView
    didSelectTabViewItem:(NSTabViewItem*)tabViewItem {
	NSInteger index = [tabView indexOfTabViewItem:tabViewItem];
	tabviewDidSelectItem( tabView, index );
}

- (void)windowDidResize:(NSNotification*)notification {
	NSWindow* window = [notification object];
	windowDidResize( window );
}

@end

void* tabviewNew( void* superview ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );

	// Make sure that we have a delegate
	static MyTabViewDelegate* delegate = 0;
	if ( !delegate ) {
		delegate = [[MyTabViewDelegate alloc] init];
	}

	// Create the button.
	// A default frame is required.  Otherwise, we get errors about negative
	// size frames when tabs are created, which have their own content views.
	// Also required for the bounds of those frames to be correctely
	// initialized.
	NSTabView* control =
	    [[NSTabView alloc] initWithFrame:NSMakeRect( 0, 0, 200, 200 )];
	[control setDelegate:delegate];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

void tabviewAddItem( void* control, char const* text ) {
	assert( control && [(id)control isKindOfClass:[NSTabView class]] );
	assert( text );

	NSString* label = [[NSString alloc] initWithUTF8String:text];
	NSTabViewItem* item = [[NSTabViewItem alloc] init];
	[item setLabel:label];
	[label release];
	[(NSTabView*)control addTabViewItem:item];
	[item release];
}

void tabviewSelectItem( void* control, int index ) {
	assert( control && [(id)control isKindOfClass:[NSTabView class]] );

	[(NSTabView*)control selectTabViewItemAtIndex:index];
}

void* tabviewContentView( void* control, int index ) {
	assert( control && [(id)control isKindOfClass:[NSTabView class]] );

	NSTabViewItem* item = [(NSTabView*)control selectedTabViewItem];
	assert( item );
	NSView* cv = [item view];
	assert( cv );
	return cv;
}

nssize_t tabviewContentInsets( void* control ) {
	assert( control && [(id)control isKindOfClass:[NSTabView class]] );

	NSRect frame = [(NSTabView*)control frame];
	NSRect cr = [(NSTabView*)control contentRect];
	nssize_t ret = {frame.size.width - cr.size.width,
	                frame.size.height - cr.size.height};
	return ret;
}

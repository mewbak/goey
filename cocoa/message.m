#include "cocoa.h"
#import <Cocoa/Cocoa.h>

void messageDialog( void* window, char const* text, char const* title, char icon ) {
	assert( !window || [(id)window isKindOfClass:[NSWindow class]] );
	assert( text );
	assert( title );

	NSAlert* alert = [[NSAlert alloc] init];

	NSString* tmp = [[NSString alloc] initWithUTF8String:title];
	[alert setMessageText:tmp];
	[tmp release];
	tmp = [[NSString alloc] initWithUTF8String:text];
	[alert setInformativeText:tmp];
	[tmp release];

	switch ( icon ) {
	case 'e':
		[alert setAlertStyle:NSCriticalAlertStyle];
		break;
	case 'w':
		[alert setAlertStyle:NSWarningAlertStyle];
		break;
	case 'i':
		[alert setAlertStyle:NSInformationalAlertStyle];
		break;
	}

	NSString* ok = [[[NSString alloc] initWithUTF8String:"OK"] autorelease];
	[alert addButtonWithTitle:ok];
	[alert runModal];
	[alert release];
}

void openPanel( void* window ) {
	assert( !window || [(id)window isKindOfClass:[NSWindow class]] );

	NSOpenPanel* panel = [NSOpenPanel openPanel];

	[panel runModal];
	[panel release];
}

void savePanel( void* window ) {
	assert( !window || [(id)window isKindOfClass:[NSWindow class]] );

	NSSavePanel* panel = [NSSavePanel savePanel];

	[panel runModal];
	[panel release];
}

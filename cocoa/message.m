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

static void setFilename( NSSavePanel* panel, char const* dir, char const* base ) {
	if ( dir ) {
		assert( base );

		NSString* tmp = [[NSString alloc] initWithUTF8String:dir];
		[panel setDirectory:tmp];
		[tmp release];
		tmp = [[NSString alloc] initWithUTF8String:base];
		[panel setNameFieldStringValue:tmp];
		[tmp release];
	}
}

char const* openPanel( void* window, char const* dir, char const* base ) {
	assert( !window || [(id)window isKindOfClass:[NSWindow class]] );

	NSOpenPanel* panel = [NSOpenPanel openPanel];
	setFilename( panel, dir, base );

	[panel runModal];
	return [[panel filename] cStringUsingEncoding:NSUTF8StringEncoding];
}

char const* savePanel( void* window, char const* dir, char const* base ) {
	assert( !window || [(id)window isKindOfClass:[NSWindow class]] );

	NSSavePanel* panel = [NSSavePanel savePanel];
	setFilename( panel, dir, base );

	[panel runModal];
	return [[panel filename] cStringUsingEncoding:NSUTF8StringEncoding];
}

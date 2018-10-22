#include "cocoa.h"
#import <Cocoa/Cocoa.h>
#include <assert.h>

void* imageNewFromRGBA( void* imageData, int width, int height ) {
	NSImageRep* imagerep = [[NSBitmapImageRep alloc]
	    initWithBitmapDataPlanes:(unsigned char**)&imageData
	                  pixelsWide:width
	                  pixelsHigh:height
	               bitsPerSample:8
	             samplesPerPixel:4
	                    hasAlpha:YES
	                    isPlanar:NO
	              colorSpaceName:NSDeviceRGBColorSpace
	                 bytesPerRow:4 * width
	                bitsPerPixel:32];
	assert( imagerep );

	NSImage* image = [[NSImage alloc]
	    initWithSize:NSMakeSize( width / 92.0, height / 92.0 )];
	assert( image );
	[image addRepresentation:imagerep];
	return image;
}

void* imageNewFromGray( void* imageData, int width, int height ) {
	NSImageRep* imagerep = [[NSBitmapImageRep alloc]
	    initWithBitmapDataPlanes:(unsigned char**)&imageData
	                  pixelsWide:width
	                  pixelsHigh:height
	               bitsPerSample:8
	             samplesPerPixel:1
	                    hasAlpha:NO
	                    isPlanar:NO
	              colorSpaceName:NSDeviceRGBColorSpace
	                 bytesPerRow:width
	                bitsPerPixel:8];
	assert( imagerep );

	NSImage* image = [[NSImage alloc]
	    initWithSize:NSMakeSize( width / 92.0, height / 92.0 )];
	assert( image );
	[image addRepresentation:imagerep];
	return image;
}

void* imageviewNew( void* window, void* image ) {
	// Create the control
	NSImageView* control = [[NSImageView alloc] init];
	[control setImage:(NSImage*)image];

	// Add the button as the view for the window
	NSView* cv = [(NSWindow*)window contentView];
	[cv addSubview:control];

	// Return handle to the control
	return control;
}

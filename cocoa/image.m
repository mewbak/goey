#include "cocoa.h"
#import <Cocoa/Cocoa.h>
#include <assert.h>

void* imageNewFromRGBA( uint8_t* imageData, int width, int height ) {
	assert( imageData );

	NSImageRep* imagerep =
	    [[NSBitmapImageRep alloc] initWithBitmapDataPlanes:&imageData
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

	NSImage* image = [[NSImage alloc] initWithSize:NSMakeSize( width, height )];
	assert( image );
	[image addRepresentation:imagerep];
	return image;
}

void* imageNewFromGray( uint8_t* imageData, int width, int height ) {
	assert( imageData );

	NSImageRep* imagerep =
	    [[NSBitmapImageRep alloc] initWithBitmapDataPlanes:&imageData
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

	NSImage* image = [[NSImage alloc] initWithSize:NSMakeSize( width, height )];
	assert( image );
	[image addRepresentation:imagerep];
	return image;
}

void imageClose( void* image ) {
	assert( image );
	assert( [(id)image isKindOfClass:[NSImage class]] );

	[(NSImage*)image release];
}

void* imageviewNew( void* window, void* image ) {
	assert( window );
	assert( [(id)window isKindOfClass:[NSWindow class]] );
	assert( image );
	assert( [(id)image isKindOfClass:[NSImage class]] );

	// Create the control
	NSImageView* control = [[NSImageView alloc] init];
	[control setImage:(NSImage*)image];
    [control setImageScaling:NSImageScaleAxesIndependently];

	// Add the button as the view for the window
	NSView* cv = [(NSWindow*)window contentView];
	[cv addSubview:control];

	// Return handle to the control
	return control;
}

void imageviewSetImage( void* control, void* image ) {
	assert( control );
	assert( image );

	[(NSImageView*)control setImage:(NSImage*)image];
}

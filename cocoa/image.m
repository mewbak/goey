#include "cocoa.h"
#import <Cocoa/Cocoa.h>
#include <assert.h>

void* imageNewFromRGBA( uint8_t* imageData, int width, int height,
                        int stride ) {
	assert( imageData );

	NSBitmapImageRep* imagerep =
	    [[NSBitmapImageRep alloc] initWithBitmapDataPlanes:NULL
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

	// Copy over the image data.
	if ( [imagerep bytesPerRow] == stride ) {
		// Check for overflow
		assert( width * height > width && width * height > height );
		assert( ( width * height * 4 ) > width * height );
		// Copy the data
		assert( [imagerep bitmapData] );
		memcpy( [imagerep bitmapData], imageData, width * height * 4 );
	} else {
		assert( false ); // not implemented
	}

	// Create the image
	NSImage* image = [[NSImage alloc] initWithSize:NSMakeSize( width, height )];
	assert( image );
	[image addRepresentation:imagerep];
	return image;
}

void* imageNewFromGray( uint8_t* imageData, int width, int height,
                        int stride ) {
	assert( imageData );

	NSBitmapImageRep* imagerep =
	    [[NSBitmapImageRep alloc] initWithBitmapDataPlanes:NULL
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

	// Copy over the image data.
	if ( [imagerep bytesPerRow] == stride ) {
		// Check for overflow
		assert( width * height > width && width * height > height );
		assert( ( width * height * 4 ) > width * height );
		// Copy the data
		assert( [imagerep bitmapData] );
		memcpy( [imagerep bitmapData], imageData, width * height * 4 );
	} else {
		assert( false ); // not implemented
	}

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

void* imageviewNew( void* superview, void* image ) {
	assert( superview && [(id)superview isKindOfClass:[NSView class]] );
	assert( image && [(id)image isKindOfClass:[NSImage class]] );

	// Create the control
	NSImageView* control = [[NSImageView alloc] init];
	[control setImage:(NSImage*)image];
	[control setImageScaling:NSImageScaleAxesIndependently];

	// Add the button as the view for the window
	[(NSView*)superview addSubview:control];

	// Return handle to the control
	return control;
}

void imageviewSetImage( void* control, void* image ) {
	assert( control );
	assert( image );

	[(NSImageView*)control setImage:(NSImage*)image];
}

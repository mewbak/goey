#ifndef GOEY_COCOA_H
#define GOEY_COCOA_H

#include <stdint.h>

// Cannot use std bool.  The builtin type _Bool does not play well with CGO.
// Need an alternate for the binding.

typedef unsigned bool_t;

typedef struct nssize_tag {
	int32_t width;
	int32_t height;
} nssize_t;

/* Event loop */
extern void init( void );
extern void run( void );
extern void do_thunk( void );
extern void stop( void );
extern bool_t isMainThread( void );

extern void trace( char const* func );

/* Window */
extern void* windowNew( char const* title, unsigned width, unsigned height );
extern void windowClose( void* handle );
extern nssize_t windowContentSize( void* handle );
extern void windowMakeFirstResponder( void* handle, void* control );
extern void windowSetMinSize( void* handle, int width, int height );

/* View */
extern void viewSetFrame( void* handle, int x, int y, int dx, int dy );

/* Control */
extern bool_t controlIsEnabled( void* handle );
extern void controlSetEnabled( void* handle, bool_t value );
extern void controlClose( void* handle );
extern nssize_t controlIntrinsicContentSize( void* handle );

/* Button */
extern void* buttonNew( void* window, char const* title );
extern void* buttonNewCheck( void* window, char const* title, bool_t value );
extern void buttonPerformClick( void* handle );
extern bool_t buttonState( void* handle );
extern void buttonSetState( void* handle, bool_t checked );
extern char const* buttonTitle( void* handle );
extern void buttonSetTitle( void* handle, char const* title );

/* Text */
extern void* textNew( void* window, char const* text );
extern void textClose( void* handle );

/* TextField */
extern void* textfieldNew( void* window, char const* text );

#endif

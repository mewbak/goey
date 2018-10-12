#include <objc/objc.h>

/* Event loop */
extern void init();
extern void run();
extern void thunkDo();
extern void stop();

/* Window */
extern void* windowNew( char const* title, unsigned width, unsigned height );
extern void windowClose( void* handle );
extern int windowContentSize( void* handle, int* h );

/* View */
extern void viewSetFrame( void* handle, int x, int y, int dx, int dy );

/* Control */
extern void controlSetEnabled( void* handle, BOOL value );
extern int controlIntrinsicContentSize( void* handle, int* h );

/* Button */
extern void* buttonNew( void* window, char const* title );
extern void buttonClose( void* handle );
extern void buttonOnClick( void* handle );
extern void buttonSetTitle( void* handle, char const* title );

/* Text */
extern void* textNew( void* window, char const* text );
extern void textClose( void* handle );

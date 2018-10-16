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
extern void windowMakeFirstResponder( void* handle, void* control );
extern void windowSetMinSize( void* handle, int width, int height );

/* View */
extern void viewSetFrame( void* handle, int x, int y, int dx, int dy );

/* Control */
extern BOOL controlIsEnabled( void* handle );
extern void controlSetEnabled( void* handle, BOOL value );
extern void controlClose( void* handle );
extern int controlIntrinsicContentSize( void* handle, int* h );

/* Button */
extern void* buttonNew( void* window, char const* title );
extern void* buttonNewCheck( void* window, char const* title );
extern void buttonPerformClick( void* handle );
extern BOOL buttonState( void* handle );
extern void buttonSetState( void* handle, BOOL checked );
extern char const* buttonTitle( void* handle );
extern void buttonSetTitle( void* handle, char const* title );

/* Text */
extern void* textNew( void* window, char const* text );
extern void textClose( void* handle );

/* TextField */
extern void* textfieldNew( void* window, char const* text );

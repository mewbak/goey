#include <objc/objc.h>

/* Event loop */
extern void init();
extern void run();
extern void thunkDo();
extern void stop();

/* Window */
extern void* windowNew(char const* title, unsigned width, unsigned height);
extern void windowClose(void* handle);

/* Control */
extern void controlSetEnabled(void* handle, BOOL value);

/* Button */
extern void* buttonNew(void* window, char const* title);
extern void buttonClose(void* handle);
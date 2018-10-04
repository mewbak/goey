/* Event loop */
extern void init();
extern void run();
extern void thunkDo();
extern void stop();

/* Window */
extern void* windowNew(char const* title, unsigned width, unsigned height);
extern void windowClose(void* handle);

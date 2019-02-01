#ifndef GOEY_INTERNAL_SYSCALL_LINUX_H
#define GOEY_INTERNAL_SYSCALL_LINUX_H

#include <gtk/gtk.h>

extern void goey_set_bounds( GtkWidget *handle, gint x, gint y, gint width, gint height );
extern void goey_set_key_info( GdkEventKey *evt, GdkWindow *window, guint r );
extern void goey_widget_send_key( GtkWidget *widget, guint r, GdkModifierType modifiers, gchar release );
extern void goey_main_context_invoke( void );
extern void goey_idle_add( void );

#endif

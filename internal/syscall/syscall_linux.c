#include <gdk/gdk.h>
#include <gtk/gtk.h>
#include "syscall_linux.h"
#include "_cgo_export.h"

void goey_set_bounds( GtkWidget *handle, gint x, gint y, gint width, gint height ) {
    GtkWidget *parent = gtk_widget_get_parent( handle );
    GtkLayout *layout = GTK_LAYOUT( parent );
    gtk_layout_move( layout, handle, x, y );
    GtkAllocation alloc = { x, y, width, height };
    gtk_widget_size_allocate( handle, &alloc );
}

void goey_set_key_info( GdkEventKey *evt, GdkWindow *window, guint r ) {
    evt->window = window;
    evt->time = GDK_CURRENT_TIME;
    evt->send_event = 1;

    switch (r) {
    case 0x1b:
        evt->keyval = GDK_KEY_Escape;
        evt->hardware_keycode = 9;
        break;
    case '\n':
        evt->keyval = GDK_KEY_Return;
        evt->hardware_keycode = 36;
        break;
    default:
        evt->keyval = r;
        break;
    }
}

void goey_widget_send_key( GtkWidget *widget, guint r, GdkModifierType modifiers, gchar release ) {
    GdkEvent *evt = gdk_event_new( release ? GDK_KEY_RELEASE : GDK_KEY_PRESS );
    goey_set_key_info( (GdkEventKey *)evt, gtk_widget_get_window(widget), r );
    gtk_widget_event( widget, evt );
}

static gboolean goey_main_context_invoke_cb(gpointer user_data) {
    // Callback into Go.
    mainContextInvokeCallback();
    // Prevent a repeat of this event.
    return FALSE;
}

void goey_main_context_invoke( void ) {
    g_main_context_invoke( NULL, goey_main_context_invoke_cb, NULL );
}

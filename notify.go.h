#ifndef GO_GLIB_H
#define GO_GLIB_H



#include <glib.h>
#include <glib-object.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

static GError* to_error(void* err) {
		return (GError*)err;
}


static inline char* to_charptr(const gchar* s) { return (char*)s; }

#endif

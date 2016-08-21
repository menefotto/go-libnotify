package notify

/*
#cgo pkg-config: libnotify

#include <libnotify/notify.h>
#include "notify.go.h"
*/
import "C"
import "unsafe"

// here to eleminate the glib-go dependency
type Error struct {
	GError *C.GError
}

func (v *Error) Error() string {
	return v.message()
}

func (v *Error) message() string {
	if unsafe.Pointer(v.GError) == nil || unsafe.Pointer(v.GError.message) == nil {
		return ""
	}
	return C.GoString(C.to_charptr(v.GError.message))
}

func ErrorFromNative(err unsafe.Pointer) error {
	return &Error{GError: C.to_error(err)}
}

package notify

/*
#cgo pkg-config: libnotify

#include <libnotify/notify.h>
#include "notify.go.h"
*/
import "C"
import "unsafe"

// part of the notify api
func Init(app_name string) bool {
	papp_name := C.CString(app_name)
	defer C.free(unsafe.Pointer(papp_name))

	return bool(C.notify_init(papp_name) != 0)
}

func UnInit() {
	C.notify_uninit()
}

func IsInitted() bool {
	return C.notify_is_initted() != 0
}

func GetAppName() string {
	return C.GoString(C.notify_get_app_name())
}

/*
func GetServerCaps() *glib.List {
	var gcaps *C.GList

	gcaps = C.notify_get_server_caps()

	return glib.ListFromNative(unsafe.Pointer(gcaps))
}
*/
func GetServerInfo(name, vendor, version, spec_version *string) bool {
	var cname *C.char
	var cvendor *C.char
	var cversion *C.char
	var cspec_version *C.char

	ret := C.notify_get_server_info(&cname,
		&cvendor,
		&cversion,
		&cspec_version) != 0
	*name = C.GoString(cname)
	*vendor = C.GoString(cvendor)
	*version = C.GoString(cversion)
	*spec_version = C.GoString(cspec_version)

	return ret
}

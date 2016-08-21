/*
 * notify.go for go-notify
 * by lenorm_f
 */

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
	return v.Message()
}

func (v *Error) Message() string {
	if unsafe.Pointer(v.GError) == nil || unsafe.Pointer(v.GError.message) == nil {
		return ""
	}
	return C.GoString(C.to_charptr(v.GError.message))
}

func ErrorFromNative(err unsafe.Pointer) *Error {
	return &Error{
		C.to_error(err)}
}

/*
 * Exported Types
 */
type Notification struct {
	_notification *C.NotifyNotification
}

const (
	NOTIFY_URGENCY_LOW      = 0
	NOTIFY_URGENCY_NORMAL   = 1
	NOTIFY_URGENCY_CRITICAL = 2
)

type Urgency int
type ActionCallback func(*Notification, string, interface{})

/*
 * Private Functions
 */
func new_notification(cnotif *C.NotifyNotification) *Notification {
	return &Notification{cnotif}
}

/*
 * Exported Functions
 */
// Pure Functions
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

func NotificationNew(title, text, image string) *Notification {
	ptitle := C.CString(title)
	ptext := C.CString(text)
	pimage := C.CString(image)
	ntext := C.g_utf8_normalize((*C.gchar)(ptext), -1, C.G_NORMALIZE_DEFAULT)
	defer func() {
		C.free(unsafe.Pointer(ptitle))
		C.free(unsafe.Pointer(ptext))
		C.free(unsafe.Pointer(pimage))

		if ntext != nil {
			C.free(unsafe.Pointer(ntext))
		}
	}()

	return new_notification(C.notify_notification_new(ptitle, (*C.char)(ntext), pimage))
}

func NotificationUpdate(notif *Notification, summary, body, icon string) bool {
	psummary := C.CString(summary)
	pbody := C.CString(body)
	picon := C.CString(icon)
	defer func() {
		C.free(unsafe.Pointer(psummary))
		C.free(unsafe.Pointer(pbody))
		C.free(unsafe.Pointer(picon))
	}()

	return C.notify_notification_update(notif._notification, psummary, pbody, picon) != 0
}

func NotificationShow(notif *Notification) *Error {
	var err *C.GError
	C.notify_notification_show(notif._notification, &err)

	return ErrorFromNative(unsafe.Pointer(err))
}

func NotificationSetTimeout(notif *Notification, timeout int32) {
	C.notify_notification_set_timeout(notif._notification, C.gint(timeout))
}

func NotificationSetCategory(notif *Notification, category string) {
	pcategory := C.CString(category)
	defer C.free(unsafe.Pointer(pcategory))

	C.notify_notification_set_category(notif._notification, pcategory)
}

func NotificationSetUrgency(notif *Notification, urgency Urgency) {
	C.notify_notification_set_urgency(notif._notification, C.NotifyUrgency(urgency))
}

func NotificationSetHintInt32(notif *Notification, key string, value int32) {
	pkey := C.CString(key)
	defer C.free(unsafe.Pointer(pkey))

	C.notify_notification_set_hint_int32(notif._notification, pkey, C.gint(value))
}

func NotificationSetHintDouble(notif *Notification, key string, value float64) {
	pkey := C.CString(key)
	defer C.free(unsafe.Pointer(pkey))

	C.notify_notification_set_hint_double(notif._notification, pkey, C.gdouble(value))
}

func NotificationSetHintString(notif *Notification, key string, value string) {
	pkey := C.CString(key)
	pvalue := C.CString(value)
	defer func() {
		C.free(unsafe.Pointer(pkey))
		C.free(unsafe.Pointer(pvalue))
	}()

	C.notify_notification_set_hint_string(notif._notification, pkey, pvalue)
}

func NotificationSetHintByte(notif *Notification, key string, value byte) {
	pkey := C.CString(key)
	defer C.free(unsafe.Pointer(pkey))

	C.notify_notification_set_hint_byte(notif._notification, pkey, C.guchar(value))
}

// FIXME: implement
func NotificationSetHintByteArray(notif *Notification, key string, value []byte, len uint32) {
	pkey := C.CString(key)
	defer C.free(unsafe.Pointer(pkey))

	// C.notify_notification_set_hint_byte_array(notif._notification, pkey, (*C.guchar)(value), C.gsize(len))
}

func NotificationSetHint(notif *Notification, key string, value interface{}) {
	switch value.(type) {
	case int32:
		NotificationSetHintInt32(notif, key, value.(int32))
	case float64:
		NotificationSetHintDouble(notif, key, value.(float64))
	case string:
		NotificationSetHintString(notif, key, value.(string))
	case byte:
		NotificationSetHintByte(notif, key, value.(byte))
	}
}

func NotificationClearHints(notif *Notification) {
	C.notify_notification_clear_hints(notif._notification)
}

// FIXME: the C function is supposed to be allowing the user to pass another function than free
func NotificationAddAction(notif *Notification, action, label string,
	callback ActionCallback, user_data interface{}) {
	// C.notify_notification_add_action(notif._notification, paction, plabel, (C.NotifyActionCallback)(callback), user_data, C.free)
}

func NotificationClearActions(notif *Notification) {
	C.notify_notification_clear_actions(notif._notification)
}

func NotificationClose(notif *Notification) *Error {
	var err *C.GError

	C.notify_notification_close(notif._notification, &err)

	return ErrorFromNative(unsafe.Pointer(err))
}

// Member Functions
func (this *Notification) Update(summary, body, icon string) bool {
	return NotificationUpdate(this, summary, body, icon)
}

func (this *Notification) Show() *Error {
	return NotificationShow(this)
}

func (this *Notification) SetTimeout(timeout int32) {
	NotificationSetTimeout(this, timeout)
}

func (this *Notification) SetCategory(category string) {
	NotificationSetCategory(this, category)
}

func (this *Notification) SetUrgency(urgency Urgency) {
	NotificationSetUrgency(this, urgency)
}

func (this *Notification) SetHintInt32(key string, value int32) {
	NotificationSetHintInt32(this, key, value)
}

func (this *Notification) SetHintDouble(key string, value float64) {
	NotificationSetHintDouble(this, key, value)
}

func (this *Notification) SetHintString(key string, value string) {
	NotificationSetHintString(this, key, value)
}

func (this *Notification) SetHintByte(key string, value byte) {
	NotificationSetHintByte(this, key, value)
}

func (this *Notification) SetHintByteArray(key string, value []byte, len uint32) {
	NotificationSetHintByteArray(this, key, value, len)
}

func (this *Notification) SetHint(key string, value interface{}) {
	NotificationSetHint(this, key, value)
}

func (this *Notification) ClearHints() {
	NotificationClearHints(this)
}

func (this *Notification) AddAction(action, label string, callback ActionCallback, user_data interface{}) {
	NotificationAddAction(this, action, label, callback, user_data)
}

func (this *Notification) ClearActions() {
	NotificationClearActions(this)
}

func (this *Notification) Close() *Error {
	return NotificationClose(this)
}

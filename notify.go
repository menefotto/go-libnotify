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

/*
 * Exported Types
 */

const (
	NOTIFY_URGENCY_LOW      = 0
	NOTIFY_URGENCY_NORMAL   = 1
	NOTIFY_URGENCY_CRITICAL = 2
)

const (
	NOTIFY_EXPIRES_NEVER   = 0
	NOTIFY_EXPIRES_DEFAULT = -1
)

type Urgency int
type ActionCallback func(*Notification, string, interface{})

// NotifyNotification api, n is what you actually need to use to send notifications
// once you hace initialized the library
type Notification struct {
	_notification *C.NotifyNotification
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

	return newNotification(C.notify_notification_new(ptitle, (*C.char)(ntext), pimage))
}

func (n *Notification) Update(summary, body, icon string) bool {
	return notificationUpdate(n, summary, body, icon)
}

func (n *Notification) Show() error {
	return notificationShow(n)
}

func (n *Notification) SetTimeout(timeout int32) {
	notificationSetTimeout(n, timeout)
}

func (n *Notification) SetCategory(category string) {
	notificationSetCategory(n, category)
}

func (n *Notification) SetUrgency(urgency Urgency) {
	notificationSetUrgency(n, urgency)
}

func (n *Notification) SetHintInt32(key string, value int32) {
	notificationSetHintInt32(n, key, value)
}

func (n *Notification) SetHintDouble(key string, value float64) {
	notificationSetHintDouble(n, key, value)
}

func (n *Notification) SetHintString(key string, value string) {
	notificationSetHintString(n, key, value)
}

func (n *Notification) SetHintByte(key string, value byte) {
	notificationSetHintByte(n, key, value)
}

func (n *Notification) SetHintByteArray(key string, value []byte, len uint32) {
	notificationSetHintByteArray(n, key, value, len)
}

func (n *Notification) SetHint(key string, value interface{}) {
	notificationSetHint(n, key, value)
}

func (n *Notification) ClearHints() {
	notificationClearHints(n)
}

func (n *Notification) AddAction(action, label string, callback ActionCallback, user_data interface{}) {
	notificationAddAction(n, action, label, callback, user_data)
}

func (n *Notification) ClearActions() {
	notificationClearActions(n)
}

func (n *Notification) Close() error {
	return notificationClose(n)
}

func newNotification(cnotif *C.NotifyNotification) *Notification {
	return &Notification{cnotif}
}

func notificationUpdate(notif *Notification, summary, body, icon string) bool {
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

func notificationShow(notif *Notification) error {
	var err *C.GError
	C.notify_notification_show(notif._notification, &err)

	return ErrorFromNative(unsafe.Pointer(err))
}

func notificationSetTimeout(notif *Notification, timeout int32) {
	C.notify_notification_set_timeout(notif._notification, C.gint(timeout))
}

func notificationSetCategory(notif *Notification, category string) {
	pcategory := C.CString(category)
	defer C.free(unsafe.Pointer(pcategory))

	C.notify_notification_set_category(notif._notification, pcategory)
}

func notificationSetUrgency(notif *Notification, urgency Urgency) {
	C.notify_notification_set_urgency(notif._notification, C.NotifyUrgency(urgency))
}

func notificationSetHintInt32(notif *Notification, key string, value int32) {
	pkey := C.CString(key)
	defer C.free(unsafe.Pointer(pkey))

	C.notify_notification_set_hint_int32(notif._notification, pkey, C.gint(value))
}

func notificationSetHintDouble(notif *Notification, key string, value float64) {
	pkey := C.CString(key)
	defer C.free(unsafe.Pointer(pkey))

	C.notify_notification_set_hint_double(notif._notification, pkey, C.gdouble(value))
}

func notificationSetHintString(notif *Notification, key string, value string) {
	pkey := C.CString(key)
	pvalue := C.CString(value)
	defer func() {
		C.free(unsafe.Pointer(pkey))
		C.free(unsafe.Pointer(pvalue))
	}()

	C.notify_notification_set_hint_string(notif._notification, pkey, pvalue)
}

func notificationSetHintByte(notif *Notification, key string, value byte) {
	pkey := C.CString(key)
	defer C.free(unsafe.Pointer(pkey))

	C.notify_notification_set_hint_byte(notif._notification, pkey, C.guchar(value))
}

// FIXME: implement
func notificationSetHintByteArray(notif *Notification, key string, value []byte, len uint32) {
	pkey := C.CString(key)
	defer C.free(unsafe.Pointer(pkey))

	// C.notify_notification_set_hint_byte_array(notif._notification, pkey, (*C.guchar)(value), C.gsize(len))
}

func notificationSetHint(notif *Notification, key string, value interface{}) {
	switch value.(type) {
	case int32:
		notificationSetHintInt32(notif, key, value.(int32))
	case float64:
		notificationSetHintDouble(notif, key, value.(float64))
	case string:
		notificationSetHintString(notif, key, value.(string))
	case byte:
		notificationSetHintByte(notif, key, value.(byte))
	}
}

func notificationClearHints(notif *Notification) {
	C.notify_notification_clear_hints(notif._notification)
}

// FIXME: the C function is supposed to be allowing the user to pass another function than free
func notificationAddAction(notif *Notification, action, label string,
	callback ActionCallback, user_data interface{}) {
	// C.notify_notification_add_action(notif._notification, paction, plabel, (C.NotifyActionCallback)(callback), user_data, C.free)
}

func notificationClearActions(notif *Notification) {
	C.notify_notification_clear_actions(notif._notification)
}

func notificationClose(notif *Notification) error {
	var err *C.GError

	C.notify_notification_close(notif._notification, &err)

	return ErrorFromNative(unsafe.Pointer(err))
}

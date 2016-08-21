package notify

import "testing"

func TestNotify(t *testing.T) {
	err := Init("AppName")
	defer UnInit()
	if !err {
		t.Errorf("could not initilized\n")
	}

	hello := NotificationNew("Hello World!",
		"This is an example notification.",
		"")

	if hello == nil {
		panic("Cannot create new notification")
	}

	hello.SetUrgency(NOTIFY_URGENCY_LOW)
	hello.SetTimeout(NOTIFY_EXPIRES_DEFAULT)
	hello.Show()
	//hello.Close()
	//	notify.NotificationClose(hello)

}

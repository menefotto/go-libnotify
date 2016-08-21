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

	//hello.SetTimeout(0)
	//notify.NotificationSetTimeout(hello, DELAY)

	hello.Show()
	//if e := notify.NotificationShow(hello); e != nil {
	//		fmt.Fprintf(os.Stderr, "%s\n", e.Message())
	//		return
	//	}

	//time.Sleep(time.Second * 3)
	//hello.Close()
	//	notify.NotificationClose(hello)

}

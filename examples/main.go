/*
 * main.go for go-notify
 * by lenorm_f
 */

package main

import notify "github.com/lenormf/go-notify"

const (
	DELAY = 3000
)

func main() {
	notify.Init("AppName")
	defer notify.UnInit()

	hello := notify.NotificationNew("Hello World!",
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

/*
 * main.go for go-notify
 * by lenorm_f
 */

package main

import notify "github.com/go-libnotify"

func main() {
	notify.Init("AppName")
	defer notify.UnInit()

	hello := notify.NotificationNew("Hello World!",
		"This is an example notification.",
		"")

	if hello == nil {
		panic("Cannot create new notification")
	}

	hello.SetTimeout(notify.NOTIFY_EXPIRES_DEFAULT)

	hello.Show()
	//hello.Close()

}

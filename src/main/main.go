package main

import (
	"JsLib/JsDispatcher"
	"JsLib/JsExit"
	"JsLib/JsNet"
	"activity"
	"bill"
	"common"
	"util"
)

func exit() int {
	JsDispatcher.Close()
	return 0
}

func main() {
	JsExit.RegisterExitCb(exit)
	JsNet.AppConf("./conf/app.conf")

	HomeInit()
	common.InitCommon()
	bill.BillInit()
	util.Mobile_init()
	activity.ActInit()
	
	JsDispatcher.Run()
}

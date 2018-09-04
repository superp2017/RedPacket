package util

import (
	"JsLib/JsDispatcher"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
)

type Return struct {
	Ret string
	Msg string
}

func Mobile_init() {
	JsDispatcher.Http("/mobileverify", mobile_veirfy)
}

func mobile_veirfy(session *JsNet.StSession) {
	type MobilePara struct {
		Mobile string //手机号码
		Expire int    //倒计时
	}

	ret := &Return{}
	para := &MobilePara{}
	e := session.GetPara(para)
	if e != nil {
		Error(e.Error())
		ret.Ret = "1"
		ret.Msg = e.Error()
		session.Forward(ret)
		return
	}

	IDAuth(para.Mobile, "君赛科技", para.Expire)

	ret.Ret = "0"
	ret.Msg = "success"
	session.Forward(ret)
}

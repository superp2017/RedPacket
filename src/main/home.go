package main

import (
	"JsLib/JsConfig"
	"JsLib/JsNet"
	"common"
	"db"
	. "util"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

func HomeInit() {
	JsNet.Http("/home", home)   //主页
	JsNet.Http("/login", Login) //后台登录
}

//开始地方
func home(session *JsNet.StSession) {
	openid := session.Get("openid")
	type home_ret struct {
		Ret string
		Msg string
	}

	ret := &home_ret{}

	if openid == "" {
		if openid == "" {
			JsNet.CheckWxAuth(session, authCb)
			return
		}
	}

	if openid != "" {

		teamwork := session.Get("teamwork")
		actlist := session.Get("actlist")
		activity := session.Get("activity")

		if teamwork != "" {
			openid = "?teamwork=true&openid=" + openid
		} else if actlist != "" {
			openid = "?actlist=true&openid=" + openid
		} else if activity != "" {
			openid = "?activity=" + activity + "&openid=" + openid
		} else {
			openid = "?openid=" + openid
		}

		JsNet.HttpRedict(session, JsConfig.CFG.WxJsApi.WeChatRedirectHome+openid)

	} else {
		ret.Ret = "1"
		ret.Msg = "openid == nil"
		session.Forward(ret)
	}
}

func authCb(u *oauth2.UserInfo, session *JsNet.StSession) {
	common.NewUser(u, 0)
}

//登陆
func Login(session *JsNet.StSession) {
	type RD_Login struct {
		UserName string //用户名
		PassWord string //用户密码
	}

	st := RD_Login{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	ip := session.RemoteAddr()

	if ip == "" {
		ForwardEx(session, "1", nil, "Login RemoteAddr is empty\n")
		return
	}

	var pwd string
	if err := db.Get(st.UserName, &pwd); err != nil {
		go common.ClearToken(ip)
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.PassWord != pwd {
		go common.ClearToken(ip)
		ForwardEx(session, "1", nil, "pwd is error! \n")
		return
	}

	Forward(session, "0", common.MapToken(ip))
}

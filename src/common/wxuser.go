package common

import (
	"JsLib/JsNet"
	"encoding/json"
	"io/ioutil"
	"net/http"
	. "util"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

type SessionKeys struct {
	OpenID     string `json:"openid"`      // 用户的唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // unionID
}

///小程序Code获取openid和unionid
func CodeGetSessionKey(session *JsNet.StSession) {
	type INFo struct {
		Code string
	}
	st := &INFo{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.Code == "" {
		ForwardEx(session, "1", nil, "CodeGetSessionKey Code isEmpty\n")
		return
	}

	url := "https://api.weixin.qq.com/sns/jscode2session?appid="
	url += "wx11cdac22d7719783"
	url += "&secret="
	url += "4c67bbf59e131fadc309c4c9a7d16059"
	url += "&js_code="
	url += st.Code
	url += "&grant_type=authorization_code"

	response, e := http.Get(url)
	defer response.Body.Close()
	if e != nil {
		ForwardEx(session, "1", nil, e.Error())
		return
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := &SessionKeys{}
	if err := json.Unmarshal(b, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	// u, err := checkUser(data.OpenID, data.UnionID)
	// if err == nil && u != nil {
	// 	Forward(session, "2", data)
	// 	return
	// }

	Forward(session, "0", data)
}

///微信授权直接转用户
func WxAccessToken2User(session *JsNet.StSession) {
	st := &WXAccessToken{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.Access_Token == "" || st.Openid == "" {
		ForwardEx(session, "1", nil, "WxAccessToken2User param failed,Access_Token=%s,Openid=%s\n",
			st.Access_Token, st.Openid)
		return
	}

	u, err := checkUser(st.Openid, st.UnionID)
	if err == nil && u != nil {
		Forward(session, "0", u)
		return
	}

	url := "https://api.weixin.qq.com/sns/userinfo?access_token="
	url += st.Access_Token
	url += "&openid="
	url += st.Openid
	url += "&lang=zh_CN"

	response, e := http.Get(url)
	defer response.Body.Close()
	if e != nil {
		ForwardEx(session, "1", nil, e.Error())
		return
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := &oauth2.UserInfo{}
	if err := json.Unmarshal(b, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	user, err := NewUser(data, 1)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", user)
}

//微信用户信息直接转系统给用户
func WX2User(session *JsNet.StSession) {
	st := &oauth2.UserInfo{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	user, err := NewUser(st, 2)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", user)
}

func checkUser(OpenID, UnionID string) (*User, error) {

	//检查unionID 是否映射UID
	UID, err := getUIDfromUnionID(UnionID)
	if err == nil && UID != "" {
		user, err1 := GetUserInfo(UID)
		if err1 == nil && user != nil {
			go OpenidMapUID(user)
			return user, nil
		}
	} else {
		uid, err1 := getUIDfromOpenID(OpenID)
		if err1 == nil && uid != "" {
			u, err2 := GetUserInfo(uid)
			if err2 == nil && u != nil {
				go UnionidMapUID(u)
				return u, nil
			}
		}
	}
	return nil, Err("not exist")
}

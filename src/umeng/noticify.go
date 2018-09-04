package umeng

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"io/ioutil"
	"net/http"
	. "util"
)

func PostUmeng() (e error) {
	// The host
	host := "http://msg.umeng.com"
	// The upload path
	// uploadPath:= "/upload"
	// The post path
	postPath := "/api/send"
	//appmaster secret
	appMasterSecret := "hqfvknbs5f9kiqt40kj7ijev5xugwlcm"
	UMENG_APPKEY := "5a20054fb27b0a344e000175"
	// UMENG_MESSAGE_SECRET:="9ae5de15d4d5049d004a4108226c00b7"
	url := host + postPath
	uInfo := &UMengAndroidInfo{}
	uInfo.Appkey = UMENG_APPKEY
	uInfo.Timestamp = CurTime()
	uInfo.Typeu = "unicast"
	payload := PayloadInfo{}
	payload.Display_type = "message"
	uInfo.Payload = payload
	uInfo.Description = "传单侠欢迎您"
	// bodyumeng:=BodyInfo{}
	// bodyumeng.Custom=

	h := md5.New()
	postBody, err := json.Marshal(uInfo)
	if err != nil {
		ErrorLog(err.Error())
		return err
	}
	signstr := "POST" + url + string(postBody) + appMasterSecret
	signbyte := []byte(signstr)
	h.Write(signbyte)
	signbkbyte := h.Sum(nil)
	urlpost := url + "?sign=" + string(signbkbyte)
	body := bytes.NewBuffer([]byte(postBody))
	res, err := http.Post(urlpost, "application/json;charset=utf-8", body)
	if err != nil {
		ErrorLog(err.Error())
		return err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		ErrorLog(err.Error())
		return err
	}
	Info("%v", result)
	return nil
}

func PostUmengNet(session *JsNet.StSession) {
	err := PostUmeng()
	if err != nil {
		ForwardEx(session, "Fail", nil, err.Error())
		return
	}
	ForwardEx(session, "Hello", nil, "Done")

}

package common
import (
	"errors"
	"net/http"
	"crypto/md5"   
	"io/ioutil"  
	"bytes"  
	"encoding/json"  
	. "JsLib/JsLogger"
	."util"
	"JsLib/JsNet"
	"encoding/hex"
	// "strconv"
	// "time"
)

func postAndroidUmeng(umengandroidInfo *UMengAndroidInfo)(str string,e error) {  
	// The host
	host:= "http://msg.umeng.com"
	postPath:= "/api/send"
	//appmaster secret
	appMasterSecret:="hqfvknbs5f9kiqt40kj7ijev5xugwlcm"
	url:=host+postPath

	h := md5.New()
	postBody, err := json.Marshal(umengandroidInfo)  
	if err != nil {  
		ErrorLog(err.Error())  
		return "", err
	} 
	signstr:="POST"+url+string(postBody)+appMasterSecret
	signbyte:=[]byte(signstr)
	h.Write(signbyte)
	signbkbyte:= h.Sum(nil) 
	urlpost:=url + "?sign=" + hex.EncodeToString(signbkbyte)
	body := bytes.NewBuffer([]byte(postBody))  
	res,err := http.Post(urlpost, "application/json;charset=utf-8", body)  
  if err != nil {  
				ErrorLog(err.Error())
        return "", err
  }  
  result, err := ioutil.ReadAll(res.Body)  
     res.Body.Close()  
     if err != nil { 
			ErrorLog(err.Error())
			return "", err 
	} 
	Info("%s", string(result)) 
	return string(result),nil
}

func PostUmengNet(session *JsNet.StSession) {
	uInfo:=NewUMengAndroidInfo()
	resumeng,err:=postAndroidUmeng(uInfo)
		if err!=nil{
			ForwardEx(session, resumeng, nil, err.Error())
			return
		}
		ForwardEx(session, "Hello", nil, resumeng)
	
}



func UOne(session *JsNet.StSession) {
	err:=PostMsgOnePerson("发送一个人的ticker","发送一个人的title","发送一个人的Message","AqD_Dv3LFxl0YyeFMOhPbM9yTpO7KfHAnP1LeeicTKno")
		if err!=nil{
			ForwardEx(session, "", nil, err.Error())
			return
		}
		ForwardEx(session, "Sucess", nil, "")
	
}

func UMulti(session *JsNet.StSession) {
	lsdeviceID:=[]string{}
	lsdeviceID=append(lsdeviceID,"AqD_Dv3LFxl0YyeFMOhPbM9yTpO7KfHAnP1LeeicTKno")
	err:=PostMessageMultiPerson("发送多个人的ticker","发送多个人的title","发送多个人的Message",lsdeviceID)
		if err!=nil{
			ForwardEx(session, "", nil, err.Error())
			return
		}
		ForwardEx(session, "Sucess", nil, "")
	
}

func UAll(session *JsNet.StSession) {
	lsdeviceID:=[]string{}
	lsdeviceID=append(lsdeviceID,"AqD_Dv3LFxl0YyeFMOhPbM9yTpO7KfHAnP1LeeicTKno")
	err:=PostMsgAllPerson("发送所有人的ticker","发送所有人的title","发送所有人的Message")
		if err!=nil{
			ForwardEx(session, "", nil, err.Error())
			return
		}
		ForwardEx(session, "Sucess", nil, "")
	
}



func PostMsgOnePerson(ticker,title,message,deviceID string)(e error){

	uInfo:=NewUMengAndroidInfo()
	device_id:=deviceID
	uInfo.Device_tokens=device_id
	uInfo.Payload.Body.Ticker=ticker
	uInfo.Payload.Body.Title=title
	uInfo.Payload.Body.Text=message
	uInfo.Description=message
	uInfo.Typeu="unicast"
	_,err:=postAndroidUmeng(uInfo)
	
		if err!=nil{
		
			return
		}

	return nil

}

func PostMessageMultiPerson(ticker,title,message string,lsdeviceID []string)(e error){

	uInfo:=NewUMengAndroidInfo()
	
	if len(lsdeviceID)==0{
		return errors.New("The ls uid length is 0")
	}
	device_id:=""

	for i:=0;i<len(lsdeviceID)-1;i++{
		device_id=device_id+lsdeviceID[i]+";"
	}
	device_id=device_id+lsdeviceID[len(lsdeviceID)-1]
	uInfo.Device_tokens=device_id
	uInfo.Payload.Body.Ticker=ticker
	uInfo.Typeu="listcast"
	uInfo.Payload.Body.Title=title
	uInfo.Payload.Body.Text=message
	uInfo.Description=message
	_,err:=postAndroidUmeng(uInfo)
		if err!=nil{
		
			return
		}


	return nil
}

func PostMsgAllPerson(ticker,title,message string)(e error){

	uInfo:=NewUMengAndroidInfo()
	uInfo.Payload.Body.Ticker=ticker
	uInfo.Payload.Body.Title=title
	uInfo.Payload.Body.Text=message
	uInfo.Description=message
	uInfo.Typeu="broadcast"
	_,err:=postAndroidUmeng(uInfo)
	
		if err!=nil{
		
			return
		}

	return nil
}








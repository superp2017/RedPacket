package common
import (
	"time"
	"strconv"

)

type PolicyInfo struct{
	Start_time     string   `json:"start_time"`
	Expire_time     string   `json:"expire_time"`
	// Max_send_num     int   `json:"max_send_num"`
	Out_biz_no     string   `json:"out_biz_no"`
}

type PolicyIOSInfo struct{
	Start_time     string   `json:"start_time"`
	Expire_time     string   `json:"expire_time"`
	// Max_send_num     int   `json:"max_send_num"`
	Apns_collapse_id     string   `json:"apns-collapse-id"`
}

type PayloadInfo struct{
	Display_type     string   `json:"display_type"`
	Body     BodyInfo   `json:"body"`
	// Extra    string    `json:"extra"`

}

type PayloadIOSInfo struct{
	Aps     AppsInfo   `json:"aps"`
	// Extra    string    `json:"extra"`
}

type AppsInfo struct{
	Alert    string    `json:"alert"`
	Badge    string    `json:"badge"`
	Sound    string    `json:"sound"`
	Contentavailable    string    `json:"content-available"`
	Category    string    `json:"category"`
}

type BodyInfo struct{
	Ticker     string   `json:"ticker"`
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	Icon     string   `json:"icon"`
	LargeIcon     string   `json:"largeIcon"`

	Img     string   `json:"img"`
	Sound     string   `json:"sound"`
	Builder_id     int   `json:"builder_id"`
	Play_vibrate     string   `json:"play_vibrate"`
	Play_lights     string   `json:"play_lights"`

	Play_sound     string   `json:"play_sound"`
	After_open     string   `json:"after_open"`
	Url     string   `json:"url"`
	Activity     string   `json:"activity"`
	Custom     string   `json:"custom"`
}

type UMengAndroidInfo struct{
	Appkey     string   `json:"appkey"`
	Timestamp     string   `json:"timestamp"`
	Typeu     string   `json:"type"`
	Device_tokens     string   `json:"device_tokens"`
	Alias_type     string   `json:"alias_type"`
	File_id     string   `json:"file_id"`
	// Filter     FilterInfo   `json:"filter"`
	Payload     PayloadInfo   `json:"payload"`
	Policy     PolicyInfo   `json:"policy"`
	Production_mode     string   `json:"production_mode"`
	Description     string   `json:"description"`
}
type FilterInfo struct{

}

type UMengIOSInfo struct{
	Appkey     string   `json:"appkey"`
	Timestamp     string   `json:"timestamp"`
	Typeu     string   `json:"type"`
	Device_tokens     string   `json:"device_tokens"`
	Alias_type     string   `json:"alias_type"`
	File_id     string   `json:"file_id"`
	// Filter     string   `json:"filter"`
	Payload     PayloadIOSInfo   `json:"payload"`
	Policy     PolicyIOSInfo   `json:"policy"`
	Production_mode     bool   `json:"production_mode"`
	Description     string   `json:"description"`
}

func NewUMengAndroidInfo()(*UMengAndroidInfo){
	UMENG_APPKEY:="5a20054fb27b0a344e000175"
	uInfo:=&UMengAndroidInfo{}
	uInfo.Appkey=UMENG_APPKEY
	payload:=PayloadInfo{}
	body:=BodyInfo{}
	//body
	body.Ticker="汇传提醒"
	body.Title="汇传友情提醒"
	body.Text="欢迎关注汇传，带给你不一样的惊喜"
	body.Builder_id=0
	body.Play_vibrate="true"
	body.Play_lights="true"
	body.Play_sound="true"
	body.After_open="go_custom"
	body.Custom="欢迎来到汇传"
	payload.Body=body
	payload.Display_type="notification"
	uInfo.Payload=payload
	t := time.Now().Unix()  
	timecurrent:= strconv.FormatInt(t, 10) 
	uInfo.Timestamp=timecurrent
	// uInfo.Typeu="broadcast"
	uInfo.Typeu="unicast"
	uInfo.Device_tokens="AqD_Dv3LFxl0YyeFMOhPbM9yTpO7KfHAnP1LeeicTKno"
	uInfo.Production_mode="false"
	uInfo.Description="传单侠欢迎您"
	return uInfo
} 




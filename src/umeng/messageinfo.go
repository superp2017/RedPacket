package umeng

// . "JsLib/JsConfig"
// . "JsLib/JsLogger"
// "JsLib/JsNet"
// "constant"
// "db"
// "ider"
// . "util"

type PolicyInfo struct {
	Start_time   string `json:"start_time"`
	Expire_time  string `json:"expire_time"`
	Max_send_num int    `json:"max_send_num"`
	Out_biz_no   string `json:"out_biz_no"`
}

type PolicyIOSInfo struct {
	Start_time       string `json:"start_time"`
	Expire_time      string `json:"expire_time"`
	Max_send_num     int    `json:"max_send_num"`
	Apns_collapse_id string `json:"apns-collapse-id"`
}

type PayloadInfo struct {
	Display_type string   `json:"display_type"`
	Body         BodyInfo `json:"body"`
	Extra        string   `json:"extra"`
}

type PayloadIOSInfo struct {
	Aps   AppsInfo `json:"aps"`
	Extra string   `json:"extra"`
}

type AppsInfo struct {
	Alert            string `json:"alert"`
	Badge            string `json:"badge"`
	Sound            string `json:"sound"`
	Contentavailable string `json:"content-available"`
	Category         string `json:"category"`
}

type BodyInfo struct {
	Ticker    string `json:"ticker"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	LargeIcon string `json:"largeIcon"`

	Img          string `json:"img"`
	Sound        string `json:"sound"`
	Builder_id   int    `json:"builder_id"`
	Play_vibrate string `json:"play_vibrate"`
	Play_lights  string `json:"play_lights"`

	Play_sound string `json:"play_sound"`
	After_open string `json:"after_open"`
	Url        string `json:"url"`
	Activity   string `json:"activity"`
	Custom     string `json:"custom"`
}

type UMengAndroidInfo struct {
	Appkey        string `json:"appkey"`
	Timestamp     string `json:"timestamp"`
	Typeu         string `json:"type"`
	Device_tokens string `json:"device_tokens"`
	Alias_type    string `json:"alias_type"`
	File_id       string `json:"file_id"`
	// Filter     FilterInfo   `json:"filter"`
	Payload         PayloadInfo `json:"payload"`
	Policy          PolicyInfo  `json:"policy"`
	Production_mode string      `json:"production_mode"`
	Description     string      `json:"description"`
}
type FilterInfo struct {
}

type UMengIOSInfo struct {
	Appkey        string `json:"appkey"`
	Timestamp     string `json:"timestamp"`
	Typeu         string `json:"type"`
	Device_tokens string `json:"device_tokens"`
	Alias_type    string `json:"alias_type"`
	File_id       string `json:"file_id"`
	// Filter     string   `json:"filter"`
	Payload         PayloadIOSInfo `json:"payload"`
	Policy          PolicyIOSInfo  `json:"policy"`
	Production_mode bool           `json:"production_mode"`
	Description     string         `json:"description"`
}

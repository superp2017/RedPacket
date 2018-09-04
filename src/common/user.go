package common

import (
	. "JsLib/JsConfig"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	. "util"

	"github.com/chanxuehong/wechat/mp/user/oauth2"
)

type WXAccessToken struct {
	Access_Token  string `json:"access_token"`
	Expires_in    int    `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
	Openid        string `json:"openid"`
	Scope         string `json:"scope"`
	UnionID       string `json:"unionid"`
}

type MoneyChangeRecord struct {
	Money  int    //金钱
	Date   string //日期
	Blance int    //余额
	Type   int    //类型(0：提现、1：抢红包、2：分享佣金:3：发红包)
}

type ManagerID struct {
	CID string   //公司id
	DID []string //网点id列表
}

type User struct {
	Name               string             //姓名
	Nickname           string             //昵称
	UID                string             //UID
	OpenId             string             //微信openid
	OpenId_App         string             //App微信openid
	OpenId_Small       string             //小程序openid
	UnionId            string             //微信统一id
	HeadImageURL       string             //头像路径
	Age                int                //年龄
	Sex                string             //性别
	Mobile             string             //手机号
	City               string             //城市
	Province           string             //省份
	Country            string             //国家
	Addr               string             //地址
	Type               string             //商家类型
	IsManager          []ManagerID        //是否是网点管理者,以此来判断是否可以消券
	Role               string             //用户的角色 0：用户 1：商家 2：Agent
	Blance             int                //余额
	Recharge           int                //总充值金额
	RecevieMoney       int                //抢到的金额
	SharedMoney        int                //分享得到的佣金
	LsCompany          []string           //下属的公司
	LsActivityLanched  []string           //代公司发布过的活动
	LsSystemQuestionId []string           //回答过的系统问题
	LsAgentFollower    []AgentFollower    //代理的下线
	ParentAgentID      string             //代理ID
	LsCompanyShortInfo []CompanyShortInfo //公司列表
	UserName           string             //登陆用户名
	PassWord           string             //登陆密码
	AgentTime          string
}

type CompanyShortInfo struct {
	CompanyID      string
	CompanyName    string
	EntityTime     string
	ActivityNumber int
}

type AgentFollower struct {
	UID         string
	Name        string
	TotalIncome int
	Address     string
	EntityTime  string
}

func NewUser(wc *oauth2.UserInfo, Type int) (*User, error) {
	if wc.OpenId == "" && wc.UnionId == "" {
		return nil, ErrorLog("oauth2.UserInfo openID&UnionId  is empty\n")
	}
	u, err := checkUser(wc.OpenId, wc.UnionId)
	// if err == nil && u != nil {
	// 	return u, nil
	// }

	user := User{
		UnionId:      wc.UnionId,
		Nickname:     wc.Nickname,
		Name:         wc.Nickname,
		City:         wc.City,
		Province:     wc.Province,
		Country:      wc.Country,
		HeadImageURL: wc.HeadImageURL,
		Age:          0,
		Blance:       0,
		Recharge:     0,
		Role:         "0",
	}
	if err == nil && u != nil {
		user.UID = u.UID
		user.Blance = u.Blance
		user.RecevieMoney = u.RecevieMoney
		user.Role = u.Role
		user.Name = u.Name
	} else {
		user.UID = ider.GenID()
	}

	switch Type {
	case 0:
		user.OpenId = wc.OpenId
	case 1:
		user.OpenId_App = wc.OpenId
	case 2:
		user.OpenId_Small = wc.OpenId
	default:
		user.OpenId = wc.OpenId
	}

	switch wc.Sex {
	case 1:
		user.Sex = "男"
	case 2:
		user.Sex = "女"
	default:
		user.Sex = "未知"
	}

	//write db
	if err := db.DirectWrite(constant.Hash_User, user.UID, &user); err != nil {
		return nil, ErrorLog("DirectWrite User failed\n")
	}

	go OpenidMapUID(&user)
	go UnionidMapUID(&user)
	return &user, nil
}

/**
 * 查询用户信息
 */
func QueryUser(session *JsNet.StSession) {
	type info struct {
		UID string
	}
	st := info{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, "QueryUser  GetPara :%s\n", err.Error())
		return
	}

	user := &User{}
	var err error = nil
	user, err = GetUserInfo(st.UID)
	if err != nil {
		ForwardEx(session, "1", nil, "ShareLock User:%s failed\n", st.UID)
		return
	}
	Forward(session, "0", user)
}

/**
 * 查询用户信息form openID
 */
func QueryUserFromOpenID(session *JsNet.StSession) {
	type info struct {
		Openid string
	}
	st := info{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, "QueryUser  GetPara :%s\n", err.Error())
		return
	}
	user := &User{}
	var err error = nil
	user, err = GetUserFromOpenID(st.Openid)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", user)
}

/**
 * 查询用户信息form UnionID
 */
func QueryUserFromUnionID(session *JsNet.StSession) {
	type info struct {
		UnionId string
	}
	st := info{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, "QueryUserFromUnionID  GetPara :%s\n", err.Error())
		return
	}

	user := &User{}
	var err error = nil
	user, err = GetUserFromUnionID(st.UnionId)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", user)
}

/**
 * 修改用户信息
 */
func ModifyUser(session *JsNet.StSession) {

	type modifyInfo struct {
		UID          string //统一id
		Name         string
		HeadImageURL string
		Age          int
		Sex          string
		Mobile       string //手机号
		City         string
		Province     string
		Country      string
	}
	st := &modifyInfo{}

	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ModifyUser GetPara :%s\n", err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "ModifyUser GetPara failed, UID is empty\n")
		return
	}

	user := &User{}
	if err := db.Modify(constant.Hash_User, st.UID, user, st); err != nil {
		ForwardEx(session, "1", nil, "Modify User:%s  failed\n", st.UID)
		return
	}

	Forward(session, "0", user)
}

///用户绑定手机号
func UserBindMobile(session *JsNet.StSession) {
	type INFO struct {
		UID    string
		Mobile string
		Code   string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if !VerifySmsCode(st.Mobile, st.Code) {
		ForwardEx(session, "1", nil, "短信校验失败\n")
		return
	}
	user := &User{}
	if err := db.WriteLock(constant.Hash_User, st.UID, user); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	user.Mobile = st.Mobile
	if err := db.WriteBack(constant.Hash_User, st.UID, user); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", user)

}

/////查询多个用户信息
func QueryMoreUserInfo(session *JsNet.StSession) {
	type INFO struct {
		UIDs []string //多个用户id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := []*User{}
	for _, v := range st.UIDs {
		if user, err := GetUserInfo(v); err == nil {
			data = append(data, user)
		}
	}
	Forward(session, "0", data)
}

///通过需求信息，更新用户信息
func UpdateUserInfo(UID, Name, Mobile, Addr string) error {
	user := &User{}
	if err := db.WriteLock(constant.Hash_User, UID, user); err != nil {

	}
	user.Name = Name
	user.Mobile = Mobile
	user.Addr = Addr
	return db.WriteBack(constant.Hash_User, UID, user)
}

/**
 *获取用户信息
 */
func GetUserInfo(UID string) (st *User, e error) {

	data := &User{}
	if UID == "" {
		return data, ErrorLog("GetUserInfo  UID:s failed!\n", UID)
	}
	err := db.ShareLock(constant.Hash_User, UID, data)
	return data, err
}

func GetUserFromUnionID(unionID string) (st *User, e error) {
	data := &User{}
	if unionID == "" {
		return data, ErrorLog("GetUserFromUnionID  unionID is empty \n")
	}
	UID, err := getUIDfromUnionID(unionID)
	if err != nil || UID == "" {
		return data, err
	}

	var ee error = nil
	data, ee = GetUserInfo(UID)
	return data, ee
}

func GetUserFromOpenID(openID string) (st *User, e error) {
	data := &User{}
	if openID == "" {
		return data, ErrorLog("GetUserFromOpenID openID is empty \n")
	}

	UID, err := getUIDfromOpenID(openID)
	if err != nil || UID == "" {
		return data, err
	}
	var ee error = nil
	data, ee = GetUserInfo(UID)
	return data, ee
}

/*
	openid映射UID
*/
func OpenidMapUID(user *User) {
	if user.OpenId != "" && user.UID != "" {
		go db.DirectWrite(constant.Hash_OpenID_UID, user.OpenId, user.UID)
	}
	if user.OpenId_App != "" && user.UID != "" {
		go db.DirectWrite(constant.Hash_OpenID_UID, user.OpenId_App, user.UID)
	}
	if user.OpenId_Small != "" && user.UID != "" {
		go db.DirectWrite(constant.Hash_OpenID_UID, user.OpenId_Small, user.UID)
	}
}

/*
	union映射UID
*/
func UnionidMapUID(user *User) {
	if user.UnionId != "" && user.UID != "" {
		db.DirectWrite(constant.Hash_UnionID_UID, user.UnionId, user.UID)
	}
}

/*
	openid获取UID
*/

func getUIDfromOpenID(openid string) (string, error) {
	if openid == "" {
		return "", ErrorLog("getUIDfromOpenID :openid is empty \n")
	}

	var UID string

	if err := db.ShareLock(constant.Hash_OpenID_UID, openid, &UID); err != nil {
		return UID, err
	}
	return UID, nil
}

/*
	union获取UID
*/
func getUIDfromUnionID(unionID string) (string, error) {
	if unionID == "" {
		return "", ErrorLog("getUIDfromUnionID :unionID is empty \n")
	}
	var UID string
	if err := db.ShareLock(constant.Hash_UnionID_UID, unionID, &UID); err != nil {
		return UID, err
	}
	return UID, nil
}

///////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////

////抢、发红包之后更新用户的账户信息
func updataUserAccount(UID string, receive, send, shared, Type int) error {
	data := &User{}
	if err := db.WriteLock(constant.Hash_User, UID, data); err != nil {
		return err
	}
	if receive > 0 {
		data.RecevieMoney += receive
		data.Blance += receive
		if data.RecevieMoney < 0 {
			data.RecevieMoney = 0
		}
	}
	if send > 0 {
		data.Recharge += send
		if data.Recharge < 0 {
			data.Recharge = 0
		}
	}
	if shared > 0 {
		data.SharedMoney += shared
		data.Blance += receive
		if data.SharedMoney < 0 {
			data.SharedMoney = 0
		}
	}

	if data.Blance < 0 {
		data.Blance = 0
	}
	if err := db.WriteBack(constant.Hash_User, UID, data); err != nil {
		return err
	}

	m := receive + shared
	if m > 0 {
		return addUserRecord(UID, m, data.Blance, Type)
	}
	return nil
}

///用户提现
func UserWithdraw(session *JsNet.StSession) {
	type INFO struct {
		UID   string //用户id
		Money int    //提现金额
		Tag   string //平台标签
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.Tag == "" || st.Money < 100 {
		ForwardEx(session, "1", nil, "UserWithdraw failed,UID =%s,Tag =%s,Money=%d\n", st.UID, st.Tag, st.Money)
		return
	}

	data := &User{}
	if err := db.WriteLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer db.WriteBack(constant.Hash_User, st.UID, data)
	if data.Blance <= 0 || data.Blance < st.Money {
		ForwardEx(session, "1", nil, "UserWithdraw failed,Blance less than %d !\n", st.Money)
		return
	}
	data.Blance -= st.Money
	if err := addUserRecord(data.UID, -1*st.Money, data.Blance, 0); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if ok := userwithdraw(data, st.Money, st.Tag); ok {
		Forward(session, "0", data)
		return
	}
	Forward(session, "1", data)
}

///添加用户资金变动记录
func GetUserMoneyRecord(session *JsNet.StSession) {
	type INFO struct {
		UID string //用户id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" {
		ForwardEx(session, "1", nil, "GetUserMoneyRecord failed,UID is empty\n")
		return
	}
	data, err := getUserRecord(st.UID)
	if err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	Forward(session, "0", data)

}

///添加用户资金变动记录
func addUserRecord(UID string, Money, Blance, Type int) error {
	record := []MoneyChangeRecord{}
	err := db.WriteLock(constant.Hash_User_MoneyRecord, UID, &record)
	record = append(record, MoneyChangeRecord{
		Money:  Money,
		Date:   CurTime(),
		Blance: Blance,
		Type:   Type,
	})
	if err != nil {
		return db.DirectWrite(constant.Hash_User_MoneyRecord, UID, &record)
	}
	return db.WriteBack(constant.Hash_User_MoneyRecord, UID, &record)
}

//获取用户的资金变动记录
func getUserRecord(UID string) ([]MoneyChangeRecord, error) {
	record := []MoneyChangeRecord{}
	err := db.ShareLock(constant.Hash_User_MoneyRecord, UID, &record)
	return record, err
}

///用户提现
func userwithdraw(user *User, Money int, Tag string) bool {
	if Tag == "" {
		return false
	}
	OpendID := ""
	appid := ""
	if Tag == "small" {
		OpendID = user.OpenId_Small
		appid = "wx11cdac22d7719783"
	} else if Tag == "wechat" {
		OpendID = user.OpenId
		appid = CFG.DirectPay.AppId
	} else if Tag == "app" {
		OpendID = user.OpenId_App
		appid = "wx995d7fad8a74a7ef"
	}
	if OpendID == "" {
		ErrorLog("userwithdraw failed,Openid is Empty,Tag=%s", Tag)
		return false
	}
	st := &ST_Transfer{
		OpenId:     OpendID,
		UserName:   user.Name,
		UserHeader: user.HeadImageURL,
		TimeStamp:  CurStamp(),
		LDate:      CurTime(),
		Desc:       "用户余额提现",
		Amount:     Money,
		AppId:      appid,
	}
	_, err := direct_transfer(st)
	return err == nil
}

///添加用户资金变动记录
func RegisiterUser(session *JsNet.StSession) {
	type INFO struct {
		UID      string //用户id
		DeviceID string //手机识别码
		UIDOK    bool   //是否有UID
	}

	type UIDDeviceIDInfo struct {
		UID        string
		DeviceID   string
		EntityTime string
	}
	uidDeviceIDInfo := &UIDDeviceIDInfo{}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" && st.UIDOK == true {
		ForwardEx(session, "1", nil, "UID is empty\n")
		return
	}

	if st.DeviceID == "" {
		ForwardEx(session, "1", nil, "DeviceID is empty\n")
		return
	}

	if st.UIDOK == true {
		err := db.WriteLock(constant.Hash_UIDDeviceID, st.UID, uidDeviceIDInfo)
		if err != nil {
			uidDeviceIDInfo.UID = st.UID
			uidDeviceIDInfo.DeviceID = st.DeviceID
			uidDeviceIDInfo.EntityTime = CurTime()
			db.DirectWrite(constant.Hash_UIDDeviceID, st.UID, uidDeviceIDInfo)
		} else {
			uidDeviceIDInfo.DeviceID = st.DeviceID
			db.WriteBack(constant.Hash_UIDDeviceID, st.UID, uidDeviceIDInfo)
		}

	}

	err := db.WriteLock(constant.Hash_DeviceIDUID, st.DeviceID, uidDeviceIDInfo)
	if err != nil {
		if st.UIDOK == false {
			uidDeviceIDInfo.UID = "NotLogin"
		} else {
			uidDeviceIDInfo.UID = st.UID
		}

		uidDeviceIDInfo.DeviceID = st.DeviceID
		uidDeviceIDInfo.EntityTime = CurTime()

		db.DirectWrite(constant.Hash_DeviceIDUID, st.DeviceID, uidDeviceIDInfo)
	} else {
		uidDeviceIDInfo.UID = st.UID
		if st.UIDOK == false {
			uidDeviceIDInfo.UID = "NotLogin"
		}
		db.WriteBack(constant.Hash_DeviceIDUID, st.DeviceID, uidDeviceIDInfo)
	}
	Forward(session, "0", uidDeviceIDInfo)

}

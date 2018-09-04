package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	. "util"
)

type Requirement struct {
	RequirementID  string   //需求id
	ActivityID     string   //关联的活动id
	Theme          string   //主题
	Type           string   //类型(备用)
	TemplateId     string   //模板id(备用)
	SellerId       string   //商家id(备用)
	AgentId        string   //小B ID(备用)
	ProxyId        string   //代理id(备用)
	UID            string   //用户id
	OrderID        string   //订单ID
	Content        string   //内容
	Logo           string   //商家logo
	Material       []string //原素材
	Pictures       []string //设计后的海报
	Name           string   //姓名
	Mobile         string   //手机号
	Addr           string   //地址
	Lat            float64  //经度
	Lon            float64  //纬度
	Range          int      //方圆几公里
	Status         string   //需求的状态
	EntityTime     string   //创建时间
	RequirementWay string   // "FurtherProcessCharge","FurtherProcessUnCharge"
}

/*
 新建一条需求
*/
func NewRequirement(session *JsNet.StSession) {
	st := &Requirement{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.UID == "" || st.RequirementWay == "" ||
		st.Theme == "" || st.Content == "" || st.Name == "" ||
		st.Mobile == "" {
		ForwardEx(session, "1", nil,
			"NewRequirement failed,UID:%s,RequirementWay:%s,Theme=%s,Content=%s,Name=%s,Mobile:%s\n",
			st.UID, st.RequirementWay, st.Theme, st.Content, st.Name, st.Mobile)
		return
	}

	st.RequirementID = ider.GenID()
	st.EntityTime = CurTime()
	st.Status = constant.Status_New
	if st.RequirementWay == "FurtherProcessCharge" {
		st.Status = constant.Status_WaitPay
	} else {
		if len(st.Pictures) == 0 {
			ForwardEx(session, "1", nil, "NewRequirement failed,Picturs is empty\n")
			return
		}
	}

	if err := db.DirectWrite(constant.Hash_Requirement, st.RequirementID, st); err != nil {
		ForwardEx(session, "1", nil, "NewRequirement DirectWrite :"+err.Error())
		return
	}

	//添加一条全局的需求记录
	go appendRequiredToGolbal(st.RequirementID)
	///添加一条用户需求记录
	go AddUserRequirement(st.UID, st.RequirementID)
	//更新用户信息
	go UpdateUserInfo(st.UID, st.Name, st.Mobile, st.Addr)

	Forward(session, "0", st)
}

/**
 * 查询单个需求信息
 */
func QueryRequitement(session *JsNet.StSession) {
	type require struct {
		RequirementID string //需求id
	}
	st := &require{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "QueryRequitement GetPara :"+err.Error())
		return
	}

	data, err := GetRequirementInfo(st.RequirementID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", data)
}

/*
 获取所有的需求列表
*/
func GetAllRequirements(session *JsNet.StSession) {
	ids, err := getGlobalRequiredID()
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := []*Requirement{}
	for _, id := range ids {
		d, e := GetRequirementInfo(id)
		if e == nil {
			data = append(data, d)
		}
	}
	Forward(session, "0", data)
}

/*
	修改需求内容
*/
func ModifyRequirement(session *JsNet.StSession) {
	type modify struct {
		RequirementID string   //需求id
		Content       string   //内容
		Pictures      []string //海报
	}
	st := &modify{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ModifyRequirement GetPara:%s\n", err.Error())
		return
	}
	if st.RequirementID == "" {
		ForwardEx(session, "1", nil, "ModifyRequirement RequirementID is empty\n")
		return
	}
	re := &Requirement{}

	if err := db.WriteLock(constant.Hash_Requirement, st.RequirementID, re); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	re.Pictures = st.Pictures
	re.Content = st.Content
	if err := db.WriteBack(constant.Hash_Requirement, st.RequirementID, re); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", re)
}

///需求提交订单
func appendOrderID2Requirement(RequirementID, OrderID string) error {
	data := &Requirement{}
	if RequirementID == "" || OrderID == "" {
		return ErrorLog("appendOrderID2Requirement failed,RequirementID=%s,OrderID=%s\n", RequirementID, OrderID)
	}
	if err := db.WriteLock(constant.Hash_Requirement, RequirementID, data); err != nil {
		return err
	}
	data.OrderID = OrderID
	return db.WriteBack(constant.Hash_Requirement, RequirementID, data)
}

///需求海报支付完成后回调
func PayRequirement(RequirementID string) error {
	data := &Requirement{}
	if RequirementID == "" {
		return ErrorLog("GetRequirementInfo failed,RequirementID=%s\n", RequirementID)
	}
	if err := db.WriteLock(constant.Hash_Requirement, RequirementID, data); err != nil {
		return err
	}
	data.Status = constant.Status_Paid
	return db.WriteBack(constant.Hash_Requirement, RequirementID, data)
}

///获取需求信息
func GetRequirementInfo(keyValue string) (*Requirement, error) {

	data := &Requirement{}
	if keyValue == "" {
		return data, ErrorLog("GetRequirementInfo failed,RequirementID=%s\n", keyValue)
	}
	if err := db.ShareLock(constant.Hash_Requirement, keyValue, data); err != nil {
		return data, ErrorLog("GetRequirementInfo fail, ShareLock(),key=%s\n", keyValue)
	}
	return data, nil
}

//添加全局的需求id列表
func appendRequiredToGolbal(id string) error {
	data := []string{}
	err := db.WriteLock(constant.Hash_Requirement, constant.KEY_Global_Requirement, &data)
	AppendUniqueString(&data, id)
	if err != nil {
		return db.DirectWrite(constant.Hash_Requirement, constant.KEY_Global_Requirement, &data)
	}
	return db.WriteBack(constant.Hash_Requirement, constant.KEY_Global_Requirement, &data)
}

///获取全局的需求id
func getGlobalRequiredID() ([]string, error) {
	data := []string{}
	e := db.ShareLock(constant.Hash_Requirement, constant.KEY_Global_Requirement, &data)
	if e != nil {
		return data, e
	}
	return data, nil
}

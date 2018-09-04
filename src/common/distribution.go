package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	. "util"
)

type Distribution struct {
	Name         string   //网点名字
	EntityTime   string   //创建时间
	DName        string   //店铺联系人姓名
	DID          string   //网点ID
	Province     string   //省份
	City         string   //所在城市
	District     string   //区
	Street       string   //街道
	StreetNumber string   //街道号
	Area         string   //区域选择时候的城市，街道标识
	Address      string   //所在地址
	Lon          float32  //经度
	Lat          float64  //纬度
	Range        float64  //领取范围
	Contacts     string   //联系方式
	CID          string   //公司ID
	LsActivity   []string //店铺活动列表
	LsArea       []string //区域列表
	DLandPhone   string   //店铺座机
	DMobile      string   //店铺手机
	DSex         string   //联系人性别
	Amount       int      //默认发放金额
	OrderInCome  int      //分支机构年营业额
	ZipCode      string   //分支机构邮编
	Fax          string   //分支机构传真
}

type DistributionSI struct {
	DID         string   //网点ID
	Name        string   //网点名字
	City        string   //所在城市
	Address     string   //所在地址
	Lon         float64  //经度
	Lat         float64  //纬度
	Range       float64  //领取范围公里
	LsArea      []string //区域列表
	District    string   //区
	DLandPhone  string   //店铺座机
	DMobile     string   //店铺手机
	DName       string   //店铺联系人姓名
	TotalMoney  int      //历史总金额
	TotlaUser   int      //历史总人数
	Money       int      //分配的总金额
	SharedMoney int      //分享的佣金
	Remaining   int      //剩余金额
}

/*
 新建一个网点
*/
func NewDistribution(session *JsNet.StSession) {
	st := &Distribution{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.Lon == 0 || st.Lat == 0 || st.CID == "" {
		ForwardEx(session, "1", nil, "New Distribution failed with Lon=%d Lat=%d CID=%s\n", st.Lon, st.Lat, st.CID)
		return
	}
	st.DID = ider.GenID()
	st.EntityTime = CurTime()
	if err := db.DirectWrite(constant.Hash_Distribution, st.DID, st); err != nil {
		ForwardEx(session, "1", nil, "NewDistribution DirectWrite :"+err.Error())
		return
	}
	_, err := appendDistribution(st.CID, st.DID)

	if err != nil {
		ForwardEx(session, "1", nil, "Append Distribution failed:"+err.Error())
		return
	}
	Forward(session, "0", st)
}

/*
修改网点信息
*/
func ModifyDistribution(session *JsNet.StSession) {
	type modifyInfo struct {
		DID      string  //网点ID
		Name     string  //网点名字
		City     string  //所在城市
		Address  string  //所在地址
		Lon      float64 //经度
		Lat      float64 //纬度
		Contacts string  //联系方式
	}
	st := &modifyInfo{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ModifyDistribution GetPara :%s\n", err.Error())
		return
	}
	if st.DID == "" {
		ForwardEx(session, "1", nil, "ModifyDistribution GetPara failed, DID is empty\n")
		return
	}
	distribution := &Distribution{}
	if err := db.Modify(constant.Hash_Distribution, st.DID, distribution, st); err != nil {
		ForwardEx(session, "1", nil, "Modify Distribution:%s  failed\n", st.DID)
		return
	}
	Forward(session, "0", distribution)
}

/*
删除一个网点
*/
func DelDistribution(session *JsNet.StSession) {
	type info struct {
		DID string //网点id
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DID == "" {
		ForwardEx(session, "1", nil, "DelDistribution DID is empty\n")
		return
	}

	if err := db.HDel(constant.Hash_Distribution, st.DID); err != nil {
		ForwardEx(session, "1", nil, "Delete Distribution error\n")
		return

	}
	Forward(session, "0", nil)
}

/*
获取网点信息内部
*/

func GetDistributionInfo(DID string) (st *Distribution, e error) {
	data := &Distribution{}
	if DID == "" {
		return data, ErrorLog("GetDistributionInfo failed,DID is empty\n")
	}
	if err := db.ShareLock(constant.Hash_Distribution, DID, data); err != nil {
		return data, err
	}
	return data, nil
}

/*
获取一个网点信息
*/
func QueryDistribution(session *JsNet.StSession) {
	type info struct {
		DID string
	}
	st := &info{}
	data := &Distribution{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.DID == "" {
		ForwardEx(session, "1", nil, "DelDistribution DID is empty\n")
		return
	}
	if err := db.ShareLock(constant.Hash_Distribution, st.DID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

/*
添加一个活动到网点
*/
func AppendActivityToDistribution(activityID, distributionID string) (e error) {
	if activityID == "" || distributionID == "" {
		return ErrorLog("ActivityID=%s,DistributionID=%s is not correct", activityID, distributionID)
	}
	st := &Distribution{}
	if err := db.WriteLock(constant.Hash_Distribution, distributionID, st); err != nil {
		return err
	}
	for _, v := range st.LsActivity {
		if v == activityID {
			db.WriteBack(constant.Hash_Distribution, distributionID, st)
			return nil
		}
	}
	st.LsActivity = append(st.LsActivity, activityID)
	return db.WriteBack(constant.Hash_Distribution, distributionID, st)
}

/*
查询一个网点下面所有的活动
*/
func GetAllActivityOfDistribution(session *JsNet.StSession) {
	type info struct {
		DID string //网点ID
	}
	lsActivity := []*Activity{}
	st := &info{}
	dis := &Distribution{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "Get the net parameter error:%s\n", err.Error())
		return
	}
	if err := db.ShareLock(constant.Hash_Distribution, st.DID, dis); err != nil {
		ForwardEx(session, "1", nil, "Get the Distribution error:%s\n", err.Error())
		return
	}

	for _, v := range dis.LsActivity {
		activity, err := GetActivityInfo(v)
		if err == nil && activity.Status != constant.Status_Del {
			lsActivity = append(lsActivity, activity)

		}
	}
	ForwardEx(session, "0", lsActivity, "sucess")
}

/*
获取公司所有网点信息
*/
func QueryCompayDistribution(session *JsNet.StSession) {
	type info struct {
		CID string //公司ID
	}

	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "Get the info error:%s\n", err.Error())
		return
	}
	lsDistribution := []*Distribution{}
	dis, err := getCompanyDistribution(st.CID)
	if err != nil {
		ForwardEx(session, "1", lsDistribution, err.Error())
		return
	}

	for i := len(dis.DistributionID); i > 0; i-- {
		disprition, err := GetDistributionInfo(dis.DistributionID[i-1])
		if err == nil {
			lsDistribution = append(lsDistribution, disprition)
		}
	}
	ForwardEx(session, "0", lsDistribution, "sucess")
}

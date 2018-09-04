package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"math"
	"strings"
	. "util"
)

type AcADP struct {
	ActivityID string
	Money      int
	Date       string
	RecevieUID string
}

type UserActivity struct {
	UID            string         //用户id
	Requirement    []string       //用户的需求id
	Send           []AcADP        //所有发出去的红包列表
	Receive        []AcADP        //所有收到的红包列表
	Shared         []AcADP        //分享的红包详情
	SendMoney      int            //发出去金额
	RecevieMoney   int            //抢到的金额
	SharedMoney    int            //分享的佣金
	Re_Detail      []*Requirement //需求详细列表
	Send_Detail    []*Activity    //发出去的活动详细列表
	Receive_Detail []*Activity    //抢到的活动详细列表
	Date           string         //创建时间
}

func (this *UserActivity) appendRequirement(id string) {
	AppendUniqueString(&this.Requirement, id)
}

func (this *UserActivity) appendSend(ActivityID string, Money int) {
	exist := false
	for i, v := range this.Send {
		if v.ActivityID == ActivityID {
			this.SendMoney -= v.Money
			this.Send[i].Money = Money
			this.SendMoney += Money
			exist = true
		}
	}
	if !exist {
		this.Send = append(this.Send, AcADP{
			ActivityID: ActivityID,
			Money:      Money,
			Date:       CurTime(),
		})
		this.SendMoney += Money
	}
}

func (this *UserActivity) appendReceive(ActivityID string, Money int) {
	this.Receive = append(this.Receive, AcADP{
		ActivityID: ActivityID,
		Money:      Money,
		Date:       CurTime(),
	})
	this.RecevieMoney += Money
}

func (this *UserActivity) appendShared(ActivityID, RecevieUID string, Money int) {
	this.Shared = append(this.Shared, AcADP{
		ActivityID: ActivityID,
		Money:      Money,
		RecevieUID: RecevieUID,
		Date:       CurTime(),
	})
	this.SharedMoney += Money
}

//获取红包
func GetRedPacket(session *JsNet.StSession) {
	type RD_ActivityRequest struct {
		UID                 string           //用户id
		ActivityID          string           //活动id
		Lat                 float64          //经度
		Lon                 float64          //纬度
		City                string           //详细地址
		IsCorrect           bool             //是否答题正确
		LTime               int              //停留时间
		LsSystemQuestion    []SystemQuestion //系统问题
		LsActivitySQuestion []SystemQuestion //活动系统问题
		Tag                 string           //平台标签
		SharedID            string           //分享者uid
	}
	DIDGen := ""
	st := &RD_ActivityRequest{}
	if err := session.GetPara(st); err != nil {
		ErrorLog("GetRedPacket GetPara :%s\n", err.Error())
		ForwardEx(session, "1", nil, "GetRedPacket GetPara :%s\n", err.Error())
		return
	}

	if st.ActivityID == "" || st.UID == "" || st.Tag == "" {
		ForwardEx(session, "1", nil, "GetRedPacket param empty, ActivityID=%s,UID=%s,Tag=%s\n", st.ActivityID, st.UID, st.Tag)
		return
	}

	isShared := (st.SharedID != "") && (st.SharedID != st.UID)

	user, err := GetUserInfo(st.UID)
	if err != nil {
		ErrorLog(err.Error())
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	AC := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, st.ActivityID, AC); err != nil {
		ErrorLog(err.Error())
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	defer db.WriteBack(constant.Hash_Activity, st.ActivityID, AC)
	if AC.Status == constant.Status_WaitPay {
		ErrorLog("GetRedPacket failed, Activity:%s  is waitPay\n")
		ForwardEx(session, "1", nil, "GetRedPacket failed, Activity:%s  is waitPay\n")
		return
	}
	DisRemaining := -1
	if AC.TestFlag == 1 || AC.TestFlag == -1 {
		///检查是否在认领范围内
		IsInRange := true
		if AC.MoneyArrangeWay == constant.MoneyArrange_TradArea {
			IsInRange = InTradingArea(st.Lat, st.Lon, AC.TardingArea)
		} else if AC.MoneyArrangeWay == constant.MoneyArrange_Coordinate {
			if AC.ArrangeType == constant.Arrange_Distribution {
				ok, DID := GetUserDistributionIDLonLat(st.Lat, st.Lon, AC.LsDistribution)
				DIDGen = DID
				IsInRange = ok
			} else {
				IsInRange = InLatLonRange(st.Lat, st.Lon, AC.LatLon)
			}
		} else if AC.MoneyArrangeWay == constant.MoneyArrange_Division {
			if AC.ArrangeType == constant.Arrange_Distribution {
				ok, DID := GetUserDistributionIDArea(st.City, AC.LsDistribution)
				DIDGen = DID
				IsInRange = ok
			} else {
				IsInRange = InAreaRange(st.City, AC.Area)
			}
		}
		if !IsInRange {
			go AppendActivityRDP(user, st.SharedID, st.ActivityID, AC.Status, st.City, st.Lat, st.Lon, AC.TestFlag, DisRemaining, st.LTime, st.IsCorrect, false)
			go appendCompanyCustomer(user, AC.CID, DIDGen, AC.ActivityID, AC.Theme, 0, st.LTime, st.Lat, st.Lon, st.City, st.IsCorrect, false)
			ErrorLog("User:%s is not in Activity:%s  range \n", user.UID, st.ActivityID)
			ForwardEx(session, "1", nil, "User:%s is not in Activity:%s  range \n", user.UID, st.ActivityID)
			return
		}
	}

	////查找网点的余额
	DisIndex := -1
	if AC.ArrangeType == constant.Arrange_Distribution {
		if DIDGen != "" {
			for i, v := range AC.LsDistribution {
				if v.DID == DIDGen {
					DisRemaining = AC.LsDistribution[i].Remaining
					DisIndex = i
					break
				}
			}
		}
	}

	//往活动认领记录添加一条，并检查是否领过
	shardMoney, money, IsStop, err := AppendActivityRDP(user, st.SharedID, st.ActivityID, AC.Status, st.City, st.Lat, st.Lon, AC.TestFlag, DisRemaining, st.LTime, st.IsCorrect, true)
	if err != nil {
		ErrorLog("AppendActivityRDP failed,User:%s,Activity:%s\n", user.UID, st.ActivityID)
		ForwardEx(session, "1", nil, "AppendActivityRDP failed,User:%s,Activity:%s\n", user.UID, st.ActivityID)
		return
	}
	if IsStop {
		AC.Status = constant.Status_Arrears
	}

	////更新网点的人数和总金额
	if AC.ArrangeType == constant.Arrange_Distribution {
		if DIDGen != "" && DisIndex != -1 &&
			len(AC.LsDistribution) > DisIndex &&
			AC.LsDistribution[DisIndex].DID == DIDGen {
			AC.LsDistribution[DisIndex].TotalMoney += money
			AC.LsDistribution[DisIndex].SharedMoney += shardMoney
			AC.LsDistribution[DisIndex].TotlaUser++
			AC.LsDistribution[DisIndex].Remaining -= money
			AC.LsDistribution[DisIndex].Remaining -= shardMoney
			if AC.LsDistribution[DisIndex].Remaining <= 0 {
				AC.LsDistribution[DisIndex].Remaining = 0
			}
		}
	}

	if st.IsCorrect {
		///更新分享佣金
		if isShared && shardMoney > 0 {
			go AddUserSharedRDP(st.SharedID, st.UID, st.ActivityID, shardMoney)
		}
		//更新抢红包的用户
		go AddUserRecevieRDP(user, st.ActivityID, st.Tag, money, AC.IsRealTime)
		//更新公司用户
		go appendCompanyCustomer(user, AC.CID, DIDGen, AC.ActivityID, AC.Theme, money, st.LTime, st.Lat, st.Lon, st.City, st.IsCorrect, true)
		//更新活动问题
		go AnswerActivitySQuestion(st.UID, st.ActivityID, st.LsActivitySQuestion)
		//更新系统活动问题
		go AnswerSystemQuestion(st.UID, st.LsSystemQuestion)
		Forward(session, "0", money)
		return
	}
	ErrorLog("问题回答错误\n")

	ForwardEx(session, "1", nil, "Question answer fail")
}

func DelUserRequirement(session *JsNet.StSession) {
	type RD_Request struct {
		UID           string //用户id
		RequirementID string //需求id
	}
	st := &RD_Request{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.UID == "" || st.RequirementID == "" {
		ForwardEx(session, "1", nil, "DelUserRequirement failed,UID:%s,RequirementID:%s\n", st.UID, st.RequirementID)
		return
	}

	d, err := GetRequirementInfo(st.RequirementID)
	if err != nil {
		ForwardEx(session, "1", nil, "DelUserRequirement failed,err:%s\n", err.Error())
		return
	}
	if d.Status == constant.Status_Paid {
		ForwardEx(session, "1", nil, "DelUserRequirement failed,RequirementID:%s is paid\n", st.RequirementID)
		return
	}

	re := &UserActivity{}
	if err := db.WriteLock(constant.Hash_UserActivity, st.UID, re); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	//删除存在的用户需求id
	DelExistString(&re.Requirement, st.RequirementID)

	if err := db.WriteBack(constant.Hash_UserActivity, st.UID, re); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", re)
}

//获取用户的活动相关信息
func GetUserActivityInfo(session *JsNet.StSession) {
	type RD_Request struct {
		UID string
	}
	st := &RD_Request{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data, err := GetUserActivity(st.UID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	for _, v := range data.Requirement {
		ac, err := GetRequirementInfo(v)
		if err == nil && ac.Status != constant.Status_Finished {
			data.Re_Detail = append(data.Re_Detail, ac)
		}
	}

	for _, v := range data.Send {
		ac, err := GetActivityInfo(v.ActivityID)
		if err == nil && ac.Status != constant.Status_Del {
			data.Send_Detail = append(data.Send_Detail, ac)
		}
	}

	for _, v := range data.Receive {
		ac, err := GetActivityInfo(v.ActivityID)
		if err == nil && ac.Status != constant.Status_Del {
			data.Receive_Detail = append(data.Receive_Detail, ac)
		}
	}

	Forward(session, "0", data)
}

///获取用户的红包相关信息
func GetUserActivity(Uuid string) (*UserActivity, error) {
	ac := &UserActivity{}
	err := db.ShareLock(constant.Hash_UserActivity, Uuid, ac)
	return ac, err
}

//添加一个需求记录
func AddUserRequirement(UID, RequirementID string) error {
	if UID == "" || RequirementID == "" {
		return ErrorLog("AddUserRequirement failed ,parame is empty\n")
	}
	ac := &UserActivity{}
	err := db.WriteLock(constant.Hash_UserActivity, UID, ac)
	ac.Requirement = append(ac.Requirement, RequirementID)
	if err != nil {
		ac.UID = UID
		ac.Date = CurTime()
		return db.DirectWrite(constant.Hash_UserActivity, UID, ac)
	}
	return db.WriteBack(constant.Hash_UserActivity, UID, ac)
}

//添加一个发出去的红包
func AddUserSendRDP(UID, ActivityID string, money int) error {
	if UID == "" || ActivityID == "" {
		return ErrorLog("AddUserSendRDP failed ,parame is empty\n")
	}
	ac := &UserActivity{}
	err := db.WriteLock(constant.Hash_UserActivity, UID, ac)
	ac.appendSend(ActivityID, money)
	if err != nil {
		ac.UID = UID
		ac.Date = CurTime()
		return db.DirectWrite(constant.Hash_UserActivity, UID, ac)
	}
	if err := db.WriteBack(constant.Hash_UserActivity, UID, ac); err != nil {
		return err
	}

	///更新用户账户信息
	return updataUserAccount(UID, 0, money, 0, 3)
}

//添加一个分享的红包
func AddUserSharedRDP(SharedID, UID, ActivityID string, money int) error {
	if SharedID == "" || ActivityID == "" || UID == "" {
		return ErrorLog("AddUserSharedRDP failed ,parame is empty\n")
	}
	go updataUserAccount(SharedID, 0, 0, money, 2)

	ac := &UserActivity{}
	err := db.WriteLock(constant.Hash_UserActivity, SharedID, ac)
	ac.appendShared(ActivityID, UID, money)
	if err != nil {
		ac.UID = SharedID
		ac.Date = CurTime()
		return db.DirectWrite(constant.Hash_UserActivity, SharedID, ac)
	}
	if err := db.WriteBack(constant.Hash_UserActivity, SharedID, ac); err != nil {
		return err
	}
	return nil
}

//添加一个收到的红包
func AddUserRecevieRDP(user *User, ActivityID, Tag string, money int, IsRealTime int) error {
	if user == nil || ActivityID == "" {
		return ErrorLog("AddUserRecevieRDP failed ,parame is empty\n")
	}
	///实时到账,直接打款
	if IsRealTime == 1 {
		if money >= 100 {
			go SendActivityRDP(user, money, Tag)
		}
	} else {
		///更新用户账户信息
		go updataUserAccount(user.UID, money, 0, 0, 1)
	}

	ac := &UserActivity{}
	err := db.WriteLock(constant.Hash_UserActivity, user.UID, ac)
	ac.appendReceive(ActivityID, money)
	if err != nil {
		ac.UID = user.UID
		ac.Date = CurTime()
		return db.DirectWrite(constant.Hash_UserActivity, user.UID, ac)
	}
	if err := db.WriteBack(constant.Hash_UserActivity, user.UID, ac); err != nil {
		return err
	}
	return nil
}

/*
 求经纬度距离
*/
func InLatLonRange(oLat, oLon float64, latlon []LATLON) bool {
	for _, v := range latlon {
		return EarthDistance(oLat, oLon, v.Lat, v.Lon) <= v.Range
	}
	return false
}

//查看是否在同一个区域
func InAreaRange(City string, Area []string) bool {
	exist := false
	for _, v := range Area {
		if City != "" && v != "" && strings.Contains(City, v) {
			exist = true
			break
		}
	}
	return exist
}

//是否在商区内
func InTradingArea(Lat, Lon float64, area []string) bool {
	exist := false
	data := getMoreTradingRrea(area)
	for _, v := range data {
		if v.Shape == 0 {
			tmpx := []float64{}
			tmpy := []float64{}
			for _, v1 := range v.Points {
				tmpx = append(tmpx, v1.Lat)
				tmpy = append(tmpy, v1.Lon)
			}
			if pnpoly(len(v.Points), tmpx, tmpy, Lat, Lon) {
				exist = true
				break
			}

		} else {
			for _, v1 := range v.Points {
				if EarthDistance(Lat, Lon, v1.Lat, v1.Lon) <= v1.Range {
					exist = true
					break
				}
			}
		}
	}
	return exist
}

//判断一个点是否在多边形内
func pnpoly(count int, vertx []float64, verty []float64, testx float64, testy float64) bool {
	exist := false
	j := count - 1
	for i := 0; i < count; i++ {
		if ((verty[i] > testy) != (verty[j] > testy)) &&
			(testx < (vertx[j]-vertx[i])*(testy-verty[i])/(verty[j]-verty[i])+vertx[i]) {
			exist = !exist
		}
		j = i
	}
	return exist
}

///计算经纬度距离
func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6371000.0 // 6378137
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return math.Abs((float64)(dist * radius))
}

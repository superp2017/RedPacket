package common

import (
	."JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	// "errors"
	. "util"
)

type ShortInfo struct {
	Key   string
	Value string
}

type CUser struct {
	CID     string                 //公司ID
	HmCUser map[string][]ShortInfo //用户资料 Key- UID
}

type UserCInfo struct {
	UID     string
	Mobile  string
	HmCUser map[string][]ShortInfo //用户资料  KEY-Company ID.
}

type CUserRecord struct {
	CID    string      //公司ID
	UID    string      //UID
	Mobile string      //用户手机号码
	Info   []ShortInfo //用户资料
}

//录入公司信息 总的Key 公司CID;   key-- 用户ID; 内容--客户具体信息

/*
客户信息录入
*/
func RecordCUserInfo(session *JsNet.StSession) {
	type CUserRecord struct {
		CID    string      //公司ID
		UID    string      //UID
		Mobile string      //用户手机号码
		Info   []ShortInfo //用户资料
	}
	st := &CUserRecord{}
	cudb := &CUser{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" || st.UID == "" || st.Mobile == "" {
		ForwardEx(session, "1", nil, "New CUser failed with CID=%s and UID=%s Mobile=%s\n", st.CID, st.UID, st.Mobile)
		return
	}
	data := &User{}
	if err := db.ShareLock(constant.Hash_User, st.UID, data); err != nil {
		ForwardEx(session, "1", nil, "There is no such an user with UID=%s\n", st.UID)
		return
	}

	if err := db.WriteLock(constant.Hash_CRecordUser, st.CID, cudb); err != nil {
		ForwardEx(session, "1", nil, err.Error(), st.UID)
		return
	}
	if cudb.HmCUser == nil {
		cudb.HmCUser = make(map[string][]ShortInfo)
	}
	cudb.HmCUser[st.UID] = st.Info

	if err := db.WriteBack(constant.Hash_CRecordUser, st.CID, cudb); err != nil {
		ForwardEx(session, "1", nil, err.Error(), st.UID)
		return
	}

	//录入客户档案 --总的key 用户UID； key 公司ID   内容 --用户
	userCInfo := &UserCInfo{}
	directWrite := false
	if err := db.WriteLock(constant.Hash_UserDesciption, st.UID, userCInfo); err != nil {
		directWrite = true
	}

	if userCInfo.HmCUser == nil {
		userCInfo.HmCUser = make(map[string][]ShortInfo)
	}
	userCInfo.HmCUser[st.CID] = st.Info
	if directWrite {
		db.DirectWrite(constant.Hash_UserDesciption, st.UID, userCInfo)
	} else {
		db.WriteBack(constant.Hash_UserDesciption, st.UID, userCInfo)
	}

	//录入客户档案 --总的key 用户Mobile； key 公司ID   内容 --用户

	//写用户表
	if err := db.DirectWrite(constant.Hash_CUser, st.CID+"@"+st.UID, st); err != nil {
		ForwardEx(session, "1", nil, err.Error(), st.UID)
		return
	}
	//写Mobile-UID对应表格
	if err := db.DirectWrite(constant.Hash_MobileCUser, st.CID+"@"+st.Mobile, st); err != nil {
		ForwardEx(session, "1", nil, err.Error(), st.UID)
		return
	}
	Forward(session, "0", nil)
}

//获取公司所有录入客户信息

func QueryCUserRecord(session *JsNet.StSession) {
	type info struct {
		CID    string
		Mobile string
	}
	st := &info{}
	data := &CUser{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" {
		ForwardEx(session, "1", nil, "DelCompany CID is empty\n")
		return
	}
	if err := db.ShareLock(constant.Hash_CRecordUser, st.CID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

//
/*
手机号获取客户信息
*/
func QueryCompanyUserInfo(session *JsNet.StSession) {
	type info struct {
		CID    string
		Mobile string
	}
	st := &info{}

	data := &CUserRecord{}

	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" || st.Mobile == "" {
		ForwardEx(session, "1", nil, "Mobile or  CID is empty\n")
		return
	}
	if err := db.ShareLock(constant.Hash_MobileCUser, st.CID+"@"+st.Mobile, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////以下是从活动添加的用户//////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////

///添加一个新用户
func appendCompanyCustomer(user *User, CID, DID, ActivityID, ActivityName string,
	Money, LTime int, Lat float64, Lon float64, Address string, IsCorrect, InRange bool) error {
	if user == nil || ActivityID == "" || CID == "" {
		return ErrorLog("appendCompanyCustomer param failed,ActivityID=%s,CID=%s\n", ActivityID, CID)
	}
	st := &ST_Customer{}
	err := db.WriteLock(constant.Hash_CompanyU, CID, st)
	if st.Increase == nil {
		st.Increase = make(map[string]int)
	}
	exist := false
	ac := LActivity{
		ActivityID:   ActivityID,
		ActivityName: ActivityName,
		Money:        Money,
		IsCoupon:     false,
		LatF:         Lat,
		LonF:         Lon,
		Address:      Address,
		Date:         CurTime(),
		IsCorrect:    IsCorrect,
		InRange:      InRange,
	}
	for i, v := range st.LsCSI {
		if v.UID == user.UID {
			ex := false
			for _, v1 := range v.LsLActivity {
				if ActivityID == v1.ActivityID {
					ex = true
					break
				}
			}
			if !ex {
				st.LsCSI[i].LsLActivity = append(st.LsCSI[i].LsLActivity, ac)
				st.LsCSI[i].TotleMoney += Money
			}
			st.LsCSI[i].Address = Address
			st.LsCSI[i].HeadImageURL = user.HeadImageURL
			st.LsCSI[i].LatF = Lat
			st.LsCSI[i].LonF = Lon
			st.LsCSI[i].Mobile = user.Mobile
			st.LsCSI[i].Name = user.Nickname
			if st.LsCSI[i].CreatDate == "" {
				st.LsCSI[i].CreatDate = CurTime()
			}
			exist = true
			break
		}
	}
	if !exist {
		st.LsCSI = append(st.LsCSI,
			CSInfo{
				UID:          user.UID,
				Name:         user.Nickname,
				Mobile:       user.Mobile,
				DID:          DID,
				HeadImageURL: user.HeadImageURL,
				Address:      Address,
				LonF:         Lon,
				LatF:         Lat,
				LsLActivity:  []LActivity{ac},
				TotleMoney:   Money,
				CreatDate:    CurTime(),
			})
		/////更新增加量
		date := CurDate()
		if v, ok := st.Increase[date]; ok {
			st.Increase[date] = v + 1
		} else {
			st.IncreaseKey = append(st.IncreaseKey, date)
			st.Increase[date] = 1
		}
	}

	if err != nil {
		st.ID = CID
		return db.DirectWrite(constant.Hash_CompanyU, CID, st)
	}
	return db.WriteBack(constant.Hash_CompanyU, CID, st)
}

//更新用户是否领了券
func updateCompanyCustomer(UID, CID, ActivityID string) error {
	if UID == "" || CID == "" || ActivityID == "" {
		return ErrorLog("updateCustomer param failed\n")
	}
	st := &ST_Customer{}
	if err := db.WriteLock(constant.Hash_CompanyU, CID, st); err != nil {
		return err
	}
	for i, v := range st.LsCSI {
		if v.UID == UID {
			for j, v1 := range v.LsLActivity {
				if v1.ActivityID == ActivityID {
					st.LsCSI[i].LsLActivity[j].IsCoupon = true
					break
				}
			}
			break
		}
	}
	return db.WriteBack(constant.Hash_CompanyU, CID, st)
}

///查询公司的用户
func queryCompanyCustomer(CID string) (*ST_Customer, error) {
	st := &ST_Customer{}
	err := db.ShareLock(constant.Hash_CompanyU, CID, st)
	return st, err
}

////查询公司用户
func QueryCompanyCustomer(session *JsNet.StSession) {
	type INFO struct {
		CID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	cu, err := queryCompanyCustomer(st.CID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", cu)
}

////获取某个网点的用户
func QueryDistributionCustomer(session *JsNet.StSession) {
	type INFO struct {
		DID string
		CID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" || st.DID == "" {
		ForwardEx(session, "1", nil, "GetDistributionCustomer param failed,CID=%s,DID=%s\n", st.CID, st.DID)
		return
	}
	customer, err := queryCompanyCustomer(st.CID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	user := []CSInfo{}
	for _, v := range customer.LsCSI {
		if v.DID == st.DID {
			user = append(user, v)
		}
	}
	Forward(session, "0", user)
}

///获取活动的顾客列表
func QueryActivityCustomer(session *JsNet.StSession) {
	type INFO struct {
		// CID        string
		ActivityID string //活动id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ActivityID == "" {
		ForwardEx(session, "1", nil, "GetActivityCustomer failed,ActivityID = %s \n", st.ActivityID)
		return
	}

	// user := []CSInfo{}
	record, err := getActivityRDP(st.ActivityID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	// if st.CID != "" {
	// 	customer, err := queryCompanyCustomer(st.CID)
	// 	if err != nil {
	// 		ForwardEx(session, "1", nil, err.Error())
	// 		return
	// 	}
	// 	for _, v := range customer.LsCSI {
	// 		for _, v1 := range record.RDP {
	// 			if v.UID == v1.UID {
	// 				user = append(user, v)
	// 			}
	// 		}
	// 	}
	// } else {
	// 	for _, v := range record.RDP {
	// 		user = append(user, CSInfo{
	// 			UID:          v.UID,
	// 			Name:         v.Name,
	// 			Mobile:       v.Mobile,
	// 			LatF:         v.Lat,
	// 			LonF:         v.Lon,
	// 			Address:      v.City,
	// 			CreatDate:    v.Date,
	// 			HeadImageURL: v.HeadImageURL,
	// 		})
	// 	}
	// }
	Forward(session, "0", record)
}

////获取某个网点的用户
func QueryActivityDistributionCustomer(session *JsNet.StSession) {
	type INFO struct {
		CID        string
		DID        string
		ActivityID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" || st.DID == "" || st.ActivityID == "" {
		ForwardEx(session, "1", nil,
			"GetDistributionActivityCustomer param failed,CID=%s,DID=%s,ActivityID=%s\n",
			st.CID, st.DID, st.ActivityID)
		return
	}
	customer, err := queryCompanyCustomer(st.CID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	user := []CSInfo{}
	for _, v := range customer.LsCSI {
		if v.DID == st.DID {
			for _, v1 := range v.LsLActivity {
				if v1.ActivityID == st.ActivityID {
					user = append(user, v)
				}
			}

		}
	}
	Forward(session, "0", user)
}

package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"

	. "util"
)

type CouponCustomer struct {
	CouponID string  //代金券ID
	CID      string  //公司ID
	DID      string  //网点id
	DisName  string  //网点名字
	UID      string  //用户id
	Name     string  //用户姓名
	Sex      string  //性别
	Lon      float64 //经度
	Lat      float64 //纬度
	Addr     string  //地址
	IsUse    bool    //是否使用
	Date     string  //领用日期
}

///添加一个代金券到用户
func AppendCouponToUser(session *JsNet.StSession) {
	type INFO struct {
		UID        string  //用户id
		CID        string  //公司ID
		CouponID   string  //代金券ID
		ActivityID string  //活动id
		DID        string  //网点id
		Lon        float64 //经度
		Lat        float64 //纬度
		Addr       string  //地址
	}

	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.CouponID == "" || st.ActivityID == "" {
		ForwardEx(session, "1", nil, "AppendCouponToUser param failed,UID=%s,CouponID=%s,ActivityID=%s\n",
			st.UID, st.CouponID, st.ActivityID)
		return
	}

	if err := AddUserCoupon(st.UID, st.CouponID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	///更新活动领券的情况
	go updateActivityRDP(st.ActivityID, st.UID)
	updateCompanyCustomer(st.UID, st.CID, st.ActivityID)
	go addNewCouponCustomer(st.UID, st.CID, st.DID, st.CouponID, st.Lon, st.Lat, st.Addr, false)
	Forward(session, "0", nil)
}

//获取用户的所有代金券
func GetUserCoupon(session *JsNet.StSession) {
	type INFO struct {
		UID string //用户id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	list, err := getUserCoupon(st.UID)
	if err != nil {
		ForwardEx(session, "1", list, err.Error())
		return
	}

	Forward(session, "0", list)
}

//扫码查找代金券
func FindValidCoupon(session *JsNet.StSession) {
	type INFO struct {
		UID string //用户id
		CID string //公司ID
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.CID == "" {
		ForwardEx(session, "1", nil, "FindValidCoupon param failed,UID=%s,CID=%s\n", st.UID, st.CID)
		return
	}

	list, err := getUserCoupon(st.UID)
	if err != nil {
		ForwardEx(session, "1", list, err.Error())
		return
	}

	data := []*Coupon{}

	for _, v := range list {
		if v.CID == st.CID {
			data = append(data, v)
		}
	}

	if len(data) == 0 {
		ForwardEx(session, "1", data, "未查到可用的券")
		return
	}

	Forward(session, "0", data)
}

////使用代金券
func UseCoupon(session *JsNet.StSession) {
	type INFO struct {
		UID      string //用户id
		CID      string //公司ID
		CouponID string //代金券ID
		DID      string //网点ID
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.UID == "" || st.CID == "" || st.CouponID == "" {
		ForwardEx(session, "1", nil, "UseCoupon param failed,UID=%s,CID=%s,CouponID=%s,DID=%s\n", st.UID, st.CID, st.CouponID, st.DID)
		return
	}

	list, err := getUserCoupon(st.UID)
	if err != nil {
		ForwardEx(session, "1", list, err.Error())
		return
	}
	exist := false

	cou := &Coupon{}
	for _, v := range list {
		if v.CouponID == st.CouponID && v.CID == st.CID {
			exist = true
			cou = v
		}
	}
	if !exist {
		ForwardEx(session, "1", nil, "未查到可用的券")
		return
	}

	cur := CurStamp()
	if cur >= cou.StartStamp && cur <= cou.StopStamp && cou.Status != -1 && cou.Status != 1 {
		go removeUserCoupon(st.UID, st.CouponID)
		go appendCompanyUsedCoupon(st.UID, st.CID, st.DID, st.CouponID)
		go addNewCouponCustomer(st.UID, st.CID, st.DID, st.CouponID, 0, 0, "", true)
		ForwardEx(session, "0", nil, "代金券使用成功\n")
		return
	}

	ForwardEx(session, "1", nil, "现金券不在使用范围,CouponID:%s,StartTime:%s,StopTime:%s\n", cou.CouponID, cou.StartTime, cou.StopTime)
}

///查询某个公司各个网点代金券使用的情况
func QueryCompanyUseedCoupon(session *JsNet.StSession) {
	type INFO struct {
		CID string //公司id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" {
		ForwardEx(session, "1", nil, "QueryCompanyUseedCoupon CID is empty\n")
		return
	}

	data, err := getCompantUsedCoupon(st.CID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

func QueryDistributionUseedCoupon(session *JsNet.StSession) {
	type INFO struct {
		CID string //公司id
		DID string //网点id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" || st.DID == "" {
		ForwardEx(session, "1", nil, "QueryDistributionUseedCoupon param failed,CID:%s,DID:%s\n", st.CID, st.DID)
		return
	}
	data, err := getCompantUsedCoupon(st.CID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	d, ok := data[st.DID]
	if !ok {
		ForwardEx(session, "1", []UseedCoupon{}, "QueryDistributionUseedCoupon not exist,DID:%s\n", st.DID)
		return
	}
	Forward(session, "0", d)
}

//添加券到用户券列表
func AddUserCoupon(UID, CouponID string) error {

	cou, err := GetCouponInfo(CouponID)
	if err != nil {
		return err
	}
	if cou.Status != 0 {
		return ErrorLog("该券已经过期或者删除,CouponID:%s\n", CouponID)
	}

	list := &[]string{}
	e := db.WriteLock(constant.Hash_UserCoupon, UID, list)
	AppendUniqueString(list, CouponID)
	if e != nil {
		return db.DirectWrite(constant.Hash_UserCoupon, UID, list)
	}
	return db.WriteBack(constant.Hash_UserCoupon, UID, list)
}

///移除用户的券
func removeUserCoupon(UID, CouponID string) error {
	list := &[]string{}
	err := db.WriteLock(constant.Hash_UserCoupon, UID, list)
	DelExistString(list, CouponID)
	if err != nil {
		return db.DirectWrite(constant.Hash_UserCoupon, UID, list)
	}
	return db.WriteBack(constant.Hash_UserCoupon, UID, list)
}

func getUserCoupon(UID string) ([]*Coupon, error) {
	list := []string{}
	data := []*Coupon{}
	if err := db.ShareLock(constant.Hash_UserCoupon, UID, &list); err != nil {
		return data, err
	}

	for _, v := range list {
		s, err := GetCouponInfo(v)
		if err == nil {
			cur := CurStamp()
			if cur <= s.StopStamp && s.Status != -1 && s.Status != 1 {
				data = append(data, s)
			}
		}
	}

	return data, nil
}

type UseedCoupon struct {
	CouponID string   //券id
	User     []string ///使用人的id列表
}

///添加用户使用过的代金券到公司
func appendCompanyUsedCoupon(UID, CID, DID, CouponID string) error {
	if UID == "" || CID == "" || DID == "" || CouponID == "" {
		return ErrorLog("appendCompanyUsedCoupon param failed,UID:%s,CID:%s,DID:%s,CouponID:%s\n", UID, CID, DID, CouponID)
	}

	data := make(map[string][]UseedCoupon)
	err := db.WriteLock(constant.Hash_CompanyCouponUsed, CID, &data)

	list := []UseedCoupon{}
	exist := false
	list, _ = data[DID]
	for _, v := range list {
		if v.CouponID == CouponID {
			exist = true
			break
		}
	}
	if !exist {
		list = append(list, UseedCoupon{
			CouponID: CouponID,
			User:     []string{UID},
		})
	}
	data[DID] = list
	if err != nil {
		return db.DirectWrite(constant.Hash_CompanyCouponUsed, CID, &data)
	}
	return db.WriteBack(constant.Hash_CompanyCouponUsed, CID, &data)
}

///获取某个公司使用过的代金券
func getCompantUsedCoupon(CID string) (map[string][]UseedCoupon, error) {
	data := make(map[string][]UseedCoupon)
	err := db.ShareLock(constant.Hash_CompanyCouponUsed, CID, &data)
	return data, err
}

///获取某个券的用户使用情况
func QueryCoupinCustomer(session *JsNet.StSession) {
	type INFO struct {
		CouponID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CouponID == "" {
		ForwardEx(session, "1", nil, "QueryCoupinCustomer CouponID is empty\n")
		return
	}
	data, err := getCouponCustomer(st.CouponID)
	if err != nil {
		ForwardEx(session, "0", data, err.Error())
		return
	}
	Forward(session, "0", data)
}

///添加一个券的领用者
func addNewCouponCustomer(UID, CID, DID, CouponID string, Lon, Lat float64, Addr string, IsUse bool) error {
	if UID == "" || CouponID == "" {
		return ErrorLog("addNewCouponCustomer param failed,UID:%s,CID:%s,DID:%s,CouponID:%s\n",
			UID, CID, DID, CouponID)
	}

	list := []*CouponCustomer{}
	err := db.WriteLock(constant.Hash_CouponCustomer, CouponID, &list)

	data := &CouponCustomer{
		UID:      UID,
		CID:      CID,
		DID:      DID,
		CouponID: CouponID,
		Lon:      Lon,
		Lat:      Lat,
		Addr:     Addr,
		IsUse:    IsUse,
		Date:     CurTime(),
	}

	exist := false
	for i, v := range list {
		if v.UID == UID {
			list[i].IsUse = IsUse
			exist = true
			break
		}
	}
	if !exist {
		if DID != "" {
			dis, e := GetDistributionInfo(DID)
			if e == nil {
				data.DisName = dis.Name
			}
		}
		user, e1 := GetUserInfo(UID)
		if e1 == nil {
			data.Name = user.Nickname
			data.Sex = user.Sex
		}
		list = append(list, data)
	}

	if err != nil {
		return db.DirectWrite(constant.Hash_CouponCustomer, CouponID, &list)
	}
	return db.WriteBack(constant.Hash_CouponCustomer, CouponID, &list)
}

///获取某个券的用户使用情况
func getCouponCustomer(CouponID string) ([]*CouponCustomer, error) {
	list := []*CouponCustomer{}
	err := db.ShareLock(constant.Hash_CouponCustomer, CouponID, &list)
	return list, err
}

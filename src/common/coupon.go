package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	. "util"
)

type Coupon struct {
	CouponID    string   //代金券ID
	CID         string   //公司ID
	Picture     string   //代金券路径
	Money       int      //代金券金额
	StartTime   string   //代金券起始日期
	StopTime    string   //代金券结束日期
	StartStamp  int64    //开始日期时间戳
	StopStamp   int64    //结束日期时间戳
	LsDis       []string //适用此券的网点id
	IsFitAll    bool     //是否使用所有的网点
	EntityTime  string   //创建日期
	Title       string   //代金券标题
	Instruction string   //代金券使用说明
	Status      int      //券的状态  0:正常 1：过期 -1：删除
}

/*
 添加一张金券
*/
func NewCoupon(session *JsNet.StSession) {
	st := &Coupon{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if st.StartTime == "" || st.StopTime == "" || st.Money <= 0 {
		ForwardEx(session, "1", nil, "NewCoupon param failed,Money:%d,StartTime:=%s,StopTime:%s,LsDis:%d",
			st.Money, st.StartTime, st.StopTime, len(st.LsDis))
		return
	}

	start, err1 := ParseTimeFromString(st.StartTime, "2006-01-02")
	end, err2 := ParseTimeFromString(st.StopTime, "2006-01-02")
	if err1 != nil || err2 != nil {
		ForwardEx(session, "1", nil, "NewCoupon Time format error,such as '2006-01-02'\n")
		return
	}

	cur := CurStamp()
	if cur > end.Unix() || start.Unix() > end.Unix() {
		ForwardEx(session, "1", nil, "不能添加已经过期的代金券\n")
		return
	}

	st.StartStamp = start.Unix()
	st.StopStamp = end.Unix()
	st.CouponID = ider.GenID()
	st.EntityTime = CurTime()
	st.Status = 0
	if err := db.DirectWrite(constant.Hash_Coupon, st.CouponID, st); err != nil {
		ForwardEx(session, "1", nil, "New Coupon DirectWrite :"+err.Error())
		return
	}
	err := AppendCouponToCompany(st.CouponID, st.CID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", st)
}

/*
修改代金券信息
*/
func ModifyCoupon(session *JsNet.StSession) {
	type modifyInfo struct {
		CouponID  string //代金券ID
		Picture   string //代金券路径
		Money     int    //代金券金额
		StartTime string //代金券起始日期
		StopTime  string //代金券结束日期
	}
	st := &modifyInfo{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ModifyCoupon GetPara :%s\n", err.Error())
		return
	}
	if st.CouponID == "" || st.StartTime == "" || st.StopTime == "" || st.Money <= 0 {
		ForwardEx(session, "1", nil,
			"ModifyCoupon para failed, CouponID:%s,StartTime:%s,StopTime:%s,Money:%d\n",
			st.CouponID, st.StartTime, st.StopTime, st.Money)
		return
	}

	start, err1 := ParseTimeFromString(st.StartTime, "2006-01-02")
	end, err2 := ParseTimeFromString(st.StopTime, "2006-01-02")
	if err1 != nil || err2 != nil {
		ForwardEx(session, "1", nil, "ModifyCoupon Time format error,such as '2006-01-02'\n")
		return
	}
	cur := CurStamp()
	if cur > end.Unix() || start.Unix() > end.Unix() {
		ForwardEx(session, "1", nil, "不能修改已经过去的时间\n")
		return
	}

	Coupon := &Coupon{}

	if err := db.WriteLock(constant.Hash_Coupon, st.CouponID, Coupon); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Coupon.StartTime = st.StartTime
	Coupon.StopTime = st.StopTime
	Coupon.StartStamp = start.Unix()
	Coupon.StopStamp = end.Unix()
	Coupon.Money = st.Money
	Coupon.Picture = st.Picture
	Coupon.Status = 0
	if err := db.WriteBack(constant.Hash_Coupon, st.CouponID, Coupon); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", Coupon)
}

/*
删除一个代金券
*/
func DelCoupon(session *JsNet.StSession) {
	type INFO struct {
		CouponID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Coupon := &Coupon{}
	if st.CouponID == "" {
		ForwardEx(session, "1", nil, "DelCoupon CouponID is empty\n")
		return
	}
	if err := db.WriteLock(constant.Hash_Coupon, st.CouponID, Coupon); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Coupon.Status = -1

	if err := db.WriteBack(constant.Hash_Coupon, st.CouponID, Coupon); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", nil)
}

/*
获取代金券信息内部
*/

func GetCouponInfo(CouponID string) (st *Coupon, e error) {
	data := &Coupon{}
	if CouponID == "" {
		return data, ErrorLog("GetCouponInfo failed,CouponID is empty\n")
	}
	if err := db.WriteLock(constant.Hash_Coupon, CouponID, data); err != nil {
		return data, err
	}

	cur := CurStamp()
	if data.Status != -1 && cur > data.StopStamp {
		data.Status = 1
	}
	err := db.WriteBack(constant.Hash_Coupon, CouponID, data)
	return data, err
}

/*
获取一个代金券信息
*/
func QueryCoupon(session *JsNet.StSession) {
	type INFO struct {
		CouponID string
	}
	data := &Coupon{}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CouponID == "" {
		ForwardEx(session, "1", nil, "QueryCoupon CouponID is empty\n")
		return
	}

	data, err := GetCouponInfo(st.CouponID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

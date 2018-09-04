package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"encoding/xml"
	"ider"
	"strconv"
	"strings"
	"time"
	. "util"
)

type ST_Transfer struct {
	Tid              string
	OpenId           string //收款人OpenId
	OrderId          string //订单ID
	UserName         string //收款人姓名
	UserHeader       string //收款人头像
	LDate            string //打款时间
	TimeStamp        int64  //时间戳
	Theme            string
	ProjectId        string
	Desc             string //描述
	MchId            string
	AppId            string
	CheckName        string
	ReUserName       string
	Amount           int //转账金额
	Spbill_create_ip string
	TransferCb       map[string]string
	LastError        string
	LastCb           map[string]string
}

////微信支付回调
type WxST_PayCb struct {
	AppId                string `xml:"appid"`
	Mch_id               string `xml:"mch_id"`
	Device_info          string `xml:"device_info"`
	Nonce_str            string `xml:"nonce_str"`
	Sign                 string `xml:"sign"`
	Sign_type            string `xml:"sign_type"`
	Result_code          string `xml:"result_code"`
	Err_code             string `xml:"err_code"`
	Err_code_des         string `xml:"err_code_des"`
	Openid               string `xml:"openid"`
	Is_subscribe         string `xml:"is_subscribe"`
	Trade_type           string `xml:"trade_type"`
	Bank_type            string `xml:"bank_type"`
	Total_fee            string `xml:"total_fee"`
	Settlement_total_fee string `xml:"settlement_total_fee"`
	Fee_type             string `xml:"fee_type"`
	Cash_fee             string `xml:"cash_fee"`
	Cash_fee_type        string `xml:"cash_fee_type"`
	Transaction_id       string `xml:"transaction_id"`
	Out_trade_no         string `xml:"out_trade_no"` //订单id
	Attach               string `xml:"attach"`
	Time_end             string `xml:"time_end"`
}

///支付信息
type ST_OrderApp struct {
	TerminalIp       string // 支付主机IP
	LocalTimeStamp   int64  // 本地时间戳
	ServiceTimeStamp int64  // 服务端时间戳
	Amount           int    // 需要支付金额
	RealPay          int    // 实际支付价格
	Desc             string // 描述
	Nonce_str        string // 随机串
	Mch_id           string // 商家ID
	AppId            string // 应用ID
	OpenId           string // opendid
	RefundId         string // 退款ID
	RefundFee        int    // 退款金额
}

//订单
type ST_Order struct {
	OrderID         string            //订单id
	OrderType       int               //1:红包 2：海报
	ActivityID      string            //活动id
	RequirementID   string            //需求id
	UID             string            //用户id
	OpenId          string            //用户openid
	UserName        string            //用户名字
	UserCell        string            //用户手机号
	UserCity        string            //用户城市
	PayWay          string            //支付途径:支付宝、微信、银联、平台余额（保留）
	PayAccount      int               //支付账号
	PayNumber       string            //第三方平台支付单号
	WxPayCb         *WxST_PayCb       //支付回调
	WxRefundCb      map[string]string //退款回调
	RefundMoney     int               //退款金额
	Charge          map[string]string //票据
	SubmitStamp     int64             //下单的时间戳
	OrderSubmitDate string            //下单时间
	PayDate         string            //支付时间
	RefundDate      string            //退款时间
	Status          string            //当前状态
	ST_OrderApp                       //支付信息
	EntityTime      string            //创建时间
}

func PaySuccess(session *JsNet.StSession) {
	body := session.Body()
	paycb := &WxST_PayCb{}
	e := xml.Unmarshal(body, paycb)
	xml := ""
	if e != nil {
		xml = `<xml>
  				<return_code><![CDATA[FAIL]]></return_code>
  				<return_msg><![CDATA[` + e.Error() + `]]></return_msg>
			   </xml>`
	} else {
		go orderPaySuccess(paycb)
		xml = `<xml>
  				<return_code><![CDATA[SUCCESS]]></return_code>
  				<return_msg><![CDATA[OK]]></return_msg>
			</xml>`
	}
	session.DirectWrite(xml)
}

///提交一个订单
func SubmitOrder(session *JsNet.StSession) {
	Error("Enter SubmitOrder....\n")
	order := &ST_Order{}
	if err := session.GetPara(order); err != nil {
		ForwardEx(session, "2", nil, err.Error())
		return
	}

	if order.UID == "" {
		ForwardEx(session, "1", nil, "SubmitOrder failed,UID is empty \n")
	}

	if order.OrderType == 1 && order.ActivityID == "" {
		ForwardEx(session, "1", nil, "SubmitOrder failed,ActivityID is empty\n")
		return
	}
	if order.OrderType == 2 && order.RequirementID == "" {
		ForwardEx(session, "1", nil, "SubmitOrder failed,RequirementID is empty\n")
		return
	}
	order.OrderID = ider.GenOrderId() //订单id
	ch, err := wx_pub_pay(order)
	if err != nil {
		ForwardEx(session, "4", nil, err.Error())
		return
	}

	if user, err := GetUserInfo(order.UID); err == nil {
		order.UserCell = user.Mobile
		order.UserName = user.Name
		order.UserCity = user.City
	}

	order.Charge = ch //收据
	addr := session.RemoteAddr()
	i := strings.Index(addr, ":")
	order.TerminalIp = addr[:i]
	order.ServiceTimeStamp = time.Now().Unix() + constant.C_TIMEAREA*3600
	order.OrderSubmitDate = CurTime()     //提交时间
	order.SubmitStamp = CurStamp()        //提交时间的时间戳
	order.Status = constant.C_PAY_WAITPAY //状态
	order.EntityTime = CurTime()
	if err := db.DirectWrite(constant.Hash_Order, order.OrderID, order); err != nil {
		ForwardEx(session, "4", nil, err.Error())
		return
	}

	if order.OrderType == 1 {
		go appendOrderID2Activity(order.ActivityID, order.OrderID)
	} else if order.OrderType == 2 {
		go appendOrderID2Requirement(order.RequirementID, order.OrderID)
	}

	go AppendUserOrder(order.OrderID, order.UID)
	go AppendToOrderGlobal(order.OrderID)

	Error("order=%v\n", order)
	Forward(session, "0", order)
}

//支付完成后
func orderPaySuccess(cb *WxST_PayCb) error {
	Info("enter orderPaySuccess.....\n")
	order := &ST_Order{}
	if err := db.WriteLock(constant.Hash_Order, cb.Out_trade_no, order); err != nil {
		return err
	}
	order.WxPayCb = cb //支付回调
	m, e := strconv.Atoi(cb.Cash_fee)
	if e == nil {
		order.RealPay = m ///支付价格
	} else {
		order.RealPay = 0 ///支付价格
	}
	order.RefundFee = order.RealPay       ///可退金钱
	order.PayNumber = cb.Out_trade_no     //第三方支付的单号
	order.PayDate = CurTime()             //支付日期
	order.Status = constant.C_PAY_SUCCESS //支付状态
	if err := db.WriteBack(constant.Hash_Order, cb.Out_trade_no, order); err != nil {
		return err
	}
	if order.OrderType == 1 {
		///更新活动状态
		go PayActivity(order.ActivityID)
		go PromoteToSeller(order.UID)
		Info("UID=%s支付活动费用成功,Money=%d\n", order.UID, order.RealPay)
	} else if order.OrderType == 2 {
		///更新需求海报
		go PromoteToSeller(order.UID)
		go PayRequirement(order.RequirementID)
		Info("UID=%s支付海报费用成功,Money=%d\n", order.UID, order.RealPay)
	}
	Info("Leave orderPaySuccess.....\n")
	return nil
}

//查询订单
func QueryOrder(session *JsNet.StSession) {
	type INFO struct {
		OrderID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	order, err := GetOrderInfo(st.OrderID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", order)
}

//获取用户的所有订单信息
func GetUserOrderList(session *JsNet.StSession) {
	type INFO struct {
		UID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", getmoreorder(getuserorderlist(st.UID)))
}

//获取多个订单
func GetMoreOrders(session *JsNet.StSession) {
	type INFO struct {
		List []string //订单id列表
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", getmoreorder(st.List))
}

//获取全局的订单信息
func GetGlobalOrders(session *JsNet.StSession) {
	Forward(session, "0", getmoreorder(getglobalorderlist()))
}

//添加一个订单id到用户
func AppendUserOrder(OrderId, Uuid string) error {
	if OrderId == "" || Uuid == "" {
		return ErrorLog("AppendUserOrder failed,OrderId:%s,Uuid:%s\n", OrderId, Uuid)
	}
	data := &[]string{}
	err := db.WriteLock(constant.Hash_UserOrder, Uuid, data)
	AppendUniqueString(data, OrderId)
	if err != nil {
		return db.DirectWrite(constant.Hash_UserOrder, Uuid, data)
	}
	return db.WriteBack(constant.Hash_UserOrder, Uuid, data)
}

///从用户列表中移除订单
func removeFormUserOrder(OrderID, UID string) error {
	if OrderID == "" || UID == "" {
		return ErrorLog("removeFormUserOrder failed,OrderID:%s,UID:%s\n", OrderID, UID)
	}
	data := &[]string{}
	err := db.WriteLock(constant.Hash_UserOrder, UID, data)

	if err != nil {
		return nil
	}
	DelExistString(data, OrderID)
	return db.WriteBack(constant.Hash_UserOrder, UID, data)
}

//添加一个订单到全局
func AppendToOrderGlobal(OrderId string) error {
	if OrderId == "" {
		return ErrorLog("AppendToOrderGlobal failed,OrderId:%s,\n", OrderId)
	}
	data := &[]string{}
	err := db.WriteLock(constant.Hash_Order, constant.KEY_Global_Order, data)
	AppendUniqueString(data, OrderId)
	if err != nil {
		return db.DirectWrite(constant.Hash_Order, constant.KEY_Global_Order, data)
	}
	return db.WriteBack(constant.Hash_Order, constant.KEY_Global_Order, data)
}

//从全局的订单列表中移除
func removeFromGlobalOrder(OrderID string) error {
	if OrderID == "" {
		return ErrorLog("removeFromGlobalOrder failed,OrderID:%s,\n", OrderID)
	}
	data := &[]string{}
	if err := db.WriteLock(constant.Hash_Order, constant.KEY_Global_Order, data); err != nil {
		return nil
	}
	DelExistString(data, OrderID)
	return db.WriteBack(constant.Hash_Order, constant.KEY_Global_Order, data)
}

///将超时的订单扔到无效的订单列表
func Append2Invalid(orderID string) error {
	if orderID == "" {
		return ErrorLog("Append2Invalid failed,orderID:%s,\n", orderID)
	}
	data := &[]string{}
	err := db.WriteLock(constant.Hash_Order, constant.KEY_Invalid_Order, data)
	AppendUniqueString(data, orderID)
	if err != nil {
		return db.DirectWrite(constant.Hash_Order, constant.KEY_Invalid_Order, data)
	}
	return db.WriteBack(constant.Hash_Order, constant.KEY_Invalid_Order, data)
}

//获取用户的所有的订单
func getuserorderlist(Uuid string) []string {
	data := []string{}
	db.ShareLock(constant.Hash_UserOrder, Uuid, &data)
	return data
}

//获取全局的订单
func getglobalorderlist() []string {
	data := []string{}
	db.ShareLock(constant.Hash_Order, constant.KEY_Global_Order, &data)
	return data
}

//获取多个订单信息
func getmoreorder(list []string) []*ST_Order {
	data := []*ST_Order{}
	if len(list) == 0 {
		return data
	}
	for _, v := range list {
		order, err := GetOrderInfo(v)
		if err != nil {
			continue
		}
		if order.Status == constant.C_PAY_INVALID {
			continue
		}
		if order.Status == constant.C_PAY_WAITPAY && time.Now().Unix()-order.SubmitStamp >= 24*3600 {
			continue
		}
		data = append(data, order)
	}
	return data
}

///获取订单信息
func GetOrderInfo(orderId string) (*ST_Order, error) {
	data := &ST_Order{}
	if orderId == "" {
		return data, ErrorLog("GetOrderInfo failed,OrderID=%s\n", orderId)
	}
	err := db.ShareLock(constant.Hash_Order, orderId, data)
	return data, err
}

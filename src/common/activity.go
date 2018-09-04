package common

import (
	. "JsLib/JsConfig"
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	. "util"
)

type LATLON struct {
	Lat   float64 //纬度
	Lon   float64 //经度
	Range float64 //领取范围公里
	Addr  string  //地址
}

// Optional1
type QuestionInfo struct {
	Question   string       //问题
	Answer     string       //正确答案
	AnswerList []AnswerInfo //答案列表
}

type Activity struct {
	Theme             string           //主题
	Html              string           //活动详情Html
	URL               string           //模板路径
	EntityTime        string           //日期
	City              string           //城市
	UID               string           //用户id
	ActivityID        string           //活动id
	RequirementID     string           //需求id
	OrderID           string           //订单ID
	TestFlag          int              //是否是测试数据 -1:测试数据 1:正式数据
	Type              string           //类型
	Name              string           //姓名
	Mobile            string           //手机号
	TemplateId        string           //模板id(备用)
	SellerId          string           //商家id(备用)
	CID               string           //发布公司ID
	AgentID           string           //代理ID
	CreateEntity      int              //0--company ;1 --agent
	Logo              string           //商家logo
	AgentId           string           //小B ID(备用)
	ProxyId           string           //代理id(备用)
	Pictures          []string         //海报
	Content           string           //内容
	Note              string           //备注
	QuestionId        []string         //问题列表
	SystemQuestionId  []string         //系统问题列表
	Lat               float64          //纬度
	Lon               float64          //经度
	Range             float64          //领取范围公里
	Addr              string           //地址
	LatLon            []LATLON         //坐标列表
	Area              []string         //城市列表
	TardingArea       []string         //商圈id
	TotalMoney        int              //总金额
	TotalCustomer     int              //分配人数
	IsAuto            int              //是否是自动模式
	AlgoType          int              //算法类型
	Probability       int              //红包概率
	MoneyArrangeWay   string           //分配方式("Coordinate","Division","TradingArea")
	Remaining         int              //剩余金额
	ShardMoney        int              //分享的给出的佣金
	LastNum           int              //剩余数量
	VisitNum          int              //访问量
	FinishedPercent   int              //完成多少了%"
	ArrangedCustomer  int              //已经抢红包人数,已经抢到钱的人
	RobNum            int              //抢红包的人，不管有没有抢到钱
	Status            string           //状态
	ReleaseDate       string           //发布时间
	SystemQuestionNum int              //系统问题数量
	CouponID          string           //代金券ID
	IsCoupon          bool             //是否使用代金券
	IsBigSeller       bool             //是否是大商家
	IsRealTime        int              //是否实时到账
	LsDistribution    []DistributionSI //多网点部署的网点简短信息
	ArrangeType       string           //活动发布形式 公司独立发布；网点发布 1.Company 2.Distribution
	QuestionTemp      []QuestionInfo   //问题列表
	LsASystemQuestion []string         //系统问题
}

//新建活动
func NewActivity(require *Requirement) (*Activity, error) {
	if require == nil {
		return nil, ErrorLog("NewActivity failed,require is nil \n")
	}

	ac := &Activity{
		ActivityID:    ider.GenID(),
		UID:           require.UID,
		RequirementID: require.RequirementID,
		Theme:         require.Theme,
		Logo:          require.Logo,
		Type:          require.Type,
		TemplateId:    require.TemplateId,
		SellerId:      require.SellerId,
		AgentId:       require.AgentId,
		ProxyId:       require.ProxyId,
		Content:       require.Content,
		Pictures:      require.Pictures,
		Lat:           require.Lat,
		Lon:           require.Lon,
		Addr:          require.Addr,
		Name:          require.Name,
		Mobile:        require.Mobile,
		IsAuto:        1,
		Probability:   0,
		AlgoType:      0,
		IsBigSeller:   false,
		Status:        constant.Status_WaitPay,
		EntityTime:    CurTime(),
	}

	if err := db.DirectWrite(constant.Hash_Activity, ac.ActivityID, ac); err != nil {
		return nil, err
	}

	go appendActivityToGolbal(ac.ActivityID)
	go AddUserSendRDP(require.UID, ac.ActivityID, 0)

	return ac, nil
}

//新建活动
func NewActivityDirectly(session *JsNet.StSession) {
	st := &Activity{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	ok, errstr := activityEvaluation(st)
	if !ok {
		ForwardEx(session, "1", nil, errstr)
		return
	}
	//to generate question
	for _, v := range st.QuestionTemp {
		q, e := CreatQuetion(v.Question, v.Answer, v.AnswerList)
		if e == nil {
			st.QuestionId = append(st.QuestionId, q.QID)
		}
	}
	st.ActivityID = ider.GenID()
	st.EntityTime = CurTime()
	st.IsBigSeller = true
	if st.TestFlag == 1 {
		st.Status = constant.Status_WaitPay
	} else if st.TestFlag == -1 {
		st.Status = constant.Status_Active
	}
	st.ReleaseDate = CurTime()
	st.IsAuto = 1
	for i, v := range st.LsDistribution {
		st.LsDistribution[i].Remaining = v.Money
	}

	if err := db.DirectWrite(constant.Hash_Activity, st.ActivityID, st); err != nil {
		ForwardEx(session, "1", nil, "New Activity error :"+err.Error())
		return
	}

	if err := NewActivityRDP(st); err != nil {
		ForwardEx(session, "1", st, err.Error())
		return
	}

	go appendActivityToGolbal(st.ActivityID)
	go AppendActivityToCompany(st.ActivityID, st.CID)
	if st.ArrangeType == constant.Arrange_Distribution {
		for _, v := range st.LsDistribution {
			AppendActivityToDistribution(st.ActivityID, v.DID)
		}
	}
	//go AddUserSendRDP(st.CID, st.ActivityID, 0)
	Forward(session, "0", st)
}

//标记删除一个活动
func DelActivity(session *JsNet.StSession) {
	type info struct {
		ActivityID string
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.Status = constant.Status_Del
	if err := db.WriteBack(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

//活动参数判断
func activityEvaluation(st *Activity) (ok bool, e string) {
	if st == nil {
		return false, "The activity is none"
	}
	if len(st.Pictures) == 0 && st.URL == "" {
		return false, "The picture length is zero"
	}
	if st.TotalCustomer <= 0 || st.TotalMoney <= 0 || st.TestFlag == 0 || st.AlgoType < 0 {
		return false, ErrorLog("ActivitySettings failed,param is empty,ActivityID=%s,TotalCustomer=%d,TotalMoney=%d,TestFlag=%f,AlgoType:%d\n",
			st.ActivityID, st.TotalCustomer, st.TotalMoney, st.TestFlag, st.AlgoType).Error()
	}
	if st.TotalMoney < st.TotalCustomer {
		return false, ErrorLog("Money:%s is less than Person:%\n", st.TotalMoney, st.TotalCustomer).Error()
	}

	if (st.MoneyArrangeWay != constant.MoneyArrange_Coordinate && st.MoneyArrangeWay != constant.MoneyArrange_Division) ||
		(st.ArrangeType != constant.Arrange_CompanyAlone && st.ArrangeType != constant.Arrange_Distribution) {
		return false, ErrorLog("MoneyArrangeWay,Arrange Way failed,MoneyArrangeWay=%s,ArrangeWay:%s",
			st.MoneyArrangeWay, st.ArrangeType).Error()
	}

	if st.MoneyArrangeWay == constant.MoneyArrange_Coordinate {
		if (st.ArrangeType == constant.Arrange_CompanyAlone && len(st.LatLon) == 0) ||
			(st.ArrangeType == constant.Arrange_Distribution && len(st.LsDistribution) == 0) {
			return false, ErrorLog("MoneyArrangeWay,Arrange Way failed,MoneyArrangeWay=%s,ArrangeWay:%s,len(LatLon)=%d,len(LsDistribution)=%d",
				st.MoneyArrangeWay, st.ArrangeType, len(st.LatLon), len(st.LsDistribution)).Error()
		}
	}

	if st.MoneyArrangeWay == constant.MoneyArrange_Division {
		if (st.ArrangeType == constant.Arrange_CompanyAlone && len(st.Area) == 0) ||
			(st.ArrangeType == constant.Arrange_Distribution && len(st.LsDistribution) == 0) {
			return false, ErrorLog("MoneyArrangeWay,Arrange Way failed,MoneyArrangeWay=%s,ArrangeWay:%s,len(LatLon)=%d,len(LsDistribution)=%d",
				st.MoneyArrangeWay, st.ArrangeType, len(st.Area), len(st.LsDistribution)).Error()
		}
	}

	if st.MoneyArrangeWay == constant.MoneyArrange_TradArea && len(st.TardingArea) == 0 {
		return false, ErrorLog("ActivitySettings failed,MoneyArrangeWay=%s,TradArea=%v\n", st.MoneyArrangeWay, st.TardingArea).Error()
	}

	if len(st.QuestionTemp) == 0 {
		return false, "ActivitySettings Question is empty"
	}
	return true, "Sucess"
}

//活动参数设置
func ActivitySettings(session *JsNet.StSession) {
	type Settings struct {
		ActivityID        string
		TotalCustomer     int              //分配人数
		TotalMoney        int              //总金额
		City              string           //城市
		TestFlag          int              //是否是测试数据 -1:测试数据 1:正式数据
		MoneyArrangeWay   string           //分配方式("Coordinate,Division")
		LatLon            []LATLON         //坐标
		Area              []string         //城市列表
		TardingArea       []string         //商圈id列表
		ReleaseDate       string           //发布时间
		QuestionId        []string         //问题id列表
		SystemQuestionId  []string         //系统问题id列表
		Status            string           //状态
		Que               []QuestionInfo   //问题列表
		AlgoType          int              //算法类型
		SystemQuestionNum int              //系统问题数量
		CouponID          string           //代金券ID
		LsDistribution    []DistributionSI //多网点部署的网点简短信息
		ArrangeType       string           //活动发布形式 公司独立发布；网点发布 1.Company 2.Distribution
		IsRealTime        int              //是否实时到账
	}

	st := &Settings{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ActivitySettings GetPara:%s\n", err.Error())
		return
	}
	if st.ActivityID == "" || st.TotalCustomer <= 0 || st.TotalMoney <= 0 || st.TestFlag == 0 || st.AlgoType < 0 {
		ForwardEx(session, "1", nil,
			"ActivitySettings failed,param is empty,ActivityID=%s,TotalCustomer=%d,TotalMoney=%d,Lat=%f,Lon=%f,Range=%f,TestFlag=%s,AlgoType:%d\n",
			st.ActivityID, st.TotalCustomer, st.TotalMoney, st.TestFlag, st.AlgoType)
		return
	}

	if st.TotalMoney < st.TotalCustomer {
		ForwardEx(session, "1", nil, "Money:%d is less than Person:%d\n", st.TotalMoney, st.TotalCustomer)
		return
	}

	if (st.MoneyArrangeWay != constant.MoneyArrange_Coordinate && st.MoneyArrangeWay != constant.MoneyArrange_Division) ||
		(st.ArrangeType != constant.Arrange_CompanyAlone && st.ArrangeType != constant.Arrange_Distribution) {
		ForwardEx(session, "1", nil, "MoneyArrangeWay,Arrange Way failed,MoneyArrangeWay=%s,ArrangeType:%s",
			st.MoneyArrangeWay, st.ArrangeType)
		return
	}

	if st.MoneyArrangeWay == constant.MoneyArrange_Coordinate {
		if (st.ArrangeType == constant.Arrange_CompanyAlone && len(st.LatLon) == 0) ||
			(st.ArrangeType == constant.Arrange_Distribution && len(st.LsDistribution) == 0) {
			ForwardEx(session, "1", nil, "MoneyArrangeWay,ArrangeType failed,MoneyArrangeWay=%s,ArrangeWay:%s,len(LatLon)=%d,len(LsDistribution)=%d",
				st.MoneyArrangeWay, st.ArrangeType, len(st.LatLon), len(st.LsDistribution))
			return
		}
	}

	if st.MoneyArrangeWay == constant.MoneyArrange_Division {
		if (st.ArrangeType == constant.Arrange_CompanyAlone && len(st.Area) == 0) ||
			(st.ArrangeType == constant.Arrange_Distribution && len(st.LsDistribution) == 0) {
			ForwardEx(session, "1", nil, "MoneyArrangeWay,ArrangeType failed,MoneyArrangeWay=%s,ArrangeWay:%s,len(LatLon)=%d,len(LsDistribution)=%d",
				st.MoneyArrangeWay, st.ArrangeType, len(st.Area), len(st.LsDistribution))
			return
		}
	}

	if st.MoneyArrangeWay == constant.MoneyArrange_TradArea && len(st.TardingArea) == 0 {
		ForwardEx(session, "1", nil, "ActivitySettings failed,MoneyArrangeWay=%s,TradArea=%v\n", st.MoneyArrangeWay, st.TardingArea)
		return
	}

	st.ReleaseDate = CurTime()

	if len(st.Que) == 0 {
		Error("ActivitySettings Question is empty...............................\n")

	}

	st.QuestionId = []string{}
	for _, v := range st.Que {
		q, e := CreatQuetion(v.Question, v.Answer, v.AnswerList)
		if e == nil {
			st.QuestionId = append(st.QuestionId, q.QID)
		}
	}

	if st.TestFlag == 1 {
		st.Status = constant.Status_WaitPay
	} else if st.TestFlag == -1 {
		st.Status = constant.Status_Active
	}

	data := &Activity{}

	if err := db.WriteLock(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.TotalCustomer = st.TotalCustomer
	data.TotalMoney = st.TotalMoney
	data.City = st.City
	data.TestFlag = st.TestFlag
	data.AlgoType = st.AlgoType
	data.MoneyArrangeWay = st.MoneyArrangeWay
	data.LatLon = st.LatLon
	data.Area = st.Area
	data.ReleaseDate = CurTime()
	data.QuestionId = st.QuestionId
	data.QuestionTemp = st.Que
	data.SystemQuestionId = st.SystemQuestionId
	data.Status = st.Status
	data.SystemQuestionNum = st.SystemQuestionNum
	data.CouponID = st.CouponID
	data.LsDistribution = st.LsDistribution
	data.ArrangeType = st.ArrangeType
	data.CID = data.UID
	data.IsRealTime = st.IsRealTime
	data.IsAuto = 1
	data.TardingArea = st.TardingArea
	if err := db.WriteBack(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if err := NewActivityRDP(data); err != nil {
		ForwardEx(session, "1", data, err.Error())
		return
	}
	Forward(session, "0", data)
}

///将订单ID和活动ID关联
func appendOrderID2Activity(ActivityID, OrderID string) error {
	if ActivityID == "" || OrderID == "" {
		return ErrorLog("appendOrderID2Activity failed,ActivityID: %s,OrderID=%s\n", ActivityID, OrderID)
	}
	data := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, ActivityID, data); err != nil {
		return err
	}
	data.OrderID = OrderID
	if err := db.WriteBack(constant.Hash_Activity, ActivityID, data); err != nil {
		return err
	}
	return nil
}

//支付成功后回调
func PayActivity(ActivityID string) error {
	data := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, ActivityID, data); err != nil {
		return err
	}
	data.Status = constant.Status_Active
	if err := db.WriteBack(constant.Hash_Activity, ActivityID, data); err != nil {
		return err
	}
	///更新用户的活动相关
	go AddUserSendRDP(data.UID, ActivityID, data.TotalMoney)
	return nil
}

//查询单个活动信息
func QueryActivity(session *JsNet.StSession) {
	type INFO struct {
		ActivityID string //活动id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ActivityID == "" {
		ForwardEx(session, "1", nil, "QueryActivity ActivityID is empty\n")
		return
	}
	ac, err := GetActivityInfo(st.ActivityID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", ac)
}

//修改活动的海报
func ModifyActivity(session *JsNet.StSession) {
	type INFO struct {
		ActivityID string   //活动id
		Pictures   []string //海报
		Theme      string   //主题
		Name       string   //姓名
		Mobile     string   //手机号
		Content    string   //内容
		URL        string   //活动路劲
	}

	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ModifyActivity GetPara:%s\n", err.Error())
		return
	}

	if st.ActivityID == "" ||
		(st.URL == "" && len(st.Pictures) == 0) ||
		st.Theme == "" || st.Name == "" ||
		st.Mobile == "" || st.Content == "" {
		ForwardEx(session, "1", nil,
			"ModifyActivity param failed, ActivityID=%s,Theme=%s,Content=%s,Name=%s,Mobile=%s,Pictures=%v,URL=%s\n",
			st.ActivityID, st.Theme, st.Content, st.Name, st.Mobile, st.Pictures, st.URL)
		return
	}
	re := &Activity{}

	if err := db.WriteLock(constant.Hash_Activity, st.ActivityID, re); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	re.Pictures = st.Pictures
	re.URL = st.URL
	re.Theme = st.Theme
	re.Content = st.Content
	re.Name = st.Name
	re.Mobile = st.Mobile

	if err := db.WriteBack(constant.Hash_Activity, st.ActivityID, re); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", re)

}

////删除全局的某一个活动（后台强制停止）
func DelGlobalActivity(session *JsNet.StSession) {
	type INFO struct {
		ActivityID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.ActivityID == "" {
		ForwardEx(session, "1", nil, "DelGlobalActivity ActivityID is empty\n")
		return
	}
	if err := delFromGolbalActivity(st.ActivityID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	go stopActivity(st.ActivityID)
	Forward(session, "0", nil)
}

//强制停止活动
func stopActivity(ActivityID string) error {
	if ActivityID == "" {
		return ErrorLog("stopActivity failed,ActivityID is empty \n")
	}
	data := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, ActivityID, data); err != nil {
		return err
	}
	data.Status = constant.Status_Stop
	return db.WriteBack(constant.Hash_Activity, ActivityID, data)
}

//获取活动信息
func GetActivityInfo(ActivityID string) (st *Activity, e error) {

	data := &Activity{}
	if ActivityID == "" {
		return data, ErrorLog("GetActivityInfo failed,ActivityID is empty\n")
	}
	if err := db.ShareLock(constant.Hash_Activity, ActivityID, data); err != nil {
		return data, err
	}
	///获取红包领取记录
	r, err := getActivityRDP(ActivityID)

	if err == nil {
		data.Remaining = r.Remaining
		data.LastNum = r.LastNum
		data.ArrangedCustomer = r.ArrangedCustomer
		data.RobNum = r.RobNum
		data.FinishedPercent = r.FinishedPercent
		data.ShardMoney = r.SharedMoney
	} else {
		if data.Status != constant.Status_WaitPay {
			ErrorLog("@@@@@@@@@@@@@@@@@@@@@@@@@@@@ GetActivityRDP failed,ActivityID=%s@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n", ActivityID)
		}
	}

	return data, nil
}

/*
 获取所有的需求列表
*/
func GetAllActivities(session *JsNet.StSession) {
	Forward(session, "0", getAllActivitys(true))
}

//获取首页活动列表
func GetHomeActivity(session *JsNet.StSession) {
	Forward(session, "0", getAllActivitys(false))
}

func getAllActivitys(isAll bool) []*Activity {
	ids := getAllActivityID()
	data := []*Activity{}
	if len(ids) == 0 {
		return data
	}
	for _, id := range ids {
		d, e := GetActivityInfo(id)
		if e == nil {
			if d.Status == constant.Status_Del {
				continue
			}
			if !isAll {
				if d.Status != constant.Status_Active {
					continue
				}
			}
			data = append(data, d)
		}
	}

	return data
}

func getAllActivityID() []string {
	data := []string{}
	err := db.ShareLock(constant.Hash_Activity, constant.KEY_Global_Activity, &data)
	if err != nil {
		Error("getAllActivityID failed !\n")
	}
	return data
}

//添加全局的活动id列表
func appendActivityToGolbal(id string) error {

	data := []string{}
	err := db.WriteLock(constant.Hash_Activity, constant.KEY_Global_Activity, &data)
	AppendUniqueString(&data, id)
	if err != nil {
		return db.DirectWrite(constant.Hash_Activity, constant.KEY_Global_Activity, &data)
	}
	return db.WriteBack(constant.Hash_Activity, constant.KEY_Global_Activity, &data)
}

func delFromGolbalActivity(id string) error {
	data := &[]string{}
	if err := db.WriteLock(constant.Hash_Activity, constant.KEY_Global_Activity, data); err != nil {
		ErrorLog(err.Error())
		return nil
	}
	DelExistString(data, id)
	return db.WriteBack(constant.Hash_Activity, constant.KEY_Global_Activity, data)
}

///发布活动
func ReleaseActivity(session *JsNet.StSession) {
	type INFO struct {
		RequirementID string //需求id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ReleaseActivity GetPara:%s\n", err.Error())
		return
	}
	re := &Requirement{}
	if err := db.WriteLock(constant.Hash_Requirement, st.RequirementID, re); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if re.Status == constant.Status_WaitPay {
		go db.WriteBack(constant.Hash_Requirement, st.RequirementID, re)
		ForwardEx(session, "1", nil, "ReleaseActivity failed,this Requirement:%s Wait Pay\n", st.RequirementID)
		return
	}

	if re.Status == constant.Status_Finished {
		go db.WriteBack(constant.Hash_Requirement, st.RequirementID, re)
		ForwardEx(session, "1", nil, "ReleaseActivity failed,this Requirement:%s is finish\n", st.RequirementID)
		return
	}

	if len(re.Pictures) == 0 {
		ForwardEx(session, "1", nil, "ReleaseActivity failed,Picturs is empty\n")
		return
	}

	re.Status = constant.Status_Finished //审核通过

	ac, e1 := NewActivity(re)
	if e1 == nil {
		re.ActivityID = ac.ActivityID
	}
	go db.WriteBack(constant.Hash_Requirement, st.RequirementID, re)
	if e1 != nil {
		ForwardEx(session, "1", nil, e1.Error())
		return
	}
	Forward(session, "0", ac)
}

///向用户发放红包
func SendActivityRDP(user *User, Money int, Tag string) bool {
	OpendID := ""
	appid := ""
	if Tag == "small" {
		OpendID = user.OpenId_Small
		appid = "wx11cdac22d7719783"
	} else if Tag == "wechat" {
		OpendID = user.OpenId
		appid = CFG.DirectPay.AppId
	} else {
		OpendID = user.OpenId_App
		appid = "wx995d7fad8a74a7ef"
	}
	if OpendID == "" {
		ErrorLog("SendActivityRDP failed,Openid is Empty,Tag=%s", Tag)
		return false
	}
	st := &ST_Transfer{
		OpenId:     OpendID,
		UserName:   user.Name,
		UserHeader: user.HeadImageURL,
		TimeStamp:  CurStamp(),
		LDate:      CurTime(),
		Desc:       "传单侠奖励红包",
		Amount:     Money,
		AppId:      appid,
	}
	_, err := direct_transfer(st)
	return err == nil
}

///更改活动
func ChangeActivitySettings(session *JsNet.StSession) {
	type INFO struct {
		ActivityID    string //活动id
		TotalCustomer int    //总人数
		TotalMoney    int    //总金额
		AlgoType      int    //红包
		Probability   int    //红包概率
		IsAuto        int    //是否是自动模式
		IsRealTime    int    //是否实时到账
		Token         string //校验token
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if ok := checkToken(session.RemoteAddr(), st.Token); !ok {
		ForwardEx(session, "1", nil, "ChangeActivitySettings,Token 校验失败\n")
		return
	}

	if st.TotalCustomer <= 0 || st.TotalMoney <= 0 || st.AlgoType < 0 || st.IsAuto == 0 {
		ForwardEx(session, "1", nil,
			"ChangeActivitySettings failed ,TotalCustomer:%d,TotalMoney:%d,AlgoType:%d,Probability:%d,IsAuto:%d\n",
			st.TotalCustomer, st.TotalMoney, st.AlgoType, st.IsAuto)
		return
	}
	if st.Probability < 0 || st.Probability > 100 {
		ForwardEx(session, "1", nil, "ChangeActivitySettings failed,Probability:%d\n", st.Probability)
		return
	}

	data := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	data.TotalCustomer = st.TotalCustomer
	data.TotalMoney = st.TotalMoney
	data.AlgoType = st.AlgoType
	data.Probability = st.Probability
	data.IsAuto = st.IsAuto
	data.IsRealTime = st.IsRealTime
	if err := db.WriteBack(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	rdp := &RedPacketGetRecord{}
	if err := db.WriteLock(constant.Hash_ActivityRDP, st.ActivityID, &rdp); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	disMoney := st.TotalMoney - rdp.TotalMoney

	rdp.TotalCustomer = st.TotalCustomer
	rdp.TotalMoney = st.TotalMoney
	rdp.AlgoType = st.AlgoType
	rdp.Probability = st.Probability
	rdp.IsAuto = st.IsAuto
	rdp.IsRealTime = st.IsRealTime
	rdp.Remaining = rdp.Remaining + disMoney
	rdp.LastNum = rdp.TotalCustomer - rdp.ArrangedCustomer
	if rdp.Remaining <= 0 {
		rdp.Remaining = 0
	}
	if rdp.LastNum <= 0 {
		rdp.LastNum = 0
	}

	data.FinishedPercent = rdp.RobNum * 100 / rdp.TotalCustomer
	if err := db.WriteBack(constant.Hash_ActivityRDP, st.ActivityID, &rdp); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", nil)
}

func ChangeActivityStatus(session *JsNet.StSession) {

	type INFO struct {
		ActivityID string //活动id
		Status     string //状态
		Token      string //校验token
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if ok := checkToken(session.RemoteAddr(), st.Token); !ok {
		ForwardEx(session, "1", nil, "ChangeActivitySettings,Token 校验失败\n")
		return
	}

	if st.ActivityID == "" {
		ForwardEx(session, "1", nil, "ChangeActivitySettings,TActivityID is empty\n")
		return
	}
	if st.Status != "WaiPay" && st.Status != "Active" && st.Status != "Arrears" && st.Status != "Stop" {
		ForwardEx(session, "1", nil, "ChangeActivitySettings,Status:%s is not exist\n", st.Status)
		return
	}

	data := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.Status = st.Status
	if err := db.WriteBack(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

///增加活动的访问量
func AddActivityVisit(session *JsNet.StSession) {
	type INFO struct {
		ActivityID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data := &Activity{}
	if err := db.WriteLock(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data.VisitNum++
	if err := db.WriteBack(constant.Hash_Activity, st.ActivityID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}



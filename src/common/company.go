package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	."util"
)

type Company struct {
	ShortName         string   //公司简称
	FullName          string   //公司全称
	CID               string   //公司ID
	UserName          string   //登陆用户名
	Password          string   //登陆密码
	UID               string   //对应微信端的UID
	openID            string   //对应微信断的OpenID
	OfficialPerson    string   //公司法人代表姓名
	RegistrationTime  string   //公司注册时间
	TaxpayerIdentity  string   //公司纳税人识别号
	RegisteredCapital string   //公司注册资金
	Companysize       string   //公司人员规模
	SubCompanyNum     int      //子公司的个数
	BusinessScope     string   //公司业务范围
	Profile           string   //公司简介
	OtherNotes        string   //公司其他补充
	LsPicture         []string //公司介绍照片
	WebSite           string   //公司网址
	Representative    string   //公司代表人姓名
	Age               int      //公司代表人年龄
	Sex               string   //公司代表人性别
	LandlineTel       string   //座机电话（021-385638827-8921）
	MobilePhone       string   //代表人电话
	City              string   //城市
	Province          string   //省份
	Country           string   //国家
	Addr              string   //地址
	Type              string   //商家类型
	HeadImageURL      string   //公司代表人头像路径
	Logo              string   //公司logo路径
	AccountID         string   //公司账户ID
	Lon               float64  //公司经度
	Lat               float64  //公司纬度
	LsActivity        []string //公司下面所有的活动
	LsCoupon          []string //代金券ID列表
	LsASystemQuestion []string //公司的问卷调查
	EntityTime        string   //公司创建时间
	Status            string   //公司状态  0：正常  -1：删除 1：停业
	LsAgent           []AgentShortInfo  //代理ID
	AgentID            string   //代理ID
	IsAgentLaunched    bool
}

type AgentShortInfo struct{
	AgentID string
	EntityTime string
}

//公司和自己的子公司和代理网店的结构
type CompanyDistribution struct {
	CID            string   //公司ID
	DistributionID []string //网点ID
}

//客户activity的信息
type LActivity struct {
	ActivityID   string  //活动ID
	ActivityName string  //活动名字
	Money        int     //在这次活动中领取的金钱
	IsCoupon     bool    //是否领券
	LatF         float64 //领取纬度
	LonF         float64 //领取经度
	Address      string  //用户带上来的地址
	Date         string  //领取时间
	LTime        int     //用户从进入到领红包的时间
	Gowil        int     //购买意愿
	IsCorrect    bool    //是否回答正确
	InRange      bool    //是否在范围内
}

//公司客户简短的信息
type CSInfo struct {
	UID          string      //最终客户的UID
	Name         string      //客户姓名
	Mobile       string      //手机，联系方式
	DID          string      //所属网点id
	TotleMoney   int         //领到的总金额
	LsLActivity  []LActivity //参加过的活动
	LatF         float64     //用户纬度
	LonF         float64     //用户经度
	Address      string      //用户带上来的地址
	CreatDate    string      //用户增加的时间
	HeadImageURL string      //头像路径
}
type ST_Customer struct {
	ID          string         //公司ID
	LsCSI       []CSInfo       //公司客户信息列表
	Increase    map[string]int //人数的增加量统计
	IncreaseKey []string       ///增长的日期列表
}

type CompanyAccount struct {
	FullName  string //公司全称
	CID       string //公司ID
	UserName  string //登陆用户名
	Password  string //登陆密码
	CreatDate string //创建日期
}

/*
 注册一个公司账号
*/
func NewCompany(session *JsNet.StSession) {
	st := &Company{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.FullName == "" || st.UserName == "" || st.Password == "" {
		ForwardEx(session, "1", nil, "New Company failed with FullName=%s and UserName=%s,Password=%s\n", st.FullName, st.UserName, st.Password)
		return
	}
	if CheckCompanyAccount(st.UserName) {
		ForwardEx(session, "1", nil, " 用户名已经存在,UserName=%s", st.UserName)
		return
	}

	st.CID = ider.GenID()
	st.EntityTime = CurTime()
	if err := db.DirectWrite(constant.Hash_Company, st.CID, st); err != nil {
		ForwardEx(session, "1", nil, "New Company DirectWrite :"+err.Error())
		return
	}
	cudb := &CUser{}
	cudb.CID = st.CID
	cudb.HmCUser = make(map[string][]ShortInfo)
	if err := db.DirectWrite(constant.Hash_CRecordUser, cudb.CID, cudb); err != nil {
		ForwardEx(session, "1", nil, "New Company DirectWrite :"+err.Error())
		return
	}
	///添加公司账号
	if err := AddCompanyAccount(st.UserName, st.Password, st.FullName, st.CID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", st)
}

/*
修改公司信息
*/
func ModifyCompany(session *JsNet.StSession) {
	type modifyInfo struct {
		ShortName         string   //公司简称
		FullName          string   //公司全称
		UserName          string   //登录账号
		CID               string   //公司ID
		Password          string   //登陆密码
		UID               string   //对应微信端的UID
		OfficialPerson    string   //公司法人代表姓名
		RegistrationTime  string   //公司注册时间
		TaxpayerIdentity  string   //公司纳税人识别号
		RegisteredCapital string   //公司注册资金
		Companysize       string   //公司人员规模
		SubCompanyNum     int      //子公司的个数
		BusinessScope     string   //公司业务范围
		Profile           string   //公司简介
		OtherNotes        string   //公司其他补充
		LsPicture         []string //公司介绍照片
		WebSite           string   //公司网址
		Representative    string   //公司代表人姓名
		Age               int      //公司代表人年龄
		Sex               string   //公司代表人性别
		LandlineTel       string   //座机电话（021-385638827-8921）
		MobilePhone       string   //代表人电话
		City              string   //城市
		Province          string   //省份
		Country           string   //国家
		Addr              string   //地址
		Type              string   //商家类型
		HeadImageURL      string   //公司代表人头像路径
		Logo              string   //公司logo路径
		AccountID         string   //公司账户ID
		Lon               float64  //公司经度
		Lat               float64  //公司纬度
		EntityTime        string   //公司创建时间
		LsDistributionID  []string //子公司的ID列表
	}
	st := &modifyInfo{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "ModifyCompany GetPara :%s\n", err.Error())
		return
	}
	if st.CID == "" {
		ForwardEx(session, "1", nil, "ModifyCompany GetPara failed, CID is empty\n")
		return
	}

	company := &Company{}
	if err := db.Modify(constant.Hash_Company, st.CID, company, st); err != nil {
		ForwardEx(session, "1", nil, "Modify Company:%s  failed\n", st.CID)
		return
	}
	go AddCompanyAccount(st.UserName, st.Password, st.FullName, st.CID)
	Forward(session, "0", company)
}

/*
删除一个公司
*/
func DelCompany(session *JsNet.StSession) {
	type info struct {
		CID string
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" {
		ForwardEx(session, "1", nil, "DelCompany CID is empty\n")
		return
	}

	if err := db.HDel(constant.Hash_Company, st.CID); err != nil {
		ForwardEx(session, "1", nil, "Delete Company error\n")
		return
	}
	Forward(session, "0", nil)
}

func GetAllCompany(session *JsNet.StSession) {
	data, err := GetGlobalCompany()
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	com := []*Company{}
	for _, v := range data {
		d, e := GetCompanyInfo(v)
		if e == nil {
			com = append(com, d)
		}
	}
	Forward(session, "0", com)
}

func GetGlobalCompany() ([]string, error) {
	data := []string{}
	err := db.ShareLock(constant.Hash_Company, constant.KEY_Global_Company, &data)
	if err != nil {
		ErrorLog(err.Error())
	}
	return data, nil
}

/*
获取公司信息内部
*/

func GetCompanyInfo(CID string) (st *Company, e error) {
	data := &Company{}
	if CID == "" {
		return data, ErrorLog("GetCompanyInfo failed,CID is empty\n")
	}
	if err := db.ShareLock(constant.Hash_Company, CID, data); err != nil {
		return data, err
	}
	return data, nil
}

/*
获取一个公司信息
*/
func QueryCompany(session *JsNet.StSession) {
	type info struct {
		CID string
	}
	st := &info{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" {
		ForwardEx(session, "1", nil, "DelCompany CID is empty\n")
		return
	}

	data, err := GetCompanyInfo(st.CID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	Forward(session, "0", data)
}

/*
获取一个公司简短信息
*/
func QueryCompanyShortInfo(session *JsNet.StSession) {
	type info struct {
		CID string
	}
	type CShortInfo struct {
		CID       string
		ShortName string //公司简称
		FullName  string //公司全称
		Logo      string //公司logo路径
	}
	st := &info{}
	data := &Company{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" {
		ForwardEx(session, "1", nil, "DelCompany CID is empty\n")
		return
	}
	if err := db.ShareLock(constant.Hash_Company, st.CID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	shortInfo := &CShortInfo{}
	shortInfo.CID = data.CID
	shortInfo.ShortName = data.ShortName
	shortInfo.FullName = data.FullName
	shortInfo.Logo = data.Logo
	Forward(session, "0", shortInfo)
}

/*
添加一个网点索引到公司
*/
func appendDistribution(CID, DID string) (dis *CompanyDistribution, e error) {
	data := &CompanyDistribution{}
	if CID == "" || DID == "" {
		return nil, ErrorLog("add distribution failed,CID:%s,DID:%s", CID, DID)
	}
	err := db.WriteLock(constant.Hash_Company_Distribution, CID, data)
	if err != nil {
		data.CID = CID
		data.DistributionID = []string{}
		data.DistributionID = append(data.DistributionID, DID)
		db.DirectWrite(constant.Hash_Company_Distribution, CID, data)
		return data, nil
	}
	for _, v := range data.DistributionID {
		if v == DID {
			db.WriteBack(constant.Hash_Company_Distribution, CID, data)
			return data, nil
		}
	}
	data.DistributionID = append(data.DistributionID, DID)
	return data, db.WriteBack(constant.Hash_Company_Distribution, CID, data)
}

///获取一个公司的所有网点
func getCompanyDistribution(CID string) (*CompanyDistribution, error) {
	data := &CompanyDistribution{}
	err := db.ShareLock(constant.Hash_Company_Distribution, CID, data)
	return data, err
}

/*
添加一个活动到公司
*/
func AppendActivityToCompany(activityID, CompanyID string) (e error) {
	if activityID == "" || CompanyID == "" {
		return ErrorLog("ActivityID=%s,CompanyID=%s is not correct", activityID, CompanyID)
	}
	st := &Company{}
	if err := db.WriteLock(constant.Hash_Company, CompanyID, st); err != nil {
		return err
	}
	for _, v := range st.LsActivity {
		if v == activityID {
			db.WriteBack(constant.Hash_Company, CompanyID, st)
			return nil
		}
	}
	st.LsActivity = append(st.LsActivity, activityID)
	return db.WriteBack(constant.Hash_Company, CompanyID, st)
}

/*
查询一个公司下面所有的活动
*/
func GetAllActivityOfCompany(session *JsNet.StSession) {
	type info struct {
		CID string //公司ID
	}
	lsActivity := []*Activity{}
	st := &info{}
	dis := &Company{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "Get the net parameter error:%s\n", err.Error())
		return
	}
	if err := db.ShareLock(constant.Hash_Company, st.CID, dis); err != nil {
		ForwardEx(session, "1", nil, "Get the Company error:%s\n", err.Error())
		return
	}

	for i := len(dis.LsActivity); i > 0; i-- {
		activity, err := GetActivityInfo(dis.LsActivity[i-1])
		if err == nil {
			lsActivity = append(lsActivity, activity)
		}
	}
	ForwardEx(session, "0", lsActivity, "sucess")

}

/*
添加一个新的代金券
*/
func AppendCouponToCompany(CouponID, CompanyID string) (e error) {
	if CouponID == "" || CompanyID == "" {
		return ErrorLog("CouponID=%s,CompanyID=%s is not correct", CouponID, CompanyID)
	}
	st := &Company{}
	if err := db.WriteLock(constant.Hash_Company, CompanyID, st); err != nil {
		return err
	}
	for _, v := range st.LsCoupon {
		if v == CouponID {
			db.WriteBack(constant.Hash_Company, CompanyID, st)
			return nil
		}
	}
	st.LsCoupon = append(st.LsCoupon, CouponID)
	return db.WriteBack(constant.Hash_Company, CompanyID, st)
}

/*
添加一个新的问卷调查
*/
func AppendSAQuestionToCompany(QuestionID, CompanyID string) (e error) {
	if QuestionID == "" || CompanyID == "" {
		return ErrorLog("QuestionID=%s,CompanyID=%s is not correct", QuestionID, CompanyID)
	}
	st := &Company{}
	if err := db.WriteLock(constant.Hash_Company, CompanyID, st); err != nil {
		return err
	}
	for _, v := range st.LsASystemQuestion {
		if v == QuestionID {
			db.WriteBack(constant.Hash_Company, CompanyID, st)
			return nil
		}
	}
	st.LsASystemQuestion = append(st.LsASystemQuestion, QuestionID)
	return db.WriteBack(constant.Hash_Company, CompanyID, st)
}

//商家登陆
func CompanyLogin(session *JsNet.StSession) {
	type RD_Login struct {
		UserName string //用户名
		PassWord string //用户密码
	}
	st := RD_Login{}
	if err := session.GetPara(&st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	data := &CompanyAccount{}
	if err := db.ShareLock(constant.Hash_Company_Account, st.UserName, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if data.UserName != st.UserName || data.Password != st.PassWord {
		ForwardEx(session, "1", nil, "The User Name or the Password is not correct")
		return
	}
	// com, err := GetCompanyInfo(data.CID)
	// if err != nil {
	// 	ForwardEx(session, "1", nil, err.Error())
	// 	return
	// }
	Forward(session, "0", data.CID)
}

func AddCompanyAccount(UserName, Password, FullName, CID string) error {

	data := &CompanyAccount{}
	err := db.WriteLock(constant.Hash_Company_Account, UserName, data)

	data.FullName = FullName
	data.Password = Password
	if err != nil {
		data.CreatDate = CurTime()
		data.CID = CID
		data.UserName = UserName
		return db.DirectWrite(constant.Hash_Company_Account, UserName, data)
	}
	return db.WriteBack(constant.Hash_Company_Account, UserName, data)

}

///检查用户名是否重复
func CheckCompanyAccount(UserName string) bool {
	data := &CompanyAccount{}
	if err := db.ShareLock(constant.Hash_Company_Account, UserName, data); err != nil {
		return false
	}
	return true
}

/*
查询一个公司下面所有的代金券
*/
func GetAllCouponOfCompany(session *JsNet.StSession) {
	type info struct {
		CID string //公司ID
	}
	lsCoupon := []*Coupon{}
	st := &info{}
	dis := &Company{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, "Get the net parameter error:%s\n", err.Error())
		return
	}
	if err := db.ShareLock(constant.Hash_Company, st.CID, dis); err != nil {
		ForwardEx(session, "1", nil, "Get the Company error:%s\n", err.Error())
		return
	}
	for i := len(dis.LsCoupon); i > 0; i-- {
		coupon, err := GetCouponInfo(dis.LsCoupon[i-1])
		if err == nil {
			if coupon.Status != -1 { //删除的不展示
				lsCoupon = append(lsCoupon, coupon)
			}
		}
	}
	ForwardEx(session, "0", lsCoupon, "sucess")
}
/*
公司挂靠代理
*/
func CompanyRegisiterAgent(companyID,agentID string)(e error){
	//获取公司信息
	companyInfo:=&Company{}
	agentInfo:=AgentShortInfo{}
	err:=db.WriteLock(constant.Hash_Company,companyID,companyInfo)
	if err!=nil{
		return err
	}
	agentInfo.AgentID=agentID
	agentInfo.EntityTime=CurTime()
	if companyInfo.LsAgent==nil{
		companyInfo.LsAgent=[]AgentShortInfo{}
	}
	companyInfo.LsAgent=append(companyInfo.LsAgent,agentInfo)
	err=db.WriteBack(constant.Hash_Company,companyID,companyInfo)
	if err!=nil{
		return err
	}
	return nil
	//挂靠代理
}
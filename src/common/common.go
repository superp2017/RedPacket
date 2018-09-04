package common

import (
	"JsLib/JsNet"
)

func InitCommon() {
	activityInit()
	requirementInit()
	userInit()
	userActivityInit()
	quetionInit()
	orderInit()
	stageInit()
	systemQuestionInit()
	systeminit()
	genertalinit()
	userCouponInit()
	initTradingArea()
	iniSystemQuestion()
}

func activityInit() {
	JsNet.Http("/queryactivity", QueryActivity)                   //查询活动
	JsNet.Http("/releaseactivity", ReleaseActivity)               //发布(需求转化为活动)
	JsNet.Http("/getallactivity", GetAllActivities)               //获取所有的活动列表
	JsNet.Http("/gethomeactivity", GetHomeActivity)               //获取首页活动列表
	JsNet.Http("/activitybindquestion", ActivityBindQuestion)     //活动绑定问题
	JsNet.Http("/activitysettings", ActivitySettings)             //活动参数设置
	JsNet.Http("/queryactivityrdp", GetActivityRDP)               //获取相关的统计信息
	JsNet.Http("/delglobalactivity", DelGlobalActivity)           //删除全局的某一个活动
	JsNet.Http("/changeactivitysettings", ChangeActivitySettings) //更改活动的配置
	JsNet.Http("/changeactivitystatus", ChangeActivityStatus)     //更改活动的状态
	JsNet.Http("/modifyactivity", ModifyActivity)                 //修改活动信息
	JsNet.Http("/delactivity", DelActivity)                       //标记删除活动
	JsNet.Http("/addactivityvisit", AddActivityVisit)             //增加活动的访问量

}

func requirementInit() {
	JsNet.Http("/newrequirement", NewRequirement)         //新建需求
	JsNet.Http("/queryrequirement", QueryRequitement)     //查询需求
	JsNet.Http("/getallquirements", GetAllRequirements)   //获取所有的需求
	JsNet.Http("/modifyrequirement", ModifyRequirement)   //修改需求
	JsNet.Http("/deluserrequirement", DelUserRequirement) //删除用户某一个需求
}

func userInit() {
	JsNet.Http("/codegetsessionkey", CodeGetSessionKey)                  //小程序code获取opeid和unionid
	JsNet.Http("/wxaccesstoken2user", WxAccessToken2User)                //微信授权信息直接转系统用户信息
	JsNet.Http("/wx2user", WX2User)                                      //微信的用户信息直接转系统用户信息
	JsNet.Http("/queryuser", QueryUser)                                  //UID查询用户信息
	JsNet.Http("/queryuserformopenid", QueryUserFromOpenID)              //openid查询用户信息
	JsNet.Http("/queryuserfromunionid", QueryUserFromUnionID)            //UnionID查询用户信息
	JsNet.Http("/modifyuser", ModifyUser)                                //修改用户基本信息
	JsNet.Http("/getallceller", GetAllCeller)                            //获取所有的商家列表
	JsNet.Http("/relationuser2distribution", RelationUserToDistribution) //商家和网点进行关联
	JsNet.Http("/bindmobile", UserBindMobile)                            //用户绑定手机号
	JsNet.Http("/newfeedback", newFeedBack)                              //新建一条反馈
	JsNet.Http("/queryfeedback", queryFeedBack)                          //查询单个反馈信息
	JsNet.Http("/getglobalfeedback", getGlobalFeedBack)                  //查询所有的反馈信息

}

func userActivityInit() {
	JsNet.Http("/getredpacket", GetRedPacket)                 //抢红包
	JsNet.Http("/getselleractivitylist", GetUserActivityInfo) //获取用户相关的活动列表
	JsNet.Http("/getusermoneyrecord", GetUserMoneyRecord)     //获取用户资金变动记录
	JsNet.Http("/userwithdraw", UserWithdraw)                 //用户提现
}

func quetionInit() {
	JsNet.Http("/newquestion", NewQuestion)                   //新建问题
	JsNet.Http("/queryquestion", QueryQuertion)               //查询问题
	JsNet.Http("/mofifyquestion", ModifyQuestion)             //修改问题
	JsNet.Http("/getactivityquestion", GetActivityQuestion)   //获取活动关联的问题
	JsNet.Http("/bindactivityquestion", ActivityBindQuestion) //活动绑定问题
}

func orderInit() {
	JsNet.Http("/queryorder", QueryOrder)             //查询单个订单
	JsNet.Http("/getmoreorders", GetMoreOrders)       //获取多个订单信息
	JsNet.Http("/getuserorderlist", GetUserOrderList) //获取用户的订单列表
	JsNet.Http("/getallorder", GetGlobalOrders)       //获取全局的订单列表
	JsNet.Http("/submitorder", SubmitOrder)           //提交订单
	JsNet.Http("/paysuccess", PaySuccess)             //订单支付回调
	Init_wx_pay()
}

func stageInit() {
	JsNet.Http("/queryactivitystatistics", QueryActivityStatistics) //查询后台的统计信息
	JsNet.Http("/backpayactivity", BackPayActity)                   ////后台支付活动
	JsNet.Http("/backpayrequirement", BackPayRequirment)            //后台支付需求
}

func systemQuestionInit() {
	JsNet.Http("/newsystemquestion", NewSystemQuetion)            //新建系统问卷
	JsNet.Http("/queryusersysquestion", GetUserNetSystemQuestion) //获取用户的当前系统问卷
}

func systeminit() {
	JsNet.Http("/getversion", GetVersion) //获取当前版本
	JsNet.Http("/setversion", SetVersion) //设置当前版本
}

func genertalinit() {
	JsNet.Http("/newcompany", NewCompany) //注册一个公司账号
	JsNet.Http("/agentnewcompany", AgentNewCompany)
	JsNet.Http("/modifycompany", ModifyCompany)                                         //修改公司信息
	JsNet.Http("/delcompany", DelCompany)                                               //删除一个公司
	JsNet.Http("/querycompany", QueryCompany)                                           //获取一个公司信息
	JsNet.Http("/getallactivityofcompany", GetAllActivityOfCompany)                     //查询一个公司下面所有的活动
	JsNet.Http("/querycompanycustomer", QueryCompanyCustomer)                           //获取一个公司的用户信息列表
	JsNet.Http("/companylogin", CompanyLogin)                                           //商家登陆
	JsNet.Http("/newdistribution", NewDistribution)                                     //新建一个网点
	JsNet.Http("/modifydistribution", ModifyDistribution)                               //修改网点信息
	JsNet.Http("/deldistribution", DelDistribution)                                     //删除一个网点
	JsNet.Http("/querydistribution", QueryDistribution)                                 //获取一个网点信息
	JsNet.Http("/getallactivityofdistribution", GetAllActivityOfDistribution)           //查询一个网点下面所有的活动
	JsNet.Http("/querycompaydistribution", QueryCompayDistribution)                     //获取公司所有网点信息
	JsNet.Http("/querydistributioncustomer", QueryDistributionCustomer)                 //获取一个网点的用户信息列表
	JsNet.Http("/queryactivitydistributioncustomer", QueryActivityDistributionCustomer) //获取一个活动在一个网点的客户
	JsNet.Http("/newcoupon", NewCoupon)                                                 //添加一张金券
	JsNet.Http("/modifycoupon", ModifyCoupon)                                           //修改代金券信息
	JsNet.Http("/delcoupon", DelCoupon)                                                 //删除一个代金券
	JsNet.Http("/querycoupon", QueryCoupon)                                             //获取一个代金券信息
	JsNet.Http("/getallcouponofcompany", GetAllCouponOfCompany)                         //获取一个公司的所有代金券信息
	JsNet.Http("/newactivitydirectly", NewActivityDirectly)                             //直接新建一个活动
	JsNet.Http("/recordcuserinfo", RecordCUserInfo)                                     //录入客户档案
	JsNet.Http("/querycuserrecord", QueryCUserRecord)                                   //获取公司所有录入客户信息
	JsNet.Http("/querycompanyshortinfo", QueryCompanyShortInfo)                         //获取一个公司简短信息
	JsNet.Http("/queryactivitycustomer", QueryActivityCustomer)                         //查询活动相关的用户
	JsNet.Http("/querydisactivitycustomer", QueryDistributionCustomer)                  //查询活动网点相关的用户
	JsNet.Http("/querycompanyuserinfo", QueryCompanyUserInfo)                           //查询一个公司某个手机号的用户
	JsNet.Http("/getallcompany", GetAllCompany)                                         ///获取所有公司列表
	JsNet.Http("/regisiteruser", RegisiterUser)                                         //登记DeviceID UID
	JsNet.Http("/umeng", PostUmengNet)
	JsNet.Http("/uone", UOne)
	JsNet.Http("/umulti", UMulti)
	JsNet.Http("/uall", UAll)
	JsNet.Http("/agentlogin", AgentLogIn)           //代理登陆
	JsNet.Http("/binduser", BindUser)               //代理发展下线
	JsNet.Http("/agentnewcompany", AgentNewCompany) //代理创建公司
	JsNet.Http("/newagent", NewAgent)               //新建代理

}

func userCouponInit() {
	JsNet.Http("/appendcoupontouser", AppendCouponToUser)                     //添加代金券到用户
	JsNet.Http("/getusercouponlist", GetUserCoupon)                           //获取用户的代金券列表
	JsNet.Http("/findvalidcoupon", FindValidCoupon)                           //查找有效的代金券
	JsNet.Http("/usecoupon", UseCoupon)                                       //使用代金券
	JsNet.Http("/querycompanyuseedCoupon", QueryCompanyUseedCoupon)           //查询公司各网点使用过的代金券
	JsNet.Http("/querydistributionuseedCoupon", QueryDistributionUseedCoupon) //查询某个网点使用过的代金券
	JsNet.Http("/querycoupincustomer", QueryCoupinCustomer)                   //查询某个代金券的认领和使用情况
}

func initTradingArea() {
	JsNet.Http("/newtradingarea", NewTradingArea)               //新建一个商圈
	JsNet.Http("/querytradingarea", QueryTradingArea)           //查询商圈
	JsNet.Http("/modifytradingarea", ModifyTradingArea)         //修改商圈信息
	JsNet.Http("/searchtradingarea", SearchTradingArea)         //搜索商圈
	JsNet.Http("/deltardingarea", DelTardingArea)               //删除一个商圈
	JsNet.Http("/getglobaltags", GetGlobalTags)                 //获取已经存在的标签
	JsNet.Http("/getcompanytradingarea", GetCompanyTradingArea) //获取公司的商圈
	JsNet.Http("/getglobaltradingarea", GetGlobalArea)          //获取全局的商圈信息
}

func iniSystemQuestion() {
	JsNet.Http("/newactivitysystemquestion", NewActivitySystemQuestion)             //新建商家系统问题
	JsNet.Http("/queryasystemquertion", QueryASystemQuertion)                       //获取商家系统问题
	JsNet.Http("/newsystemquetion", NewSystemQuetion)                               //新建后台系统问题
	JsNet.Http("/querysystemquertion", QuerySystemQuertion)                         //获取后台系统问题
	JsNet.Http("/deletesystemquestion", DeleteSystemQuestion)                       //删除后台系统问题
	JsNet.Http("/querymultiasystemquertion", QueryMultiASystemQuertion)             //获取多个商家系统问题
	JsNet.Http("/getusernetsystemquestion", GetUserNetSystemQuestion)               //获取一个用户所有的系统问题
	JsNet.Http("/queryallcustomrize", QueryAllCustomrize)                           //获取系统所有的用户化列表
	JsNet.Http("/querydedicatecompanycustomrize", QueryDedicateCompanyCustomrize)   //获取一个公司所有的问题人群列表
	JsNet.Http("/querydedicateactivitycustomrize", QueryDedicateActivityCustomrize) //获取一个活动所有的问题人群列表
	JsNet.Http("/queryallcompanycustomrize", QueryAllCompanyCustomrize)             //获取所有公司的问题人群列表
	JsNet.Http("/restsys", Restsys)                                                 //获取所有公司的问题人群列表
	JsNet.Http("/queryallsystemquertion", QueryAllSystemQuertion)                   //获取所有系统消息
	JsNet.Http("/queryallcompanysystemaquestion", QueryAllCompanySystemAQuestion)   //获取所有公司系统消息

}

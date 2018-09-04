package constant

//全局配置
const (
	C_DB = "rdpdb"
)

//Key
const (
	Hash_User                 = "User"                //前端用户
	Hash_Activity             = "Activity"            //海报
	Hash_ActivityRDP          = "ActivityRDP"         //活动红包领取记录
	Hash_OpenID_UID           = "OpenID_UID"          //openID映射UID
	Hash_UnionID_UID          = "Union_UID"           //UnionID 映射UID
	Hash_Cell_UID             = "Cell_UID"            //cell-uid
	Hash_Order                = "Order"               //订单
	Hash_UserOrder            = "UserOrder"           //用户订单
	Hash_Requirement          = "Requirement"         //需求
	Hash_UserActivity         = "UserActivity"        //用户相关的活动表
	Hash_UserFromAgent        = "UserFromAgent"       //小B邀请用户统计
	Hash_ActivityQuestion     = "ActivityQuestion"    //活动问题
	Hash_SystemQuestion       = "SystemQuestion"      //平台活动问题
	Hash_ActivitySQuestion    = "ActivitySQuestion"   //活动的问卷调查问题
	Hash_Customerize          = "Customrize"          //用户特征表格
	Hash_ActivityCustomerize  = "ActivityCustomrize"  //用户特征表格
	Hash_SystemQuestionID     = "SystemQuestionID"    //系统ID列表
	Hash_UserSystemQuestion   = "UserSystemQuestion"  //用户分配的当前问题列表
	Hash_User_MoneyRecord     = "UserMoneyRRecord"    //用户资金变动记录
	Hash_Seller_Store         = "SellerStore"         //商家门店
	Hash_Company_Distribution = "CompanyDistribution" //卖家网点
	Hash_Distribution         = "Distribution"        //网点
	Hash_Company              = "Company"             //公司，商家
	Hash_DistributionU        = "DistributionU"       //网点，客户
	Hash_CompanyU             = "CompanyU"            //公司客户
	Hash_Coupon               = "Coupon"              //代金券
	Hash_UserCoupon           = "UserCoupon"          //用户的代金券
	Hash_ActivityDISU         = "ActivityDISU"        //活动+网点 组成的客户信息表格
	Hash_ActivityU            = "ActivityU"           //活动用户表格
	Hash_CRecordUser          = "CRecordUser"         //售后服务公司用户
	Hash_UserDesciption       = "UserDesciption"      //用户对公司的售后统计
	Hash_CUser                = "CUser"               //公司反馈用户列表
	Hash_MobileCUser          = "MobileCUser"         //手机用户信息列表
	Hash_CompanyCouponUsed    = "CompanyCouponUsed"   //公司抵用券的使用情况
	Hash_CouponCustomer       = "CouponCustomer"      //领券的人的信息
	Hash_Company_Account      = "Company_Account"     //公司，商家
	Hash_TradingArea          = "TradingArea"         //商圈
	Hash_TradingAreaTag       = "TradingAreaTag"      //商圈标签
	Hash_FeedBack             = "FeedBack"            //用户反馈
	Hash_UIDDeviceID       = "UIDDeviceID"      //UID DeviceID
	Hash_DeviceIDUID       = "DeviceIDUID"      //DeviceID UID
	Hash_Agent_Account     ="AgentAccount"
)

const (
	KEY_Global_Requirement = "GlobalRequirement"  //全局的需求key
	KEY_Global_Activity    = "GlobalActivity"     //全局的活动key
	KEY_Global_Seller      = "GlobalSeller"       //全局的所有的商家
	KEY_Global_Order       = "GlobalOrder"        //全局的订单
	KEY_Invalid_Order      = "InvalidOrder"       //全局的订单
	KEY_SYSQuestion_Normal = "SYSQuestion"        //全局的问题ID
	KEY_Global_Company     = "GlobalCompany"      //全局的公司
	KEY_CityTradingArea    = "CityTradingArea"    //分城市的商圈
	KEY_GlobalTags         = "GlobalTag"          //所有标签
	KEY_GlobalArea         = "GlobalArea"         //所有商圈id
	KEY_CompanyTradingArea = "CompanyTradingArea" //公司创建的商圈
	KEY_GlobalFeedBack     = "GlobalFeedBack"     //所有的用户反馈的id
)

//系统数据库
const (
	C_JUNSIEDB  = "junsiedb"
	C_CONSULTDB = "servicedb"
)

///网络返回字段
const (
	CT_Ret    = "Ret"
	CT_Msg    = "Msg"
	CT_Entity = "Entity"
)

const (
	Status_New      = "Require_New"  ///新建需求，不需要设计的情况
	Status_Paid     = "Require_Paid" //已经支付//已经支付海报费用
	Status_Finished = "Finished"     //审核通过，已经发布为活动
)

const (
	Status_FurtherProcessCharge   = "FurtherProcessCharge"   //需要设计
	Status_FurtherProcessUnCharge = "FurtherProcessUnCharge" //不需要设计
)

const (
	Status_WaitPay = "WaiPay"     //新建的活动,待支付，不可抢
	Status_Active  = "Active"     //充值完成，活动状态,可转发，可抢
	Status_Arrears = "Arrears"    //余额不足，欠费状态,没法抢红包
	Status_Stop    = "Stop"       //停止状态(后台强制停止)
	Status_Del     = "Status_Del" //删除活动（标记删除）
)

const (
	C_PAY_WAITPAY = "WaitPay"   //待支付
	C_PAY_REFUND  = "Refunded"  //已退款
	C_PAY_SUCCESS = "PaySucces" //支付成功
	C_PAY_INVALID = "INVALID"   //无效的订单

	C_TIMEAREA  = 8
	C_SENDLIMIT = 75
)

const (
	Question_Choise      = "QuestionChoise"
	Question_Input       = "QuestionInput"
	Question_ChoiseInput = "QuestionChoiseInput"
)
const (
	Arrange_CompanyAlone = "Company"      //公司独立部署（可能多点，可能单点，看经纬度长度)
	Arrange_Distribution = "Distribution" //分布网店部署，这个属于多点部署
)
const (
	MoneyArrange_Coordinate = "Coordinate"  //以经纬度的方式
	MoneyArrange_Division   = "Division"    //以地理位置的方式
	MoneyArrange_TradArea   = "TradingArea" //以商圈商圈
)

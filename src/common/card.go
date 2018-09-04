package common

// "card": {
//   "card_type": "GROUPON",
//   "groupon": {
//       "base_info": {
//           "logo_url": "http://mmbiz.qpic.cn/mmbiz/iaL1LJM1mF9aRKPZJkmG8xXhiaHqkKSVMMWeN3hLut7X7hicFNjakmx ibMLGWpXrEXB33367o7zHN0CwngnQY7zb7g/0",
//           "brand_name":"微信餐厅",
//           "code_type":"CODE_TYPE_TEXT",
//           "title": "132元双人火锅套餐",
//           "sub_title": "周末狂欢必备",
//           "color": "Color010",
//           "notice": "使用时向服务员出示此券",
//           "service_phone": "020-88888888",
//           "description": "不可与其他优惠同享\n如需团购券发票，请在消费时向商户提出\n店内均可使用，仅限堂食",
//           "date_info": {
//               "type": "DATE_TYPE_FIX_TERM",
//               "fixed_term": 15 ,
//               "fixed_begin_term": 0
//           },
//           "sku": {
//               "quantity": 500000
//           },
//           "get_limit": 3,
//           "use_custom_code": false,
//           "bind_openid": false,
//           "can_share": true,
//           "can_give_friend": true,
//           "location_id_list" : [123, 12321, 345345],
//           "custom_url_name": "立即使用",
//           "custom_url": "http://www.qq.com",
//           "custom_url_sub_title": "6个汉字tips",
//           "promotion_url_name": "更多优惠",
//           "promotion_url": "http://www.qq.com"
//},
//        	  "deal_detail": "以下锅底2选1（有菌王锅、麻辣锅、大骨锅、番茄锅、清补凉锅、酸 菜鱼锅可选）：\n大锅1份 12元\n小锅2份 16元 "
//    }
// }

type ST_DateInfo struct {
	Type      string `json:"type"`
	FixedTerm int    `json:"fixed_term"`
	BeginTerm int    `json:"fixed_begin_term"`
}

type ST_SKU struct {
	Quantity int `json:"quantity"`
}
type CardBaseInfo struct {
	Logo             string      `json:"logo_url"`
	CompanyName      string      `json:"brand_name"`
	CodeType         string      `json:"code_type"`
	Title            string      `json:"title"`
	SubTitle         string      `json:"sub_title"`
	Color            string      `json:"color"`
	Notice           string      `json:"notice"`
	Mobile           string      `json:"service_phone"`
	description      string      `json:"description"`
	DateInfo         ST_DateInfo `json:"date_info"`
	SKU              ST_SKU      `json:"sku"`
	GetLimit         int         `json:"get_limit"`
	IsUseCustomCode  bool        `json:""use_custom_code"`
	IsBindOpenid     bool        `json:""bind_openid"`
	CanShare         bool        `json:""can_share"`
	CanGiveFriend    bool        `json:""can_give_friend"`
	LocationIDs      []int       `json:""location_id_list"`
	CustomUrlName    string      `json:""custom_url_name"`
	CustomUrlSubName string      `json:""custom_url_sub_title"`
	CustomURL        string      `json:""custom_url"`
	PromotionUrlName string      `json:""promotion_url_name"`
	PromotionUrl     string      `json:""promotion_url"`
}

type CardGrupon struct {
	BaseInfo    CardBaseInfo `json:"base_info"`
	Deal_detail string       `json:"deal_detail"`
}

type Card struct {
	CouponID   string     //代金券ID
	CID        string     //公司ID
	Money      int        //代金券金额
	StartTime  string     //代金券起始日期
	StopTime   string     //代金券结束日期
	StartStamp int64      //开始日期时间戳
	StopStamp  int64      //结束日期时间戳
	CardType   string     `json:"card_type"`
	Grupon     CardGrupon `json:"groupon"`
	Status     int        `json:"status"` //券的状态  0:正常 1：过期 -1：删除
	EntityTime string     //创建日期
}

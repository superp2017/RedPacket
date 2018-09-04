package common

import (
	. "JsLib/JsLogger"
	"JsLib/JsNet"
	"constant"
	"db"
	"ider"
	. "util"
)

type PointF struct {
	Lat   float64 //纬度
	Lon   float64 //经度
	Range float64 //范围
}

type TardingArea struct {
	AreaID   string   //商圈id
	CID      string   //公司id
	AreaName string   //商圈名称
	City     string   //城市
	Province string   //省份
	Points   []PointF //坐标点
	Shape    int      //形状 0:多边形 1:圆
	Center   PointF   //中心点
	Tags     []string //标签
	Direct   int      //0:后台创建 1:商家创建
	Date     string   //创建日期
}

///新建一个商圈
func NewTradingArea(session *JsNet.StSession) {
	st := &TardingArea{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.AreaName == "" || st.City == "" {
		ForwardEx(session, "1", nil, "NewTradingArea param failed,AreaName=%s,City=%s,Tags lenght=%d",
			st.AreaName, st.City, len(st.Tags))
		return
	}
	if (st.Shape == 0 && len(st.Points) < 3) || (st.Shape == 1 && len(st.Points) < 1) {
		ForwardEx(session, "1", nil, "NewTradingArea param failed,Shape=%d,Points lenght=%d",
			st.Shape, len(st.Points))
		return
	}

	if st.Shape == 0 && st.Center.Lat == 0 && st.Center.Lon == 0 {
		ForwardEx(session, "1", nil,
			"NewTradingArea param failed,Shape=%d,Center=%v\n", st.Shape, st.Center)
		return
	}
	if st.Direct == 1 && st.CID == "" {
		ForwardEx(session, "1", nil, "NewTradingArea param failed,CID=%s,Direct=%d\n", st.CID, st.Direct)
		return
	}

	st.AreaID = ider.GenID()
	st.Date = CurTime()
	if err := db.DirectWrite(constant.Hash_TradingArea, st.AreaID, st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	//标签映射
	go TagMapArea(st.Tags, st.City, st.AreaID)
	go addGlobalTags(st.Tags)
	if st.CID != "" {
		go addCmpanyTradingArea(st.CID, st.AreaID)
	}
	Forward(session, "0", st)
}

//查询商圈信息
func QueryTradingArea(session *JsNet.StSession) {
	type INFO struct {
		AreaID string //商圈id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	data, err := getTradingArea(st.AreaID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

//根据标签、城市查询商圈
func SearchTradingArea(session *JsNet.StSession) {
	type INFO struct {
		City []string //城市
		Tags []string //标签列表
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if len(st.City) == 0 && len(st.Tags) == 0 {
		ForwardEx(session, "1", nil, "SearchTradingArea param failed,City:%v,Tags:%v\n", st.City, st.Tags)
		return
	}

	tmp := []string{}
	citylist := []string{}
	taglist := []string{}
	if len(st.City) != 0 {
		for _, v := range st.City {
			l := getCityArea(v)
			for _, v1 := range l {
				AppendUniqueString(&citylist, v1)
			}
		}

	}
	for _, v := range st.Tags {
		tt := getTradingAreaTag(v)
		for _, v := range tt {
			AppendUniqueString(&taglist, v)
		}
	}
	if len(st.City) != 0 && len(st.Tags) == 0 {
		tmp = citylist
	} else if len(st.City) == 0 && len(st.Tags) != 0 {
		tmp = taglist
	} else {
		for _, v := range citylist {
			for _, v1 := range taglist {
				if v == v1 {
					AppendUniqueString(&tmp, v)
				}
			}
		}
	}
	data := []*TardingArea{}
	for _, v := range tmp {
		t, err := getTradingArea(v)
		if err == nil {
			data = append(data, t)
		}
	}

	Forward(session, "0", data)
}

//修改商圈信息
func ModifyTradingArea(session *JsNet.StSession) {
	type INFO struct {
		CID      string   //公司id
		AreaID   string   //商圈id
		AreaName string   //商圈名称
		City     string   //城市
		Province string   //省份
		Points   []PointF //坐标点
		Tags     []string //标签
		Shape    int      //形状 0:多边形 1:圆
		Center   PointF   //中心点
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.AreaID == "" || st.AreaName == "" || st.City == "" {
		ForwardEx(session, "1", nil,
			"ModifyTradingArea param failed,AreaID=%s,AreaName:%s,City:%s,\n",
			st.AreaID, st.AreaName, st.City)
		return
	}

	if (st.Shape == 0 && len(st.Points) < 3) || (st.Shape == 1 && len(st.Points) < 1) {
		ForwardEx(session, "1", nil, "ModifyTradingArea param failed,Shape=%d,Points lenght=%d",
			st.Shape, len(st.Points))
		return
	}

	if st.Shape == 0 && st.Center.Lat == 0 && st.Center.Lon == 0 {
		ForwardEx(session, "1", nil,
			"ModifyTradingArea param failed,Shape=%d,Center=%v\n", st.Shape, st.Center)
		return
	}

	data := &TardingArea{}
	if err := db.WriteLock(constant.Hash_TradingArea, st.AreaID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if data.CID != "-1" && data.Direct == 1 && data.CID != st.CID {
		go db.WriteBack(constant.Hash_TradingArea, st.AreaID, data)
		ForwardEx(session, "1", nil, "ModifyTradingArea param failed,CID is empty\n")
		return
	}

	data.AreaName = st.AreaName
	data.Shape = st.Shape
	data.Center = st.Center
	data.Province = st.Province
	data.Points = st.Points
	if data.City != st.City {
		removeCityArea(data.City, st.AreaID)
		addCityArea(st.City, st.AreaID)
		data.City = st.City
	}
	go changeTardingAreaTags(data.Tags, st.Tags, st.AreaID)
	data.Tags = st.Tags
	if err := db.WriteBack(constant.Hash_TradingArea, st.AreaID, data); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	Forward(session, "0", data)
}

///删除商圈
func DelTardingArea(session *JsNet.StSession) {
	type INFO struct {
		AreaID string
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.AreaID == "" {
		ForwardEx(session, "1", nil, "DelTardingArea param is empty\n")
		return
	}

	data, err := getTradingArea(st.AreaID)
	if err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	if err := db.HDel(constant.Hash_TradingArea, st.AreaID); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}

	go TagUnMapArea(data.Tags, data.City, data.AreaID)

	Forward(session, "0", nil)
}

///获取全部的标签
func GetGlobalTags(session *JsNet.StSession) {
	Forward(session, "0", getGlobalTags())
}

///获取公司的商圈列表
func GetCompanyTradingArea(session *JsNet.StSession) {
	type INFO struct {
		CID string //公司id
	}
	st := &INFO{}
	if err := session.GetPara(st); err != nil {
		ForwardEx(session, "1", nil, err.Error())
		return
	}
	if st.CID == "" {
		ForwardEx(session, "1", nil, "GetCompanyTradingArea param failed,CID is empty\n")
		return
	}

	Forward(session, "0", getMoreTradingRrea(getCompanyArea(st.CID)))
}

//获取全局的商圈信息
func GetGlobalArea(session *JsNet.StSession) {
	Forward(session, "0", getMoreTradingRrea(getGlobalTArea()))
}

//增加公司的商圈
func addCmpanyTradingArea(CID, AreaID string) error {
	data := make(map[string][]string)
	err := db.WriteBack(constant.Hash_TradingArea, constant.KEY_CompanyTradingArea, &data)
	if v, ok := data[CID]; ok {
		AppendUniqueString(&v, AreaID)
		data[CID] = v
	} else {
		data[CID] = []string{AreaID}
	}
	if err != nil {
		return db.DirectWrite(constant.Hash_TradingArea, constant.KEY_CompanyTradingArea, &data)
	}
	return db.WriteBack(constant.Hash_TradingArea, constant.KEY_CompanyTradingArea, &data)
}

//添加全局的商圈
func addGlobalArea(AreaID string) error {
	data := []string{}
	err := db.WriteLock(constant.Hash_TradingArea, constant.KEY_GlobalArea, &data)
	AppendUniqueString(&data, AreaID)
	if err != nil {
		return db.DirectWrite(constant.Hash_TradingArea, constant.KEY_GlobalArea, &data)
	}
	return db.WriteBack(constant.Hash_TradingArea, constant.KEY_GlobalArea, &data)
}

//获取全局的商圈
func getGlobalTArea() []string {
	data := []string{}
	db.ShareLock(constant.Hash_TradingArea, constant.KEY_GlobalArea, &data)
	return data
}

//获取公司的商圈列表
func getCompanyArea(CID string) []string {
	data := make(map[string][]string)
	if err := db.ShareLock(constant.Hash_TradingArea, constant.KEY_CompanyTradingArea, &data); err != nil {
		return []string{}
	}
	l, _ := data[CID]
	return l
}

//获取商圈信息
func getTradingArea(AreaID string) (*TardingArea, error) {
	data := &TardingArea{}
	err := db.ShareLock(constant.Hash_TradingArea, AreaID, data)
	return data, err
}

///获取多个商圈信息
func getMoreTradingRrea(ids []string) []*TardingArea {
	data := []*TardingArea{}
	for _, v := range ids {
		a, err := getTradingArea(v)
		if err == nil {
			data = append(data, a)
		}
	}
	return data
}

//标签映射商圈
func TagMapArea(Tags []string, City, AreaID string) {
	if AreaID == "" {
		return
	}
	go addCityArea(City, AreaID)
	for _, v := range Tags {
		if v != "" {
			go appendTradingAreaTag(v, AreaID)
		}
	}
}

//取消商圈的映射
func TagUnMapArea(Tags []string, City, AreaID string) {
	if AreaID == "" {
		return
	}
	removeCityArea(City, AreaID)
	for _, v := range Tags {
		if v != "" {
			removeTradingAreaTag(v, AreaID)
		}
	}
}

///修改标签，更改对应的映射
func changeTardingAreaTags(old, cur []string, AreaID string) {
	for _, v := range old {
		exist := false
		for _, v1 := range cur {
			if v == v1 {
				exist = true
				break
			}
		}
		if !exist {
			go removeTradingAreaTag(v, AreaID)
		}
	}

	for _, v := range cur {
		exist := false
		for _, v1 := range old {
			if v == v1 {
				exist = true
				break
			}
		}
		if !exist {
			go appendTradingAreaTag(v, AreaID)
		}
	}
	return
}

//添加城市商圈id
func addCityArea(City, AreaID string) error {
	if City == "" || AreaID == "" {
		return ErrorLog("addCityArea failed,City%s,AreaID%s", City, AreaID)
	}

	data := make(map[string][]string)
	err := db.WriteLock(constant.Hash_TradingArea, constant.KEY_CityTradingArea, &data)

	if k, ok := data[City]; ok {
		AppendUniqueString(&k, AreaID)
		data[City] = k
	} else {
		data[City] = []string{AreaID}
	}
	if err != nil {
		return db.DirectWrite(constant.Hash_TradingArea, constant.KEY_CityTradingArea, &data)
	}
	return db.WriteBack(constant.Hash_TradingArea, constant.KEY_CityTradingArea, &data)
}

///获取城市商圈id
func getCityArea(City string) []string {
	list := []string{}
	if City == "" {
		return list
	}
	data := make(map[string][]string)
	if err := db.ShareLock(constant.Hash_TradingArea, constant.KEY_CityTradingArea, &data); err != nil {
		ErrorLog(err.Error())
		return list
	}

	list, _ = data[City]
	return list
}

//移除城市商圈id
func removeCityArea(City, AreaID string) error {
	if City == "" || AreaID == "" {
		return ErrorLog("removeCityArea failed,City%s,AreaID%s", City, AreaID)
	}
	data := make(map[string][]string)
	if err := db.WriteLock(constant.Hash_TradingArea, constant.KEY_CityTradingArea, &data); err != nil {
		return err
	}

	if k, ok := data[City]; ok {
		DelExistString(&k, AreaID)
		data[City] = k
	}
	return db.WriteBack(constant.Hash_TradingArea, constant.KEY_CityTradingArea, &data)
}

//添加一个商圈到某个标签
func appendTradingAreaTag(Tag, AreaID string) error {
	data := []string{}
	err := db.WriteLock(constant.Hash_TradingAreaTag, Tag, &data)
	AppendUniqueString(&data, AreaID)
	if err != nil {
		return db.DirectWrite(constant.Hash_TradingAreaTag, Tag, &data)
	}
	return db.WriteBack(constant.Hash_TradingAreaTag, Tag, &data)
}

//获取某个标签的所有商圈
func getTradingAreaTag(Tag string) []string {
	data := []string{}
	db.ShareLock(constant.Hash_TradingAreaTag, Tag, &data)
	return data
}

//移除某个标签下的商圈
func removeTradingAreaTag(Tag, AreaID string) error {
	data := []string{}
	if err := db.WriteLock(constant.Hash_TradingAreaTag, Tag, &data); err != nil {
		return err
	}
	DelExistString(&data, AreaID)
	return db.WriteBack(constant.Hash_TradingAreaTag, Tag, &data)
}

//添加到全局的标签
func addGlobalTags(Tags []string) error {
	data := []string{}
	err := db.WriteLock(constant.Hash_TradingAreaTag, constant.KEY_GlobalTags, &data)

	for _, v := range Tags {
		exist := false
		for _, v1 := range data {
			if v == v1 {
				exist = true
				break
			}
		}
		if !exist {
			data = append(data, v)
		}
	}
	if err != nil {
		return db.DirectWrite(constant.Hash_TradingAreaTag, constant.KEY_GlobalTags, &data)
	}
	return db.WriteBack(constant.Hash_TradingAreaTag, constant.KEY_GlobalTags, &data)
}

///获取全局所有的标签
func getGlobalTags() []string {
	data := []string{}
	db.ShareLock(constant.Hash_TradingAreaTag, constant.KEY_GlobalTags, &data)
	return data
}
